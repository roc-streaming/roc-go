package roc

import "testing"

func Test_NewAddress(t *testing.T) {
	testcases := []struct {
		name    string
		ip      string
		wantErr bool
	}{
		{
			name:    "correct ip is parsed with no error",
			ip:      "127.0.0.1",
			wantErr: false,
		},
		{
			name:    "incorrect ip causes negative error",
			ip:      "invalid ip string",
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewAddress(AfAuto, tt.ip, 23456)
			if (err != nil) != tt.wantErr {
				t.Errorf("WantErr: %v, got %v", tt.wantErr, err)
			}
			if err == nil && a == nil {
				t.Errorf("Address is nil, no error reported")
			}
		})
	}
}
