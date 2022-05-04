package taps

import (
	"io"
)

type Connection interface {
	io.ReadWriteCloser
	Preconnection() *Preconnection
}
