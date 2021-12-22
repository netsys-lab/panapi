package taps

import "errors"

var (
	ExpiredError = errors.New("Message expired")
	SendError    = errors.New("Could not send")
	ReceiveError = errors.New("Could not receive")
)

type EstablishmentError error

func NewEstablishmentError(reason string) EstablishmentError {
	return errors.New(reason)
}
