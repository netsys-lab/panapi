package taps

import (
	"errors"
	"strconv"
)

var (
	errUnknownServiceType      = errors.New("unknown service type")
	errInvalidPort             = errors.New("invalid port")
	errInvalidIPAddress        = errors.New("invalid ip address")
	errInvalidEndpointType     = errors.New("invalid endpoint type")
	errInvalidArgument         = errors.New("invalid argument")
	errUnknownSetName          = errors.New("unknown name for set target")
	errReadOnClosedConnection  = errors.New("read on closed connection")
	errWriteOnClosedConnection = errors.New("write on closed connection")
)

type tapsError struct {
	Op string

	ServiceType int
	Ipv4address string
	Port        string

	ServiceTypeInvalid string
	ArgNum             int
	SetName            string

	Err error
}

func (e *tapsError) Unwrap() error { return e.Err }

func (e *tapsError) Error() string {
	if e == nil {
		return "<nil>"
	}
	s := e.Op + ": "
	if e.ServiceType != 0 {
		s += "service type: " + SERV_NAMES[e.ServiceType] + ", "
	}
	if e.ServiceTypeInvalid != "" {
		s += "service sype: " + e.ServiceTypeInvalid + ", "
	}
	if e.Ipv4address != "" {
		s += "IPv4 address: " + e.Ipv4address + ", "
	}
	if e.Port != "" {
		s += "port: " + e.Port
	}
	if e.ArgNum > 0 {
		s += "argument number: " + strconv.Itoa(e.ArgNum)
	}
	if e.SetName != "" {
		s += "target: " + e.SetName
	}
	s += " -> " + e.Err.Error()
	return s
}
