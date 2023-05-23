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
	"errors"
	"fmt"
	"sync"
)

// Receiver peer.
//
// Receiver gets the network packets from multiple senders, decodes audio streams
// from them, mixes multiple streams into a single stream, and returns it to the user.
//
// # Context
//
// Receiver is automatically attached to a context when opened and detached from it when
// closed. The user should not close the context until the receiver is closed.
//
// Receiver work consists of two parts: packet reception and stream decoding. The
// decoding part is performed in the receiver itself, and the reception part is
// performed in the context network worker threads.
//
// # Life cycle
//
// - A receiver is created using OpenReceiver().
//
//   - Optionally, the receiver parameters may be fine-tuned using Receiver.Set*()
//     functions.
//
//   - The receiver either binds local endpoints using Receiver.Bind(), allowing senders
//     connecting to them, or itself connects to remote sender endpoints using
//     Receiver.Connect(). What approach to use is up to the user.
//
//   - The audio stream is iteratively read from the receiver using Receiver.Read*().
//     Receiver returns the mixed stream from all connected senders.
//
// - The receiver is destroyed using Receiver.Close().
//
// The user is responsible for closing any opened receiver before exiting the program.
//
// # Slots, interfaces, and endpoints
//
// Receiver has one or multiple slots, which may be independently bound or connected.
// Slots may be used to bind receiver to multiple addresses. Slots are numbered from
// zero and are created automatically. In simple cases just use SlotDefault.
//
// Each slot has its own set of interfaces, one per each type defined in Interface
// type. The interface defines the type of the communication with the remote peer
// and the set of the protocols supported by it.
//
// Supported actions with the interface:
//
//   - Call Receiver.Bind() to bind the interface to a local Endpoint. In this
//     case the receiver accepts connections from senders mixes their streams into the
//     single output stream.
//
//   - Call Receiver.Connect() to connect the interface to a remote Endpoint.
//     In this case the receiver initiates connection to the sender and requests it
//     to start sending media stream to the receiver.
//
// Supported interface configurations:
//
//   - Bind InterfaceConsolidated to a local endpoint (e.g. be an RTSP server).
//
//   - Connect InterfaceConsolidated to a remote endpoint (e.g. be an RTSP
//     client).
//
//   - Bind InterfaceAudioSource, InterfaceAudioRepair (optionally,
//     for FEC), and InterfaceAudioControl (optionally, for control messages)
//     to local endpoints (e.g. be an RTP/FECFRAME/RTCP receiver).
//
// # FEC scheme
//
// If InterfaceConsolidated is used, it automatically creates all necessary
// transport interfaces and the user should not bother about them.
//
// Otherwise, the user should manually configure InterfaceAudioSource and
// InterfaceAudioRepair interfaces:
//
//   - If FEC is disabled (FecEncodingDisable), only
//     InterfaceAudioSource should be configured. It will be used to transmit
//     audio packets.
//
//   - If FEC is enabled, both InterfaceAudioSource and
//     InterfaceAudioRepair interfaces should be configured. The second interface
//     will be used to transmit redundant repair data.
//
// The protocols for the two interfaces should correspond to each other and to the FEC
// scheme. For example, if FecEncodingRs8m is used, the protocols should be
// ProtoRtpRs8mSource and ProtoRs8mRepair.
//
// # Sessions
//
// Receiver creates a session object for every sender connected to it. Sessions can appear
// and disappear at any time. Multiple sessions can be active at the same time.
//
// A session is identified by the sender address. A session may contain multiple packet
// streams sent to different receiver ports. If the sender employs FEC, the session will
// contain source and repair packet streams. Otherwise, the session will contain a single
// source packet stream.
//
// A session is created automatically on the reception of the first packet from a new
// address and destroyed when there are no packets during a timeout. A session is also
// destroyed on other events like a large latency underrun or overrun or broken playback,
// but if the sender continues to send packets, it will be created again shortly.
//
// # Mixing
//
// Receiver mixes audio streams from all currently active sessions into a single output
// stream.
//
// The output stream continues no matter how much active sessions there are at the moment.
// In particular, if there are no sessions, the receiver produces a stream with all zeros.
//
// Sessions can be added and removed from the output stream at any time, probably in the
// middle of a frame.
//
// # Sample rate
//
// Every session may have a different sample rate. And even if nominally all of them are
// of the same rate, device frequencies usually differ by a few tens of Hertz.
//
// Receiver compensates these differences by adjusting the rate of every session stream to
// the rate of the receiver output stream using a per-session resampler. The frequencies
// factor between the sender and the receiver clocks is calculated dynamically for every
// session based on the session incoming packet queue size.
//
// Resampling is a quite time-consuming operation. The user can choose between completely
// disabling resampling (at the cost of occasional underruns or overruns) or several
// resampler profiles providing different compromises between CPU consumption and quality.
//
// # Clock source
//
// Receiver should decode samples at a constant rate that is configured when the receiver
// is created. There are two ways to accomplish this:
//
//   - If the user enabled internal clock (ClockInternal), the receiver employs a
//     CPU timer to block reads until it's time to decode the next bunch of samples
//     according to the configured sample rate.
//
//     This mode is useful when the user passes samples to a non-realtime destination,
//     e.g. to an audio file.
//
//   - If the user enabled external clock (ClockExternal), the samples read from
//     the receiver are decoded immediately and hence the user is responsible to call
//     read operation according to the sample rate.
//
//     This mode is useful when the user passes samples to a realtime destination with its
//     own clock, e.g. to an audio device. Internal clock should not be used in this case
//     because the audio device and the CPU might have slightly different clocks, and the
//     difference will eventually lead to an underrun or an overrun.
//
// # Thread safety
//
// Can be used concurrently.
type Receiver struct {
	mu   sync.RWMutex
	cPtr *C.roc_receiver
}

