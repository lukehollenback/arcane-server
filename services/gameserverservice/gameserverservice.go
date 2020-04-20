package gameserverservice

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/lukehollenback/arcane-server/models"
	"github.com/lukehollenback/arcane-server/models/msgmodels"
	"github.com/lukehollenback/arcane-server/services/msghandlerservice"
	"github.com/lukehollenback/packet-server/tcp"
)

var (
	o    *GameServerService
	once sync.Once
)

//
// GameServerService represents an instance of the Game Server Service, which is responsible for
// communicated with game clients over TCP/IP and UDP protocols.
//
type GameServerService struct {
	mu          *sync.Mutex            // Mutex to protect against concurrent modification of the client table.
	config      *Config                // Structure with the service's configuration parameters.
	tcpServer   *tcp.Server            // Instance of a TCP/IP packet server used for interacting with clients.
	clients     map[int]*models.Client // Table of known connected clients.
	chHBKill    chan bool              // Channel that can be used to send a kill signal to the heartbeat watchdog goroutine.
	chHBStopped chan bool              // Channel upon which the heartbeat watchdog goroutine will send a signal upon completing its shut-down process.
}

//
// Config represents a struct of configuration settings for the Game Server Service.
//
type Config struct {
	TCPAddr                    string
	ClientHeartbeatTimeoutSecs int
}

//
// Instance provides a singleton instance of the message handler service.
//
func Instance() *GameServerService {
	once.Do(func() {
		o = &GameServerService{
			mu: &sync.Mutex{},
		}
	})

	return o
}

//
// Config allows for the Game Server Service to be configured. It is up to the caller to execute
// this method when the service is NOT running. Failing to do so may result in a corrupt program
// state.
//
func (o *GameServerService) Config(config *Config) {
	o.config = config
}

//
// Start fires up the Game Server Service in its own goroutine.
//
func (o *GameServerService) Start() (<-chan bool, error) {
	log.Printf("The Game Server Service is starting...")

	//
	// (Re)-initialize some of the service's structures.
	//
	o.clients = make(map[int]*models.Client, 0)
	o.chHBKill = make(chan bool)
	o.chHBStopped = make(chan bool)

	//
	// Create a new TCP server instance.
	//
	o.tcpServer = tcp.CreateServer(&tcp.ServerConfig{
		Address: o.config.TCPAddr,
		OnNewClient: func(tcpClient *tcp.Client) {
			//
			// Create a new client instance and add it to the service's client table.
			//
			client := models.CreateClient(tcpClient)

			o.addClient(client)
		},
		OnNewMessage: func(tcpClient *tcp.Client, msg string) {
			//
			// Locate the client in the client table and update its "last received message" timestamp.
			//
			client := o.clients[tcpClient.ID()]

			client.UpdateLastMsgTimestamp()

			//
			// Clean the message.
			//
			msg = strings.TrimRight(msg, "\n")

			//
			// Deserialize the message. If this fails, it is a bogus message.
			//
			// TODO: If too many bogus messages are recieved from the same client, we should kick that
			//  client off. Such a scenario could be a potential attack.
			//
			var m msgmodels.Msg

			unmarshallErr := json.Unmarshal([]byte(msg), &m)
			if unmarshallErr != nil {
				log.Printf(
					"%sAn error occured while attempting to unmarshal a recieved message. (Message: %s) "+
						"(Error: %s) (Hint: Are they sending bogus messages?)",
					tcpClient.RcvLogPrefix(),
					msg,
					unmarshallErr,
				)

				return
			}

			//
			// Log the unmarshalled messaged (in case we need to go back and debug something).
			//
			log.Printf("%s%+v", tcpClient.RcvLogPrefix(), m)

			//
			// Attempt to execute a registered handler for the message.
			//
			handlerErr := msghandlerservice.Instance().ExecuteMsgHandler(client, &m)
			if handlerErr != nil {
				log.Printf(
					"%sCould not handle message type. (Error: %s)",
					tcpClient.LogPrefix(),
					handlerErr,
				)
			}
		},
		OnClientConnectionClosed: func(tcpClient *tcp.Client) {
			o.forgetClient(tcpClient.ID())
		},
	})

	//
	// Start listening for connections to the TCP/IP packet server (which will spin up its own
	// goroutine), and start watching conncted clients' heartbeats
	//
	chTCPServerStarted, err := o.tcpServer.Start()
	if err != nil {
		return nil, err
	}

	go o.monitorClientHeartbeats()

	//
	// Return the "started" channel from the TCP/IP server because, in this case, that is the only
	// concurrent process that we might be waiting on for start-up to complete
	//
	return chTCPServerStarted, nil
}

