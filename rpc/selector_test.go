package rpc

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
)

func TestSelectorMsgEncoding(t *testing.T) {
	b := bytes.Buffer{}
	enc := gob.NewEncoder(&b)
	msg := SelectorMsg{
		Local:         new(pan.UDPAddr),
		Remote:        new(pan.UDPAddr),
		Fingerprint:   new(pan.PathFingerprint),
		PathInterface: new(pan.PathInterface),
		Paths:         []*Path{},
	}
	err := enc.Encode(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", msg)
}
