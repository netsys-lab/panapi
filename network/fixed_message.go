package network

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/textproto"
)

//TODOs
//1- implement get
//2- implement set
//3- implement add?
//use it with http or concurrent server
//implement EOF message type

type FixedMessage struct {
	b          []byte
	readindex  int
	writeindex int
	// header     *textproto.MIMEHeader
}

func NewFixedMessageString(s string) *FixedMessage {
	b := []byte(s)
	return &FixedMessage{b, 0, len(b)}
}

func NewFixedMessage(size int) *FixedMessage {
	return &FixedMessage{make([]byte, size), 0, 0}
}

func (m FixedMessage) String() string {
	return string(m.b)
}

//TODO, add proper locking and reusability of buffer
func (m FixedMessage) Read(p []byte) (n int, err error) {
	for n = m.readindex; n < len(m.b) && n < len(p); n++ {
		p[n] = m.b[n]
	}
	m.readindex = n
	if n == len(m.b) {
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

func (m *FixedMessage) SetHeader(header *textproto.MIMEHeader) {

	b := new(bytes.Buffer)

	for key, val := range *header {
		// might need to add carriage return at
		// the end of each header line
		fmt.Fprintf(b, "%s: %s\r\n", key, val[0])
	}
	fmt.Fprintf(b, "\r\n")

	m.clearHeader()

	m.b = append(b.Bytes(), m.b...)

}

func (m *FixedMessage) clearHeader() {
	headersEnd := 0

	for n := 0; n < len(m.b); n++ {
		if m.b[n] == '\n' && n+1 > len(m.b) && m.b[n+1] == '\n' {
			headersEnd = n + 1
		}
	}

	if headersEnd > 0 && headersEnd+1 > len(m.b) {
		m.b = m.b[headersEnd+1:]
	}
}

func (m *FixedMessage) GetHeader() (*textproto.MIMEHeader, error) {
	m.b = []byte("HTTP 1.1 200 OK \r\n content-type: text\r\n\r\nthis is the body")
	// h := extractHeaderBytes()
	bufReader := bufio.NewReader(m)
	mimeReader := textproto.NewReader(bufReader)

	header, err := mimeReader.ReadMIMEHeader()

	if err != nil {
		return nil, err
	}

	return &header, nil
}

// func (m *FixedMessage) extractHeaderBytes()

func (m *FixedMessage) ToHTTPMessage() {
	//TODO get http start-line from user and check for validity

	b := []byte("GET / HTTP/1.0\r\n")
	m.b = append(b, m.b...)

}

// should the user call this?
// or should get called by the message before read?
// if user is going to call this name should be more obvious
// updateMessageWithMimeHeader maybe?

// func (m FixedMessage) AddMIMEHeaderToMesaage() error {

// 	var s string
// 	for key, val := range *m.Header {
// 		// might need to add carriage return at
// 		// the end of each header line
// 		s = fmt.Sprintf("%s: %s\n", key, val)
// 	}

// 	b := []byte(s)
// 	b = append(m.b, b...)

// 	m.b = b

// 	return nil
// }
