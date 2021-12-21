package taps

type Event interface{}

type ErrorEvent struct {
	Error error
}

type SentEvent struct {
	MessageContext *MessageContext
}

type ExpiredEvent struct {
	MessageContext *MessageContext
}

type SendErrorEvent struct {
	Error          error
	MessageContext *MessageContext
}

type SoftErrorEvent ErrorEvent

type ReadyEvent struct{}

type PathChangeEvent struct{}

type ReceivedEvent struct {
	MessageContext *MessageContext
	MessageData    []byte
}

type ReceivedPartialEvent struct {
	MessageContext *MessageContext
	MessageData    []byte
	EndOfMessage   bool
}

type ReceiveErrorEvent struct {
	Error          error
	MessageContext *MessageContext
}

// ConnectionReceivedEvent is raised with a new Connection when a
// Remote Endpoint has established a transport-layer connection to
// this Listener (for Connection-oriented transport protocols), or
// when the first Message has been received from the Remote Endpoint
// (for Connectionless protocols), causing a new Connection to be
// created.
type ConnectionReceivedEvent struct {
	Connection *Connection
}

type StoppedEvent struct{}

type RendezvousDoneEvent struct{}

type ClosedEvent struct{}

type EstablishmentErrorEvent ErrorEvent

type ConnectionErrorEvent ErrorEvent
