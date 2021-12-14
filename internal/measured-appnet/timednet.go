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
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sock/reliable"
)

func NewNetwork(ia addr.IA, dispatcher reliable.Dispatcher,
	revHandler snet.RevocationHandler) *snet.SCIONNetwork {

	return &snet.SCIONNetwork{
		LocalIA: ia,
		Dispatcher: &TimedPacketDispatcherService{
			Dispatcher: dispatcher,
			SCMPHandler: &snet.DefaultSCMPHandler{
				RevocationHandler: revHandler,
			},
		},
	}
}
