package handlers

import (
	"log"
	"reflect"
	"time"

	"github.com/lukehollenback/arcane-server/models"
	"github.com/lukehollenback/arcane-server/models/msgmodels"
	"github.com/lukehollenback/arcane-server/service/gameserverservice"
	"github.com/lukehollenback/arcane-server/service/msghandlerservice"
	"github.com/mitchellh/mapstructure"
)

const nanosInMilli = 1000000

func init() {
	msghandlerservice.Instance().RegisterMsgHandler(reflect.TypeOf(new(msgmodels.Ping)).Elem().Name(),
		handlePingMsg)
}

//
// handlePingMsg is intended to be registered with the Message Handler Service to be used to
// actually processes a recieved message.
//
func handlePingMsg(client *models.Client, rcvMsg *msgmodels.Msg) error {
	//
	// Deserialize the data payload in the message.
	//
	rcvMsgData := &msgmodels.Ping{}

	mapstructure.Decode(rcvMsg.Data, rcvMsgData)

	//
	// Pull the timestamp out and turn it into a usable time struct.
	//
	rcvMsgDataSentTimeNanos := (rcvMsgData.SentTime * nanosInMilli)
	rcvMsgDataSentTime := time.Unix(0, rcvMsgDataSentTimeNanos)

	//
	// Generate a timestamp that can be used as the origination time of the pong we are going to send
	// back.
	//
	now := time.Now()

	//
	// Log some debug info.
	//
	log.Printf("%sRecieved ping that originated at %s.", client.LogPrefix(), rcvMsgDataSentTime)
	log.Printf("%sLatency appears to be %d nanoseconds.", client.LogPrefix(),
		(now.UnixNano() - rcvMsgDataSentTimeNanos))
	log.Printf("%sSending pong with origination of %s...", client.LogPrefix(), now)

	//
	// Generate and send a pong message back.
	//
	sndMsgData := &msgmodels.Ping{
		SentTime: (now.UnixNano() / nanosInMilli),
	}

	sndMsg := msgmodels.CreateMsg(sndMsgData)

	gameserverservice.Instance().SendMessage(client, sndMsg)

	return nil
}
