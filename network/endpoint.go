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

type Endpoint struct {
	Local         bool
	LocalAddress  string
	RemoteAddress string
	Transport     string
	Network       string
}

func NewRemoteEndpoint() *Endpoint {
	return &Endpoint{Local: false}
}

func NewLocalEndpoint() *Endpoint {
	return &Endpoint{Local: true}
}

func (e *Endpoint) WithNetwork(network string) {
	e.Network = network
}

func (e *Endpoint) WithTransport(transport string) {
	e.Transport = transport
}

func (e *Endpoint) WithAddress(addr string) {
	if e.Local {
		e.LocalAddress = addr
	} else {
		e.RemoteAddress = addr
	}
}
