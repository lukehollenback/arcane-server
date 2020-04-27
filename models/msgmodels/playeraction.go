package msgmodels

//
// CharacterAction represents the structure of a message that tells clients to execute an action on
// behalf of an existing character.
//
type CharacterAction struct {
	ClientID int                    // The client ID of the character instance that should execute the action.
	Action   string                 // The name of the action to execute.
	Details  map[string]interface{} // Special details about the action to execute.
}
