package taps

import (
	"github.com/netsys-lab/panapi/internal/enum"
)

type TransportProperties struct {
	Reliability              bool
	PreserveMsgBoundaries    bool
	PerMsgReliability        bool
	PreserveOrder            bool
	ZeroRTTMsg               bool
	Multistreaming           bool
	FullChecksumSend         bool
	FullChecksumRecv         bool
	CongestionControl        bool
	KeepAlive                bool
	Interface                map[string]enum.Preference
	PvD                      map[string]enum.Preference
	Multipath                enum.MultipathPreference
	UseTemporaryLocalAddress bool
	AdvertisesAltAddr        bool
	Direction                enum.Directionality
	SoftErrorNotify          bool
	ActiveReadBeforeSend     bool
}

// Copy returns a new TransportProperties struct with its values deeply copied from tp
func (tp *TransportProperties) Copy() *TransportProperties {
	var (
		newInterface = make(map[string]enum.Preference)
		newPvD       = make(map[string]enum.Preference)
	)
	for key, value := range tp.Interface {
		newInterface[key] = value
	}
	for key, value := range tp.PvD {
		newPvD[key] = value
	}
	return &TransportProperties{
		Reliability:              tp.Reliability,
		PreserveMsgBoundaries:    tp.PreserveMsgBoundaries,
		PerMsgReliability:        tp.PerMsgReliability,
		PreserveOrder:            tp.PreserveOrder,
		ZeroRTTMsg:               tp.ZeroRTTMsg,
		Multistreaming:           tp.Multistreaming,
		FullChecksumSend:         tp.FullChecksumSend,
		FullChecksumRecv:         tp.FullChecksumRecv,
		CongestionControl:        tp.CongestionControl,
		KeepAlive:                tp.KeepAlive,
		Interface:                newInterface,
		PvD:                      newPvD,
		Multipath:                tp.Multipath,
		UseTemporaryLocalAddress: tp.UseTemporaryLocalAddress,
		AdvertisesAltAddr:        tp.AdvertisesAltAddr,
		Direction:                tp.Direction,
		SoftErrorNotify:          tp.SoftErrorNotify,
		ActiveReadBeforeSend:     tp.ActiveReadBeforeSend,
	}
}
