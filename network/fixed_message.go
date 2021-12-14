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
	"errors"
	"io"
)

type FixedMessage struct {
	b          []byte
	readindex  int
	writeindex int
}

func NewFixedMessageString(s string) *FixedMessage {
	b := []byte(s)
	return &FixedMessage{b, 0, len(b)}
}

func NewFixedMessage(size int) *FixedMessage {
	return &FixedMessage{make([]byte, size), 0, 0}
}

func (m FixedMessage) String() string {
	return string(m.b)
}

//TODO, add proper locking and reusability of buffer
func (m FixedMessage) Read(p []byte) (i int, err error) {
	for i = m.readindex; i < len(m.b) && i < len(p); i++ {
		p[i] = m.b[i]
	}
	m.readindex = i
	if i == len(m.b) {
		err = io.EOF
	}
	return
}

// Write appends the contents of p to the message m, return error network.EOM when the capacity is met or exceeded
func (m FixedMessage) Write(p []byte) (i int, err error) {
	if m.writeindex == len(m.b) {
		errors.New("No space left in message")
	}
	for i = m.writeindex; i < len(m.b) && i < len(p); i++ {
		m.b[i] = p[i]
	}
	m.writeindex = i
	if m.writeindex == len(m.b) {
		err = EOM
	}
	return
}