//
// Stop shuts down the Game Server Service.
//
func (o *GameServerService) Stop() (<-chan bool, error) {
	log.Printf("The Game Server Service is stopping...")

	//
	// Kill the heartbeat monitor goroutine. We must wait for it to gracefully stop.
	//
	o.chHBKill <- true

	<-o.chHBStopped

	//
	// Kill the TCP server. We must wait for it to finish gracefully shutdown.
	//
	chTCPServerStopped, err := o.tcpServer.Stop()
	if err != nil {
		return nil, err
	}

	//
	// Return the "stopped" channel from the TCP/IP server because, in this case, that is the only
	// concurrent process that we might be waiting on for shut-down to complete.
	//
	return chTCPServerStopped, nil
}

//
// SendAllMessage sends the provided message to all connected clients.
//
func (o *GameServerService) SendAllMessage(msg *msgmodels.Msg) {
	o.sendAllMessage(msg)
}

//
// SendMessage sends the provided message to the provided client.
//
func (o *GameServerService) SendMessage(client *models.Client, msg *msgmodels.Msg) {
	//
	// Serialize the message.
	//
	rawMsg, err := msg.JSON()
	if err != nil {
		log.Fatal("Failed to serialize message intended for client (%s) into JSON. (Message: %+v) "+
			"(Error: %s)", client.String(), msg, err)
	}

	//
	// Log the message.
	//
	log.Printf("%s%s", client.SndLogPrefix(), rawMsg)

	//
	// Add a delimiter to the end of the serialized message so that our network protocol knows how to
	// buffer it properly.
	//
	rawMsg = append(rawMsg, '\n')

	//
	// Fire off the message to the client.
	//
	client.TCPClient().SendBytes(rawMsg)
}

//
// SendAllMessage sends the provided message to all connected clients.
//
func (o *GameServerService) sendAllMessage(msg *msgmodels.Msg) {
	//
	// Serialize the message.
	//
	rawMsg, err := msg.JSON()
	if err != nil {
		log.Fatal("Failed to serialize message intended for all connected clients into JSON. "+
			"(Message: %+v) (Error: %s)", msg, err)
	}

	//
	// Log the message.
	//
	log.Printf("<~>           %-21s <~ %s", "All Connected Clients", rawMsg)

	//
	// Add a delimiter to the end of the serialized message so that our network protocol knows how to
	// buffer it properly.
	//
	rawMsg = append(rawMsg, '\n')

	//
	// Fire off the raw message to all connected clients.
	//
	o.tcpServer.SendBytesAll(rawMsg)
}

//
// Kick forcefully disconnects the specified client and sends a message to the game world stating
// the specified reason for the kick.
//
func (o *GameServerService) kick(client *models.Client, reason string) {
	//
	// Send a message to the world explaining that the client is being kicked.
	//
	msgData := &msgmodels.Chat{
		Author:  "Server",
		Content: fmt.Sprintf("Kicking player %s. (Reason: %s)", *client.Username(), reason),
		Color:   msgmodels.ChatColDef,
	}

	msg := msgmodels.CreateMsg(msgData)

	o.sendAllMessage(msg)

	//
	// Actually disconnect the client.
	//
	client.TCPClient().Close()
}

//
// monitorClientHeartbeats loops every 60 seconds and checks if there are any clients from which no
// message has been recieved within the last minute. Intended to be run in its own goroutine.
//
func (o *GameServerService) monitorClientHeartbeats() {
	log.Printf("Client heartbeat monitoring has started.")

	for cont := true; cont; {
		select {
		case <-o.chHBKill:
			cont = false
		case <-time.After(60 * time.Second):
			o.kickTimedOutClients()
		}
	}

	log.Printf("Client heartbeat monitoring has stopped.")

	o.chHBStopped <- true
}

//
// kickTimedOutClients forcefully disconnets any clients from which a message has not been recieved
// within the last minute.
//
func (o *GameServerService) kickTimedOutClients() {
	// NOTE: We must lock because we are going to scroll through, and possibly mutate, the client
	//  table. There may be multiple goroutines attempting to modify the client table around the same
	//  time that this is occurring.

	o.mu.Lock()
	defer o.mu.Unlock()

	cutoff := (time.Now().Add(-60 * time.Second))

	log.Printf("Checking for clients that have not beat their heart since before %s...", cutoff)

	for _, client := range o.clients {
		if client.LastMsgTimestamp().Before(cutoff) {
			o.kick(client, "no message recieved in the last minute")
		}
	}
}

//
// addClient adds the provided client to the client table.
//
func (o *GameServerService) addClient(client *models.Client) {
	// NOTE: We must lock because we are going to mutate the client table. Multiple goroutines may be
	//  attempting to do the same around the same time.

	o.mu.Lock()
	defer o.mu.Unlock()

	o.clients[client.TCPClient().ID()] = client
}

//
// forgetClient removes the provided client from the client table.
//
func (o *GameServerService) forgetClient(id int) {
	// NOTE: We must lock because we are going to mutate the client table. Multiple goroutines may be
	//  attempting to do the same around the same time.

	o.mu.Lock()
	defer o.mu.Unlock()

	delete(o.clients, id)
}
