package panapi

import (
	"crypto"
	"crypto/tls"
	"time"
)

type Preference uint8

const (
	Require Preference = iota
	Prefer
	Ignore
	Avoid
	Prohibit
)

var preferenceNames = [...]string{
	"Require",
	"Prefer",
	"Ignore",
	"Avoid",
	"Prohibit",
}

func (p Preference) String() string {
	return preferenceNames[p-Require]
}

type TransportProperties struct {
	properties  map[string]string
	preferences map[string]Preference
}

// NewTransportProperties creates TransportProperties with certain defaults set
func NewTransportProperties() *TransportProperties {
	tp := TransportProperties{
		map[string]string{},
		map[string]Preference{},
	}
	// Defaults from
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#name-specifying-transport-proper
	tp.Require("reliability")
	tp.Ignore("preserveMsgBoundaries")
	tp.Ignore("perMsgReliability")
	tp.Require("preserveOrder")
	tp.Ignore("zeroRttMsg")
	tp.Prefer("multistreaming")
	tp.Require("FullChecksumSend")
	tp.Require("FullChecksumRecv")
	tp.Require("congestionControl")
	tp.Ignore("keepAlive")
	// skip pvd and interface
	tp.Prefer("useTemporaryLocalAddress")
	// skip "multipath", "advertises-altaddr", "direction"
	tp.Ignore("softErrorNotify")
	tp.Ignore("activeReadBeforeSend")

	return &tp
}

func (tp *TransportProperties) SetProperty(property, value string) {
	tp.properties[property] = value
}

// SetPreference stores preference
func (tp *TransportProperties) SetPreference(property string, preference Preference) {
	tp.preferences[property] = preference
}

/* GetPreference returns the stored Preference for property, or Ignore if preference is not found

   Deprecated: GetX is not idiomatic Go, use X instead, i.e., "Preference(property)"
*/
func (tp *TransportProperties) GetPreference(property string) Preference {
	return tp.Preference(property)
}

// Preference returns the stored Preference for property, or Ignore if preference is not found
func (tp *TransportProperties) Preference(property string) Preference {
	p, ok := tp.preferences[property]
	if !ok {
		return Ignore
	}
	return p
}

// Require has the effect of selecting only protocols/paths providing the property and failing otherwise
func (tp *TransportProperties) Require(property string) {
	tp.preferences[property] = Require

}

// Prefer has the effect of prefering protocols/paths providing the property and proceeding otherwise
func (tp *TransportProperties) Prefer(property string) {
	tp.preferences[property] = Prefer
}

// Ignore has the effect of expressing no preference for the given property
func (tp *TransportProperties) Ignore(property string) {
	tp.preferences[property] = Ignore

}

// Avoid has the effect of avoiding protocols/paths with the property and proceeding otherwise
func (tp *TransportProperties) Avoid(property string) {
	tp.preferences[property] = Avoid
}

// Prohibit has the effect of failing if the property can not be avoided
func (tp *TransportProperties) Prohibit(property string) {
	tp.preferences[property] = Prohibit
}

// SecurityParameters is a structure used to configure security for a Preconnection
type SecurityParameters struct {
	Identity              string //? []byte?
	KeyPair               KeyPair
	SupportedGroup        tls.CurveID
	CipherSuite           tls.CipherSuite
	SignatureAlgorithm    tls.SignatureScheme
	MaxCachedSessions     uint
	CachedSessionLifetime time.Duration
	//PSK not supported
}

// SetTrustVerificationCallback is not yet implemented
func (sp SecurityParameters) SetTrustVerificationCallback() {
}

// SetIdentityChallengeCallback is not yet implemented
func (sp SecurityParameters) SetIdentityChallengeCallback() {
}

type CapacityProfile uint8

const (
	Default CapacityProfile = iota
	Scavenger
	LowLatencyInteractive
	LowLatencyNonInteractive
	ConstantRateStreaming
	CapacitySeeking
)

var profileNames = [...]string{
	"Default",
	"Scavenger",
	"Low Latency/Interactive",
	"Low Latency/Non-Interactive",
	"Constant-Rate Streaming",
	"Capacity-Seeking",
}

func (p CapacityProfile) String() string {
	return profileNames[p-Default]
}

type ConnectionProperties struct {
	RecvChecksumLen uint
	ConnPrio        uint
	ConnTimeout     time.Duration
	// ConnScheduler not yet implemented
	ConnCapacityProfile CapacityProfile
	/*
	   MultipathPolicy
	   MinSendRate
	   MinRecvRate
	   MaxSendRate
	   MaxRecvRate
	   GroupConnLimit
	   IsolateSession
	   not yet implemented
	*/
}

// KeyPair clearly associates a Private and Public Key into a pair
type KeyPair struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}

type Connection struct {
	Ready      chan bool
	SoftError  chan error
	PathChange chan bool
}
