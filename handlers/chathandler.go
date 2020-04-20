package handlers

import (
	"reflect"

	"github.com/lukehollenback/arcane-server/models"
	"github.com/lukehollenback/arcane-server/models/msgmodels"
	"github.com/lukehollenback/arcane-server/services/gameserverservice"
	"github.com/lukehollenback/arcane-server/services/msghandlerservice"
	"github.com/mitchellh/mapstructure"
)

func init() {
	msghandlerservice.Instance().RegisterMsgHandler(reflect.TypeOf(new(msgmodels.Chat)).Elem().Name(),
		handleChat)
}

//
// handle is intended to be registered with the Message Handler Service to be used to actually
// processes a recieved message.
//
func handleChat(client *models.Client, rcvMsg *msgmodels.Msg) error {
	//
	// Deserialize the data payload in the message.
	//
	rcvMsgData := new(msgmodels.Chat)

	mapstructure.Decode(rcvMsg.Data, rcvMsgData)

	//
	// Generate a "ChatMsg"-type message and send it to all connected players. If any fields were left
	// out of the recieved payload, attempt to default them to the best possible value.
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
