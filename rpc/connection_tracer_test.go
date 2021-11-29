package rpc

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/lucas-clemente/quic-go/logging"
	"github.com/netsec-ethz/scion-apps/pkg/pan"
)

func TestConnectionTracerMsgEncoding(t *testing.T) {
	b := bytes.Buffer{}
	enc := gob.NewEncoder(&b)
	msg := ConnectionTracerMsg{
		Local:           new(pan.UDPAddr),
		Remote:          new(pan.UDPAddr),
		SrcConnID:       new(logging.ConnectionID),
		DestConnID:      new(logging.ConnectionID),
		Chosen:          new(logging.VersionNumber),
		Versions:        []logging.VersionNumber{},
		ClientVersions:  []logging.VersionNumber{},
		ServerVersions:  []logging.VersionNumber{},
		ErrorMsg:        new(string),
		Parameters:      &logging.TransportParameters{},
		ByteCount:       new(logging.ByteCount),
		Cwnd:            new(logging.ByteCount),
		CongestionState: new(logging.CongestionState),
		Packets:         new(int),
		Header:          new(logging.Header),
		ExtendedHeader:  &logging.ExtendedHeader{},
		Frames:          []logging.Frame{},
		AckFrame:        new(logging.AckFrame),
		PacketType:      new(logging.PacketType),
		DropReason:      new(logging.PacketDropReason),
		LossReason:      new(logging.PacketLossReason),
		EncryptionLevel: new(logging.EncryptionLevel),
		PacketNumber:    new(logging.PacketNumber),
		TimerType:       new(logging.TimerType),
	}
	err := enc.Encode(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+s", msg.String())
}
