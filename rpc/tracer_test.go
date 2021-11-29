package rpc

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/lucas-clemente/quic-go/logging"
)

func TestTracerMsgEncoding(t *testing.T) {
	b := bytes.Buffer{}
	enc := gob.NewEncoder(&b)
	msg := TracerMsg{
		ID:           new(int),
		TracingID:    new(uint64),
		Perspective:  new(logging.Perspective),
		ConnectionID: new(logging.ConnectionID),
		//Addr:         new(pan.UDPAddr),
		Header:     new(logging.Header),
		ByteCount:  new(logging.ByteCount),
		Frames:     []logging.Frame{},
		PacketType: new(logging.PacketType),
		DropReason: new(logging.PacketDropReason),
	}
	err := enc.Encode(msg)
	if err != nil {
		t.Fatal(err)
	}

	id := IDMsg{new(int)}
	err = enc.Encode(&id)
	if err != nil {
		t.Fatal(err)
	}
	//t.Logf("%+v", msg)
}
