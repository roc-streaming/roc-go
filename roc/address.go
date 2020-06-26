package roc

/*
#include <roc/address.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// Address represents network endpoint address.
//
// Consists of IP address plus UDP or TCP port number.
// Similar to net.Addr in Go and struct sockaddr in C.
//
// Should not be used concurrently
type Address struct {
	raw *C.roc_address
	mem []byte
}

// NewAddress creates and initializes a new Address.
//
// The IP address is parsed from a string representation. If family is AfAuto, the
// address family is auto-detected from string format. Otherwise, the string format
// should correspond to the family specified.
//
// The port number should be in range [0; 65536).
//
// When Address is used to bind a sender or receiver port, the "0.0.0.0" (for IPv4)
// or "::" (for IPv6) may be used to bind the port to all network interfaces, and
// zero port number may be used to bind the port to a randomly chosen ephemeral port.
func NewAddress(family Family, ip string, port int) (*Address, error) {
	a := new(Address)
	a.mem = make([]byte, C.sizeof_roc_address)
	a.raw = (*C.roc_address)(unsafe.Pointer(&a.mem[0]))

	cfamily := (C.roc_family)(family)
	cip := toCStr(ip)
	cport := (C.int)(port)
	errCode := C.roc_address_init(a.raw, cfamily, (*C.char)(unsafe.Pointer(&cip[0])), cport)

	if errCode == 0 {
		return a, nil
	}
	if errCode < 0 {
		return nil, ErrInvalidArguments
	}
	return nil, ErrInvalidApi
}

// Family returns address family.
//
// If AfAuto was used to construct address, the actually selected family, i.e.
// either AfIPv4 or AfIPv6, is reported.
func (a *Address) Family() (Family, error) {
	f := C.roc_address_family(a.raw)
	family := (Family)(f)
	if family == afInvalid {
		return family, ErrInvalidArguments
	}
	return family, nil
}

// IP returns IP address formatted to string.
func (a *Address) IP() (string, error) {
	const buflen = 255
	sIP := make([]byte, buflen)
	res := C.roc_address_ip(a.raw, (*C.char)(unsafe.Pointer(&sIP[0])), buflen)
	if res == nil {
		return "", ErrInvalidArguments
	}
	ret := C.GoString(res)
	runtime.KeepAlive(sIP)
	return ret, nil
}

// Port return UDP or TCP port number.
//
// If Address was passed to sender or receiver bind and the initial port number was
// zero, which means "use random port", this function will return the actually
// selected port number.
func (a *Address) Port() (int, error) {
	res := C.roc_address_port(a.raw)
	if res < 0 {
		return (int)(res), ErrInvalidArguments
	}
	return (int)(res), nil
}
