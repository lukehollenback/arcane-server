package handlers

import (
	"reflect"

	"github.com/lukehollenback/arcane-server/models"
	"github.com/lukehollenback/arcane-server/models/msgmodels"
	"github.com/lukehollenback/arcane-server/service/gameserverservice"
	"github.com/lukehollenback/arcane-server/service/msghandlerservice"
	"github.com/mitchellh/mapstructure"
)

func init() {
	msghandlerservice.Instance().RegisterMsgHandler(reflect.TypeOf(new(models.NewChatMsg)).Elem().Name(),
		handleNewChatMsg)
}

//
// handle is intended to be registered with the Message Handler Service to be used to actually
// processes a recieved message.
//
func handleNewChatMsg(client *models.Client, rcvMsg *msgmodels.Msg) error {
	//
	// Deserialize the data payload in the message.
	//
	rcvMsgData := new(models.NewChatMsg)

	mapstructure.Decode(rcvMsg.Data, rcvMsgData)

	//
	// Generate a "ChatMsg"-type message and send it to all connected players.
	//
	sndMsgData := &msgmodels.Chat{
		Author:  *client.Username(),
		Content: rcvMsgData.Content,
		Color:   msgmodels.ChatColDef,
	}

	sndMsg := msgmodels.CreateMsg(sndMsgData)

	gameserverservice.Instance().SendAllMessage(sndMsg)

	return nil
}
