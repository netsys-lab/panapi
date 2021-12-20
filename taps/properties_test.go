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

package taps

import (
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
		{"Interface", map[string]Preference{"eth0": Ignore}, false},
		{"PvD", map[string]uint8{"fnord": 1}, true},
		{"AdvertisesAltAddr", false, false},
		{"Direction", UnidirectionalSend, false},
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
