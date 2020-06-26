package roc

// Family defines network address family
type Family int

const (
	afInvalid Family = -1

	// AfAuto means automatically detect address family from string format
	AfAuto Family = 0

	// AfIPv4 indicates IPv4 address
	AfIPv4 Family = 1

	// AfIPv6 indicates IPv6 address
	AfIPv6 Family = 2
)
