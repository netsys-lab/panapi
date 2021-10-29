package network

import "errors"

var (
	NetTypeError  = errors.New("invalid network type")
	AddrTypeError = errors.New("invalid address type")
)
