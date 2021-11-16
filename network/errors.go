package network

import "errors"

var (
	NetTypeError  = errors.New("invalid network type")
	AddrTypeError = errors.New("invalid address type")
	EOM           = errors.New("End of Message")
	NewlineError  = errors.New("Invalid String for Message (contains Newlines)")
)
