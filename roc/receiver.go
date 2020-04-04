package roc

/*
#include <roc/receiver.h>
*/
import "C"

func OpenReceiver(rocContext *Context, receiverConfig *ReceiverConfig) (*Receiver, error) {
	receiverConfigC := C.struct_roc_receiver_config{
		frame_sample_rate:         (C.uint)(receiverConfig.FrameSampleRate),
		frame_channels:            (C.roc_channel_set)(receiverConfig.FrameChannels),
		frame_encoding:            (C.roc_frame_encoding)(receiverConfig.FrameEncoding),
		automatic_timing:          (C.uint)(receiverConfig.AutomaticTiming),
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
		return nil, ErrInvalidArguments
	}
	return (*Receiver)(receiver), nil
}

func (r *Receiver) Close() error {
	errCode := C.roc_receiver_close((*C.roc_receiver)(r))
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return ErrInvalidArguments
	}
	return ErrInvalidApi
}
