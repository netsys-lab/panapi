// Copyright 2021 Clemens Beck (clemens.beck@ovgu.de)
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

package oracle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lucas-clemente/quic-go/logging"
	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsec-ethz/scion-apps/pkg/shttp"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type stats struct {
	Bandwidth float64 `json:"bandwidth"`
}

// ConnectionStats tracks lifetime attributes of a connection
type ConnectionStats struct {
	PacketsSent     int64
	PacketsReceived int64
	BytesSent       logging.ByteCount
	BytesReceived   logging.ByteCount
	Opened          *time.Time
	Closed          *time.Time
	ActivePath      *pan.Path
	Local           *pan.UDPAddr
	Remote          *pan.UDPAddr
}

func (cd *ConnectionStats) LifeTime() time.Duration {
	if cd.Opened == nil {
		return 0
	}
	if cd.Closed == nil {
		return time.Since(*cd.Opened)
	}
	return cd.Closed.Sub(*cd.Opened)
}

// oracleReport holds properties to be reported to a path oracle
type oracleReport struct {
	Stats stats    `json:"stats,omitempty"`
	Src   pan.IA   `json:"-"`
	Dst   pan.IA   `json:"-"`
	Path  pan.Path `json:"path"`
}

// TODO: dicover path oracle
const oracleLocation = "19-ffaa:1:e4b,[127.0.0.1]:443"
const oracleStatsEndpointFormat = "http://" + oracleLocation + "/stats/%s/%s/%s/%s/%s"

const jsonContentType = "application/json"

var instance *Client

type Client struct {
	*http.Client
	stats map[string]map[string]*ConnectionStats // local -> remote -> stats
}

var once sync.Once

// GetInstance returns the singelton instance
func GetInstance() *Client {
	once.Do(func() {
		instance = create()
	})
	return instance
}

func create() *Client {
	return &Client{
		Client: &http.Client{Transport: shttp.DefaultTransport},
		stats:  make(map[string]map[string]*ConnectionStats)}
}

func (c *Client) OnConnectionStarted(local, remote *pan.UDPAddr) {
	// log.Println("OracleClient - OnConnectionStarted")
	now := time.Now()
	if c.stats[local.String()] == nil {
		c.stats[local.String()] = make(map[string]*ConnectionStats)
	}
	c.stats[local.String()][remote.String()] = &ConnectionStats{Opened: &now, Local: local, Remote: remote}
}

func (c *Client) OnConnectionClosed(local, remote *pan.UDPAddr) {
	// log.Println("OracleClient - OnConnectionClosed")
	s := c.stats[local.String()][remote.String()]
	now := time.Now()
	s.Closed = &now
	c.reportToOracle(local, remote)
}

func (c *Client) OnPacketSent(local, remote *pan.UDPAddr, pktSize logging.ByteCount) {
	// log.Println("OracleClient - OnPacketSent")
	s := c.stats[local.String()][remote.String()]
	s.PacketsSent++
	s.BytesSent += pktSize
}

func (c *Client) OnPacketReceived(local, remote *pan.UDPAddr, pktSize logging.ByteCount) {
	// log.Println("OracleClient - OnPacketReceived")
	s := c.stats[local.String()][remote.String()]
	s.BytesReceived++
	s.BytesReceived += pktSize
}

func (c *Client) OnPathEvaluated(local, remote *pan.UDPAddr, newPath *pan.Path) {
	// log.Println("OracleClient - OnPathEvaluated")
	// TODO: report to oracle on path change ?
	s := c.stats[local.String()][remote.String()]
	s.ActivePath = newPath
}

func (c *Client) reportToOracle(local, remote *pan.UDPAddr) error {
	return c.postToOracle(c.getOracleReport(local, remote))
}

func (c *Client) postToOracle(report *oracleReport) error {
	reqUrl := fmt.Sprintf(oracleStatsEndpointFormat, report.Src.I.String(),
		strconv.FormatInt(int64(report.Src.A), 10),
		report.Dst.I.String(),
		strconv.FormatInt(int64(report.Dst.A), 10),
		report.Path.Fingerprint)

	log.Printf("posting to oracle, endpoint: %s", reqUrl)
	body, err := json.Marshal(report.Stats)
	if err != nil {
		log.Printf("error posting to path oracle stats: %s", err)
		return err
	}

	_, err = c.Post(shttp.MangleSCIONAddrURL(reqUrl), jsonContentType, bytes.NewReader(body))
	if err != nil {
		log.Printf("error posting to path oracle %stats", err)
		return err
	}
	log.Println("successfully reported to path oracle")
	return nil
}

func (c *Client) getOracleReport(local, remote *pan.UDPAddr) *oracleReport {
	s := c.stats[local.String()][remote.String()]
	report := oracleReport{
		Src:   s.Local.IA,
		Dst:   s.Remote.IA,
		Path:  *s.ActivePath,
		Stats: stats{},
	}
	lifetime := s.LifeTime()
	if lifetime > 0 {
		report.Stats.Bandwidth = float64(s.BytesSent) / lifetime.Seconds()
	}
	return &report
}
