package roc

/*
#include "roc/config.h"
*/
import "C"

// PortType as declared in roc/config.h:35
type PortType int32

// PortType enumeration from roc/config.h:35
const (
	PortAudioSource PortType = 1
	PortAudioRepair PortType = 2
)

// Protocol as declared in roc/config.h:54
type Protocol int32

// Protocol enumeration from roc/config.h:54
const (
	ProtoRtp           Protocol = 1
	ProtoRtpRs8mSource Protocol = 2
	ProtoRs8mRepair    Protocol = 3
	ProtoRtpLdpcSource Protocol = 4
	ProtoLdpcRepair    Protocol = 5
)

// FecCode as declared in roc/config.h:81
type FecCode int32

// FecCode enumeration from roc/config.h:81
const (
	FecDisable       FecCode = -1
	FecDefault       FecCode = 0
	FecRs8m          FecCode = 1
	FecLdpcStaircase FecCode = 2
)

// PacketEncoding as declared in roc/config.h:91
type PacketEncoding int32

// PacketEncoding enumeration from roc/config.h:91
const (
	PacketEncodingAvpL16 PacketEncoding = 2
)

// FrameEncoding as declared in roc/config.h:100
type FrameEncoding int32

// FrameEncoding enumeration from roc/config.h:100
const (
	FrameEncodingPcmFloat FrameEncoding = 1
)

// ChannelSet as declared in roc/config.h:108
type ChannelSet int32

// ChannelSet enumeration from roc/config.h:108
const (
	ChannelSetStereo ChannelSet = 2
)

// ResamplerProfile as declared in roc/config.h:128
type ResamplerProfile int32

// ResamplerProfile enumeration from roc/config.h:128
const (
	ResamplerDisable ResamplerProfile = -1
	ResamplerDefault ResamplerProfile = 0
	ResamplerHigh    ResamplerProfile = 1
	ResamplerMedium  ResamplerProfile = 2
	ResamplerLow     ResamplerProfile = 3
)

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
