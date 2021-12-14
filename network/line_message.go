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
	"bytes"
	"strings"
)

type LineMessage struct {
	b *bytes.Buffer
}

func NewLineMessageString(s string) (*LineMessage, error) {
	if strings.ContainsAny(s, "\r\n") {
		return nil, NewlineError
	}
	b := bytes.NewBufferString(s + "\n")
	return &LineMessage{b}, nil
}

func NewLineMessage() *LineMessage {
	return &LineMessage{new(bytes.Buffer)}
}

func (m LineMessage) String() string {
	return m.b.String()
}

func (m LineMessage) Read(p []byte) (int, error) {
	return m.b.Read(p)
}

// Write the contents of p into this message m, return error network.EOM when a newline is found
func (m LineMessage) Write(p []byte) (int, error) {
	// new temporary buffer to hold the
	buf := new(bytes.Buffer)
	n, _ := buf.Write(p)
	line, err := buf.ReadString('\n')
	if err != nil {
		//no newline was found (yet)
		return n, err
	} else {
		//a newline was found
		m.b.WriteString(line)
		err = EOM
	}
	return n, err
}
