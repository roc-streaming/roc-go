package roc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndpoint(t *testing.T) {
	tests := []struct {
		name       string
		uri        string
		protocol   Protocol
		host       string
		port       int
		resource   string
		parseErr   error
		composeErr error
	}{
		// protocols
		{
			name:       "rtsp protocol",
			uri:        "rtsp://192.168.0.1:12345/path?query1=query1&query2=query2",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "/path?query1=query1&query2=query2",
			parseErr:   nil,
			composeErr: nil,
		},
		{
			name:       "rtp protocol",
			uri:        "rtp://192.168.0.1:12345",
			protocol:   ProtoRtp,
			host:       "192.168.0.1",
			port:       12345,
			parseErr:   nil,
			composeErr: nil,
		},
		{
			name:       "rtp+rs8m protocol",
			uri:        "rtp+rs8m://192.168.0.1:12345",
			protocol:   ProtoRtpRs8mSource,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "",
			parseErr:   nil,
			composeErr: nil,
		},
		{
			name:       "rs8m protocol",
			uri:        "rs8m://192.168.0.1:12345",
			protocol:   ProtoRs8mRepair,
			host:       "192.168.0.1",
			port:       12345,
			parseErr:   nil,
			composeErr: nil,
		},
		{
			name:       "rtp+ldpc protocol",
			uri:        "rtp+ldpc://192.168.0.1:12345",
			protocol:   ProtoRtpLdpcSource,
			host:       "192.168.0.1",
			port:       12345,
			parseErr:   nil,
			composeErr: nil,
		},
		{
			name:       "ldpc protocol",
			uri:        "ldpc://192.168.0.1:12345",
			protocol:   ProtoLdpcRepair,
			host:       "192.168.0.1",
			port:       12345,
			parseErr:   nil,
			composeErr: nil,
		},
		{
			name:       "rtcp protocol",
			uri:        "rtcp://192.168.0.1:12345",
			protocol:   ProtoRtcp,
			host:       "192.168.0.1",
			port:       12345,
			parseErr:   nil,
			composeErr: nil,
		},
		// components
		{
			name:       "use default rtsp port",
			uri:        "rtsp://192.168.0.1",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       -1,
			resource:   "",
			parseErr:   nil,
			composeErr: nil,
		},
		{
			name:       "rtsp without resource",
			uri:        "rtsp://192.168.0.1:12345",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "",
			parseErr:   nil,
			composeErr: nil,
		},
		{
			name:       "rtsp without query params",
			uri:        "rtsp://192.168.0.1:12345/path",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "/path",
			parseErr:   nil,
			composeErr: nil,
		},
		{
			name:       "use zero port (for bind)",
			uri:        "rtsp://192.168.0.1:0",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       0,
			resource:   "",
			parseErr:   nil,
			composeErr: nil,
		},
		// errors
		{
			name:       "empty uri",
			uri:        "",
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_protocol()", -1),
		},
		{
			name:       "missing host and port",
			uri:        "rtsp://",
			protocol:   ProtoRtsp,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_get_uri()", -1),
		},
		{
			name:       "missing host",
			uri:        "rtsp://:12345",
			protocol:   ProtoRtsp,
			port:       12345,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_get_uri()", -1),
		},
		{
			name:       "port out of range",
			uri:        "rtsp://192.168.0.1:65536",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       655356,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_port()", -1),
		},
		{
			name:       "port out of range - negative",
			uri:        "rtsp://192.168.0.1:-2",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       -2,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_port()", -1),
		},
		{
			name:       "invalid resource",
			uri:        "rtsp://192.168.0.1/??",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       -1,
			resource:   "??",
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_resource()", -1),
		},
		{
			name:       "invalid protocol",
			protocol:   Protocol(1),
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_protocol()", -1),
		},
		{
			name:       "resource not allowed for protocol",
			uri:        "rtp://192.168.0.1:12345/path",
			protocol:   ProtoRtp,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "/path",
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_get_uri()", -1),
		},
		{
			name:       "default port not defined for protocol",
			uri:        "rtp://192.168.0.1",
			protocol:   ProtoRtp,
			host:       "192.168.0.1",
			port:       -1,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_get_uri()", -1),
		},
		{
			name:       "zero byte in uri",
			uri:        "rtsp://192.168.0.1:12345\x00",
			parseErr:   errors.New("invalid uri: "),
			composeErr: newNativeErr("roc_endpoint_set_protocol()", -1),
		},
		{
			name:       "zero byte in host",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1\x00",
			port:       12345,
			resource:   "",
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: errors.New("invalid host: "),
		},
		{
			name:       "zero byte in resource",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "/path\x00",
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: errors.New("invalid resource: "),
		},
	}

	for _, tt := range tests {
		t.Run("parse/"+tt.name, func(t *testing.T) {
			endp, err := ParseEndpoint(tt.uri)

			if tt.parseErr == nil {
				require.NoError(t, err)
				require.NotNil(t, endp)

				assert.Equal(t, tt.protocol, endp.Protocol)
				assert.Equal(t, tt.host, endp.Host)
				assert.Equal(t, tt.port, endp.Port)
				assert.Equal(t, tt.resource, endp.Resource)
			} else {
				require.Contains(t, err.Error(), tt.parseErr.Error())
				require.Nil(t, endp)
			}
		})

		t.Run("compose/"+tt.name, func(t *testing.T) {
			endp := Endpoint{
				Protocol: tt.protocol,
				Host:     tt.host,
				Port:     tt.port,
				Resource: tt.resource,
			}

			uri, err := endp.URI()

			if tt.composeErr == nil {
				require.NoError(t, err)
				require.Equal(t, tt.uri, uri)
			} else {
				require.Contains(t, err.Error(), tt.composeErr.Error())
				require.Empty(t, uri)
			}
		})
	}
}
