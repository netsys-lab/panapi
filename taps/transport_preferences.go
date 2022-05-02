package taps

type TransportPreferences struct {
	Reliability       Preference
	PreserveOrder     Preference
	CongestionControl Preference
	Interface         map[string]Preference
	Multipath         MultipathPreference
}

// Copy returns a new TransportProperties struct with its values deeply copied from tp
func (tp *TransportPreferences) Copy() *TransportPreferences {
	var (
		newInterface = make(map[string]Preference)
	)
	for key, value := range tp.Interface {
		newInterface[key] = value
	}
	return &TransportPreferences{
		Reliability:       tp.Reliability,
		PreserveOrder:     tp.PreserveOrder,
		CongestionControl: tp.CongestionControl,
		Interface:         newInterface,
		Multipath:         tp.Multipath,
	}
}

// NewTransportPreferences creates TransportPreferences with the
// recommended defaults from
// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2
func NewTransportPreferences() *TransportPreferences {
	return &TransportPreferences{
		Reliability:       Require,
		PreserveOrder:     Require,
		CongestionControl: Require,
		Interface:         map[string]Preference{},
		Multipath:         dynamic,
	}
}
