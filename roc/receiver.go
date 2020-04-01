package roc

/*
#cgo LDFLAGS: -lroc
#include <roc/receiver.h>
#include <roc/address.h>
#include <roc/sender.h>
#include <roc/log.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

func OpenReceiver(rocContext *Context, receiverConfig *ReceiverConfig) (*Receiver, error) {
	rocContextCPtr := (*C.roc_context)(unsafe.Pointer(rocContext))
	receiverConfigC := C.struct_roc_receiver_config{
		frame_sample_rate:         (C.uint)(receiverConfig.FrameSampleRate),
		frame_channels:            (C.roc_channel_set)(receiverConfig.FrameChannels),
		frame_encoding:            (C.roc_frame_encoding)(receiverConfig.FrameEncoding),
		automatic_timing:          (C.uint)(receiverConfig.AutomaticTiming),
		resampler_profile:         (C.roc_resampler_profile)(receiverConfig.ResamplerProfile),
		target_latency:            (C.ulonglong)(receiverConfig.TargetLatency),
		max_latency_overrun:       (C.ulonglong)(receiverConfig.MaxLatencyOverrun),
		no_playback_timeout:       (C.longlong)(receiverConfig.NoPlaybackTimeout),
		broken_playback_timeout:   (C.longlong)(receiverConfig.BrokenPlaybackTimeout),
		breakage_detection_window: (C.ulonglong)(receiverConfig.BreakageDetectionWindow),
	}

	receiver := C.roc_receiver_open(rocContextCPtr, &receiverConfigC)
	if receiver == nil {
		return nil, ErrInvalidArguments
	}
	return (*Receiver)(receiver), nil
}

func (r *Receiver) Close() error {
	errCode := C.roc_receiver_close((*C.roc_receiver)(unsafe.Pointer(r)))
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return ErrInvalidArguments
	}
	return ErrInvalidApi
}
