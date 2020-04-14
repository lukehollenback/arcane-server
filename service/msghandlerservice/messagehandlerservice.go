package msghandlerservice

import (
	"fmt"
	"log"
	"sync"

	"github.com/lukehollenback/arcane-server/model"
)

var (
	o    *MessageHandlerService
	once sync.Once
)

//
// MessageHandlerService represents an instance of the message handler service.
//
type MessageHandlerService struct {
	handlers map[string]func(*model.Client, *model.Msg) error
}

//
// Instance provides a singleton instance of the message handler service.
//
func Instance() *MessageHandlerService {
	once.Do(func() {
		o = new(MessageHandlerService)
		o.handlers = make(map[string]func(*model.Client, *model.Msg) error)
	})

	return o
}

//
// RegisterMsgHandler registers a handler function to be executed when messages of the specified key
// are recieved from clients.
//
func (o *MessageHandlerService) RegisterMsgHandler(key string, callback func(*model.Client, *model.Msg) error) {
	o.handlers[key] = callback

	log.Printf("Registered new message handler for the message type key \"%s\".", key)
}

//
// ExecuteMsgHandler attempts to execute the appropriate registered handler function for the
// provided message.
//
func (o *MessageHandlerService) ExecuteMsgHandler(client *model.Client, msg *model.Msg) error {
	//
	// Attempt to retrieve the handler callback from the map of those that are registered.
	//
	handler, prs := o.handlers[msg.Key]

	if !prs {
		return fmt.Errorf("no message handler is known for the message type key \"%s\"", msg.Key)
	}

	//
	// Execute the handler callback.
	//
	return handler(client, msg)
}
