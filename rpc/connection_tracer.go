package rpc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go/logging"
	"github.com/netsec-ethz/scion-apps/pkg/pan"
)

type ConnectionTracerMsg struct {
	Local, Remote                            *pan.UDPAddr
	SrcConnID, DestConnID                    *logging.ConnectionID
	Chosen                                   *logging.VersionNumber
	Versions, ClientVersions, ServerVersions []logging.VersionNumber
	ErrorMsg, Key, Value                     *string
	Parameters                               *logging.TransportParameters
	ByteCount, Cwnd                          *logging.ByteCount
	Packets                                  *int
	Header                                   *logging.Header
	ExtendedHeader                           *logging.ExtendedHeader
	Frames                                   []logging.Frame
	AckFrame                                 *logging.AckFrame
	PacketType                               *logging.PacketType
	DropReason                               *logging.PacketDropReason
	LossReason                               *logging.PacketLossReason
	EncryptionLevel                          *logging.EncryptionLevel
	PacketNumber                             *logging.PacketNumber
	CongestionState                          *logging.CongestionState
	PTOCount                                 *uint32
	Perspective                              *logging.Perspective
	Bool                                     *bool
	Generation                               *logging.KeyPhase
	TimerType                                *logging.TimerType
	Time                                     *time.Time
}

func non_nil_string(name string, i interface{}) string {
	if i != nil {
		return fmt.Sprintf("%s: %+v\n", name, i)
	}
	return ""
}

func (m *ConnectionTracerMsg) String() string {
	s := ""
	s += non_nil_string("Local", m.Local)
	s += non_nil_string("Remote", m.Remote)
	s += non_nil_string("Chosen", m.Chosen)
	s += non_nil_string("Versions", m.Versions)
	s += non_nil_string("ClientVersions", m.ClientVersions)
	s += non_nil_string("ServerVersions", m.ServerVersions)
	s += non_nil_string("ErrorMsg", m.ErrorMsg)
	s += non_nil_string("Parameters", m.Parameters)
	s += non_nil_string("ByteCount", m.ByteCount)
	s += non_nil_string("Cwnd", m.Cwnd)
	s += non_nil_string("Packets", m.Packets)
	s += non_nil_string("Header", m.Header)
	s += non_nil_string("ExtendedHeader", m.ExtendedHeader)
	s += non_nil_string("Frames", m.Frames)
	s += non_nil_string("AckFrame", m.AckFrame)
	s += non_nil_string("PacketType", m.PacketType)
	s += non_nil_string("DropReason", m.DropReason)
	s += non_nil_string("TimerType", m.TimerType)
	return s
}

type ConnectionTracerClient struct {
	rpc *Client
	l   *log.Logger
}

func NewConnectionTracerClient(client *Client, p logging.Perspective, odcid logging.ConnectionID) logging.ConnectionTracer {

	return &ConnectionTracerClient{client, client.l}
}

