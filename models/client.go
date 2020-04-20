package models

import (
	"fmt"
	"sync"
	"time"

	"github.com/lukehollenback/packet-server/tcp"
)

//
// Client represents a connected player.
//
type Client struct {
	mu            *sync.Mutex
	tcpClient     *tcp.Client
	authenticated bool
	username      string
	lastMsg       time.Time
}

//
// CreateClient constructes a new client structure instance (to represent a connected player) and
// returns a pointer to it.
//
func CreateClient(tcpClient *tcp.Client) *Client {
	client := &Client{
		mu:            &sync.Mutex{},
		tcpClient:     tcpClient,
		authenticated: false,
		username:      "Unknown",
		lastMsg:       time.Now(),
	}

	return client
}

//
// String returns a string explanation of the client.
//
func (o *Client) String() string {
	return fmt.Sprintf("username: %s, authed: %t, lastMsg: %s, tcpRemoteAddr: %s, tcpLocalAddr: %s",
		o.username, o.authenticated, o.lastMsg, o.TCPRemoteAddr(), o.TCPLocalAddr())
}

//
// TCPClient returns a pointer to the actual TCP/IP client object that can be used to communicate
// with the client.
//
func (o *Client) TCPClient() *tcp.Client {
	return o.tcpClient
}

//
// Username returns a pointer to the client's authenticated username.
//
func (o *Client) Username() *string {
	return &o.username
}

//
// TCPRemoteAddr returns the remote address string for the TCP connection to the client.
//
func (o *Client) TCPRemoteAddr() string {
	return o.tcpClient.RemoteAddr()
}

//
// TCPLocalAddr returns the local address string for the TCP connection to the client.
//
func (o *Client) TCPLocalAddr() string {
	return o.tcpClient.LocalAddr()
}

//
// LogPrefix generates a prefix string that can be used in log messages about the client.
//
func (o *Client) LogPrefix() string {
	return o.tcpClient.LogPrefix()
}

//
// SndLogPrefix generates a prefix string that can be used in log messages about messages sent to
// the client.
//
func (o *Client) SndLogPrefix() string {
	return o.tcpClient.SndLogPrefix()
}

//
// RcvLogPrefix generates a prefix string that can be used in log messages about messages recieved
// from the client.
//
func (o *Client) RcvLogPrefix() string {
	return o.tcpClient.RcvLogPrefix()
}

//
// LastMsgTimestamp returns the timestamp of the last time a message was recieved from the client.
// Can be used to check if the client is still connected and responding as expected.
//
func (o *Client) LastMsgTimestamp() time.Time {
	return o.lastMsg
}

//
// UpdateLastMsgTimestamp sets the timestamp of the client's last known recieved message to the
// current time.
//
func (o *Client) UpdateLastMsgTimestamp() {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.lastMsg = time.Now()
}
