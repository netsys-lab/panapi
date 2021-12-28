package taps

type Message struct {
	// MessageContext is an optional pointer to a context object
	// for this Message
	*MessageContext

	// Data provides access to the bytes that are to be sent OR
	// that have been received for this Message
	Data []byte

	// EndOfMessage is flag that indicates whether, for the
	// purposes of the underlying transport, this message has been
	// completely submitted for transmission OR whether it has
	// been completely received from the underlying connection
	EndOfMessage bool
}
