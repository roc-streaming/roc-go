package roc

import "testing"

func Test_roc_address_init(t *testing.T) {
	tests := []struct {
		f    Family
		ip   string
		port int
	}{
		{
			f:    AfIpv4,
			ip:   "192.168.0.1",
			port: 4567,
		},
		{
			f:    AfIpv6,
			ip:   "2001:db8:85a3:1:2:8a2e:37:7334",
			port: 9858,
		},
	}

	for _, tt := range tests {
		a, err := NewAddress(tt.f, tt.ip, tt.port)

		if err != nil {
			fail(nil, err, t)
		}

		if a == nil {
			fail(nil, "Address is nil", t)
		}

		fam, err := a.Family()
		if err != nil {
			fail(nil, err, t)
		}
		if fam != tt.f {
			fail(tt.f, fam, t)
		}

		ip, err := a.IP()
		if err != nil {
			fail(nil, err, t)
		}
		if ip != tt.ip {
			fail(tt.ip, ip, t)
		}

		port, err := a.Port()
		if err != nil {
			fail(nil, err, t)
		}
		if port != tt.port {
			fail(tt.port, port, t)
		}
	}
}
