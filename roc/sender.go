package roc

/*
#include <roc/sender.h>
int rocGoSenderWriteFloats(roc_sender* sender, float* samples, unsigned long samples_size) {
    roc_frame frame = {(void*)samples, samples_size*sizeof(float)};
    return roc_sender_write(sender, &frame);
 }
*/
import "C"

import (
	"fmt"
)

// Sender as declared in roc/sender.h:96
type Sender C.roc_sender

func OpenSender(rocContext *Context, senderConfig *SenderConfig) (*Sender, error) {
	senderConfigC := C.struct_roc_sender_config{
		frame_sample_rate:        (C.uint)(senderConfig.FrameSampleRate),
		frame_channels:           (C.roc_channel_set)(senderConfig.FrameChannels),
		frame_encoding:           (C.roc_frame_encoding)(senderConfig.FrameEncoding),
		packet_sample_rate:       (C.uint)(senderConfig.PacketSampleRate),
		packet_channels:          (C.roc_channel_set)(senderConfig.PacketChannels),
		packet_encoding:          (C.roc_packet_encoding)(senderConfig.PacketEncoding),
		packet_length:            (C.ulonglong)(senderConfig.PacketLength),
		packet_interleaving:      boolToUint(senderConfig.PacketInterleaving),
		automatic_timing:         boolToUint(senderConfig.AutomaticTiming),
		resampler_profile:        (C.roc_resampler_profile)(senderConfig.ResamplerProfile),
		fec_code:                 (C.roc_fec_code)(senderConfig.FecCode),
		fec_block_source_packets: (C.uint)(senderConfig.FecBlockSourcePackets),
		fec_block_repair_packets: (C.uint)(senderConfig.FecBlockRepairPackets),
	}
	sender := C.roc_sender_open((*C.roc_context)(rocContext), &senderConfigC)
	if sender == nil {
		return nil, ErrInvalidArgs
	}
	return (*Sender)(sender), nil
}

func (s *Sender) Bind(a *Address) error {
	errCode := C.roc_sender_bind((*C.roc_sender)(s), a.raw)
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return ErrInvalidArgs
	}
	panic(fmt.Sprintf(
		"unexpected return code %d from roc_sender_bind()", errCode))
}

func (s *Sender) Connect(portType PortType, proto Protocol, address *Address) error {
	errCode := C.roc_sender_connect(
		(*C.roc_sender)(s),
		(C.roc_port_type)(portType),
		(C.roc_protocol)(proto),
		address.raw,
	)
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return ErrInvalidArgs
	}
	panic(fmt.Sprintf(
		"unexpected return code %d from roc_sender_connect()", errCode))
}

func (s *Sender) WriteFloats(frame []float32) error {
	if frame == nil {
		return ErrInvalidArgs
	}
	if len(frame) == 0 {
		return nil
	}
	errCode := C.rocGoSenderWriteFloats(
		(*C.roc_sender)(s), (*C.float)(&frame[0]), (C.ulong)(len(frame)))
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return ErrInvalidArgs
	}
	panic(fmt.Sprintf(
		"unexpected return code %d from roc_sender_write()", errCode))
}

func (s *Sender) Close() error {
	errCode := C.roc_sender_close((*C.roc_sender)(s))
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return ErrInvalidArgs
	}
	panic(fmt.Sprintf(
		"unexpected return code %d from roc_sender_close()", errCode))
}
