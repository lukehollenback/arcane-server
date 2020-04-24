package msgmodels

//
// Disc represents the data payload of a "Disc"-type message, which indicates that the server is
// forcefully disconnecting a client.
//
type Disc struct {
	Reason string // A message explaining the reason for disconnecting the client.
}
