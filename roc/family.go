package roc

// Family as declared in roc/address.h:36
type Family int32

// Family enumeration from roc/address.h:36
const (
	AfInvalid Family = -1
	AfAuto    Family = 0
	AfIpv4    Family = 1
	AfIpv6    Family = 2
)
