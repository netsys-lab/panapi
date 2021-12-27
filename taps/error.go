package taps

import "errors"

var (
	StoppedError            = errors.New("Listener has stopped")
	SendError               = errors.New("Could not send")
	ReceiveError            = errors.New("Could not receive")
	NotYetImplementendError = errors.New("Not yet implemented")
	ExpiredError            = errors.New("The message could not be sent before its lifetime")
	PartialMessageError     = errors.New("Message not yet fully delivered")
)

type EstablishmentError error

func NewEstablishmentError(reason string) EstablishmentError {
	return errors.New(reason)
}
