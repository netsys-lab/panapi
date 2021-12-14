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
package network

type TransportProperties struct {
	Properties      map[string]string
	RequireProhibit map[string]bool
	PreferAvoid     map[string]bool
}

func NewTransportProperties() *TransportProperties {
	return &TransportProperties{
		map[string]string{},
		map[string]bool{},
		map[string]bool{},
	}
}

func (tp *TransportProperties) Set(property, value string) {
	tp.Properties[property] = value
}

func (tp *TransportProperties) Require(property string) {
	delete(tp.PreferAvoid, property)
	tp.RequireProhibit[property] = true

}
func (tp *TransportProperties) Prefer(property string) {
	delete(tp.RequireProhibit, property)
	tp.PreferAvoid[property] = true

}
func (tp *TransportProperties) Ignore(property string) {
	delete(tp.RequireProhibit, property)
	delete(tp.PreferAvoid, property)

}
func (tp *TransportProperties) Avoid(property string) {
	delete(tp.RequireProhibit, property)
	tp.PreferAvoid[property] = false
}
func (tp *TransportProperties) Prohibit(property string) {
	delete(tp.PreferAvoid, property)
	tp.RequireProhibit[property] = false
}
