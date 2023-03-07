package roc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndpoint(t *testing.T) {
	tests := []struct {
		uri        string
		protocol   Protocol
		host       string
		port       int
		resource   string
		parseErr   error
		composeErr error
	}{
		{
			uri:        "rtsp://192.168.0.1",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       -1,
			resource:   "",
			parseErr:   nil,
			composeErr: nil,
		},
		{
			uri:        "rtsp://192.168.0.1:12345",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "",
			parseErr:   nil,
			composeErr: nil,
		},
		{
			uri:        "rtsp://192.168.0.1:12345/path",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "/path",
			parseErr:   nil,
			composeErr: nil,
		},
		{
			uri:        "rtsp://192.168.0.1:12345/path?query1=query1&query2=query2",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "/path?query1=query1&query2=query2",
			parseErr:   nil,
			composeErr: nil,
		},
		{
			uri:        "",
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_protocol()", -1),
		},
		{
			uri:        "rtsp://",
			protocol:   ProtoRtsp,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_get_uri()", -1),
		},
		{
			uri:        "rtsp://:12345",
			protocol:   ProtoRtsp,
			port:       12345,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_get_uri()", -1),
		},
		{
			uri:        "rtsp://192.168.0.1:65536",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       655356,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_port()", -1),
		},
		{
			uri:        "rtsp://192.168.0.1:-2",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       -2,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_port()", -1),
		},
		{
			uri:        "rtsp://192.168.0.1/??",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       -1,
			resource:   "??",
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_resource()", -1),
		},
		{
			protocol:   Protocol(1),
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_protocol()", -1),
		},
		{
			uri:        "rtp://192.168.0.1:12345",
			protocol:   ProtoRtp,
			host:       "192.168.0.1",
			port:       12345,
			parseErr:   nil,
			composeErr: nil,
		},
		{
			uri:        "rtp://192.168.0.1:12345/path",
			protocol:   ProtoRtp,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "/path",
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_get_uri()", -1),
		},
		{
			uri:        "rtp://192.168.0.1",
			protocol:   ProtoRtp,
			host:       "192.168.0.1",
			port:       -1,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_get_uri()", -1),
		},
		{
			uri:        "rtp+rs8m://192.168.0.1:12345",
			protocol:   ProtoRtpRs8mSource,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "",
			parseErr:   nil,
			composeErr: nil,
		},
		{
			uri:        "rs8m://192.168.0.1:12345",
			protocol:   ProtoRs8mRepair,
			host:       "192.168.0.1",
			port:       12345,
			parseErr:   nil,
			composeErr: nil,
		},
		{
			uri:        "rtp+ldpc://192.168.0.1:12345",
			protocol:   ProtoRtpLdpcSource,
			host:       "192.168.0.1",
			port:       12345,
			parseErr:   nil,
			composeErr: nil,
		},
		{
			uri:        "ldpc://192.168.0.1:12345",
			protocol:   ProtoLdpcRepair,
			host:       "192.168.0.1",
			port:       12345,
			parseErr:   nil,
			composeErr: nil,
		},
		{
			uri:        "rtcp://192.168.0.1:12345",
			protocol:   ProtoRtcp,
			host:       "192.168.0.1",
			port:       12345,
			parseErr:   nil,
			composeErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run("parse/"+tt.uri, func(t *testing.T) {
			endp, err := ParseEndpoint(tt.uri)

			if tt.parseErr == nil {
				require.NoError(t, err)
				require.NotNil(t, endp)

				assert.Equal(t, tt.protocol, endp.Protocol)
				assert.Equal(t, tt.host, endp.Host)
				assert.Equal(t, tt.port, endp.Port)
				assert.Equal(t, tt.resource, endp.Resource)
			} else {
				require.Equal(t, tt.parseErr, err)
				require.Nil(t, endp)
			}
		})

		t.Run("compose/"+tt.uri, func(t *testing.T) {
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
				require.Equal(t, tt.composeErr, err)
				require.Empty(t, uri)
			}
		})
	}
}
