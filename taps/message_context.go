package taps

import (
	"math"
	"time"

	"github.com/netsys-lab/panapi/internal/enum"
)

// MessageContext carries extra information that an application might
// wish to attach to a Message
//
// As laid out in the TAPS spec, this type of Object is reminiscent of
// https://pkg.go.dev/context, which could conceivably also be
// used instead of this dedicated MessageContext Type. Feedback welcome.
type MessageContext struct {
	// MsgLifetime specifies how long a particular Message can
	// wait to be sent to the Remote Endpoint before it is
	// irrelevant and no longer needs to be (re-)transmitted. This
	// is a hint to the Transport Services system - it is not
	// guaranteed that a Message will not be sent when its
	// Lifetime has expired.
	//
	// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.1.3.1
	MsgLifetime time.Duration

	// MsgPrio specifies the priority of a Message, relative to
	// other Messages sent over the same Connection.A Message with
	// Priority 0 will yield to a Message with Priority 1, which
	// will yield to a Message with Priority 2, and so
	// on. Priorities may be used as a sender-side scheduling
	// construct only, or be used to specify priorities on the
	// wire for Protocol Stacks supporting prioritization.
	//
	// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.1.3.2
	MsgPrio uint

	// MsgOrdered, if true, preserves the order on delivery in
	// which Messages were originally submitted.
	//
	// Deprecated: The TAPS spec sets the default of this value to
	// "the queried Boolean value of the Selection Property
	// preserveOrder", but the latter is never actually in
	// scope/queriable when the MessageContext object is
	// created. Not setting the parameter defaults to "false",
	// which is indistinguishable from actively setting it to
	// "false". Feedback welcome.
	//
	// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.1.3.3
	MsgOrdered bool

	// SafelyReplayable, if true, specifies that a Message is safe
	// to send to the Remote Endpoint more than once for a single
	// Send Action. It marks the data as safe for certain 0-RTT
	// establishment techniques, where retransmission of the 0-RTT
	// data may cause the remote application to receive the
	// Message multiple times.
	//
	// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.1.3.4
	SafelyReplayable bool

	// Final, if true, this indicates a Message is the last that
	// the application will send on a Connection. This allows
	// underlying protocols to indicate to the Remote Endpoint
	// that the Connection has been effectively closed in the
	// sending direction. For example, TCP-based Connections can
	// send a FIN once a Message marked as Final has been
	// completely sent, indicated by marking
	// endOfMessage. Protocols that do not support signalling the
	// end of a Connection in a given direction will ignore this
	// property.
	//
	// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.1.3.5
	Final bool

	// MsgChecksumLen specifies the minimum length of the section
	// of a sent Message, starting from byte 0, that the
	// application requires to be delivered without corruption due
	// to lower layer errors. It is used to specify options for
	// simple integrity protection via checksums. A value of 0
	// means that no checksum needs to be calculated, and Full
	// Coverage means that the entire Message needs to be
	// protected by a checksum. Only Full Coverage is guaranteed,
	// any other requests are advisory, which may result in Full
	// Coverage being applied.
	//
	// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.1.3.6
	MsgChecksumLen uint

	// MsgReliable, if true, specifies that a Message should be
	// sent in such a way that the transport protocol ensures all
	// data is received on the other side without
	// corruption. Changing the Reliable Data Transfer property on
	// Messages is only possible for Connections that were
	// established enabling the Selection Property Configure
	// Per-Message Reliability. When this is not the case,
	// changing msgReliable will generate an error.
	//
	// Deprecated: The TAPS spec sets the default of this value to
	// "the queried Boolean value of the Selection Property
	// reliability", but the latter is never actually in
	// scope/queriable when the MessageContext object is
	// created. Not setting the parameter defaults to "false",
	// which is indistinguishable from actively setting it to
	// "false". Feedback welcome.
	//
	// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.1.3.7
	MsgReliable bool

	// MsgCapacityProfile specifies the application's preferred
	// tradeoffs for sending this Message; it is a per-Message
	// override of the CapacityProfile ConnectionProperty
	//
	// Implementation note: The TAPS spec sets the default of this
	// value to "inherited from the Connection Property
	// connCapacityProfile", but the latter is never actually in
	// scope/queriable when the MessageContext object is
	// created. However, the default "Default" CapacityProfile is
	// defined as "no information about the expected capacity
	// profile", which this implementation interprets as "no
	// override" for the purposes of MessageContext. Feedback
	// welcome.
	//
	// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.1.3.8
	MsgCapacityProfile enum.CapacityProfile

	// NoFragmentation specifies that a message should be sent and
	// received without network-layer fragmentation, if
	// possible. It can be used to avoid network layer
	// fragmentation when transport segmentation is prefered.
	//
	// This only takes effect when the transport uses a network
	// layer that supports this functionality. When it does take
	// effect, setting this property to true will cause the sender
	// to avoid network-layer source frgementation. When using
	// IPv4, this will result in the Don't Fragment bit being set
	// in the IP header.
	//
	// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.1.3.9
	NoFragmentation bool

	// NoSegmentation, if set, requests the transport layer to not
	// provide segmentation of messages larger than the maximum
	// size permitted by the network layer, and also to avoid
	// network-layer source fragmentation of messages. When
	// running over IPv4, setting this property to true can result
	// in a sending endpoint setting the Don't Fragment bit in the
	// IPv4 header of packets generated by the transport layer. An
	// attempt to send a message that results in a size greater
	// than the transport's current estimate of its maximum packet
	// size (singularTransmissionMsgMaxLen) will result in a
	// SendError. This only takes effect when the transport and
	// network layer support this functionality.
	//
	// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-9.1.3.10
	NoSegmentation bool

	localEndpoint  *LocalEndpoint
	remoteEndpoint *RemoteEndpoint
}

func NewMessageContext() *MessageContext {
	return &MessageContext{
		MsgPrio:        100,
		MsgChecksumLen: math.MaxUint,
	}

}

/*
// GetLocalEndpoint can be called by the application to query information about the Local Endpoint
//
// Deprecated: Per https://go.dev/doc/effective_go#Getters, it is not idiomatic Go to put "Get" into a getter's name
func (mc *MessageContext) GetLocalEndpoint() *LocalEndpoint {
	return mc.LocalEndpoint()
}
*/

// LocalEndpoint can be called by the application to query information about the Local Endpoint
func (mc *MessageContext) LocalEndpoint() *LocalEndpoint {
	return mc.localEndpoint
}

/*
// GetRemoteEndpoint can be called by the application to query information about the Remote Endpoint
//
// Deprecated: Per https://go.dev/doc/effective_go#Getters, it is not idiomatic Go to put "Get" into a getter's name
func (mc *MessageContext) GetRemoteEndpoint() *RemoteEndpoint {
	return mc.RemoteEndpoint()
}
*/

// RemoteEndpoint can be called by the application to query information about the Remote Endpoint
func (mc *MessageContext) RemoteEndpoint() *RemoteEndpoint {
	return mc.remoteEndpoint
}
