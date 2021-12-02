package lua

import (
	"fmt"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go/logging"
	"github.com/netsys-lab/panapi/rpc"
	lua "github.com/yuin/gopher-lua"
)

func strhlpr(s fmt.Stringer) lua.LString {
	return lua.LString(fmt.Sprintf("%s", s))
}

func new_lua_parameters(p *logging.TransportParameters) *lua.LTable {
	t := new(lua.LTable)
	if p != nil {
		t.RawSetString("InitialMaxStreamDataBidiLocal", lua.LNumber(p.InitialMaxStreamDataBidiLocal))
		t.RawSetString("InitialMaxStreamDataBidiRemote", lua.LNumber(p.InitialMaxStreamDataBidiRemote))
		t.RawSetString("InitialMaxStreamDataUni", lua.LNumber(p.InitialMaxStreamDataUni))
		t.RawSetString("InitialMaxData", lua.LNumber(p.InitialMaxData))
		t.RawSetString("MaxAckDelay", lua.LNumber(p.MaxAckDelay))
		t.RawSetString("AckDelayExponent", lua.LNumber(p.AckDelayExponent))
		t.RawSetString("DisableActiveMigration", lua.LBool(p.DisableActiveMigration))
		t.RawSetString("MaxUDPPayloadSize", lua.LNumber(p.MaxUDPPayloadSize))
		t.RawSetString("MaxUniStreamNum", lua.LNumber(p.MaxUniStreamNum))
		t.RawSetString("MaxBidiStreamNum", lua.LNumber(p.MaxBidiStreamNum))
		t.RawSetString("MaxIdleTimeout", lua.LNumber(p.MaxIdleTimeout))
		a := p.PreferredAddress
		if a != nil {
			t.RawSetString("PreferredAddress", lua.LString(
				fmt.Sprintf(
					"IPv4: %s:%d, IPv6: %s:%d, ConnectionID: %s, Token: %x",
					a.IPv4, a.IPv4Port, a.IPv6, a.IPv6Port, a.ConnectionID, a.StatelessResetToken,
				),
			))
		}
		t.RawSetString("OriginalDestinationConnectionID", strhlpr(p.OriginalDestinationConnectionID))
		t.RawSetString("InitialSourceConnectionID", strhlpr(p.InitialSourceConnectionID))
		t.RawSetString("RetrySourceConnectionID", strhlpr(p.RetrySourceConnectionID))
		t.RawSetString("StatelessResetToken", lua.LString(fmt.Sprintf("%x", p.StatelessResetToken)))
		t.RawSetString("ActiveConnectionIDLimit", lua.LNumber(p.ActiveConnectionIDLimit))
		t.RawSetString("MaxDatagramFrameSize", lua.LNumber(p.MaxDatagramFrameSize))
		return t
	}
	return nil
}

func new_lua_rtt_stats(stats *rpc.RTTStats) *lua.LTable {
	t := new(lua.LTable)
	if stats != nil {
		t.RawSetString("LatestRTT", lua.LNumber(stats.LatestRTT.Seconds()))
		t.RawSetString("MaxAckDelay", lua.LNumber(stats.MaxAckDelay.Seconds()))
		t.RawSetString("MeanDeviation", lua.LNumber(stats.MeanDeviation.Seconds()))
		t.RawSetString("MinRTT", lua.LNumber(stats.MinRTT.Seconds()))
		t.RawSetString("PTO", lua.LNumber(stats.PTO.Seconds()))
		t.RawSetString("SmoothedRTT", lua.LNumber(stats.SmoothedRTT.Seconds()))
		return t
	}
	return nil
}

type Stats struct {
	*State
	mod *lua.LTable
}

