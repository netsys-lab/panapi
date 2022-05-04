package udp

import (
	"errors"

	"github.com/netsys-lab/panapi/taps"
)

type UDP struct {
}

func (u *UDP) Satisfy(sp taps.SelectionProperties) (taps.TransportProperties, error) {
	sp = *sp.Copy()
	var err error
	if sp.Reliability == taps.Require ||
		sp.PreserveMsgBoundaries == taps.Require ||
		sp.PerMsgReliability == taps.Require ||
		sp.PreserveOrder == taps.Require ||
		sp.ZeroRTTMsg == taps.Prohibit ||
		sp.Multistreaming == taps.Require ||
		sp.FullChecksumSend == taps.Require ||
		sp.FullChecksumRecv == taps.Require ||
		sp.CongestionControl == taps.Require ||
		sp.KeepAlive == taps.Require {
		err = errors.New("Can't satisfy all constraints")
	}
	return taps.TransportProperties{
		Reliability: false,
		//PreserveMsgBoundaries:    false,
		//PerMsgReliability:        false,
		PreserveOrder: false,
		//ZeroRTTMsg:               true,
		//Multistreaming:           false,
		//FullChecksumSend:         false,
		//FullChecksumRecv:         false,
		CongestionControl: false,
		//KeepAlive:                false,
		//Interface:                sp.Interface,
		//PvD:                      sp.PvD,
		Multipath: taps.Disabled,
		//UseTemporaryLocalAddress: sp.UseTemporaryLocalAddress == taps.Require || sp.UseTemporaryLocalAddress == taps.Prefer,
		//AdvertisesAltAddr:        false,
		//Direction:                sp.Direction,
		//SoftErrorNotify:          sp.SoftErrorNotify == taps.Require || sp.SoftErrorNotify == taps.Prefer,
		//ActiveReadBeforeSend:     sp.ActiveReadBeforeSend == taps.Require || sp.ActiveReadBeforeSend == taps.Prefer,
	}, err
}
