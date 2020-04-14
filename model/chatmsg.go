package model

const (
	//
	// ChatMsgColDef is the color of default (e.g. player-sent) messages.
	//
	ChatMsgColDef = "default"

	//
	// ChatMsgColGame is the color of game-related (e.g. NPC conversation) messages.
	//
	ChatMsgColGame = "game"

	//
	// ChatMsgColMod is the color of moderator-sent messages.
	//
	ChatMsgColMod = "moderator"

	//
	// ChatMsgColSvr is the color of server-sent messages.
	//
	ChatMsgColSvr = "Server"
)

//
// ChatMsg represents the data payload of a message holding a new item for each connected player's
// chat module.
//
type ChatMsg struct {
	Author  string
	Content string
	Color   string
}
