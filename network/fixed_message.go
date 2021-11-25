package network

import (
	"errors"
	"fmt"
	"io"
	"net/textproto"
)

type FixedMessage struct {
	b          []byte
	readindex  int
	writeindex int
	Header     *textproto.MIMEHeader
}

func NewFixedMessageString(s string) *FixedMessage {
	b := []byte(s)
	h := make(textproto.MIMEHeader)
	return &FixedMessage{b, 0, len(b), &h}
}

func NewFixedMessage(size int) *FixedMessage {
	h := make(textproto.MIMEHeader)
	return &FixedMessage{make([]byte, size), 0, 0, &h}
}

func (m FixedMessage) String() string {
	return string(m.b)
}

//TODO, add proper locking and reusability of buffer
func (m FixedMessage) Read(p []byte) (i int, err error) {
	for i = m.readindex; i < len(m.b) && i < len(p); i++ {
		p[i] = m.b[i]
	}
	m.readindex = i
	if i == len(m.b) {
		err = io.EOF
	}
	return
}

// Write appends the contents of p to the message m, return error network.EOM when the capacity is met or exceeded
func (m FixedMessage) Write(p []byte) (i int, err error) {
	if m.writeindex == len(m.b) {
		errors.New("No space left in message")
	}
	for i = m.writeindex; i < len(m.b) && i < len(p); i++ {
		m.b[i] = p[i]
	}
	m.writeindex = i
	if m.writeindex == len(m.b) {
		err = EOM
	}
	return
}

// should the user call this?
// or should get called by the message before read?
// if user is going to call this name should be more obvious
// updateMessageWithMimeHeader maybe?

func (m FixedMessage) AddMIMEHeaderToMesaage() error {

	var s string
	for key, val := range *m.Header {
		// might need to add carriage return at
		// the end of each header line
		s = fmt.Sprintf("%s: %s\n", key, val)
	}

	b := []byte(s)
	b = append(m.b, b...)

	m.b = b

	return nil
}
