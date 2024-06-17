package roc

/*
#include <roc/receiver.h>

int rocGoSetOutgoingAddress(roc_interface_config* config, const char* value);
int rocGoSetMulticastGroup(roc_interface_config* config, const char* value);
int rocGoReceiverReadFloats(roc_receiver* receiver, float* samples, unsigned long samples_size);
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
)

// Receiver peer.
//
// Receiver gets the network packets from multiple senders, decodes audio
// streams from them, mixes multiple streams into a single stream, and returns
// it to the user.
//
// # Context
//
// Receiver is automatically attached to a context when opened and detached from
// it when closed. The user should not close the context until the receiver is
// closed.
//
// Receiver work consists of two parts: packet reception and stream decoding.
// The decoding part is performed in the receiver itself, and the reception part
// is performed in the context network worker threads.
//
// # Life cycle
//
//  - A receiver is created using OpenReceiver().
//  - Optionally, the receiver parameters may be fine-tuned using
//    Receiver.Configure().
//  - The receiver either binds local endpoints using Receiver.Bind(), allowing
//    senders connecting to them, or itself connects to remote sender endpoints
//    using Receiver.Connect(). What approach to use is up to the user.
//  - The audio stream is iteratively read from the receiver using
//    Receiver.ReadFloats(). Receiver returns the mixed stream from all connected
//    senders.
//  - The receiver is destroyed using Receiver.Close().
//
// # Slots, interfaces, and endpoints
//
// Receiver has one or multiple slots, which may be independently bound or
// connected. Slots may be used to bind receiver to multiple addresses. Slots
// are numbered from zero and are created automatically. In simple cases just
// use SlotDefault.
//
// Each slot has its own set of interfaces, one per each type defined in
// Interface. The interface defines the type of the communication with the
// remote peer and the set of the protocols supported by it.
//
// Supported actions with the interface:
//
//  - Call Receiver.Bind() to bind the interface to a local Endpoint. In this
//    case the receiver accepts connections from senders mixes their streams
//    into the single output stream.
//  - Call Receiver.Connect() to connect the interface to a remote Endpoint.
//    In this case the receiver initiates connection to the sender and requests
//    it to start sending media stream to the receiver.
//
// Supported interface configurations:
//
//  - Bind InterfaceConsolidated to a local endpoint (e.g. be an RTSP server).
//  - Connect InterfaceConsolidated to a remote endpoint (e.g. be an RTSP
//    client).
//  - Bind InterfaceAudioSource, InterfaceAudioRepair (optionally, for FEC), and
//    InterfaceAudioControl (optionally, for control messages) to local
//    endpoints (e.g. be an RTP/FECFRAME/RTCP receiver).
//
// Slots can be removed using Receiver.Unlink(). Removing a slot also removes all
// its interfaces and terminates all associated connections.
//
// Slots can be added and removed at any time on fly and from any thread. It is
// safe to do it from another thread concurrently with reading frames.
// Operations with slots won't block concurrent reads.
//
// # FEC schemes
//
// If InterfaceConsolidated is used, it automatically creates all necessary
// transport interfaces and the user should not bother about them.
//
// Otherwise, the user should manually configure InterfaceAudioSource and
// InterfaceAudioRepair interfaces:
//
//  - If FEC is disabled (FecEncodingDisable), only InterfaceAudioSource should
//    be configured. It will be used to transmit audio packets.
//  - If FEC is enabled, both InterfaceAudioSource and InterfaceAudioRepair
//    interfaces should be configured. The second interface will be used to
//    transmit redundant repair data.
//
// The protocols for the two interfaces should correspond to each other and to
// the FEC scheme. For example, if FecEncodingRs8m is used, the protocols should
// be ProtoRtpRs8mSource and ProtoRs8mRepair.
//
// # Connections
//
// Receiver creates a connection object for every sender connected to it.
// Connections can appear and disappear at any time. Multiple connections can be
// active at the same time.
//
// A connection may contain multiple streams sent to different receiver ports.
// If the sender employs FEC, connection usually has source, repair, and control
// streams. Otherwise, connection usually has source and control streams.
//
// Connection is created automatically on the reception of the first packet from
// a new sender, and terminated when there are no packets during a timeout.
// Connection can also be terminated on other events like a large latency
// underrun or overrun or continuous stuttering, but if the sender continues to
// send packets, connection will be created again shortly.
//
// # Mixing
//
// Receiver mixes audio streams from all currently active connections into a
// single output stream.
//
// The output stream continues no matter how much active connections there are
// at the moment. In particular, if there are no connections, the receiver
// produces a stream with all zeros.
//
// Connections can be added and removed from the output stream at any time,
// probably in the middle of a frame.
//
// # Transcoding
//
// Every connection may have a different sample rate, channel layout, and
// encoding.
//
// Before mixing, receiver automatically transcodes all incoming streams to the
// format of receiver frames.
//
// # Latency tuning and bounding
//
// If latency tuning is enabled (which is by default enabled on receiver),
// receiver monitors latency of each connection and adjusts per-connection clock
// to keep latency close to the target value. The user can configure how the
// latency is measured, how smooth is the tuning, and the target value.
//
// If latency bounding is enabled (which is also by default enabled on
// receiver), receiver also ensures that latency lies within allowed boundaries,
// and terminates connection otherwise. The user can configure those boundaries.
//
// To adjust connection clock, receiver uses resampling with a scaling factor
// slightly above or below 1.0. Since resampling may be a quite time-consuming
// operation, the user can choose between several resampler backends and
// profiles providing different compromises between CPU consumption, quality,
// and precision.
//
// # Clock source
//
// Receiver should decode samples at a constant rate that is configured when the
// receiver is created. There are two ways to accomplish this:
//
//  - If the user enabled internal clock (ClockSourceInternal), the receiver
//    employs a CPU timer to block reads until it's time to decode the next
//    bunch of samples according to the configured sample rate. This mode is
//    useful when the user passes samples to a non-realtime destination, e.g. to
//    an audio file.
//  - If the user enabled external clock (ClockSourceExternal), the samples read
//    from the receiver are decoded immediately and hence the user is
//    responsible to call read operation according to the sample rate. This mode
//    is useful when the user passes samples to a realtime destination with its
//    own clock, e.g. to an audio device. Internal clock should not be used in
//    this case because the audio device and the CPU might have slightly
//    different clocks, and the difference will eventually lead to an underrun
//    or an overrun.
//
// # Thread safety
//
// Can be used concurrently.
type Receiver struct {
	mu   sync.RWMutex
	cPtr *C.roc_receiver
}

// Open a new receiver.
//
// Allocates and initializes a new receiver, and attaches it to the context.
func OpenReceiver(context *Context, config ReceiverConfig) (receiver *Receiver, err error) {
	logWrite(LogDebug, "entering OpenReceiver(): context=%p config=%+v", context, config)
	defer func() {
		logWrite(LogDebug,
			"leaving OpenReceiver(): context=%p receiver=%p err=%#v", context, receiver, err,
		)
	}()

	checkVersionFn()

	if context == nil {
		return nil, errors.New("context is nil")
	}

	context.mu.RLock()
	defer context.mu.RUnlock()

	if context.cPtr == nil {
		return nil, errors.New("context is closed")
	}

	cTargetLatency, err := go2cUnsignedDuration(config.TargetLatency)
	if err != nil {
		return nil, fmt.Errorf("invalid config.TargetLatency: %w", err)
	}

	cLatencyTolerance, err := go2cUnsignedDuration(config.LatencyTolerance)
	if err != nil {
		return nil, fmt.Errorf("invalid config.LatencyTolerance: %w", err)
	}

	cConfig := C.struct_roc_receiver_config{
		frame_encoding: C.struct_roc_media_encoding{
			rate:     C.uint(config.FrameEncoding.Rate),
			format:   C.roc_format(config.FrameEncoding.Format),
			channels: C.roc_channel_layout(config.FrameEncoding.Channels),
			tracks:   C.uint(config.FrameEncoding.Tracks),
		},
		clock_source:            C.roc_clock_source(config.ClockSource),
		latency_tuner_backend:   C.roc_latency_tuner_backend(config.LatencyTunerBackend),
		latency_tuner_profile:   C.roc_latency_tuner_profile(config.LatencyTunerProfile),
		resampler_backend:       C.roc_resampler_backend(config.ResamplerBackend),
		resampler_profile:       C.roc_resampler_profile(config.ResamplerProfile),
		target_latency:          cTargetLatency,
		latency_tolerance:       cLatencyTolerance,
		no_playback_timeout:     go2cSignedDuration(config.NoPlaybackTimeout),
		choppy_playback_timeout: go2cSignedDuration(config.ChoppyPlaybackTimeout),
	}

	var cRecv *C.roc_receiver
	errCode := C.roc_receiver_open(context.cPtr, &cConfig, &cRecv)
	if errCode != 0 {
		return nil, newNativeErr("roc_receiver_open()", errCode)
	}
	if cRecv == nil {
		panic("roc_receiver_open() returned nil")
	}

	receiver = &Receiver{
		cPtr: cRecv,
	}

	return receiver, nil
}

// Set receiver interface configuration.
//
// Updates configuration of specified interface of specified slot. If called,
// the call should be done before calling Receiver.Bind() or
// Receiver.Connect() for the same interface.
//
// Automatically initializes slot with given index if it's used first time.
//
// If an error happens during configure, the whole slot is disabled and marked
// broken. The slot index remains reserved. The user is responsible for removing
// the slot using Receiver.Unlink(), after which slot index can be reused.
func (r *Receiver) Configure(slot Slot, iface Interface, config InterfaceConfig) (err error) {
	logWrite(LogDebug,
		"entering Receiver.Configure(): receiver=%p slot=%+v iface=%+v config=%+v",
		r, slot, iface, config,
	)
	defer func() {
		logWrite(LogDebug, "leaving Receiver.Configure(): receiver=%p err=%#v", r, err)
	}()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.cPtr == nil {
		return errors.New("receiver is closed")
	}

	cOutgoingAddress, err := go2cStr(config.OutgoingAddress)
	if err != nil {
		return fmt.Errorf("invalid config.OutgoingAddress: %w", err)
	}

	cMulticastGroup, err := go2cStr(config.MulticastGroup)
	if err != nil {
		return fmt.Errorf("invalid config.MulticastGroup: %w", err)
	}

	var cConfig C.struct_roc_interface_config
	var errCode C.int

	errCode = C.rocGoSetOutgoingAddress(&cConfig, (*C.char)(&cOutgoingAddress[0]))
	if errCode != 0 {
		return fmt.Errorf("invalid config.OutgoingAddress: too long")
	}

	errCode = C.rocGoSetMulticastGroup(&cConfig, (*C.char)(&cMulticastGroup[0]))
	if errCode != 0 {
		return fmt.Errorf("invalid config.MulticastGroup: too long")
	}

	cConfig.reuse_address = (C.int)(go2cBool(config.ReuseAddress))

	errCode = C.roc_receiver_configure(
		r.cPtr,
		(C.roc_slot)(slot),
		(C.roc_interface)(iface),
		&cConfig)
	if errCode != 0 {
		return newNativeErr("roc_receiver_configure()", errCode)
	}

	return nil
}

// Bind the receiver interface to a local endpoint.
//
// Checks that the endpoint is valid and supported by the interface, allocates a
// new ingoing port, and binds it to the local endpoint.
//
// Each slot's interface can be bound or connected only once. May be called
// multiple times for different slots or interfaces.
//
// Automatically initializes slot with given index if it's used first time.
//
// If an error happens during bind, the whole slot is disabled and marked
// broken. The slot index remains reserved. The user is responsible for removing
// the slot using Receiver.Unlink(), after which slot index can be reused.
//
// If Endpoint has explicitly set zero port, the receiver is bound to a randomly
// chosen ephemeral port. If the function succeeds, the actual port to which the
// receiver was bound is written back to Endpoint.
func (r *Receiver) Bind(slot Slot, iface Interface, endpoint *Endpoint) (err error) {
	logWrite(LogDebug,
		"entering Receiver.Bind(): receiver=%p slot=%v iface=%v endpoint=%+v", r, slot, iface, endpoint,
	)
	defer func() {
		logWrite(LogDebug,
			"leaving Receiver.Bind(): receiver=%p endpoint=%+v err=%#v", r, endpoint, err,
		)
	}()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.cPtr == nil {
		return errors.New("receiver is closed")
	}

	if endpoint == nil {
		return errors.New("endpoint is nil")
	}

	var errCode C.int

	var cEndp *C.roc_endpoint
	errCode = C.roc_endpoint_allocate(&cEndp)
	if errCode != 0 {
		panic(fmt.Sprintf("roc_endpoint_allocate() failed with code %v", errCode))
	}
	if cEndp == nil {
		panic("roc_endpoint_allocate() returned nil")
	}

	defer func() {
		errCode = C.roc_endpoint_deallocate(cEndp)
		if errCode != 0 {
			panic(fmt.Sprintf("roc_endpoint_deallocate() failed with code %v", errCode))
		}
	}()

	if err = endpoint.toC(cEndp); err != nil {
		return err
	}

	errCode = C.roc_receiver_bind(
		r.cPtr,
		(C.roc_slot)(slot),
		(C.roc_interface)(iface),
		cEndp)
	if errCode != 0 {
		return newNativeErr("roc_receiver_bind()", errCode)
	}

	if err = endpoint.fromC(cEndp); err != nil {
		return err
	}

	return nil
}

// Delete receiver slot.
//
// Disconnects, unbinds, and removes all slot interfaces and removes the slot.
// All associated connections to remote peers are properly terminated.
//
// After unlinking the slot, it can be re-created again by re-using slot index.
func (r *Receiver) Unlink(slot Slot) (err error) {
	logWrite(LogDebug,
		"entering Receiver.Unlink(): receiver=%p slot=%+v", r, slot,
	)
	defer func() {
		logWrite(LogDebug, "leaving Receiver.Unlink(): receiver=%p err=%#v", r, err)
	}()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.cPtr == nil {
		return errors.New("receiver is closed")
	}

	var errCode C.int

	errCode = C.roc_receiver_unlink(
		r.cPtr,
		(C.roc_slot)(slot))
	if errCode != 0 {
		return newNativeErr("roc_receiver_unlink()", errCode)
	}

	return nil
}

// Read samples from the receiver.
//
// Reads retrieved network packets, decodes packets, repairs losses, extracts
// samples, adjusts sample rate and channel layout, compensates clock drift,
// mixes samples from all connections, and finally stores samples into the
// provided frame.
//
// If ClockSourceInternal is used, the function blocks until it's time to decode
// the samples according to the configured sample rate.
//
// Until the receiver is connected to at least one sender, it produces silence.
// If the receiver is connected to multiple senders, it mixes their streams into
// one.
func (r *Receiver) ReadFloats(frame []float32) (err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.cPtr == nil {
		return errors.New("receiver is closed")
	}

	if frame == nil {
		return errors.New("frame is nil")
	}

	if len(frame) == 0 {
		return nil
	}

	errCode := C.rocGoReceiverReadFloats(
		r.cPtr, (*C.float)(&frame[0]), (C.ulong)(len(frame)))
	if errCode != 0 {
		return newNativeErr("roc_receiver_read()", errCode)
	}

	return nil
}

// Close the receiver.
//
// Deinitializes and deallocates the receiver, and detaches it from the context.
// The user should ensure that nobody uses the receiver during and after this
// call. If this function fails, the receiver is kept opened and attached to the
// context.
func (r *Receiver) Close() (err error) {
	logWrite(LogDebug, "entering Receiver.Close(): receiver=%p", r)
	defer func() {
		logWrite(LogDebug, "leaving Receiver.Close(): receiver=%p err=%#v", r, err)
	}()

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cPtr != nil {
		errCode := C.roc_receiver_close(r.cPtr)
		if errCode != 0 {
			return newNativeErr("roc_receiver_close()", errCode)
		}

		r.cPtr = nil
	}

	return nil
}
