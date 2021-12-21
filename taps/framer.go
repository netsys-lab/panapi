package taps

import (
	"fmt"
	"io"
	"net/http"
)

type Framer interface {
	Encode(kv map[string]interface{}, dst io.Writer, src io.Reader) (written int64, err error)
	Decode(kv map[string]interface{}, dst io.Writer, src io.Reader) (written int64, err error)
}

//
type HTTPMessageFramer struct {
	KeyValue map[string]string
}

func (h *HTTPMessageFramer) Encode(kv map[string]interface{}, dst io.Writer, src io.Reader) (written int64, err error) {
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
	} else if _, ok := kv["Request"]; ok {
		// This is a request
		err = fmt.Errorf("Request not yet implemented")
	} else {
		err = fmt.Errorf("Neither Request nor Response key found")
	}
	return
}

func (h *HTTPMessageFramer) Decode(kv map[string]interface{}, dst io.Writer, src io.Reader) (written int64, err error) {
	return
}
