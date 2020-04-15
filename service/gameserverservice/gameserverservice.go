package gameserverservice

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/lukehollenback/arcane-server/model"
	"github.com/lukehollenback/arcane-server/service/msghandlerservice"
	tcpserver "github.com/lukehollenback/tcp-server"
)

var (
	o    *GameServerService
	once sync.Once
)

//
// GameServerService represents an instance of the Game Server Service.
//
type GameServerService struct {
	mu        *sync.Mutex
	running   bool
	config    *Config
	tcpServer *tcpserver.Server
	clients   map[string]*model.Client
	hbStop    chan bool
	hbDone    chan bool
}

//
// Config represents a struct of configuration settings for the Game Server
// Service.
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
// Config allows for the Game Server Service to be configured. Must be called while the service is
// NOT running.
//
func (o *GameServerService) Config(config *Config) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.running {
		log.Fatal("Cannot configure the Game Server Service while it is running!")
	}

	o.config = config
}

//
// Start fires up the Game Server Service in its own goroutine.
//
func (o *GameServerService) Start() {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.running {
		log.Fatal("Cannot start the Game Server Service while it is already running!")
	}

	//
	// (Re)-initialize some of the service's structures.
	//
	o.clients = make(map[string]*model.Client, 0)
	o.hbStop = make(chan bool)
	o.hbDone = make(chan bool)

	//
	// Create a new TCP server instance.
	//
	o.tcpServer = tcpserver.New(o.config.TCPAddr)

	//
	// Configure event callbacks for the TCP server instance.
	//
	o.tcpServer.OnNewClient(func(tcpClient *tcpserver.Client) {
		// NOTE: We must obtain a lock because we will be modifying the client table.
		o.mu.Lock()
		defer o.mu.Unlock()

		//
		// Create the client and add it to the client table.
		//
		client := model.CreateClient(tcpClient)

		o.clients[tcpClient.RemoteAddr()] = client

		//
		// Log some debug info.
		//
		log.Printf("%sClient connection has been established.", client.LogPrefix())
	})

	o.tcpServer.OnClientConnectionClosed(func(tcpClient *tcpserver.Client, err error) {
		// NOTE: We must obtain a lock because we will be modifying the client table.
		o.mu.Lock()
		defer o.mu.Unlock()

		//
		// Grab the client from the lookup table so that the proper log prefix can be obtained.
		//
		logPrefix := o.clients[tcpClient.RemoteAddr()].LogPrefix()

		//
		// Remove the client from the client table.
		//
		delete(o.clients, tcpClient.RemoteAddr())

		//
		// Log some debug info.
		//
		log.Printf("%sClient connection has closed.", logPrefix)
	})

	o.tcpServer.OnNewMessage(func(tcpClient *tcpserver.Client, msg string) {
		//
		// Locate the client in the client table and update its "last received message" timestamp.
		//
		client := o.clients[tcpClient.RemoteAddr()]

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
		var m model.Msg

		unmarshallErr := json.Unmarshal([]byte(msg), &m)
		if unmarshallErr != nil {
			log.Printf("%sAn error occured while attempting to unmarshal a recieved "+
				"message. (Message: %s) (Error: %s) (Hint: Are they sending bogus messages?)",
				tcpClient.RcvLogPrefix(), msg, unmarshallErr)

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
			log.Printf("Could not handle message type. Kicking client.")
			tcpClient.Close()
		}
	})

	//
	// Start listening for connections to the TCP server, and start watching conncted clients'
	// heartbeats.
	//

	go o.tcpServer.Start()
	go o.monitorClientHeartbeats(o.hbStop, o.hbDone)

	//
	// Set the running sentinel.
	//
	o.running = true
}

//
// Stop shuts down the Game Server Service (thus causing its goroutine to also shut down).
//
func (o *GameServerService) Stop() {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.running {
		log.Fatal("Cannot stop the Game Server Service when it has not been already started!")
	}

	//
	// Kill the heartbeat monitor goroutine.
	//
	o.hbStop <- true

	<-o.hbDone

	//
	// Kill the TCP server.
	//
	o.tcpServer.Stop()

	//
	// Unset the running sentinel.
	//
	o.running = false
}

//
// SendAllMessage sends the provided message to all connected clients.
//
func (o *GameServerService) SendAllMessage(msg *model.Msg) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.sendAllMessage(msg)
}

//
// SendAllMessage sends the provided message to all connected clients.
//
func (o *GameServerService) sendAllMessage(msg *model.Msg) {
	o.verifyRunning("Cannot send message.")

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
	log.Printf("<~> %-21s <~ %s", "All Connected Clients", rawMsg)

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
// SendMessage sends the provided message to the provided client.
//
func (o *GameServerService) SendMessage(client *model.Client, msg *model.Msg) {
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

	client.TCPClient().SendBytes(rawMsg)
}

//
// Kick forcefully disconnects the specified client and sends a message to the game world stating
// the specified reason for the kick.
//
func (o *GameServerService) kick(client *model.Client, reason string) {
	//
	// Send a message to the world explaining that the client is being kicked.
	//
	msgData := &model.ChatMsg{
		Author:  "Server",
		Content: fmt.Sprintf("Kicking player %s. (Reason: %s)", *client.Username(), reason),
		Color:   model.ChatMsgColDef,
	}

	msg := model.CreateMsg(msgData)

	o.sendAllMessage(msg)

	//
	// Actually disconnect the client.
	//
	client.TCPClient().Close()
}

//
// verifyRunning ensures that the Game Server Service is running. If it is not running, a fatal will
// be invoked.
//
func (o *GameServerService) verifyRunning(errMsg string) {
	if !o.running {
		log.Fatalf("%s Game Server Service is not started!", errMsg)
	}
}

//
// monitorClientHeartbeats loops every 60 seconds and checks if there are any clients from which no
// message has been recieved within the last minute. Intended to be run in its own goroutine.
//
func (o *GameServerService) monitorClientHeartbeats(stop <-chan bool, done chan<- bool) {
	log.Printf("Client heartbeat monitoring has started.")

	monitor := true

	for monitor {
		select {
		case _ = <-stop:
			monitor = false
		case <-time.After(60 * time.Second):
			o.kickTimedOutClients()
		}
	}

	log.Printf("Client heartbeat monitoring has stopped.")

	done <- true
}

//
// kickTimedOutClients forcefully disconnets any clients from which a message has not been recieved
// within the last minute.
//
func (o *GameServerService) kickTimedOutClients() {
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
