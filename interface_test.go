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
package panapi

import (
	"crypto/tls"
	"fmt"
	"testing"
)

func TestTransportProperties(t *testing.T) {
	proptests := []struct {
		property   string
		preference interface{}
		err        bool
	}{
		{"Reliability", Ignore, false},
		{"preserve-msg-boundaries", Require, false},
		{"PerMsgReliability", Bidirectional, true},
		{"Typo", Require, true},
		{"PreserveOrder", nil, true},
		{"ZeroRTTMsg", unset, false},
		//Multistreaming:           Prefer,
		//FullChecksumSend:         Require,
		//FullChecksumRecv:         Require,
		//CongestionControl:        Require,
		//KeepAlive:                Ignore,
		{"Interface", map[string]Preference{"eth0": Ignore}, false},
		{"PvD", map[string]uint8{"fnord": 1}, true},
		//UseTemporaryLocalAddress: unset,   // Needs to be resolved at runtime: Avoid for Listeners and Rendezvous Connections, else Prefer
		//Multipath:                dynamic, // Needs to be resolved at runtime: Disabled for Initiated and Rendezvous Connections, else Passive
		{"AdvertisesAltAddr", false, false},
		{"Direction", UnidirectionalSend, false},
		//SoftErrorNotify:          Ignore,
		//ActiveReadBeforeSend:     Ignore,*/

	}

	p := NewTransportProperties()

	for _, tt := range proptests {
		err := p.Set(tt.property, tt.preference)
		t.Logf("%s: %v (Error: %v)", tt.property, tt.preference, err)
		if tt.err {
			if err == nil {
				t.Fatal("expected error, got no error")
			}
		} else if err != nil {
			t.Fatal(err)
		}
	}
	//t.Logf("%+v", p)
}

func ExampleTransportProperties_Set() {
	tp := NewTransportProperties()

	// Calling Set() is possible, because the TAPS (draft)
	// specifies a Set function on the TransportProperties Object
	// (i.e., struct).
	err := tp.Set("preserve-msg-boundaries", Require)
	if err != nil {
		panic(err)
	}

	// Idiomatic Go would be to instead directly access the Field:
	tp.PreserveMsgBoundaries = Ignore

	fmt.Println(tp.PreserveMsgBoundaries)
	// Output: Ignore

}

func ExampleSecurityParameters_Set() {
	sp := NewSecurityParameters()

	var suite *tls.CipherSuite
	// find CipherSuite
	for _, suite = range tls.CipherSuites() {
		if suite.ID == tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256 {
			break
		}
	}

	// Calling Set() is possible, because the TAPS (draft)
	// specifies a Set function on the SecurityParameters Object
	// (i.e., struct).
	err := sp.Set("ciphersuite", suite)
	if err != nil {
		panic(err)
	}

	// Idiomatic Go would be to instead directly access the Field:
	sp.CipherSuite = suite

	fmt.Printf("%s", sp.CipherSuite.Name)
	// Output: TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256

}