func (c *ConnectionTracerClient) StartedConnection(local, remote net.Addr, srcConnID, destConnID logging.ConnectionID) {
	c.l.Printf("StartedConnection")
	err := c.rpc.Call("ConnectionTracerServer.StartedConnection",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) NegotiatedVersion(chosen logging.VersionNumber, clientVersions, serverVersions []logging.VersionNumber) {
	c.l.Printf("NegotiatedVersion")
	err := c.rpc.Call("ConnectionTracerServer.NegotiatedVersion",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) ClosedConnection(e error) {
	c.l.Printf("ClosedConnection")
	s := e.Error()
	err := c.rpc.Call("ConnectionTracerServer.ClosedConnection",
		&ConnectionTracerMsg{
			ErrorMsg: &s,
		},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) SentTransportParameters(parameters *logging.TransportParameters) {
	c.l.Printf("SentTransportParameters")
	err := c.rpc.Call("ConnectionTracerServer.SentTransportParameters",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) ReceivedTransportParameters(parameters *logging.TransportParameters) {
	c.l.Printf("ReceivedTransportParameters")
	err := c.rpc.Call("ConnectionTracerServer.ReceivedTransportParameters",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) RestoredTransportParameters(parameters *logging.TransportParameters) {
	c.l.Printf("RestoredTransportParameters")
	err := c.rpc.Call("ConnectionTracerServer.RestoredTransportParameters",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) SentPacket(hdr *logging.ExtendedHeader, size logging.ByteCount, ack *logging.AckFrame, frames []logging.Frame) {
	c.l.Printf("SentPacket")
	err := c.rpc.Call("ConnectionTracerServer.SentPacket",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) ReceivedVersionNegotiationPacket(hdr *logging.Header, versions []logging.VersionNumber) {
	c.l.Printf("ReceivedVersionNegotiationPacket")
	err := c.rpc.Call("ConnectionTracerServer.ReceivedVersionNegotiationPacket",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) ReceivedRetry(hdr *logging.Header) {
	c.l.Printf("ReceivedRetry")
	err := c.rpc.Call("ConnectionTracerServer.ReceivedRetry",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) ReceivedPacket(hdr *logging.ExtendedHeader, size logging.ByteCount, frames []logging.Frame) {
	c.l.Printf("ReceivedPacket")
	err := c.rpc.Call("ConnectionTracerServer.ReceivedPacket",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) BufferedPacket(ptype logging.PacketType) {
	c.l.Printf("BufferedPacket")
	err := c.rpc.Call("ConnectionTracerServer.BufferedPacket",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) DroppedPacket(ptype logging.PacketType, size logging.ByteCount, reason logging.PacketDropReason) {
	c.l.Printf("DroppedPacket")
	err := c.rpc.Call("ConnectionTracerServer.DroppedPacket",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) UpdatedMetrics(rttStats *logging.RTTStats, cwnd, bytesInFlight logging.ByteCount, packetsInFlight int) {
	c.l.Printf("UpdatedMetrics")
	err := c.rpc.Call("ConnectionTracerServer.UpdatedMetrics",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) AcknowledgedPacket(logging.EncryptionLevel, logging.PacketNumber) {
	c.l.Printf("AcknowledgedPacket")
	err := c.rpc.Call("ConnectionTracerServer.AcknowledgedPacket",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) LostPacket(logging.EncryptionLevel, logging.PacketNumber, logging.PacketLossReason) {
	c.l.Printf("LostPacket")
	err := c.rpc.Call("ConnectionTracerServer.LostPacket",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) UpdatedCongestionState(logging.CongestionState) {
	c.l.Printf("UpdatedCongestionState")
	err := c.rpc.Call("ConnectionTracerServer.UpdatedCongestionState",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) UpdatedPTOCount(value uint32) {
	c.l.Printf("UpdatedPTOCount")
	err := c.rpc.Call("ConnectionTracerServer.UpdatedPTOCount",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) UpdatedKeyFromTLS(logging.EncryptionLevel, logging.Perspective) {
	c.l.Printf("UpdatedKeyFromTLS")
	err := c.rpc.Call("ConnectionTracerServer.UpdatedKeyFromTLS",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) UpdatedKey(generation logging.KeyPhase, remote bool) {
	c.l.Printf("UpdatedKey")
	err := c.rpc.Call("ConnectionTracerServer.UpdatedKey",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) DroppedEncryptionLevel(logging.EncryptionLevel) {
	c.l.Printf("DroppedEncryptionLevel")
	err := c.rpc.Call("ConnectionTracerServer.DroppedEncryptionLevel",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) DroppedKey(generation logging.KeyPhase) {
	c.l.Printf("DroppedKey")
	err := c.rpc.Call("ConnectionTracerServer.DroppedKey",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) SetLossTimer(logging.TimerType, logging.EncryptionLevel, time.Time) {
	c.l.Printf("SetLossTimer")
	err := c.rpc.Call("ConnectionTracerServer.SetLossTimer",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) LossTimerExpired(logging.TimerType, logging.EncryptionLevel) {
	c.l.Printf("LossTimerExpired")
	err := c.rpc.Call("ConnectionTracerServer.LossTimerExpired",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) LossTimerCanceled() {
	c.l.Printf("LossTimerCanceled")
	err := c.rpc.Call("ConnectionTracerServer.LossTimerCanceled",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) Close() {
	c.l.Printf("Close")
	err := c.rpc.Call("ConnectionTracerServer.Close",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) Debug(name, msg string) {
	c.l.Printf("Debug")
	err := c.rpc.Call("ConnectionTracerServer.Debug",
		&ConnectionTracerMsg{},
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}

type NilMsg struct{}

type ConnectionTracerServer struct {
	l      *log.Logger
	tracer logging.ConnectionTracer
}

func NewConnectionTracerServer(tracer logging.ConnectionTracer, l *log.Logger) *ConnectionTracerServer {
	return &ConnectionTracerServer{l, tracer}
}

func (c *ConnectionTracerServer) StartedConnection(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("StartedConnection called")
	if args.Local == nil || args.Remote == nil || args.SrcConnID == nil || args.DestConnID == nil {
		return ErrDeref
	}
	c.tracer.StartedConnection(*args.Local, *args.Remote, *args.SrcConnID, *args.DestConnID)
	return nil
}
func (c *ConnectionTracerServer) NegotiatedVersion(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("NegotiatedVersion called")
	if args.Chosen == nil {
		return ErrDeref
	}
	c.tracer.NegotiatedVersion(*args.Chosen, args.ClientVersions, args.ServerVersions)
	return nil
}
func (c *ConnectionTracerServer) ClosedConnection(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("ClosedConnection called")
	if args.ErrorMsg == nil {
		c.tracer.ClosedConnection(nil)
	} else {
		c.tracer.ClosedConnection(errors.New(*args.ErrorMsg))
	}
	return nil
}
func (c *ConnectionTracerServer) SentTransportParameters(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("SentTransportParameters called")
	c.tracer.SentTransportParameters(args.Parameters)
	return nil
}
func (c *ConnectionTracerServer) ReceivedTransportParameters(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("ReceivedTransportParameters called")
	c.tracer.ReceivedTransportParameters(args.Parameters)
	return nil
}
func (c *ConnectionTracerServer) RestoredTransportParameters(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("RestoredTransportParameters called")
	c.tracer.RestoredTransportParameters(args.Parameters)
	return nil
}
func (c *ConnectionTracerServer) SentPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("SentPacket called")
	if args.ByteCount == nil {
		return ErrDeref
	}
	c.tracer.SentPacket(args.ExtendedHeader, *args.ByteCount, args.AckFrame, args.Frames)
	return nil
}
func (c *ConnectionTracerServer) ReceivedVersionNegotiationPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("ReceivedVersionNegotiationPacket called")
	c.tracer.ReceivedVersionNegotiationPacket(args.Header, args.Versions)
	return nil
}
func (c *ConnectionTracerServer) ReceivedRetry(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("ReceivedRetry called")
	c.tracer.ReceivedRetry(args.Header)
	return nil
}
func (c *ConnectionTracerServer) ReceivedPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("ReceivedPacket called")
	if args.ByteCount == nil {
		return ErrDeref
	}
	c.tracer.ReceivedPacket(args.ExtendedHeader, *args.ByteCount, args.Frames)
	return nil
}
func (c *ConnectionTracerServer) BufferedPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("BufferedPacket called")
	if args.PacketType == nil {
		return ErrDeref
	}
	c.tracer.BufferedPacket(*args.PacketType)
	return nil
}
func (c *ConnectionTracerServer) DroppedPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("DroppedPacket called")
	if args.PacketType == nil || args.ByteCount == nil || args.DropReason == nil {
		return ErrDeref
	}
	c.tracer.DroppedPacket(*args.PacketType, *args.ByteCount, *args.DropReason)
	return nil
}
func (c *ConnectionTracerServer) UpdatedMetrics(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("UpdatedMetrics called")
	if args.Cwnd == nil || args.ByteCount == nil || args.Packets == nil {
		return ErrDeref
	}
	c.tracer.UpdatedMetrics(nil, *args.Cwnd, *args.ByteCount, *args.Packets)
	return nil
}
func (c *ConnectionTracerServer) AcknowledgedPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("AcknowledgedPacket called")
	if args.EncryptionLevel == nil || args.PacketNumber == nil {
		return ErrDeref
	}
	c.tracer.AcknowledgedPacket(*args.EncryptionLevel, *args.PacketNumber)
	return nil
}
func (c *ConnectionTracerServer) LostPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("LostPacket called")
	if args.EncryptionLevel == nil || args.PacketNumber == nil || args.LossReason == nil {
		return ErrDeref
	}
	c.tracer.LostPacket(*args.EncryptionLevel, *args.PacketNumber, *args.LossReason)
	return nil
}
func (c *ConnectionTracerServer) UpdatedCongestionState(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("UpdatedCongestionState called")
	if args.CongestionState == nil {
		return ErrDeref
	}
	c.tracer.UpdatedCongestionState(*args.CongestionState)
	return nil
}
func (c *ConnectionTracerServer) UpdatedPTOCount(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("UpdatedPTOCount called")
	if args.PTOCount == nil {
		return ErrDeref
	}
	c.tracer.UpdatedPTOCount(*args.PTOCount)
	return nil
}
func (c *ConnectionTracerServer) UpdatedKeyFromTLS(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("UpdatedKeyFromTLS called")
	if args.EncryptionLevel == nil || args.Perspective == nil {
		return ErrDeref
	}
	c.tracer.UpdatedKeyFromTLS(*args.EncryptionLevel, *args.Perspective)
	return nil
}
func (c *ConnectionTracerServer) UpdatedKey(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("UpdatedKey called")
	if args.Generation == nil || args.Bool == nil {
		return ErrDeref
	}
	c.tracer.UpdatedKey(*args.Generation, *args.Bool)
	return nil
}
func (c *ConnectionTracerServer) DroppedEncryptionLevel(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("DroppedEncryptionLevel called")
	if args.EncryptionLevel == nil {
		return ErrDeref
	}
	c.tracer.DroppedEncryptionLevel(*args.EncryptionLevel)
	return nil
}
func (c *ConnectionTracerServer) DroppedKey(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("DroppedKey called")
	if args.Generation == nil {
		return ErrDeref
	}
	c.tracer.DroppedKey(*args.Generation)
	return nil
}
func (c *ConnectionTracerServer) SetLossTimer(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("SetLossTimer called")
	if args.TimerType == nil || args.EncryptionLevel == nil || args.Time == nil {
		return ErrDeref
	}
	c.tracer.SetLossTimer(*args.TimerType, *args.EncryptionLevel, *args.Time)
	return nil
}
func (c *ConnectionTracerServer) LossTimerExpired(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("LossTimerExpired called")
	if args.TimerType == nil || args.EncryptionLevel == nil {
		return ErrDeref
	}
	c.tracer.LossTimerExpired(*args.TimerType, *args.EncryptionLevel)
	return nil
}
func (c *ConnectionTracerServer) LossTimerCanceled(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("LossTimerCanceled called")
	c.tracer.LossTimerCanceled()
	return nil
}
func (c *ConnectionTracerServer) Close(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("Close called")
	c.tracer.Close()
	return nil
}
func (c *ConnectionTracerServer) Debug(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("Debug called")
	if args.Key == nil || args.Value == nil {
		return ErrDeref
	}
	c.tracer.Debug(*args.Key, *args.Value)
	return nil
}
