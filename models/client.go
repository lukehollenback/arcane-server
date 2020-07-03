package models

import (
  "fmt"
  "github.com/google/uuid"
  "sync"
  "time"

  "github.com/lukehollenback/packet-server/tcp"
)

//
// Client represents a connected player.
//
type Client struct {
  mu        *sync.Mutex // Mutex to prevent concurrent modification issues when mutating struct members.
  tcpClient *tcp.Client // The actual TCP/IP packet server client instance that is interacting with the client.
  authed    bool        // Whether or not the client has successfully authenticated yet. Some message handlers will fail until this is true.
  authedID  string      // The Player ID that the client authenticated themselves to be.
  objectID  uuid.UUID   // The unique identifier for the object instance representing the client.
  lastMsg   time.Time   // Timestamp of when the last known message was received from the client.
}

//
// CreateClient constructs a new client structure instance (to represent a connected player) and
// returns a pointer to it.
//
func CreateClient(tcpClient *tcp.Client) *Client {
  client := &Client{
    mu:        &sync.Mutex{},
    tcpClient: tcpClient,
    authed:    false,
    authedID:  "Unknown",
    objectID:  uuid.New(),
    lastMsg:   time.Now(),
  }

  return client
}

//
// String returns a string explanation of the client.
//
func (o *Client) String() string {
  return fmt.Sprintf("username: %s, authed: %t, lastMsg: %s, tcpRemoteAddr: %s, tcpLocalAddr: %s",
    o.authedID, o.authed, o.lastMsg, o.TCPRemoteAddr(), o.TCPLocalAddr())
}

//
// TCPClient returns a pointer to the actual TCP/IP client object that can be used to communicate
// with the client.
//
func (o *Client) TCPClient() *tcp.Client {
  return o.tcpClient
}

//
// Authed returns whether or not the client has successfully authenticated yet.
//
func (o *Client) Authed() bool {
  return o.authed
}

//
// SetAuthed modifies the client's "authenticated" sentinel.
//
func (o *Client) SetAuthed(authed bool) {
  o.mu.Lock()
  defer o.mu.Unlock()

  o.authed = authed
}

//
// PlayerID returns a pointer to the client's authenticated Player ID. This is typically the unique
// identifier of the player's record in the database.
//
func (o *Client) PlayerID() string {
  return o.authedID
}

//
// SetPlayerID modifies the player ID against which client is "authenticated". This is typically the
// unique identifier of the player's record in the database.
//
func (o *Client) SetPlayerID(authedID string) {
  o.mu.Lock()
  defer o.mu.Unlock()

  o.authedID = authedID
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
// UpdateLastMsgTimestamp sets the timestamp of the client's last known received message to the
// current time.
//
func (o *Client) UpdateLastMsgTimestamp() {
  o.mu.Lock()
  defer o.mu.Unlock()

  o.lastMsg = time.Now()
}

//
// Implementation of Object.ObjectID().
//
func (o *Client) ObjectID() string {
  return o.objectID.String()
}

//
// Implementation of Object.AreaID().
//
func (o *Client) AreaID() string {
  return "SecretMountain"
}

//
// Implementation of Object.X().
//
func (o *Client) X() int {
  return 224
}

//
// Implementation of Object.Y().
//
func (o *Client) Y() int {
  return 160
}

//
// Implementation of Object.Depth().
//
func (o *Client) Depth() int {
  return 400
}
