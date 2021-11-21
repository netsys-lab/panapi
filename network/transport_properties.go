package network

type TransportProperties struct {
	Properties      map[string]string
	RequireProhibit map[string]bool
	PreferAvoid     map[string]bool
}

func NewTransportProperties() *TransportProperties {
	return &TransportProperties{
		map[string]string{},
		map[string]bool{},
		map[string]bool{},
	}
}

func (tp *TransportProperties) Set(property, value string) {
	tp.Properties[property] = value
}

func (tp *TransportProperties) Require(property string) {
	delete(tp.PreferAvoid, property)
	tp.RequireProhibit[property] = true

}
func (tp *TransportProperties) Prefer(property string) {
	delete(tp.RequireProhibit, property)
	tp.PreferAvoid[property] = true

}
func (tp *TransportProperties) Ignore(property string) {
	delete(tp.RequireProhibit, property)
	delete(tp.PreferAvoid, property)

}
func (tp *TransportProperties) Avoid(property string) {
	delete(tp.RequireProhibit, property)
	tp.PreferAvoid[property] = false
}
func (tp *TransportProperties) Prohibit(property string) {
	delete(tp.PreferAvoid, property)
	tp.RequireProhibit[property] = false
}
