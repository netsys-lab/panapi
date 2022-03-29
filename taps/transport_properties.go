package taps

import (
	"github.com/netsys-lab/panapi/internal/enum"
)

type TransportProperties struct {
	Reliability       bool
	PreserveOrder     bool
	CongestionControl bool
	Interface         string
	Multipath         enum.MultipathPreference
}

// Copy returns a new TransportProperties struct with its values deeply copied from tp
func (tp *TransportProperties) Copy() *TransportProperties {
	return &TransportProperties{
		Reliability:       tp.Reliability,
		PreserveOrder:     tp.PreserveOrder,
		CongestionControl: tp.CongestionControl,
		Interface:         tp.Interface,
		Multipath:         tp.Multipath,
	}
}