func NewStats(state *State) rpc.ServerConnectionTracer {
	state.Lock()
	defer state.Unlock()
	mod := map[string]lua.LGFunction{}
	for _, fn := range []string{
		"TracerForConnection",
		"StartedConnection",
		"NegotiatedVersion",
		"ClosedConnection",
		"SentTransportParameters",
		"ReceivedTransportParameters",
		"RestoredTransportParameters",
		"SentPacket",
		"ReceivedVersionNegotiationPacket",
		"ReceivedRetry",
		"ReceivedPacket",
		"BufferedPacket",
		"DroppedPacket",
		"UpdatedMetrics",
		"AcknowledgedPacket",
		"LostPacket",
		"UpdatedCongestionState",
		"UpdatedPTOCount",
		"UpdatedKeyFromTLS",
		"UpdatedKey",
		"DroppedEncryptionLevel",
		"DroppedKey",
		"SetLossTimer",
		"LossTimerExpired",
		"LossTimerCanceled",
		"Debug",
	} {
		s := fmt.Sprintf("function %s not implemented in script", fn)
		mod[fn] = func(L *lua.LState) int {
			state.Logger.Println(s)
			return 0
		}
	}

	stats := state.RegisterModule("stats", mod).(*lua.LTable)
	return &Stats{state, stats}
}

