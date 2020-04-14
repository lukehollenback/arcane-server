package model

//
// PingMsg represents the data payload of a "PingMsg"-type message.
//
type PingMsg struct {
	SentTime int64 // An epoch milliseconds timestamp of when the ping was fired off.
}
