package roc

// Network interface.
//
// Interface is a way to access the peer (sender or receiver) via network.
//
// Each peer slot has multiple interfaces, one of each type. The user interconnects
// peers by binding one of the first peer's interfaces to an URI and then connecting the
// corresponding second peer's interface to that URI.
//
// A URI is represented by Endpoint object.
//
// The interface defines the type of the communication with the remote peer and the
// set of protocols (URI schemes) that can be used with this particular interface.
//
// InterfaceConsolidated is an interface for high-level protocols which
// automatically manage all necessary communication: transport streams, control messages,
// parameter negotiation, etc. When a consolidated connection is established, peers may
// automatically setup lower-level interfaces like InterfaceAudioSource,
// InterfaceAudioRepair, and InterfaceAudioControl.
//
// InterfaceConsolidated is mutually exclusive with lower-level interfaces.
// In most cases, the user needs only InterfaceConsolidated. However, the
// lower-level interfaces may be useful if an external signaling mechanism is used or for
// compatibility with third-party software.
//
// InterfaceAudioSource and InterfaceAudioRepair are lower-level
// unidirectional transport-only interfaces. The first is used to transmit audio stream,
// and the second is used to transmit redundant repair stream, if FEC is enabled.
//
// InterfaceAudioControl is a lower-level interface for control streams.
// If you use InterfaceAudioSource and InterfaceAudioRepair, you
// usually also need to use InterfaceAudioControl to enable carrying additional
// non-transport information.
//
//go:generate stringer -type Interface -trimprefix Interface
type Interface int

const (
	// Interface that consolidates all types of streams (source, repair, control).
	//
	// Allowed operations:
	//  - bind    (sender, receiver)
	//  - connect (sender, receiver)
	//
	// Allowed protocols:
	//  - ProtoRtsp
	//
	InterfaceConsolidated Interface = 1

	// Interface for audio stream source data.
	//
	// Allowed operations:
	//  - bind    (receiver)
	//  - connect (sender)
	//
	// Allowed protocols:
	//  -  ProtoRtp
	//  -  ProtoRtpRs8mSource
	//  -  ProtoRtpLdpcSource
	//
	InterfaceAudioSource Interface = 11

	// Interface for audio stream repair data.
	//
	// Allowed operations:
	//  - bind    (receiver)
	//  - connect (sender)
	//
	// Allowed protocols:
	//  -  ProtoRs8mRepair
	//  -  ProtoLdpcRepair
	//
	InterfaceAudioRepair Interface = 12

	// Interface for audio control messages.
	//
	// Allowed operations:
	//  - bind    (sender, receiver)
	//  - connect (sender, receiver)
	//
	// Allowed protocols:
	//  -  ProtoRrcp
	//
	InterfaceAudioControl Interface = 13
)
