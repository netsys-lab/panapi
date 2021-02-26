package taps

const (
	NETWORK_IP    = "IP"
	NETWORK_IPV4  = "IPv4"
	NETWORK_IPV6  = "IPv6"
	NETWORK_SCION = "SCION"

	TRANSPORT_UDP  = "UDP"
	TRANSPORT_TCP  = "TCP"
	TRANSPORT_QUIC = "QUIC"
)

type Message string

func (m Message) String() string {
	return string(m)
}
