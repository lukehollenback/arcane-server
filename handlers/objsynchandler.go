package handlers

import (
  "github.com/lukehollenback/arcane-server/models"
  "github.com/lukehollenback/arcane-server/models/msgmodels"
  "github.com/lukehollenback/arcane-server/services/gameserverservice"
  "github.com/lukehollenback/arcane-server/services/msghandlerservice"
  "github.com/mitchellh/mapstructure"
  "reflect"
)

func init() {
  msghandlerservice.Instance().RegisterMsgHandler(
    reflect.TypeOf(new(msgmodels.ObjSync)).Elem().Name(),
    true,
    handleObjSync,
  )
}

//
// handle is intended to be registered with the Message Handler Service to be used to actually
// processes a received message.
//
func handleObjSync(client *models.Client, rcvMsg *msgmodels.Msg) error {
  //
  // Deserialize the data payload in the message and perform anti-cheat validation on the
  // synchronized variable's values.
  //
  data := new(msgmodels.ObjSync)

  if err := mapstructure.Decode(rcvMsg.Data, data); err != nil {
    return err
  }

  // TODO ~> Perform anti-cheat validation.s

  //
  // Fire off the synchronization message to all other connected clients.
  //
  // TODO ~> Localize this to only clients in the relevant area.
  //
  gameserverservice.Instance().SendAllMessage(rcvMsg, []int{client.TCPClient().ID()})

  return nil
}