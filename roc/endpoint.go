package roc

/*
#include <roc/endpoint.h>
*/
import "C"

import (
	"fmt"
)

// Network endpoint.
//
// Endpoint is a network entry point of a peer. The definition includes the
// protocol being used, network host and port, and, for some protocols, a
// resource. All these parts together are unambiguously represented
// by a URI. The user may set or get the entire URI or its individual parts.
//
// # Endpoint URI
//
// Endpoint URI syntax is a subset of the syntax defined in RFC 3986:
//
//	protocol://host[:port][/path][?query]
//
// Examples:
//
//	rtsp://localhost:123/path?query
//	rtp+rs8m://localhost:123
//	rtp://127.0.0.1:123
//	rtp://[::1]:123
//
// The following protocols (schemes) are supported:
//
//	rtp://       (ProtoRtp)
//	rtp+rs8m://  (ProtoRtpRs8mSource)
//	rs8m://      (ProtoRs8mRepair)
//	rtp+ldpc://  (ProtoRtpLdpcSource)
//	ldpc://      (ProtoLdpcRepair)
//
// The host field should be either FQDN (domain name), or IPv4 address, or
// IPv6 address in square brackets.
//
// The port field can be omitted if the protocol defines standard port. Otherwise,
// the port can not be omitted. For example, RTSP defines standard port,
// but RTP doesn't.
//
// The path and query fields are allowed only for protocols that support them.
// For example, they're supported by RTSP, but not by RTP.
//
// # Thread safety
//
// Should not be used concurrently.
type Endpoint struct {
	// Protocol of the endpoint (URI scheme).
	// Should be set.
	Protocol Protocol

	// Host IP address or domain.
	// Should be set.
	// To bind to all interfaces, use "0.0.0.0" or "[::]".
	Host string

	// TCP or UDP port number.
	// To bind to random port, use 0.
	// To use default port for specified protocol, use -1.
	// Some protocols don't have default port.
	Port int

	// Resource path.
	// Can be empty.
	// Some protocols don't have default resource component.
	Resource string
}

// ParseEndpoint decomposes URI string into Endpoint instance.
func ParseEndpoint(uri string) (*Endpoint, error) {
	checkVersionFn()

	var errCode C.int

	var cEndp *C.roc_endpoint
	errCode = C.roc_endpoint_allocate(&cEndp)
	if errCode != 0 {
		panic(fmt.Sprintf("roc_endpoint_allocate() failed with code %v", errCode))
	}
	if cEndp == nil {
		panic("roc_endpoint_allocate() returned nil")
	}

	defer func() {
		errCode = C.roc_endpoint_deallocate(cEndp)
		if errCode != 0 {
			panic(fmt.Sprintf("roc_endpoint_deallocate() failed with code %v", errCode))
		}
	}()

	cURI, err := go2cStr(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid uri: %w", err)
	}
	errCode = C.roc_endpoint_set_uri(cEndp, (*C.char)(&cURI[0]))
	if errCode != 0 {
		return nil, newNativeErr("roc_endpoint_set_uri()", errCode)
	}

	endp := new(Endpoint)

	if err := endp.fromC(cEndp); err != nil {
		return nil, err
	}

	return endp, nil
}

// URI composes Endpoint instance into URI string.
func (endp *Endpoint) URI() (string, error) {
	var errCode C.int

	var cEndp *C.roc_endpoint
	errCode = C.roc_endpoint_allocate(&cEndp)
	if errCode != 0 || cEndp == nil {
		panic("roc_endpoint_allocate() failed")
	}

	defer func() {
		errCode = C.roc_endpoint_deallocate(cEndp)
		if errCode != 0 {
			panic("roc_endpoint_deallocate() failed")
		}
	}()

	if err := endp.toC(cEndp); err != nil {
		return "", err
	}

	var cURISize C.size_t
	errCode = C.roc_endpoint_get_uri(cEndp, nil, &cURISize)
	if errCode != 0 {
		return "", newNativeErr("roc_endpoint_get_uri()", errCode)
	}

	cURI := make([]C.char, cURISize)
	errCode = C.roc_endpoint_get_uri(cEndp, (*C.char)(&cURI[0]), &cURISize)
	if errCode != 0 {
		return "", newNativeErr("roc_endpoint_get_uri()", errCode)
	}

	uri := c2goStr(cURI)

	return uri, nil
}

// fills endp from cEndp
func (endp *Endpoint) fromC(cEndp *C.roc_endpoint) error {
	var errCode C.int

	var cProto C.roc_protocol
	errCode = C.roc_endpoint_get_protocol(cEndp, &cProto)
	if errCode != 0 {
		return newNativeErr("roc_endpoint_get_protocol()", errCode)
	}
	endp.Protocol = Protocol(cProto)

	var cHostSize C.size_t
	errCode = C.roc_endpoint_get_host(cEndp, nil, &cHostSize)
	if errCode != 0 {
		return newNativeErr("roc_endpoint_get_host()", errCode)
	}

	cHost := make([]C.char, cHostSize)
	errCode = C.roc_endpoint_get_host(cEndp, (*C.char)(&cHost[0]), &cHostSize)
	if errCode != 0 {
		return newNativeErr("roc_endpoint_get_host()", errCode)
	}
	endp.Host = c2goStr(cHost)

	var cPort C.int
	errCode = C.roc_endpoint_get_port(cEndp, &cPort)
	if errCode == 0 {
		endp.Port = int(cPort)
	} else {
		endp.Port = -1
	}

	var cResourceSize C.size_t
	errCode = C.roc_endpoint_get_resource(cEndp, nil, &cResourceSize)
	if errCode == 0 {
		cResource := make([]C.char, cResourceSize)
		errCode = C.roc_endpoint_get_resource(cEndp, (*C.char)(&cResource[0]), &cResourceSize)
		if errCode != 0 {
			return newNativeErr("roc_endpoint_get_resource()", errCode)
		}
		endp.Resource = c2goStr(cResource)
	}

	return nil
}

// fills cEndp from endp
func (endp *Endpoint) toC(cEndp *C.roc_endpoint) error {
	var errCode C.int

	errCode = C.roc_endpoint_set_protocol(cEndp, C.roc_protocol(endp.Protocol))
	if errCode != 0 {
		return newNativeErr("roc_endpoint_set_protocol()", errCode)
	}

	if endp.Host != "" {
		cHost, err := go2cStr(endp.Host)
		if err != nil {
			return fmt.Errorf("invalid host: %w", err)
		}
		errCode = C.roc_endpoint_set_host(cEndp, (*C.char)(&cHost[0]))
		if errCode != 0 {
			return newNativeErr("roc_endpoint_set_host()", errCode)
		}
	}

	if endp.Port != -1 {
		errCode = C.roc_endpoint_set_port(cEndp, C.int(endp.Port))
		if errCode != 0 {
			return newNativeErr("roc_endpoint_set_port()", errCode)
		}
	}

	if endp.Resource != "" {
		cResource, err := go2cStr(endp.Resource)
		if err != nil {
			return fmt.Errorf("invalid resource: %w", err)
		}
		errCode = C.roc_endpoint_set_resource(cEndp, (*C.char)(&cResource[0]))
		if errCode != 0 {
			return newNativeErr("roc_endpoint_set_resource()", errCode)
		}
	}

	return nil
}
