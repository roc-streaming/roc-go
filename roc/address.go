package roc

/*
#include <roc/address.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// NewAddress parses the `ip`, `port` and `family` and initializes the Address object
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

func (a *Address) Family() (Family, error) {
	f := C.roc_address_family(a.raw)
	family := (Family)(f)
	if family == AfInvalid {
		return family, ErrInvalidArguments
	}
	return family, nil
}

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

func (a *Address) Port() (int, error) {
	res := C.roc_address_port(a.raw)
	if res < 0 {
		return (int)(res), ErrInvalidArguments
	}
	return (int)(res), nil
}
