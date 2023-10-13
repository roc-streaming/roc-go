package roc

// Forward Error Correction encoding.
//
//go:generate stringer -type FecEncoding -trimprefix FecEncoding
type FecEncoding int

const (
	// No FEC encoding.
	// Compatible with ProtoRtp protocol.
	FecEncodingDisable FecEncoding = -1

	// Default FEC encoding.
	// Current default is FecEncodingRs8m.
	FecEncodingDefault FecEncoding = 0

	// Reed-Solomon FEC encoding (RFC 6865) with m=8.
	// Good for small block sizes (below 256 packets).
	// Compatible with ProtoRtpRs8mSource and ProtoRs8mRepair
	// protocols for source and repair endpoints.
	FecEncodingRs8m FecEncoding = 1

	// LDPC-Staircase FEC encoding (RFC 6816).
	// Good for large block sizes (above 1024 packets).
	// Compatible with ProtoRtpLdpcSource and ProtoLdpcRepair
	// protocols for source and repair endpoints.
	FecEncodingLdpcStaircase FecEncoding = 2
)
