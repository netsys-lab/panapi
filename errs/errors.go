package errs

import "errors"

var (
	NetworkType   = errors.New("invalid network type")
	TransportType = errors.New("invalid address type")
)
