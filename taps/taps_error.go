package taps

import (
	"errors"
	"strconv"
)

var (
	errUnknownServiceType      = errors.New("unknown service type")
	errNoServiceType           = errors.New("no service type set")
	errInvalidPort             = errors.New("invalid port")
	errInvalidIPAddress        = errors.New("invalid ip address")
	errInvalidEndpointType     = errors.New("invalid endpoint type")
	errInvalidArgument         = errors.New("invalid argument")
	errUnknownSetName          = errors.New("unknown name for set")
	errUnknownRequireName      = errors.New("unknown name for require")
	errReadOnClosedConnection  = errors.New("read on closed connection")
	errWriteOnClosedConnection = errors.New("write on closed connection")
	errNoClientAddr            = errors.New("no address for client")
	errInvalidAddressType      = errors.New("invalid address type")
)

type tapsError struct {
	Op     string
	Endp   interface{}
	ArgNum int
	Desc   string
	Err    error
}

func (e *tapsError) Unwrap() error {
	return e.Err
}

func (e *tapsError) Error() string {
	if e == nil {
		return "<nil>"
	}
	s := "@" + e.Op + " -> " + e.Err.Error() + "\n"
	var ep *Endpoint
	switch e.Endp.(type) {
	case *LocalEndpoint:
		ep = &e.Endp.(*LocalEndpoint).Endpoint
	case *RemoteEndpoint:
		ep = &e.Endp.(*RemoteEndpoint).Endpoint
	case *Endpoint:
		ep = e.Endp.(*Endpoint)
	}
	if ep != nil {
		if ep.interfaceName != "" {
			s += "\t" + "interface name: " + ep.interfaceName + "\n"
		}
		if ep.serviceType != SERV_INVALID {
			s += "\t" + "service type: " + SERV_NAMES[ep.serviceType] + "\n"
		}
		if ep.address != "" {
			s += "\t" + "address: " + ep.address + "\n"
		}
		if ep.port != "" {
			s += "\t" + "port: " + ep.port + "\n"
		}
	}
	if e.ArgNum > 0 {
		s += "\t" + strconv.Itoa(e.ArgNum) + ". argument" + "\n"
	}
	if e.Desc != "" {
		s += "\t" + "given value: " + e.Desc + "\n"
	}
	return s
}
