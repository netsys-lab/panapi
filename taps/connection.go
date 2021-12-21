package taps

import "github.com/netsys-lab/panapi/internal/enum"

// Connection
//
type Connection struct {
	Events chan Event
	State  enum.ConnectionState
	// Need to keep track of ConnPrio Property separately from the
	// other ConnectionProperties. Feedback welcome.
	ConnPrio uint
	cp       *ConnectionProperties
}

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

//
func (c *Connection) Send(messageData []byte) {
	c.SendContext(nil, messageData, true)
}

//
func (c *Connection) SendContext(messageContext *MessageContext, messageData []byte, endOfMessage bool) {

}

func (c *Connection) Receive() {

}

func (c *Connection) Close() {

}

func (c *Connection) Abort() {

}

/*


func (c *Connection) Clone() *ConnectionGroup {
	return &ConnectionGroup{}
}

type ConnectionGroup struct {
	// TODO
        }*/
