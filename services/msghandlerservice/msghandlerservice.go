package msghandlerservice

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/lukehollenback/arcane-server/models"
	"github.com/lukehollenback/arcane-server/models/msgmodels"
)

var (
	o    *MsgHandlerService
	once sync.Once
)

//
// MsgHandlerService represents an instance of the message handler service.
//
type MsgHandlerService struct {
	handlers map[string]*registeredMsgHandler
}

//
// registeredMsgHandler represents a registered message handler.
//
type registeredMsgHandler struct {
	requiresAuth bool                                       // Whether or not the client must be authenticated in order for the message to be handled.
	callback     func(*models.Client, *msgmodels.Msg) error // The actual handler method to execute upon recieving the message.
}

//
// Instance provides a singleton instance of the message handler service.
//
func Instance() *MsgHandlerService {
	once.Do(func() {
		o = new(MsgHandlerService)
		o.handlers = make(map[string]*registeredMsgHandler)
	})

	return o
}

//
// RegisterMsgHandler registers a handler function to be executed when messages of the specified key
// are recieved from clients.
//
func (o *MsgHandlerService) RegisterMsgHandler(
	key string,
	requiresAuth bool,
	callback func(*models.Client, *msgmodels.Msg) error,
) {
	o.handlers[key] = &registeredMsgHandler{
		requiresAuth: requiresAuth,
		callback:     callback,
	}

	log.Printf("Registered new message handler for the message type key \"%s\".", key)
}

//
// ExecuteMsgHandler attempts to execute the appropriate registered handler function for the
// provided message.
//
func (o *MsgHandlerService) ExecuteMsgHandler(client *models.Client, msg *msgmodels.Msg) error {
	//
	// Attempt to retrieve the handler callback from the map of those that are registered.
	//
	handler, prs := o.handlers[msg.Key]

	if !prs {
		return fmt.Errorf("no message handler is known for the message type key \"%s\"", msg.Key)
	}

	//
	// Make sure the client has authenticated already if the handler requires it.
	//
	if handler.requiresAuth && !client.Authed() {
		return errors.New("message handling requires authentication")
	}

	//
	// Execute the handler callback.
	//
	return handler.callback(client, msg)
}
