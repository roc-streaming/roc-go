package roc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringer(t *testing.T) {
	for i := -10; i <= 1000; i++ {
		assert.NotEmpty(t, ChannelLayout(i).String())
		assert.NotEmpty(t, ClockSource(i).String())
		assert.NotEmpty(t, ClockSyncBackend(i).String())
		assert.NotEmpty(t, ClockSyncProfile(i).String())
		assert.NotEmpty(t, FecEncoding(i).String())
		assert.NotEmpty(t, Format(i).String())
		assert.NotEmpty(t, Interface(i).String())
		assert.NotEmpty(t, LogLevel(i).String())
		assert.NotEmpty(t, PacketEncoding(i).String())
		assert.NotEmpty(t, Protocol(i).String())
		assert.NotEmpty(t, ResamplerBackend(i).String())
		assert.NotEmpty(t, ResamplerProfile(i).String())
	}
}
