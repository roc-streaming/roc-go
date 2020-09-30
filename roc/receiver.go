package roc

/*
 #include <roc/receiver.h>
 #include <roc/config.h>
int rocGoReceiverReadFloats(roc_receiver* receiver, float* samples, unsigned long samples_size) {
    roc_frame frame = {(void*)samples, samples_size*sizeof(float)};
    return roc_receiver_read(receiver, &frame);
}
*/
import "C"

import (
	"fmt"
)

// Receiver as declared in roc/receiver.h:117
type Receiver C.roc_receiver

func OpenReceiver(rocContext *Context, receiverConfig *ReceiverConfig) (*Receiver, error) {
	receiverConfigC := C.struct_roc_receiver_config{
		frame_sample_rate:         (C.uint)(receiverConfig.FrameSampleRate),
		frame_channels:            (C.roc_channel_set)(receiverConfig.FrameChannels),
		frame_encoding:            (C.roc_frame_encoding)(receiverConfig.FrameEncoding),
		automatic_timing:          boolToUint(receiverConfig.AutomaticTiming),
		resampler_profile:         (C.roc_resampler_profile)(receiverConfig.ResamplerProfile),
		target_latency:            (C.ulonglong)(receiverConfig.TargetLatency),
		max_latency_overrun:       (C.ulonglong)(receiverConfig.MaxLatencyOverrun),
		max_latency_underrun:      (C.ulonglong)(receiverConfig.MaxLatencyUnderrun),
		no_playback_timeout:       (C.longlong)(receiverConfig.NoPlaybackTimeout),
		broken_playback_timeout:   (C.longlong)(receiverConfig.BrokenPlaybackTimeout),
		breakage_detection_window: (C.ulonglong)(receiverConfig.BreakageDetectionWindow),
	}

	receiver := C.roc_receiver_open((*C.roc_context)(rocContext), &receiverConfigC)
	if receiver == nil {
		return nil, ErrInvalidArgs
	}
	return (*Receiver)(receiver), nil
}

func (r *Receiver) Bind(portType PortType, proto Protocol, a *Address) error {
	errCode := C.roc_receiver_bind(
		(*C.roc_receiver)(r),
		(C.roc_port_type)(portType),
		(C.roc_protocol)(proto),
		a.raw)
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return ErrInvalidArgs
	}
	panic(fmt.Sprintf(
		"unexpected return code %d from roc_receiver_bind()", errCode))
}

func (r *Receiver) ReadFloats(frame []float32) error {
	if frame == nil {
		return ErrInvalidArgs
	}
	if len(frame) == 0 {
		return nil
	}
	errCode := C.rocGoReceiverReadFloats((*C.roc_receiver)(r), (*C.float)(&frame[0]), (C.ulong)(len(frame)))
	if errCode == 0 {
		return nil
	}

	if errCode < 0 {
		return ErrInvalidArgs
	}
	panic(fmt.Sprintf(
		"unexpected return code %d from roc_receiver_read()", errCode))
}

func (r *Receiver) Close() error {
	errCode := C.roc_receiver_close((*C.roc_receiver)(r))

	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return ErrInvalidArgs
	}

	panic(fmt.Sprintf(
		"unexpected return code %d from roc_receiver_close()", errCode))
}
