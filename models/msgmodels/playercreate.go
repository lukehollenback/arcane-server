package msgmodels

//
// CharacterCreate represents the structure of a message that tells clients to create a new instance
// of a character.
//
type CharacterCreate struct {
	Type     string // The name of the character object to create.
	ClientID int    // The client ID of the character instance. For use in things like client tables.
	X        int    // The initial horizontal location of the relevant character instance.
	Y        int    // The initial vertical location of the relevant character instance.
}
