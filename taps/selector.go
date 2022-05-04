package taps

import "github.com/netsec-ethz/scion-apps/pkg/pan"

type Selector interface {
	pan.Selector
	SetPreferences(*ConnectionPreferences) error
}

type DefaultSelector struct {
	pan.DefaultSelector
}

func (s *DefaultSelector) SetPreferences(*ConnectionPreferences) error {
	return nil
}
