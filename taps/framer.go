package taps

import (
	"bufio"
	"bytes"
	"io"
)

/*type MessageFramer struct {
}

func (f *MessageFramer) Start() (*Connection, error) {

}

func (f *MessageFramer) Stop() (*Connection, error) {

  }

type Framer interface {
	Encode(kv map[string]interface{}, dst io.Writer, src io.Reader) (written int64, err error)
	Decode(kv map[string]interface{}, dst io.Writer, src io.Reader) (written int64, err error)
        }*/

type FrameSender interface {
	SendFrame(messageData []byte, messageContext *MessageContext) error
}

type FrameReceiver interface {
	ReceiveFrame() (messageData []byte, messageContext *MessageContext, err error)
}

/*type SendFramer interface {
	FrameSender
	ChainedSender(lower FrameSender) SendFramer
}

type ReceiveFramer interface {
	FrameReceiver
	ChainedReceiver(lower FrameReceiver)
        }*/

type Framer interface {
	FrameSender
	FrameReceiver
}

type NewlineFramer struct {
	lower Framer
	r     *bufio.Reader
	ctx   *MessageContext
}

func NewNewlineFramer(lower Framer) *NewlineFramer {
	return &NewlineFramer{lower, nil, nil}
}

func (nf *NewlineFramer) SendFrame(messageData []byte, messageContext *MessageContext) error {
	r := bufio.NewReader(bytes.NewReader(messageData))
	for {
		bs, err := r.ReadBytes('\n')
		if err == nil || err == io.EOF {
			err2 := nf.lower.SendFrame(bs, messageContext)
			if err2 != nil {
				return err2
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (nf *NewlineFramer) ReceiveFrame() ([]byte, *MessageContext, error) {
	var (
		bs  []byte
		err error
	)
	if nf.r == nil {
		bs, nf.ctx, err = nf.lower.ReceiveFrame()
		nf.r = bufio.NewReader(bytes.NewReader(bs))
	}
	messageData, err2 := nf.r.ReadBytes('\n')
	if err2 != nil {
		nf.r = nil
		return messageData, nf.ctx, err2
	}
	return messageData, nf.ctx, err

}

/*
   //
type HTTPServerFramer struct {
	KeyValue map[string]string
}

func (h *HTTPServerFramer) Encode(kv map[string]interface{}, dst io.Writer, src io.Reader) (written int64, err error) {
	status, ok := kv["Statuscode"]
	if ok {
		// it's a response
		var statuscode, n int
		if statuscode, ok = status.(int); !ok {
			err = fmt.Errorf("Could not parse %v to int", status)
			return
		}
		n, err = dst.Write([]byte(fmt.Sprintf(
			"HTTP/1.1 %d %s\r\n",
			statuscode,
			http.StatusText(statuscode),
		)))

		for key, value := range kv {
			if value == status {
				continue
			}
			n, err = dst.Write([]byte(fmt.Sprintf(
				"%s: %v\r\n",
				key,
				value,
			)))

			written += int64(n)
			if err != nil {
				return
			}
		}
		n, err = dst.Write([]byte("\r\n"))
		written += int64(n)
		if err != nil {
			return
		}
		n2, err2 := io.Copy(dst, src)
		written += n2
		err = err2
	} else {
		err = fmt.Errorf("Statuscode for response not supplied")
	}
	return
}

func (h *HTTPServerFramer) Decode(kv map[string]interface{}, dst io.Writer, src io.Reader) (written int64, err error) {
	return
}
*/
