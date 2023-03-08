package roc

import (
	"errors"
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
		// protocols
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
			uri:        "rtp://192.168.0.1:12345",
			protocol:   ProtoRtp,
			host:       "192.168.0.1",
			port:       12345,
			parseErr:   nil,
			composeErr: nil,
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
		// components
		{
			uri:        "rtsp://192.168.0.1",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       -1, // use default rtsp port
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
		/* FIXME: uncomment after https://github.com/roc-streaming/roc-toolkit/issues/519
			{
				uri:        "rtsp://192.168.0.1:0",
				protocol:   ProtoRtsp,
				host:       "192.168.0.1",
				port:       0, // use zero port (for bind)
				resource:   "",
				parseErr:   nil,
				composeErr: nil,
			},
		*/
		// errors
		{
			uri:        "", // empty uri
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_protocol()", -1),
		},
		{
			uri:        "rtsp://", // missing host and port
			protocol:   ProtoRtsp,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_get_uri()", -1),
		},
		{
			uri:        "rtsp://:12345", // missing host
			protocol:   ProtoRtsp,
			port:       12345,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_get_uri()", -1),
		},
		{
			uri:        "rtsp://192.168.0.1:65536", // port out of range
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       655356,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_port()", -1),
		},
		{
			uri:        "rtsp://192.168.0.1:-2", // port out of range
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       -2,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_port()", -1),
		},
		{
			uri:        "rtsp://192.168.0.1/??", // invalid resource
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       -1,
			resource:   "??",
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_resource()", -1),
		},
		{
			protocol:   Protocol(1), // invalid protocol
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_set_protocol()", -1),
		},
		{
			uri:        "rtp://192.168.0.1:12345/path", // resource not allowed for protocol
			protocol:   ProtoRtp,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "/path",
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_get_uri()", -1),
		},
		{
			uri:        "rtp://192.168.0.1", // default port not defined for protocol
			protocol:   ProtoRtp,
			host:       "192.168.0.1",
			port:       -1,
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: newNativeErr("roc_endpoint_get_uri()", -1),
		},
		{
			uri:        "rtsp://192.168.0.1:12345\x00", // zero byte in uri
			parseErr:   errors.New("invalid uri: "),
			composeErr: newNativeErr("roc_endpoint_set_protocol()", -1),
		},
		{
			uri:        "1",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1\x00", // zero byte in host
			port:       12345,
			resource:   "",
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: errors.New("invalid host: "),
		},
		{
			uri:        "2",
			protocol:   ProtoRtsp,
			host:       "192.168.0.1",
			port:       12345,
			resource:   "/path\x00", // zero byte in resource
			parseErr:   newNativeErr("roc_endpoint_set_uri()", -1),
			composeErr: errors.New("invalid resource: "),
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
				require.Contains(t, err.Error(), tt.parseErr.Error())
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
				require.Contains(t, err.Error(), tt.composeErr.Error())
				require.Empty(t, uri)
			}
		})
	}
}
