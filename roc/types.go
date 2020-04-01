// MIT
package roc

/*
#cgo LDFLAGS: -lroc
#include "roc/receiver.h"
#include "roc/sender.h"
#include "roc/log.h"
#include <stdlib.h>
*/
import "C"

// Address as declared in roc/address.h:59
type Address struct {
	raw *C.roc_address
	mem []byte
}

// ContextConfig as declared in roc/config.h:147
type ContextConfig struct {
	MaxPacketSize uint32
	MaxFrameSize  uint32
}

// SenderConfig as declared in roc/config.h:234
type SenderConfig struct {
	FrameSampleRate       uint32
	FrameChannels         ChannelSet
	FrameEncoding         FrameEncoding
	PacketSampleRate      uint32
	PacketChannels        ChannelSet
	PacketEncoding        PacketEncoding
	PacketLength          uint64
	PacketInterleaving    uint32
	AutomaticTiming       uint32
	ResamplerProfile      ResamplerProfile
	FecCode               FecCode
	FecBlockSourcePackets uint32
	FecBlockRepairPackets uint32
}

// ReceiverConfig as declared in roc/config.h:316
type ReceiverConfig struct {
	FrameSampleRate         uint32
	FrameChannels           ChannelSet
	FrameEncoding           FrameEncoding
	AutomaticTiming         uint32
	ResamplerProfile        ResamplerProfile
	TargetLatency           uint64
	MaxLatencyOverrun       uint64
	MaxLatencyUnderrun      uint64
	NoPlaybackTimeout       int64
	BrokenPlaybackTimeout   int64
	BreakageDetectionWindow uint64
}

// Context as declared in roc/context.h:41
type Context C.roc_context

// LogHandler type as declared in roc/log.h:64
type LogHandler func(level LogLevel, component string, message string)

// Receiver as declared in roc/receiver.h:117
type Receiver C.roc_receiver

// Sender as declared in roc/sender.h:96
type Sender C.roc_sender
