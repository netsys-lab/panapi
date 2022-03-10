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

const (
	// TODO: discover path oracle via RAINS
	oracleLocation = "19-ffaa:1:e4b,127.0.0.1:443"
	reportingURL   = "http://" + oracleLocation + "/stats/%s/%s/%s/"
	scoringURL     = "http://" + oracleLocation + "/scorings/"
)

const jsonContentType = "application/json"

var instance *Client

type Client struct {
	*http.Client
	stats            map[string]map[string]*ConnectionStats // local -> remote -> stats
	subscribedScores subscribedScores
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
		Client:           &http.Client{Transport: shttp.DefaultTransport},
		stats:            make(map[string]map[string]*ConnectionStats),
		subscribedScores: make(subscribedScores)}
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
	c.reportStats(local, remote)
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

func (c *Client) Subscribe(dst pan.IA, service scoringServiceName) (*scoreResponse, error) {
	c.subscribedScores[dst.String()] = append(c.subscribedScores[dst.String()], service)
	return c.fetchScores()
}

func (c *Client) reportStats(local, remote *pan.UDPAddr) error {
	return c.postReport(c.createReport(local, remote))
}

func (c *Client) postReport(report *oracleReport) error {
	reqUrl := fmt.Sprintf(reportingURL,
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

func (c *Client) createReport(local, remote *pan.UDPAddr) *oracleReport {
	s := c.stats[local.String()][remote.String()]
	report := oracleReport{
		Src:   s.Local.IA,
		Dst:   s.Remote.IA,
		Path:  *s.ActivePath,
		Stats: stats{},
	}

	if lifetime := s.LifeTime(); lifetime > 0 {
		report.Stats.Bandwidth = float64(s.BytesSent) / lifetime.Seconds()
		log.Printf("reporting bandwidth of: %.2f bytes/s\n", report.Stats.Bandwidth)
	}
	return &report
}

func (c *Client) fetchScores() (*scoreResponse, error) {
	scoreReq, err := json.Marshal(scoreRequest{c.subscribedScores})
	if err != nil {
		return nil, err
	}

	res, err := c.Post(scoringURL, jsonContentType, bytes.NewBuffer(scoreReq))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	scores := &scoreResponse{}
	if err = parseJson(res, scores); err != nil {
		return nil, err
	}
	return scores, nil
}

func parseJson(response *http.Response, target interface{}) error {
	return json.NewDecoder(response.Body).Decode(target)
}
