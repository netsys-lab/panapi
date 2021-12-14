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

import (
	"io"
	"net"
)

type Message interface {
	String() string
	io.ReadWriter
}

type Connection interface {
	io.ReadWriteCloser
	Send(Message) error
	Receive(Message) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetError(error)
	GetError() error
}

type Dialer interface {
	Dial() (Connection, error)
}

type Listener interface {
	Listen() (Connection, error)
	Stop() error
}

type Preconnection interface {
	Listen() (Listener, error)
	Initiate() (Connection, error)
}

type Network interface {
	NewListener(*Endpoint) (Listener, error)
	NewDialer(*Endpoint) (Dialer, error)
}
