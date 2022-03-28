package tcp

import (
	"errors"
	"net"

	"github.com/netsys-lab/panapi/taps"
)

type TCP struct {
	conn net.TCPConn
}

func (t *TCP) Satisfy(sp taps.SelectionProperties) (taps.TransportProperties, error) {
	sp = *sp.Copy()
	var err error
	if sp.Reliability == taps.Prohibit ||
		sp.PreserveMsgBoundaries == taps.Require ||
		sp.PerMsgReliability == taps.Require ||
		sp.PreserveOrder == taps.Prohibit ||
		sp.ZeroRTTMsg == taps.Require ||
		sp.Multistreaming == taps.Require ||
		sp.FullChecksumSend == taps.Require ||
		sp.FullChecksumRecv == taps.Require ||
		sp.CongestionControl == taps.Prohibit {
		err = errors.New("Can't satisfy all constraints")
	}
	return taps.TransportProperties{
		Reliability:              true,
		PreserveMsgBoundaries:    false,
		PerMsgReliability:        false,
		PreserveOrder:            true,
		ZeroRTTMsg:               false,
		Multistreaming:           false,
		FullChecksumSend:         false,
		FullChecksumRecv:         false,
		CongestionControl:        true,
		KeepAlive:                sp.KeepAlive == taps.Require || sp.KeepAlive == taps.Prefer,
		Interface:                sp.Interface,
		PvD:                      sp.PvD,
		Multipath:                taps.Disabled,
		UseTemporaryLocalAddress: sp.UseTemporaryLocalAddress == taps.Require || sp.UseTemporaryLocalAddress == taps.Prefer,
		AdvertisesAltAddr:        false,
		Direction:                sp.Direction,
		SoftErrorNotify:          sp.SoftErrorNotify == taps.Require || sp.SoftErrorNotify == taps.Prefer,
		ActiveReadBeforeSend:     sp.ActiveReadBeforeSend == taps.Require || sp.ActiveReadBeforeSend == taps.Prefer,
	}, err
}

func (t *TCP) SendFrame(messageData []byte, messageContext *taps.MessageContext) error {
	n, err := t.conn.Write(messageData)
	if err != nil {
		return err
	}
	if n != len(messageData) {
		return errors.New("Short write")
	}
	return nil
}

// not sure if we can actually call this "ReceiveFrame" here, because we just get any amount of data, not necessarily a whole frame
func (t *TCP) ReceiveFrame() ([]byte, *taps.MessageContext, error) {
	//FIXME TODO, what is a good buffer size here?
	messageData := make([]byte, 1500)
	n, err := t.conn.Read(messageData)
	return messageData[:n], taps.NewMessageContext(), err
}
