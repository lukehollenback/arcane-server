package msgmodels

//
// Auth represents the data payload of a "Auth"-type message. Used when a client needs to tell the
// server that it has successfully authenticated itself over HTTPS and would like to provide the
// token that it recieved during that process.
//
type Auth struct {
	Token string // Token recieved during the client's HTTPS authentication handshake.
}
