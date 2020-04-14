package model

import (
	"encoding/json"
	"reflect"
)

//
// Msg represents a generic message that contains a command key (e.g. "ping", or "sendChatMessage")
// and a map of the remaining message payload for use by the appropriate handler implementation.
//
// NOTE: We intentionally make all members of this class public to help with both serialization and
//  with logging.
//
type Msg struct {
	Key  string
	Data interface{}
}

//
// CreateMsg constructs a new message instance that the message handler service understands.
//
func CreateMsg(data interface{}) *Msg {
	msg := &Msg{
		Key:  reflect.TypeOf(data).Elem().Name(),
		Data: data,
	}

	return msg
}

//
// JSON serializes the message.
//
func (o *Msg) JSON() ([]byte, error) {
	return json.Marshal(o)
}
