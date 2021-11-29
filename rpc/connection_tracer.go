package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/logging"
	"github.com/netsec-ethz/scion-apps/pkg/pan"
)

type ConnectionTracerMsg struct {
	Local, Remote                            *pan.UDPAddr
	OdcID, SrcConnID, DestConnID             *logging.ConnectionID
	Chosen                                   logging.VersionNumber
	Versions, ClientVersions, ServerVersions []logging.VersionNumber
	ErrorMsg, Key, Value                     *string
	Parameters                               *logging.TransportParameters
	ByteCount, Cwnd                          logging.ByteCount
	Packets, ID                              int
	Header                                   *logging.Header
	ExtendedHeader                           *logging.ExtendedHeader
	Frames                                   []logging.Frame
	AckFrame                                 *logging.AckFrame
	PacketType                               logging.PacketType
	DropReason                               logging.PacketDropReason
	LossReason                               logging.PacketLossReason
	EncryptionLevel                          logging.EncryptionLevel
	PacketNumber                             logging.PacketNumber
	CongestionState                          logging.CongestionState
	PTOCount                                 uint32
	TracingID                                uint64
	Perspective                              logging.Perspective
	Bool                                     bool
	Generation                               logging.KeyPhase
	TimerType                                logging.TimerType
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
	s += non_nil_string("OdcID", m.OdcID)
	s += non_nil_string("Perspective", m.Perspective)
	s += non_nil_string("ID", m.ID)
	s += non_nil_string("TracingID", m.TracingID)
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
	s += non_nil_string("CongestionState", m.CongestionState)
	return s
}

type ConnectionTracerClient struct {
	rpc        *Client
	l          *log.Logger
	p          logging.Perspective
	odcid      logging.ConnectionID
	tracing_id uint64
}

func (c *ConnectionTracerClient) new_msg() *ConnectionTracerMsg {
	return &ConnectionTracerMsg{
		Perspective: c.p,
		OdcID:       &c.odcid,
		ID:          c.rpc.id,
		TracingID:   c.tracing_id,
	}
}

func NewConnectionTracerClient(client *Client, id uint64, p logging.Perspective, odcid logging.ConnectionID) logging.ConnectionTracer {
	client.l.Printf("NewConnectionTracerClient %v", odcid)
	err := client.Call("ConnectionTracerServer.NewTracerForConnection",
		&ConnectionTracerMsg{
			Perspective: p,
			OdcID:       &odcid,
			ID:          client.id,
			TracingID:   id,
		},
		&NilMsg{},
	)
	if err != nil {
		client.l.Fatalln(err)
	}

	return &ConnectionTracerClient{client, client.l, p, odcid, id}
}

