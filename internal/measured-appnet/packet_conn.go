// Copyright 2021 Thorben Krüger (thorben.krueger@ovgu.de)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package measured_appnet

import (
	"log"
	"net"
	"time"

	"github.com/scionproto/scion/go/lib/serrors"
	"github.com/scionproto/scion/go/lib/slayers"
	"github.com/scionproto/scion/go/lib/snet"
)

func timeTrack(start time.Time, name string) {
	log.Printf("%s:  %s", name, time.Since(start))
}

type TimedSCIONPacketConn struct {
	conn        net.PacketConn
	scmpHandler snet.SCMPHandler
	timed       time.Time
}

type timeStore struct {
	cnt   int
	times []time.Time
}

func (c *TimedSCIONPacketConn) SetDeadline(d time.Time) error {
	return c.conn.SetDeadline(d)
}

func (c *TimedSCIONPacketConn) Close() error {
	return c.conn.Close()
}

func (c *TimedSCIONPacketConn) WriteTo(pkt *snet.Packet, ov *net.UDPAddr) error {
	defer timeTrack(time.Now(), "WriteTo")
	if err := pkt.Serialize(); err != nil {
		return serrors.WrapStr("serialize SCION packet", err)
	}

	_, err := c.conn.WriteTo(pkt.Bytes, ov)

	c.timed = time.Now()

	if err != nil {
		return serrors.WrapStr("Reliable socket write error", err)
	}
	return nil
}

func (c *TimedSCIONPacketConn) SetWriteDeadline(d time.Time) error {
	return c.conn.SetWriteDeadline(d)
}

func (c *TimedSCIONPacketConn) ReadFrom(pkt *snet.Packet, ov *net.UDPAddr) error {
	for {
		if err := c.readFrom(pkt, ov); err != nil {
			return err
		}

		log.Printf("RTT (decoded): %s", time.Since(c.timed))

		if scmp, ok := pkt.Payload.(snet.SCMPPayload); ok {
			if c.scmpHandler == nil {
				return serrors.New("scmp packet received, but no handler found",
					"type_code", slayers.CreateSCMPTypeCode(scmp.Type(), scmp.Code()),
					"src", pkt.Source)
			}
			if err := c.scmpHandler.Handle(pkt); err != nil {
				return err
			}
			continue
		}

		return nil
	}
}

func (c *TimedSCIONPacketConn) readFrom(pkt *snet.Packet, ov *net.UDPAddr) error {
	pkt.Prepare()
	n, lastHopNetAddr, err := c.conn.ReadFrom(pkt.Bytes)
	if err != nil {
		return serrors.WrapStr("Reliable socket read error", err)
	}

	pkt.Bytes = pkt.Bytes[:n]
	var lastHop *net.UDPAddr

	var ok bool
	lastHop, ok = lastHopNetAddr.(*net.UDPAddr)
	if !ok {
		return serrors.New("Invalid lastHop address Type",
			"Actual", lastHopNetAddr)
	}

	if err := pkt.Decode(); err != nil {
		return serrors.WrapStr("decoding packet", err)
	}

	if ov != nil {
		*ov = *lastHop
	}
	return nil
}

func (c *TimedSCIONPacketConn) SetReadDeadline(d time.Time) error {
	return c.conn.SetReadDeadline(d)
}