func (s *Stats) TracerForConnection(tracing_id uint64, p logging.Perspective, odcid logging.ConnectionID) error {
	s.Printf("TracerForConnection")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("TracerForConnection"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		lua.LNumber(p),
		strhlpr(odcid),
	)
}
func (s *Stats) StartedConnection(tracing_id uint64, local, remote net.Addr, srcConnID, destConnID logging.ConnectionID) error {
	s.Printf("StartedConnection")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("StartedConnection"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		strhlpr(local),
		strhlpr(remote),
		strhlpr(srcConnID),
		strhlpr(destConnID),
	)

}
func (s *Stats) NegotiatedVersion(tracing_id uint64, chosen logging.VersionNumber, clientVersions, serverVersions []logging.VersionNumber) error {
	s.Printf("NegotiatedVersion")
	s.Lock()
	defer s.Unlock()
	var (
		c_vs = lua.LTable{}
		s_vs = lua.LTable{}
	)
	for _, v := range clientVersions {
		c_vs.Append(strhlpr(v))
	}
	for _, v := range serverVersions {
		s_vs.Append(strhlpr(v))
	}

	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("NegotiatedVersion"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		strhlpr(chosen),
		&c_vs,
		&s_vs,
	)

}
func (s *Stats) ClosedConnection(tracing_id uint64, err error) error {
	s.Printf("ClosedConnection")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("ClosedConnection"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		lua.LString(err.Error()),
	)

}
func (s *Stats) SentTransportParameters(tracing_id uint64, parameters *logging.TransportParameters) error {
	s.Printf("SentTransportParameters")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("SentTransportParameters"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		new_lua_parameters(parameters),
	)

}
func (s *Stats) ReceivedTransportParameters(tracing_id uint64, parameters *logging.TransportParameters) error {
	s.Printf("ReceivedTransportParameters")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("ReceivedTransportParameters"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		new_lua_parameters(parameters),
	)

}
func (s *Stats) RestoredTransportParameters(tracing_id uint64, parameters *logging.TransportParameters) error {
	s.Printf("RestoredTransportParameters")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("RestoredTransportParameters"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		new_lua_parameters(parameters),
	)

}
func (s *Stats) SentPacket(tracing_id uint64, hdr *logging.ExtendedHeader, size logging.ByteCount, ack *logging.AckFrame, frames []logging.Frame) error {
	s.Printf("SentPacket: only stub implementation")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("SentPacket"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
	)

}
func (s *Stats) ReceivedVersionNegotiationPacket(tracing_id uint64, hdr *logging.Header, versions []logging.VersionNumber) error {
	s.Printf("ReceivedVersionNegotiationPacket: only stub implementation")
	s.Lock()
	defer s.Unlock()
	var vs = new(lua.LTable)
	for _, v := range versions {
		vs.Append(strhlpr(v))
	}

	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("ReceivedVersionNegotiationPacket"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		vs,
	)

}
func (s *Stats) ReceivedRetry(tracing_id uint64, hdr *logging.Header) error {
	s.Printf("ReceivedRetry: only stub implementation")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("ReceivedRetry"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
	)

}
func (s *Stats) ReceivedPacket(tracing_id uint64, hdr *logging.ExtendedHeader, size logging.ByteCount, frames []logging.Frame) error {
	s.Printf("ReceivedPacket: only stub implementation")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("ReceivedPacket"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
	)

}
func (s *Stats) BufferedPacket(tracing_id uint64, ptype logging.PacketType) error {
	s.Printf("BufferedPacket")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("BufferedPacket"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		lua.LNumber(ptype),
	)

}
func (s *Stats) DroppedPacket(tracing_id uint64, ptype logging.PacketType, size logging.ByteCount, reason logging.PacketDropReason) error {
	s.Printf("DroppedPacket")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("DroppedPacket"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		lua.LNumber(ptype),
		lua.LNumber(size),
		lua.LNumber(reason),
	)

}
func (s *Stats) UpdatedMetrics(tracing_id uint64, rttStats *rpc.RTTStats, cwnd, bytesInFlight logging.ByteCount, packetsInFlight int) error {
	s.Printf("UpdatedMetrics")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("UpdatedMetrics"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		new_lua_rtt_stats(rttStats),
		lua.LNumber(cwnd),
		lua.LNumber(bytesInFlight),
		lua.LNumber(packetsInFlight),
	)

}
func (s *Stats) AcknowledgedPacket(tracing_id uint64, level logging.EncryptionLevel, num logging.PacketNumber) error {
	s.Printf("AcknowledgedPacket")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("AcknowledgedPacket"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		strhlpr(level),
		lua.LNumber(num),
	)

}
func (s *Stats) LostPacket(tracing_id uint64, level logging.EncryptionLevel, num logging.PacketNumber, reason logging.PacketLossReason) error {
	s.Printf("LostPacket")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("LostPacket"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		strhlpr(level),
		lua.LNumber(num),
		lua.LNumber(reason),
	)

}
func (s *Stats) UpdatedCongestionState(tracing_id uint64, state logging.CongestionState) error {
	s.Printf("UpdatedCongestionState")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("UpdatedCongestionState"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		lua.LNumber(state),
	)

}
func (s *Stats) UpdatedPTOCount(tracing_id uint64, value uint32) error {
	s.Printf("UpdatedPTOCount")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("UpdatedPTOCount"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		lua.LNumber(value),
	)

}
func (s *Stats) UpdatedKeyFromTLS(tracing_id uint64, level logging.EncryptionLevel, p logging.Perspective) error {
	s.Printf("UpdatedKeyFromTLS")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("UpdatedKeyFromTLS"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		strhlpr(level),
		lua.LNumber(p),
	)

}
func (s *Stats) UpdatedKey(tracing_id uint64, generation logging.KeyPhase, remote bool) error {
	s.Printf("UpdatedKey")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("UpdatedKey"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		lua.LNumber(generation),
		lua.LBool(remote),
	)

}
func (s *Stats) DroppedEncryptionLevel(tracing_id uint64, level logging.EncryptionLevel) error {
	s.Printf("DroppedEncryptionLevel")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("DroppedEncryptionLevel"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		strhlpr(level),
	)

}
func (s *Stats) DroppedKey(tracing_id uint64, generation logging.KeyPhase) error {
	s.Printf("DroppedKey")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("DroppedKey"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		lua.LNumber(generation),
	)

}
func (s *Stats) SetLossTimer(tracing_id uint64, ttype logging.TimerType, level logging.EncryptionLevel, t time.Time) error {
	s.Printf("SetLossTimer")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("SetLossTimer"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		lua.LNumber(ttype),
		strhlpr(level),
		strhlpr(t),
	)

}
func (s *Stats) LossTimerExpired(tracing_id uint64, ttype logging.TimerType, level logging.EncryptionLevel) error {
	s.Printf("LossTimerExpired")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("LossTimerExpired"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		lua.LNumber(ttype),
		strhlpr(level),
	)

}
func (s *Stats) LossTimerCanceled(tracing_id uint64) error {
	s.Printf("LossTimerCanceled")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("LossTimerCanceled"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
	)

}
func (s *Stats) Close(tracing_id uint64) error {
	s.Printf("Close")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("Close"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
	)

}
func (s *Stats) Debug(tracing_id uint64, name, msg string) error {
	s.Printf("Debug")
	s.Lock()
	defer s.Unlock()
	return s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("Debug"),
			NRet:    0,
			Protect: true,
		},
		lua.LNumber(tracing_id),
		lua.LString(name),
		lua.LString(msg),
	)

}
