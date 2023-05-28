package roc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringer(t *testing.T) {
	for i := -10; i <= 1000; i++ {
		assert.NotEmpty(
			t,
			Interface(i).String(),
			"Interface String func must always return something",
		)
		assert.NotEmpty(
			t,
			Protocol(i).String(),
			"Protocol String func must always return something",
		)
		assert.NotEmpty(
			t,
			FecEncoding(i).String(),
			"FecEncoding String func must always return something",
		)
		assert.NotEmpty(
			t,
			PacketEncoding(i).String(),
			"PacketEncoding String func must always return something",
		)
		assert.NotEmpty(
			t,
			FrameEncoding(i).String(),
			"FrameEncoding String func must always return something",
		)
		assert.NotEmpty(
			t,
			ChannelSet(i).String(),
			"ChannelSet String func must always return something",
		)
		assert.NotEmpty(
			t,
			ResamplerBackend(i).String(),
			"ResamplerBackend String func must always return something",
		)
		assert.NotEmpty(
			t,
			ResamplerProfile(i).String(),
			"ResamplerProfile String func must always return something",
		)
		assert.NotEmpty(
			t,
			ClockSource(i).String(),
			"ClockSource String func must always return something",
		)
		assert.NotEmpty(
			t,
			LogLevel(i).String(),
			"LogLevel String func must always return something",
		)
	}
}
