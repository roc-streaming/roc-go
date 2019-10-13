package roc

import "testing"

func Test_convertErr(t *testing.T) {
	testcases := []struct {
		name    string
		code    int32
		wantErr bool
	}{
		{
			name:    "code 0 means no error",
			code:    0,
			wantErr: false,
		},
		{
			name:    "negative code means error",
			code:    -1,
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			err := convertErr(tt.code, "This error is not expected")
			if (err != nil) != tt.wantErr {
				t.Errorf("Want err: %v, err: %v", tt.wantErr, err)
			}
		})
	}
}