// Open a new receiver.
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

	cConfig := C.struct_roc_receiver_config{
		frame_sample_rate:         (C.uint)(config.FrameSampleRate),
		frame_channels:            (C.roc_channel_set)(config.FrameChannels),
		frame_encoding:            (C.roc_frame_encoding)(config.FrameEncoding),
		clock_source:              (C.roc_clock_source)(config.ClockSource),
		resampler_backend:         (C.roc_resampler_backend)(config.ResamplerBackend),
		resampler_profile:         (C.roc_resampler_profile)(config.ResamplerProfile),
		target_latency:            (C.ulonglong)(config.TargetLatency),
		max_latency_overrun:       (C.ulonglong)(config.MaxLatencyOverrun),
		max_latency_underrun:      (C.ulonglong)(config.MaxLatencyUnderrun),
		no_playback_timeout:       (C.longlong)(config.NoPlaybackTimeout),
		broken_playback_timeout:   (C.longlong)(config.BrokenPlaybackTimeout),
		breakage_detection_window: (C.ulonglong)(config.BreakageDetectionWindow),
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

// Set receiver interface multicast group.
//
// Optional.
//
// Multicast group should be set only when binding receiver interface to an endpoint with
// multicast IP address. If present, it defines an IP address of the OS network interface
// on which to join the multicast group. If not present, no multicast group is joined.
//
// It's possible to receive multicast traffic from only those OS network interfaces, on
// which the process has joined the multicast group. When using multicast, the user should
// either call this function, or join multicast group manually using OS-specific API.
//
// It is allowed to set multicast group to `0.0.0.0` (for IPv4) or to `::` (for IPv6),
// to be able to receive multicast traffic from all available interfaces. However, this
// may not be desirable for security reasons.
//
// Each slot's interface can have only one multicast group. The function should be called
// before calling roc_receiver_bind() for the interface. It should not be called when
// calling Receiver.Connect() for the interface.
//
// Automatically initializes slot with given index if it's used first time.
func (r *Receiver) SetMulticastGroup(slot Slot, iface Interface, ip string) (err error) {
	logWrite(LogDebug,
		"entering Receiver.SetMulticastGroup(): receiver=%p slot=%v iface=%v ip=%v", r, slot, iface, ip,
	)
	defer func() {
		logWrite(LogDebug, "leaving Receiver.SetMulticastGroup(): receiver=%p err=%#v", r, err)
	}()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.cPtr == nil {
		return errors.New("receiver is closed")
	}

	cIP, parseErr := go2cStr(ip)
	if parseErr != nil {
		return fmt.Errorf("invalid ip: %w", parseErr)
	}
	errCode := C.roc_receiver_set_multicast_group(
		r.cPtr,
		(C.roc_slot)(slot),
		(C.roc_interface)(iface),
		(*C.char)(&cIP[0]))
	if errCode != 0 {
		return newNativeErr("roc_receiver_set_multicast_group()", errCode)
	}

	return nil
}

// Set receiver interface address reuse option.
//
// Optional.
//
// When set to true, SO_REUSEADDR is enabled for interface socket, regardless of socket
// type, unless binding to ephemeral port (port explicitly set to zero).
//
// When set to false, SO_REUSEADDR is enabled only for multicast sockets, unless binding
// to ephemeral port (port explicitly set to zero).
//
// By default set to false.
//
// For TCP-based protocols, SO_REUSEADDR allows immediate reuse of recently closed socket
// in TIME_WAIT state, which may be useful you want to be able to restart server quickly.
//
// For UDP-based protocols, SO_REUSEADDR allows multiple processes to bind to the same
// address, which may be useful if you're using socket activation mechanism.
//
// Automatically initializes slot with given index if it's used first time.
func (r *Receiver) SetReuseaddr(slot Slot, iface Interface, enabled bool) (err error) {
	logWrite(LogDebug,
		"entering Receiver.SetReuseaddr(): receiver=%p slot=%v iface=%v enabled=%v",
		r, slot, iface, enabled,
	)
	defer func() {
		logWrite(LogDebug, "leaving Receiver.SetReuseaddr(): receiver=%p err=%#v", r, err)
	}()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.cPtr == nil {
		return errors.New("receiver is closed")
	}

	cEnabled := go2cBool(enabled)

	errCode := C.roc_receiver_set_reuseaddr(
		r.cPtr,
		(C.roc_slot)(slot),
		(C.roc_interface)(iface),
		(C.int)(cEnabled),
	)

	if errCode != 0 {
		return newNativeErr("roc_receiver_set_reuseaddr()", errCode)
	}

	return nil
}

// Bind the receiver interface to a local endpoint.
//
// Checks that the endpoint is valid and supported by the interface, allocates
// a new ingoing port, and binds it to the local endpoint.
//
// Each slot's interface can be bound or connected only once.
// May be called multiple times for different slots or interfaces.
//
// Automatically initializes slot with given index if it's used first time.
//
// If endpoint has explicitly set zero port, the receiver is bound to a randomly
// chosen ephemeral port. If the function succeeds, the actual port to which the
// receiver was bound is written back to endpoint.
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

// Read samples from the receiver.
//
// Reads network packets received on bound ports, routes packets to sessions, repairs lost
// packets, decodes samples, resamples and mixes them, and finally stores samples into the
// provided frame.
//
// If ClockInternal is used, the function blocks until it's time to decode the
// samples according to the configured sample rate.
//
// Until the receiver is connected to at least one sender, it produces silence.
// If the receiver is connected to multiple senders, it mixes their streams into one.
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
// Deinitializes and deallocates the receiver, and detaches it from the context. The user
// should ensure that nobody uses the receiver during and after this call. If this
// function fails, the receiver is kept opened and attached to the context.
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