func (c *ConnectionTracerClient) StartedConnection(local, remote net.Addr, srcConnID, destConnID logging.ConnectionID) {
	c.l.Printf("StartedConnection")
	msg := c.new_msg()
	l := local.(pan.UDPAddr)
	r := remote.(pan.UDPAddr)

	msg.Local = &l
	msg.Remote = &r
	msg.SrcConnID = &srcConnID
	msg.DestConnID = &destConnID

	err := c.rpc.Call("ConnectionTracerServer.StartedConnection",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) NegotiatedVersion(chosen logging.VersionNumber, clientVersions, serverVersions []logging.VersionNumber) {
	c.l.Printf("NegotiatedVersion")
	msg := c.new_msg()
	msg.Chosen = chosen
	msg.ClientVersions = clientVersions
	msg.ServerVersions = serverVersions

	err := c.rpc.Call("ConnectionTracerServer.NegotiatedVersion",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) ClosedConnection(e error) {
	c.l.Printf("ClosedConnection")
	s := e.Error()
	msg := c.new_msg()
	msg.ErrorMsg = &s
	err := c.rpc.Call("ConnectionTracerServer.ClosedConnection",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) SentTransportParameters(parameters *logging.TransportParameters) {
	c.l.Printf("SentTransportParameters")
	msg := c.new_msg()
	msg.Parameters = parameters
	err := c.rpc.Call("ConnectionTracerServer.SentTransportParameters",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) ReceivedTransportParameters(parameters *logging.TransportParameters) {
	c.l.Printf("ReceivedTransportParameters")
	msg := c.new_msg()
	msg.Parameters = parameters
	err := c.rpc.Call("ConnectionTracerServer.ReceivedTransportParameters",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) RestoredTransportParameters(parameters *logging.TransportParameters) {
	c.l.Printf("RestoredTransportParameters")
	msg := c.new_msg()
	msg.Parameters = parameters
	err := c.rpc.Call("ConnectionTracerServer.RestoredTransportParameters",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) SentPacket(hdr *logging.ExtendedHeader, size logging.ByteCount, ack *logging.AckFrame, frames []logging.Frame) {
	c.l.Printf("SentPacket")
	msg := c.new_msg()
	msg.ExtendedHeader = hdr
	msg.ByteCount = size
	msg.AckFrame = ack
	//msg.Frames = frames
	err := c.rpc.Call("ConnectionTracerServer.SentPacket",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) ReceivedVersionNegotiationPacket(hdr *logging.Header, versions []logging.VersionNumber) {
	c.l.Printf("ReceivedVersionNegotiationPacket")
	msg := c.new_msg()
	msg.Header = hdr
	msg.Versions = versions

	err := c.rpc.Call("ConnectionTracerServer.ReceivedVersionNegotiationPacket",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) ReceivedRetry(hdr *logging.Header) {
	c.l.Printf("ReceivedRetry")
	msg := c.new_msg()
	msg.Header = hdr
	err := c.rpc.Call("ConnectionTracerServer.ReceivedRetry",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) ReceivedPacket(hdr *logging.ExtendedHeader, size logging.ByteCount, frames []logging.Frame) {
	c.l.Printf("ReceivedPacket")
	msg := c.new_msg()
	msg.ExtendedHeader = hdr
	msg.ByteCount = size
	//msg.Frames = frames
	err := c.rpc.Call("ConnectionTracerServer.ReceivedPacket",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) BufferedPacket(ptype logging.PacketType) {
	c.l.Printf("BufferedPacket")
	msg := c.new_msg()
	msg.PacketType = ptype
	err := c.rpc.Call("ConnectionTracerServer.BufferedPacket",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) DroppedPacket(ptype logging.PacketType, size logging.ByteCount, reason logging.PacketDropReason) {
	c.l.Printf("DroppedPacket")
	msg := c.new_msg()
	msg.PacketType = ptype
	msg.ByteCount = size
	msg.DropReason = reason
	err := c.rpc.Call("ConnectionTracerServer.DroppedPacket",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) UpdatedMetrics(rttStats *logging.RTTStats, cwnd, bytesInFlight logging.ByteCount, packetsInFlight int) {
	c.l.Printf("UpdatedMetrics")
	msg := c.new_msg()
	msg.Cwnd = cwnd
	msg.ByteCount = bytesInFlight
	msg.Packets = packetsInFlight
	err := c.rpc.Call("ConnectionTracerServer.UpdatedMetrics",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) AcknowledgedPacket(level logging.EncryptionLevel, pnum logging.PacketNumber) {
	c.l.Printf("AcknowledgedPacket")
	msg := c.new_msg()
	msg.EncryptionLevel = level
	msg.PacketNumber = pnum
	err := c.rpc.Call("ConnectionTracerServer.AcknowledgedPacket",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) LostPacket(level logging.EncryptionLevel, pnum logging.PacketNumber, reason logging.PacketLossReason) {
	c.l.Printf("LostPacket")
	msg := c.new_msg()
	msg.EncryptionLevel = level
	msg.PacketNumber = pnum
	msg.LossReason = reason
	err := c.rpc.Call("ConnectionTracerServer.LostPacket",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) UpdatedCongestionState(state logging.CongestionState) {
	msg := c.new_msg()
	msg.CongestionState = state
	c.l.Printf("UpdatedCongestionState")
	err := c.rpc.Call("ConnectionTracerServer.UpdatedCongestionState",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) UpdatedPTOCount(value uint32) {
	c.l.Printf("UpdatedPTOCount")
	msg := c.new_msg()
	msg.PTOCount = value
	err := c.rpc.Call("ConnectionTracerServer.UpdatedPTOCount",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) UpdatedKeyFromTLS(level logging.EncryptionLevel, p logging.Perspective) {
	c.l.Printf("UpdatedKeyFromTLS")
	msg := c.new_msg()
	msg.EncryptionLevel = level
	msg.Perspective = p
	err := c.rpc.Call("ConnectionTracerServer.UpdatedKeyFromTLS",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) UpdatedKey(generation logging.KeyPhase, remote bool) {
	c.l.Printf("UpdatedKey")
	msg := c.new_msg()
	msg.Generation = generation
	msg.Bool = remote
	err := c.rpc.Call("ConnectionTracerServer.UpdatedKey",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) DroppedEncryptionLevel(level logging.EncryptionLevel) {
	c.l.Printf("DroppedEncryptionLevel")
	msg := c.new_msg()
	msg.EncryptionLevel = level
	err := c.rpc.Call("ConnectionTracerServer.DroppedEncryptionLevel",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) DroppedKey(generation logging.KeyPhase) {
	c.l.Printf("DroppedKey")
	msg := c.new_msg()
	msg.Generation = generation
	err := c.rpc.Call("ConnectionTracerServer.DroppedKey",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) SetLossTimer(ttype logging.TimerType, level logging.EncryptionLevel, t time.Time) {
	c.l.Printf("SetLossTimer")
	msg := c.new_msg()
	msg.TimerType = ttype
	msg.EncryptionLevel = level
	msg.Time = &t
	err := c.rpc.Call("ConnectionTracerServer.SetLossTimer",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) LossTimerExpired(ttype logging.TimerType, level logging.EncryptionLevel) {
	c.l.Printf("LossTimerExpired")
	msg := c.new_msg()
	msg.TimerType = ttype
	msg.EncryptionLevel = level
	err := c.rpc.Call("ConnectionTracerServer.LossTimerExpired",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) LossTimerCanceled() {
	c.l.Printf("LossTimerCanceled")
	msg := c.new_msg()
	err := c.rpc.Call("ConnectionTracerServer.LossTimerCanceled",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) Close() {
	c.l.Printf("Close")
	msg := c.new_msg()
	err := c.rpc.Call("ConnectionTracerServer.Close",
		msg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}
func (c *ConnectionTracerClient) Debug(name, msg string) {
	c.l.Printf("Debug")
	mesg := c.new_msg()
	mesg.Key = &name
	mesg.Value = &msg
	err := c.rpc.Call("ConnectionTracerServer.Debug",
		mesg,
		&NilMsg{},
	)
	if err != nil {
		c.l.Fatalln(err)
	}
}

type NilMsg struct{}

type ConnectionTracerServer struct {
	l        *log.Logger
	tracer   logging.Tracer
	ctracers map[uint64]logging.ConnectionTracer
}

func NewConnectionTracerServer(tracer logging.Tracer, l *log.Logger) *ConnectionTracerServer {
	return &ConnectionTracerServer{l, tracer, map[uint64]logging.ConnectionTracer{}}
}

func (c *ConnectionTracerServer) NewTracerForConnection(args *ConnectionTracerMsg, resp *NilMsg) error {
	if args.OdcID == nil {
		return ErrDeref
	}
	tracing_id := args.TracingID
	c.ctracers[tracing_id] = c.tracer.TracerForConnection(context.WithValue(context.Background(), quic.SessionTracingKey, tracing_id), args.Perspective, *args.OdcID)

	return nil
}

func (c *ConnectionTracerServer) StartedConnection(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("StartedConnection called")
	if args.Local == nil || args.Remote == nil || args.SrcConnID == nil || args.DestConnID == nil {
		return ErrDeref
	}
	tracing_id := args.TracingID
	c.ctracers[tracing_id].StartedConnection(args.Local, args.Remote, *args.SrcConnID, *args.DestConnID)
	return nil
}
func (c *ConnectionTracerServer) NegotiatedVersion(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("NegotiatedVersion called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].NegotiatedVersion(args.Chosen, args.ClientVersions, args.ServerVersions)
	return nil
}
func (c *ConnectionTracerServer) ClosedConnection(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("ClosedConnection called")
	tracing_id := args.TracingID

	if args.ErrorMsg == nil {
		c.ctracers[tracing_id].ClosedConnection(nil)
	} else {
		c.ctracers[tracing_id].ClosedConnection(errors.New(*args.ErrorMsg))
	}
	return nil
}
func (c *ConnectionTracerServer) SentTransportParameters(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("SentTransportParameters called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].SentTransportParameters(args.Parameters)
	return nil
}
func (c *ConnectionTracerServer) ReceivedTransportParameters(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("ReceivedTransportParameters called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].ReceivedTransportParameters(args.Parameters)
	return nil
}
func (c *ConnectionTracerServer) RestoredTransportParameters(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("RestoredTransportParameters called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].RestoredTransportParameters(args.Parameters)
	return nil
}
func (c *ConnectionTracerServer) SentPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("SentPacket called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].SentPacket(args.ExtendedHeader, args.ByteCount, args.AckFrame, args.Frames)
	return nil
}
func (c *ConnectionTracerServer) ReceivedVersionNegotiationPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("ReceivedVersionNegotiationPacket called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].ReceivedVersionNegotiationPacket(args.Header, args.Versions)
	return nil
}
func (c *ConnectionTracerServer) ReceivedRetry(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("ReceivedRetry called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].ReceivedRetry(args.Header)
	return nil
}
func (c *ConnectionTracerServer) ReceivedPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("ReceivedPacket called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].ReceivedPacket(args.ExtendedHeader, args.ByteCount, args.Frames)
	return nil
}
func (c *ConnectionTracerServer) BufferedPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("BufferedPacket called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].BufferedPacket(args.PacketType)
	return nil
}
func (c *ConnectionTracerServer) DroppedPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("DroppedPacket called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].DroppedPacket(args.PacketType, args.ByteCount, args.DropReason)
	return nil
}
func (c *ConnectionTracerServer) UpdatedMetrics(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("UpdatedMetrics called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].UpdatedMetrics(&logging.RTTStats{}, args.Cwnd, args.ByteCount, args.Packets)
	return nil
}
func (c *ConnectionTracerServer) AcknowledgedPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("AcknowledgedPacket called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].AcknowledgedPacket(args.EncryptionLevel, args.PacketNumber)
	return nil
}
func (c *ConnectionTracerServer) LostPacket(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("LostPacket called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].LostPacket(args.EncryptionLevel, args.PacketNumber, args.LossReason)
	return nil
}
func (c *ConnectionTracerServer) UpdatedCongestionState(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Printf("UpdatedCongestionState called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].UpdatedCongestionState(args.CongestionState)
	return nil
}
func (c *ConnectionTracerServer) UpdatedPTOCount(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("UpdatedPTOCount called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].UpdatedPTOCount(args.PTOCount)
	return nil
}
func (c *ConnectionTracerServer) UpdatedKeyFromTLS(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("UpdatedKeyFromTLS called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].UpdatedKeyFromTLS(args.EncryptionLevel, args.Perspective)
	return nil
}
func (c *ConnectionTracerServer) UpdatedKey(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("UpdatedKey called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].UpdatedKey(args.Generation, args.Bool)
	return nil
}
func (c *ConnectionTracerServer) DroppedEncryptionLevel(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("DroppedEncryptionLevel called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].DroppedEncryptionLevel(args.EncryptionLevel)
	return nil
}
func (c *ConnectionTracerServer) DroppedKey(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("DroppedKey called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].DroppedKey(args.Generation)
	return nil
}
func (c *ConnectionTracerServer) SetLossTimer(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("SetLossTimer called")
	if args.Time == nil {
		return ErrDeref
	}
	tracing_id := args.TracingID
	c.ctracers[tracing_id].SetLossTimer(args.TimerType, args.EncryptionLevel, *args.Time)
	return nil
}
func (c *ConnectionTracerServer) LossTimerExpired(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("LossTimerExpired called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].LossTimerExpired(args.TimerType, args.EncryptionLevel)
	return nil
}
func (c *ConnectionTracerServer) LossTimerCanceled(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("LossTimerCanceled called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].LossTimerCanceled()
	return nil
}
func (c *ConnectionTracerServer) Close(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("Close called")
	tracing_id := args.TracingID
	c.ctracers[tracing_id].Close()
	return nil
}
func (c *ConnectionTracerServer) Debug(args *ConnectionTracerMsg, resp *NilMsg) error {
	c.l.Println("Debug called")
	if args.Key == nil || args.Value == nil {
		return ErrDeref
	}
	tracing_id := args.TracingID
	c.ctracers[tracing_id].Debug(*args.Key, *args.Value)
	return nil
}
