package msgmodels

//
// CharacterDestroy represents the structure of a message that tells clients to destroy an existing
// instance of a character.
//
type CharacterDestroy struct {
	ClientID int // The client ID of the character instance to destroy.
}
