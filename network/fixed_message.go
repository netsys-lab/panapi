package network

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"strings"
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
	header     *textproto.MIMEHeader
	httpHeader []byte
}

func NewFixedMessageString(s string) *FixedMessage {
	b := []byte(s)
	h := make(textproto.MIMEHeader)
	httpH := make([]byte, 0)
	return &FixedMessage{b, 0, len(b), &h, httpH}
}

func NewFixedMessage(size int) *FixedMessage {
	h := make(textproto.MIMEHeader)
	httpH := make([]byte, 0)
	return &FixedMessage{make([]byte, size), 0, 0, &h, httpH}
}

func (m FixedMessage) String() string {
	httpH := string(m.httpHeader)
	h := m.headerAsString()
	return httpH + h + string(m.b)
}

//TODO, add proper locking and reusability of buffer
func (m *FixedMessage) Read(p []byte) (n int, err error) {
	h := []byte(m.headerAsString())

	headers := append(m.httpHeader, h...)
	completeMessage := append(headers, m.b...)

	for n = m.readindex; n < len(completeMessage) && n < len(p); n++ {
		p[n] = completeMessage[n]
	}

	m.readindex = n
	if n == len(completeMessage) {
		err = io.EOF
	}
	return
}

// Write appends the contents of p to the message m, return error network.EOM when the capacity is met or exceeded
func (m *FixedMessage) Write(p []byte) (i int, err error) {
	//TODO should it write byte by byte?
	//Or should it write all then parse?
	//what to do with possible header reading erros
	httpHeader, mimeHeader, body, err := parseMessage(p)
	m.httpHeader = httpHeader
	if err == nil {
		m.header = &mimeHeader
	}

	if m.writeindex == len(m.b) {
		errors.New("No space left in message")
	}
	for i = m.writeindex; i < len(m.b) && i < len(body); i++ {
		m.b[i] = body[i]
	}
	m.writeindex = i
	if m.writeindex == len(m.b) {
		err = EOM
	}

	i = len(p)
	return
}

func (m *FixedMessage) SetHeader(header *textproto.MIMEHeader) {
	m.header = header
}

func (m *FixedMessage) GetHeader() *textproto.MIMEHeader {
	return m.header
}

func (m *FixedMessage) SetHttpHeader(header []byte) {
	m.httpHeader = header
}

func (m *FixedMessage) GetHttpHeader() []byte {
	return m.httpHeader
}

func (m FixedMessage) headerAsString() string {

	b := new(bytes.Buffer)

	//TODO handle multi-valued headers
	for key, val := range *m.header {
		fmt.Fprintf(b, "%s: %s\r\n", key, val[0])
	}
	fmt.Fprintf(b, "\r\n")

	return b.String()
}

func parseMessage(p []byte) (httpHeader []byte, mimeHeader textproto.MIMEHeader, body []byte, err error) {
	//TODO handle missing parts of the message like no mimeHeader, no body etc.
	httpHeaderAndMessage := strings.SplitAfterN(string(p), "\r\n", 2)
	if len(httpHeaderAndMessage) == 2 {
		httpHeader = []byte(httpHeaderAndMessage[0])
		mimeHeaderAndMessage := strings.SplitAfterN(httpHeaderAndMessage[1], "\r\n\r\n", 2)

		if len(mimeHeaderAndMessage) == 2 {
			body = []byte(mimeHeaderAndMessage[1])

			headerReader := strings.NewReader(mimeHeaderAndMessage[0])
			bufReader := bufio.NewReader(headerReader)
			mimeReader := textproto.NewReader(bufReader)

			mimeHeader, err = mimeReader.ReadMIMEHeader()

			if err != nil {
				return
			}
			return

		} else {
			body = []byte(mimeHeaderAndMessage[0])
		}
	}

	return
}
