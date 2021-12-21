package taps

type Event interface{}

type SentEvent struct {
	MessageContext *MessageContext
}

type ExpiredEvent struct {
	MessageContext *MessageContext
}

type SendErrorEvent struct {
	MessageContext *MessageContext
	Error          error
}

type SoftErrorEvent struct {
	Error error
}

type ReadyEvent struct{}

type PathChangeEvent struct{}

type ReceivedEvent struct {
	MessageContext *MessageContext
	MessageData    []byte
}

type ReceivedPartialEvent struct {
	ReceivedEvent
	EndOfMessage bool
}

type ReceiveErrorEvent struct {
	MessageContext *MessageContext
	Error          error
}

type ConnectionReceivedEvent struct{}

type RendezvousDoneEvent struct{}

type ClosedEvent struct{}

type EstablishmentErrorEvent struct {
	Error error
}

type ConnectionErrorEvent struct {
	Error error
}
