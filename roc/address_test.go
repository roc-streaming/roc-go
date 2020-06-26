package roc

import "testing"

func Test_roc_address_init(t *testing.T) {
	tests := []struct {
		f    Family
		ip   string
		port int
		err  error
	}{
		{
			f:    AfIPv4,
			ip:   "192.168.0.1",
			port: 4567,
			err:  nil,
		},
		{
			f:    AfIPv6,
			ip:   "2001:db8:85a3:1:2:8a2e:37:7334",
			port: 9858,
			err:  nil,
		},
		{
			f:    AfIPv6,
			ip:   "192.168.0.1",
			port: 9858,
			err:  ErrInvalidArgs,
		},
	}

	for _, tt := range tests {
		a, err := NewAddress(tt.f, tt.ip, tt.port)

		if err != tt.err {
			fail(tt.err, err, t)
		}

		if err != nil {
			continue
		} // negative test, err is not nil, skipping next tests

		if a == nil {
			fail("Address initialized", "Address is nil", t)
		}

		fam := a.Family()
		if fam != tt.f {
			fail(tt.f, fam, t)
		}

		ip := a.IP()
		if ip != tt.ip {
			fail(tt.ip, ip, t)
		}

		port := a.Port()
		if port != tt.port {
			fail(tt.port, port, t)
		}
	}
}
