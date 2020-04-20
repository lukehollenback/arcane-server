package msgmodels

const (
	//
	// ChatColDef is the color of default (e.g. player-sent) messages.
	//
	ChatColDef = "default"

	//
	// ChatColGame is the color of game-related (e.g. NPC conversation) messages.
	//
	ChatColGame = "game"

	//
	// ChatColMod is the color of moderator-sent messages.
	//
	ChatColMod = "moderator"

	//
	// ChatColSvr is the color of server-sent messages.
	//
	ChatColSvr = "Server"
)

//
// Chat represents the data payload of a message holding a new item for each connected player's
// chat module.
//
type Chat struct {
	Author  string
	Content string
	Color   string
}