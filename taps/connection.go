package taps

import (
	"bytes"
	"net"

	"github.com/netsys-lab/panapi/internal/enum"
)

// Connection
//
type Connection struct {
	State enum.ConnectionState
	// Need to keep track of ConnPrio Property separately from the
	// other ConnectionProperties. Feedback welcome.
	ConnPrio uint
	cp       *ConnectionProperties
	conn     net.Conn
}

// SimpleSend sends a message, blocks until sending has succeeded or
// returns an error if sending was not successful. This error
// represents either the "Expired" or "SendError" Events from
// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.2.2
func (c *Connection) SimpleSend(messageData []byte) error {
	return c.Send(messageData, nil, true)
}

// Send sends a message with a specific, optional,
// messageContext, and an endOfMessage flag that indicates whether,
// for the purposes of the underlying transport, this message is
// complete and can be sent. When endOfMessage is set, a call to this
// function blocks until sending of the underlying data either
// succeeded or returns an error if sending was not successful. This
// error represents either the "Expired" or "SendError" Events from
// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.2.2
// If endOfMessage is unset, this function will return as soon as the
// data could be handed off to an underlying queue and return an error
// if not enough resources are available.
func (c *Connection) Send(messageData []byte, messageContext *MessageContext, endOfMessage bool) error {
	_, err := c.conn.Write(messageData)
	return err
}

// TODO
func (c *Connection) Receive(messageData []byte) (n int, messageContext *MessageContext, endOfMessage bool, err error) {
	n, err = c.conn.Read(messageData)
	if err != nil {
		return
	}
}

func (c *Connection) Close() error {
	return NotYetImplementendError
}

func (c *Connection) Abort() error {
	return NotYetImplementendError
}

// PathChange blocks until the underlying transport has signalled a Path Change.
func (c *Connection) PathChange() error {
	return NotYetImplementendError
}

// Blocks until the underlying transport signals an ICMP error related to the Connection c.
func (c *Connection) SoftError() error {
	return NotYetImplementendError
}

/*
// Sent blocks until a message is sent, returning the messageContext
// or an error if the message could not be delivered, e.g. when the message was not
func (c *Connection) Sent() (messageContext *MessageContext, err error) {
	return
}
*/

/*


func (c *Connection) Clone() *ConnectionGroup {
	return &ConnectionGroup{}
}

type ConnectionGroup struct {
	// TODO
        }*/
