package taps

import (
	"crypto"
	"crypto/tls"
	"time"
)

// SecurityParameters is a structure used to configure security for a
// Preconnection
type SecurityParameters struct {
	// Local identity and private keys: Used to perform private
	// key operations and prove one's identity to the Remote
	// Endpoint.
	Identity string
	KeyPair  KeyPair

	// Supported algorithms: Used to restrict what parameters are
	// used by underlying transport security protocols. When not
	// specified, these algorithms should use known and safe
	// defaults for the system. Parameters include: ciphersuites,
	// supported groups, and signature algorithms. These
	// parameters take a collection of supported algorithms as
	// parameter.
	SupportedGroup        tls.CurveID
	CipherSuite           *tls.CipherSuite
	SignatureAlgorithm    tls.SignatureScheme
	MaxCachedSessions     uint
	CachedSessionLifetime time.Duration

	// Unsupported
	PSK crypto.PrivateKey
}

// NewSecurityParameters TODO
func NewSecurityParameters() *SecurityParameters {
	return &SecurityParameters{}
}

// NewDisabledSecurityParameters is intended for compatibility with
// endpoints that do not support transport security protocols (such as
// TCP without support for TLS)
func NewDisabledSecurityParameters() *SecurityParameters {
	return &SecurityParameters{}
}

// NewOpportunisticSecurityParameters() is not yet implemented
func NewOpportunisticSecurityParameters() *SecurityParameters {
	return &SecurityParameters{}
}

// SetTrustVerificationCallback is not yet implemented
func (sp SecurityParameters) SetTrustVerificationCallback() {
}

// SetIdentityChallengeCallback is not yet implemented
func (sp SecurityParameters) SetIdentityChallengeCallback() {
}

/*
// Set stores value for parameter, which is stripped of case and
// non-alphabetic characters before being matched against the (equally
// stripped) exported Field names of sp. The type of value must be
// assignable to type of the targeted parameter Field, otherwise an
// error is returned.
//
// For the sake of respecting the TAPS (draft) spec as closely as
// possible, this function allows you to say:
//  err := sp.Set("supported-group", tls.CurveP521)
//  if err != nil {
//    ... // handle runtime error
//  }
//
// In idiomatic Go, you would (and should) instead say:
//  sp.SupportedGroup = tls.CurveP521
//
// Deprecated: Use func sp.Set only if you must. Direct access of the
// SecurityParameters struct Fields is usually preferred. This
// function is implemented using reflection and dynamic string
// matching, which is inherently inefficient and prone to bugs
// triggered at runtime.
func (sp *SecurityParameters) Set(parameter string, value interface{}) error {
	return set(sp, parameter, value)
}
*/

// Copy returns a new SecurityParameters struct with its values deeply copied from sp
func (sp *SecurityParameters) Copy() *SecurityParameters {
	return &SecurityParameters{
		Identity:              sp.Identity,
		KeyPair:               sp.KeyPair,
		SupportedGroup:        sp.SupportedGroup,
		CipherSuite:           sp.CipherSuite,
		SignatureAlgorithm:    sp.SignatureAlgorithm,
		MaxCachedSessions:     sp.MaxCachedSessions,
		CachedSessionLifetime: sp.CachedSessionLifetime,
		PSK:                   sp.PSK,
	}
}

// KeyPair clearly associates a Private and Public Key into a pair
type KeyPair struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}
