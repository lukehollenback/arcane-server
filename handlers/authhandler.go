package handlers

import (
	"fmt"
	"reflect"

	"github.com/lukehollenback/arcane-server/models"
	"github.com/lukehollenback/arcane-server/models/msgmodels"
	"github.com/lukehollenback/arcane-server/services/gameserverservice"
	"github.com/lukehollenback/arcane-server/services/msghandlerservice"
	"github.com/lukehollenback/arcane-server/services/playerinfoservice"
	"github.com/mitchellh/mapstructure"
)

func init() {
	msghandlerservice.Instance().RegisterMsgHandler(
		reflect.TypeOf(new(msgmodels.Auth)).Elem().Name(),
		false,
		handleAuth,
	)
}

//
// handleAuth is intended to be registered with the Message Handler Service to be used to
// actually processes a recieved message.
//
func handleAuth(client *models.Client, rcvMsg *msgmodels.Msg) error {
	//
	// Deserialize the data payload in the message.
	//
	rcvMsgData := &msgmodels.Auth{}

	mapstructure.Decode(rcvMsg.Data, rcvMsgData)

	//
	// Set the client's "authenticated" sentinel.
	//
	client.SetAuthed(true)

	//
	// Generate and send a pong message back.
	//
	username := playerinfoservice.Instance().GetUsername(client.AuthedID())
	content := fmt.Sprintf("Welcome, %s!", username)
	sndMsgData := &msgmodels.Chat{
		Author:  "Server",
		Content: content,
		Color:   msgmodels.ChatColSvr,
	}

	sndMsg := msgmodels.CreateMsg(sndMsgData)

	gameserverservice.Instance().SendAllMessage(sndMsg)

	return nil
}
