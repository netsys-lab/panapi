package taps

type TransportProtocol interface {
	FrameSender
	FrameReceiver
	Satisfy(SelectionProperties) (TransportProperties, error)
}
