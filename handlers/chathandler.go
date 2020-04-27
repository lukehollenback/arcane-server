package handlers

import (
	"reflect"

	"github.com/lukehollenback/arcane-server/models"
	"github.com/lukehollenback/arcane-server/models/msgmodels"
	"github.com/lukehollenback/arcane-server/services/gameserverservice"
	"github.com/lukehollenback/arcane-server/services/msghandlerservice"
	"github.com/lukehollenback/arcane-server/services/playerinfoservice"
	"github.com/lukehollenback/arcane-server/util"
	"github.com/mitchellh/mapstructure"
)

func init() {
	msghandlerservice.Instance().RegisterMsgHandler(
		reflect.TypeOf(new(msgmodels.Chat)).Elem().Name(),
		true,
		handleChat,
	)
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
	// Generate a "ChatMsg"-type message and send it to all connected players. To prevent the ability
	// for any players to be weird and spoof their username, said field is always looked up – even if
	// it was provided. If a color was optionally provided, it will be used.
	//
	// TODO: Validate everything – content (for excessive whitespace, illegal characters, and so on),
	//  color (to be allowed according to the senders permissions), and so on.
	//
	sndMsgAuthor := playerinfoservice.Instance().GetUsername(client.AuthedID())
	sndMsgColor := util.GetStrVal(rcvMsgData.Color, msgmodels.ChatColDef)
	sndMsgData := &msgmodels.Chat{
		Author:  sndMsgAuthor,
		Content: rcvMsgData.Content,
		Color:   sndMsgColor,
	}

	sndMsg := msgmodels.CreateMsg(sndMsgData)

	gameserverservice.Instance().SendAllMessage(sndMsg, nil)

	return nil
}
