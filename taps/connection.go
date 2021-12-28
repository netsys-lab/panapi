package taps

import "github.com/netsys-lab/panapi/internal/enum"

// Connection
//
type Connection struct {
	State enum.ConnectionState
	// Need to keep track of ConnPrio Property separately from the
	// other ConnectionProperties. Feedback welcome.
	ConnPrio uint
	cp       *ConnectionProperties
}

/*
// GetProperties can be called at any time by the application to query ConnectionProperties
//
// Deprecated: Per https://go.dev/doc/effective_go#Getters, it is not
// idiomatic Go to put "Get" into a getter's name
func (c *Connection) GetProperties() *ConnectionProperties {
	return c.Properties()
}

// Properties can be called at any time by the application to query ConnectionProperties
//
// Deprecated: Property could instead simply be an exported Field of
// Connection c. Feedback welcome.
func (c *Connection) Properties() *ConnectionProperties {
	return c.cp
}

// SetProperty stores value for property, which is stripped of case
// and non-alphabetic characters before being matched against the
// (equally stripped) exported Field names of c. The type of value
// must be assignable to type of the targeted property Field,
// otherwise an error is returned.
//
// For the sake of respecting the TAPS (draft) spec as closely as
// possible, this function allows you to say:
//  err := c.SetProperty("groupConnLimit", 100)
//  if err != nil {
//    ... // handle runtime error
//  }
//
// In idiomatic Go, you would (and should) instead say:
//  c.Properties().GroupConnLimit = 100
//
// Deprecated: Use func c.SetProperty only if you must. Direct access
// of the underlying ConnectionProperties struct Fields is usually
// preferred. This function is implemented using reflection and
// dynamic string matching, which is inherently inefficient and prone
// to bugs triggered at runtime.
func (c *Connection) SetProperty(property string, value interface{}) error {
	// ConnPrio apparently is an oddball, because it is not shared
	// between all Connections of a group. See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-7.4
	// We need to keep track of it separately
	tmp := c.cp.ConnPrio
	err := set(c.cp, property, value)
	if c.cp.ConnPrio != tmp {
		// ConnPrio changed, but only for this Connection, not
		// the rest in the group
		c.ConnPrio = c.cp.ConnPrio
		c.cp.ConnPrio = tmp
	}
	return err
}
*/

// Send sends a message, blocks until sending has succeeded or
// returns an error if sending was not successful. This error
// represents either the "Expired" or "SendError" Events from
// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.2.2
func (c *Connection) Send(message *Message) error {
	return NotYetImplementendError
}

/*
// SendContext sends a message with a specific, optional,
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
func (c *Connection) SendContext(messageContext *MessageContext, messageData []byte, endOfMessage bool) error {
	return nil
}
*/

func (c *Connection) Receive() (message *Message, err error) {
	err = NotYetImplementendError
	return
}

func (c *Connection) Close() error {
	return NotYetImplementendError
}

/*// Ready blocks until a Connection created with Initiate() or
// InitiateWithSend() transitions to Established state.
func (c *Connection) Ready() error {
	return NotYetImplementendError
}
*/

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
