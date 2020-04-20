package msgmodels

//
// Ping represents the data payload of a "Ping"-type message.
//
type Ping struct {
	SentTime int64 // An epoch milliseconds timestamp of when the ping was fired off.
}
