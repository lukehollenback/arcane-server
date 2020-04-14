package handler

import (
	"reflect"

	"github.com/lukehollenback/arcane-server/model"
	"github.com/lukehollenback/arcane-server/service/gameserverservice"
	"github.com/lukehollenback/arcane-server/service/msghandlerservice"
	"github.com/mitchellh/mapstructure"
)

func init() {
	msghandlerservice.Instance().RegisterMsgHandler(reflect.TypeOf(new(model.NewChatMsg)).Elem().Name(),
		handleNewChatMsg)
}

//
// handle is intended to be registered with the Message Handler Service to be used to actually
// processes a recieved message.
//
func handleNewChatMsg(client *model.Client, rcvMsg *model.Msg) error {
	//
	// Deserialize the data payload in the message.
	//
	rcvMsgData := new(model.NewChatMsg)

	mapstructure.Decode(rcvMsg.Data, rcvMsgData)

	//
	// Generate a "ChatMsg"-type message and send it to all connected players.
	//
	sndMsgData := &model.ChatMsg{
		Author:  *client.Username(),
		Content: rcvMsgData.Content,
		Color:   model.ChatMsgColDef,
	}

	sndMsg := model.CreateMsg(sndMsgData)

	gameserverservice.Instance().SendAllMessage(sndMsg)

	return nil
}
