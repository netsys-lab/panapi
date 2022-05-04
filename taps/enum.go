package taps

//go:generate stringer -type=Preference,MultipathPreference,MultipathPolicy,Directionality,CapacityProfile,StreamScheduler,ConnectionState -output enum_string.go

type Preference uint8

const (
	// (Implementation detail: Indicate that recommended default
	// value for Property should be used)
	unset Preference = iota

	// No preference
	Ignore

	// Select only protocols/paths providing the property, fail
	// otherwise
	Require

	// Prefer protocols/paths providing the property, proceed
	// otherwise
	Prefer

	// Prefer protocols/paths not providing the property, proceed
	// otherwise
	Avoid

	// Select only protocols/paths not providing the property,
	// fail otherwise
	Prohibit
)

type MultipathPreference uint8

const (
	// (Implementation detail: need to use different defaults
	// depending on endpoint)
	dynamic MultipathPreference = iota

	// The connection will not use multiple paths once
	// established, even if the chosen transport supports using
	// multiple paths.
	Disabled

	// The connection will negotiate the use of multiple paths if
	// the chosen transport supports this.
	Active

	// The connection will support the use of multiple paths if
	// the Remote Endpoint requests it.
	Passive
)

type MultipathPolicy uint8

const (
	// The connection ought only to attempt to migrate between
	// different paths when the original path is lost or becomes
	// unusable.
	Handover MultipathPolicy = iota

	// The connection ought only to attempt to minimize the
	// latency for interactive traffic patterns by transmitting
	// data across multiple paths when this is beneficial. The
	// goal of minimizing the latency will be balanced against the
	// cost of each of these paths. Depending on the cost of the
	// lower-latency path, the scheduling might choose to use a
	// higher-latency path. Traffic can be scheduled such that
	// data may be transmitted on multiple paths in parallel to
	// achieve a lower latency.
	Interactive

	// The connection ought to attempt to use multiple paths in
	// parallel to maximize available capacity and possibly
	// overcome the capacity limitations of the individual paths.
	Aggregate
)

type Directionality uint8

const (
	// The connection must support sending and receiving data
	Bidirectional Directionality = iota

	// The connection must support sending data, and the application cannot use the connection to receive any data
	UnidirectionalSend

	// The connection must support receiving data, and the application cannot use the connection to send any data
	UnidirectionalReceive
)

type CapacityProfile uint8

const (
	// The application provides no information about its expected
	// capacity profile.
	Default CapacityProfile = iota

	// The application is not interactive. It expects to send
	// and/or receive data without any urgency. This can, for
	// example, be used to select protocol stacks with scavenger
	// transmission control and/or to assign the traffic to a
	// lower-effort service.
	Scavenger

	// The application is interactive, and prefers loss to
	// latency. Response time should be optimized at the expense
	// of delay variation and efficient use of the available
	// capacity when sending on this connection. This can be used
	// by the system to disable the coalescing of multiple small
	// Messages into larger packets (Nagle's algorithm); to prefer
	// immediate acknowledgment from the peer endpoint when
	// supported by the underlying transport; and so on.
	LowLatencyInteractive

	// The application prefers loss to latency, but is not
	// interactive. Response time should be optimized at the
	// expense of delay variation and efficient use of the
	// available capacity when sending on this connection.
	LowLatencyNonInteractive

	// The application expects to send/receive data at a constant
	// rate after Connection establishment. Delay and delay
	// variation should be minimized at the expense of efficient
	// use of the available capacity. This implies that the
	// Connection might fail if the Path is unable to maintain the
	// desired rate.
	ConstantRateStreaming

	// The application expects to send/receive data at the maximum
	// rate allowed by its congestion controller over a relatively
	// long period of time.
	CapacitySeeking
)

type StreamScheduler uint8

const (
	SCTP_SS_FCFS   StreamScheduler = iota // First-Come, First-Served Scheduler
	SCTP_SS_RR                            // Round-Robin Scheduler
	SCTP_SS_RR_PKT                        // Round-Robin Scheduler per Packet
	SCTP_SS_PRIO                          // Priority-Based Scheduler
	SCTP_SS_FC                            // Fair Capacity Scheduler
	SCTP_SS_WFQ                           // Weighted Fair Queueing Scheduler
)

type ConnectionState uint8

const (
	Establishing ConnectionState = iota
	Established
	Closing
	Closed
)
