package models

//
// NewChatMsg represents the data payload of a "NewChatMsg"-type message that can be sent by a
// client to indicate that they have created a new message.
//
type NewChatMsg struct {
	Content string // The content of the chat message.
}
