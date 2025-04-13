package roc

/*
#include <roc/sender.h>

int rocGoSetOutgoingAddress(roc_interface_config* config, const char* value);
int rocGoSetMulticastGroup(roc_interface_config* config, const char* value);
int rocGoSenderWriteFloats(roc_sender* sender, float* samples, unsigned long samples_size);
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
)

// Sender peer.
//
// Sender gets an audio stream from the user, encodes it into network packets,
// and transmits them to a remote receiver.
//
// # Context
//
// Sender is automatically attached to a context when opened and detached from
// it when closed. The user should not close the context until the sender is
// closed.
//
// Sender work consists of two parts: stream encoding and packet transmission.
// The encoding part is performed in the sender itself, and the transmission
// part is performed in the context network worker threads.
//
// # Life cycle
//
//   - A sender is created using OpenSender().
//   - Optionally, the sender parameters may be fine-tuned using
//     Sender.Configure().
//   - The sender either binds local endpoints using Sender.Bind(), allowing
//     receivers connecting to them, or itself connects to remote receiver
//     endpoints using Sender.Connect(). What approach to use is up to the user.
//   - The audio stream is iteratively written to the sender using Sender.WriteFloats().
//     The sender encodes the stream into packets and send to connected
//     receiver(s).
//   - The sender is destroyed using Sender.Close().
//
// # Slots, interfaces, and endpoints
//
// Sender has one or multiple slots, which may be independently bound or
// connected. Slots may be used to connect sender to multiple receivers. Slots
// are numbered from zero and are created automatically. In simple cases just
// use SlotDefault.
//
// Each slot has its own set of interfaces, one per each type defined in
// Interface. The interface defines the type of the communication with the
// remote peer and the set of the protocols supported by it.
//
// Supported actions with the interface:
//
//   - Call Sender.Bind() to bind the interface to a local Endpoint. In this
//     case the sender accepts connections from receivers and sends media stream
//     to all connected receivers.
//   - Call Sender.Connect() to connect the interface to a remote Endpoint. In
//     this case the sender initiates connection to the receiver and starts
//     sending media stream to it.
//
// Supported interface configurations:
//
//   - Connect InterfaceConsolidated to a remote endpoint (e.g. be an RTSP
//     client).
//   - Bind InterfaceConsolidated to a local endpoint (e.g. be an RTSP server).
//   - Connect InterfaceAudioSource, InterfaceAudioRepair (optionally, for FEC),
//     and InterfaceAudioControl (optionally, for control messages) to remote
//     endpoints (e.g. be an RTP/FECFRAME/RTCP sender).
//
// Slots can be removed using Sender.Unlink(). Removing a slot also removes all
// its interfaces and terminates all associated connections.
//
// Slots can be added and removed at any time on fly and from any thread. It is
// safe to do it from another thread concurrently with writing frames.
// Operations with slots won't block concurrent writes.
//
// # FEC schemes
//
// If InterfaceConsolidated is used, it automatically creates all necessary
// transport interfaces and the user should not bother about them.
//
// Otherwise, the user should manually configure InterfaceAudioSource and
// InterfaceAudioRepair interfaces:
//
//   - If FEC is disabled (FecEncodingDisable), only InterfaceAudioSource should
//     be configured. It will be used to transmit audio packets.
//   - If FEC is enabled, both InterfaceAudioSource and InterfaceAudioRepair
//     interfaces should be configured. The second interface will be used to
//     transmit redundant repair data.
//
// The protocols for the two interfaces should correspond to each other and to
// the FEC scheme. For example, if FecEncodingRs8m is used, the protocols should
// be ProtoRtpRs8mSource and ProtoRs8mRepair.
//
// # Transcoding
//
// If encoding of sender frames and network packets are different, sender
// automatically performs all necessary transcoding.
//
// # Latency tuning and bounding
//
// Usually latency tuning and bounding is done on receiver side, but it's
// possible to disable it on receiver and enable on sender. It is useful if
// receiver is does not support it or does not have enough CPU to do it with
// good quality. This feature requires use of ProtoRtcp to deliver necessary
// latency metrics from receiver to sender.
//
// If latency tuning is enabled (which is by default disabled on sender), sender
// monitors latency and adjusts connection clock to keep latency close to the
// target value. The user can configure how the latency is measured, how smooth
// is the tuning, and the target value.
//
// If latency bounding is enabled (which is also by default disabled on sender),
// sender also ensures that latency lies within allowed boundaries, and restarts
// connection otherwise. The user can configure those boundaries.
//
// To adjust connection clock, sender uses resampling with a scaling factor
// slightly above or below 1.0. Since resampling may be a quite time-consuming
// operation, the user can choose between several resampler backends and
// profiles providing different compromises between CPU consumption, quality,
// and precision.
//
// # Clock source
//
// Sender should encode samples at a constant rate that is configured when the
// sender is created. There are two ways to accomplish this:
//
//   - If the user enabled internal clock (ClockSourceInternal), the sender
//     employs a CPU timer to block writes until it's time to encode the next
//     bunch of samples according to the configured sample rate. This mode is
//     useful when the user gets samples from a non-realtime source, e.g. from an
//     audio file.
//   - If the user enabled external clock (ClockSourceExternal), the samples
//     written to the sender are encoded and sent immediately, and hence the user
//     is responsible to call write operation according to the sample rate. This
//     mode is useful when the user gets samples from a realtime source with its
//     own clock, e.g. from an audio device. Internal clock should not be used in
//     this case because the audio device and the CPU might have slightly
//     different clocks, and the difference will eventually lead to an underrun
//     or an overrun.
//
// # Thread safety
//
// Can be used concurrently.
type Sender struct {
	mu   sync.RWMutex
	cPtr *C.roc_sender
}

// Open a new sender.
//
// Allocates and initializes a new sender, and attaches it to the context.
func OpenSender(context *Context, config SenderConfig) (sender *Sender, err error) {
	logWrite(LogDebug, "entering OpenSender(): context=%p config=%+v", context, config)
	defer func() {
		logWrite(LogDebug,
			"leaving OpenSender(): context=%p sender=%p err=%#v", context, sender, err,
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

	cPacketLength, err := go2cUnsignedDuration(config.PacketLength)
	if err != nil {
		return nil, fmt.Errorf("invalid config.PacketLength: %w", err)
	}

	cTargetLatency, err := go2cUnsignedDuration(config.TargetLatency)
	if err != nil {
		return nil, fmt.Errorf("invalid config.TargetLatency: %w", err)
	}

	cLatencyTolerance, err := go2cUnsignedDuration(config.LatencyTolerance)
	if err != nil {
		return nil, fmt.Errorf("invalid config.LatencyTolerance: %w", err)
	}

	cConfig := C.struct_roc_sender_config{
		frame_encoding: C.struct_roc_media_encoding{
			rate:     C.uint(config.FrameEncoding.Rate),
			format:   C.roc_format(config.FrameEncoding.Format),
			channels: C.roc_channel_layout(config.FrameEncoding.Channels),
			tracks:   C.uint(config.FrameEncoding.Tracks),
		},
		packet_encoding:          C.roc_packet_encoding(config.PacketEncoding),
		packet_length:            cPacketLength,
		packet_interleaving:      go2cBool(config.PacketInterleaving),
		fec_encoding:             C.roc_fec_encoding(config.FecEncoding),
		fec_block_source_packets: C.uint(config.FecBlockSourcePackets),
		fec_block_repair_packets: C.uint(config.FecBlockRepairPackets),
		clock_source:             C.roc_clock_source(config.ClockSource),
		latency_tuner_backend:    C.roc_latency_tuner_backend(config.LatencyTunerBackend),
		latency_tuner_profile:    C.roc_latency_tuner_profile(config.LatencyTunerProfile),
		resampler_backend:        C.roc_resampler_backend(config.ResamplerBackend),
		resampler_profile:        C.roc_resampler_profile(config.ResamplerProfile),
		target_latency:           cTargetLatency,
		latency_tolerance:        cLatencyTolerance,
	}

	var cSender *C.roc_sender
	errCode := C.roc_sender_open(context.cPtr, &cConfig, &cSender)
	if errCode != 0 {
		return nil, newNativeErr("roc_sender_open()", errCode)
	}
	if cSender == nil {
		panic("roc_sender_open() returned nil")
	}

	sender = &Sender{
		cPtr: cSender,
	}

	return sender, nil
}

// Set sender interface configuration.
//
// Updates configuration of specified interface of specified slot. If called,
// the call should be done before calling Sender.Bind() or Sender.Connect()
// for the same interface.
//
// Automatically initializes slot with given index if it's used first time.
//
// If an error happens during configure, the whole slot is disabled and marked
// broken. The slot index remains reserved. The user is responsible for removing
// the slot using Sender.Unlink(), after which slot index can be reused.
func (s *Sender) Configure(slot Slot, iface Interface, config InterfaceConfig) (err error) {
	logWrite(LogDebug,
		"entering Sender.Configure(): sender=%p slot=%+v iface=%+v config=%+v", s, slot, iface, config,
	)
	defer func() {
		logWrite(LogDebug, "leaving Sender.Configure(): sender=%p err=%#v", s, err)
	}()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.cPtr == nil {
		return errors.New("sender is closed")
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

	errCode = C.roc_sender_configure(
		s.cPtr,
		(C.roc_slot)(slot),
		(C.roc_interface)(iface),
		&cConfig)
	if errCode != 0 {
		return newNativeErr("roc_sender_configure()", errCode)
	}

	return nil
}

// Connect the sender interface to a remote receiver endpoint.
//
// Checks that the endpoint is valid and supported by the interface, allocates a
// new outgoing port, and connects it to the remote endpoint.
//
// Each slot's interface can be bound or connected only once. May be called
// multiple times for different slots or interfaces.
//
// Automatically initializes slot with given index if it's used first time.
//
// If an error happens during connect, the whole slot is disabled and marked
// broken. The slot index remains reserved. The user is responsible for removing
// the slot using Sender.Unlink(), after which slot index can be reused.
func (s *Sender) Connect(slot Slot, iface Interface, endpoint *Endpoint) (err error) {
	logWrite(LogDebug,
		"entering Sender.Connect(): sender=%p slot=%+v iface=%+v endpoint=%+v", s, slot, iface, endpoint,
	)
	defer func() {
		logWrite(LogDebug, "leaving Sender.Connect(): sender=%p err=%#v", s, err)
	}()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.cPtr == nil {
		return errors.New("sender is closed")
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

	errCode = C.roc_sender_connect(
		s.cPtr,
		(C.roc_slot)(slot),
		(C.roc_interface)(iface),
		cEndp)
	if errCode != 0 {
		return newNativeErr("roc_sender_connect()", errCode)
	}

	return nil
}

// Delete sender slot.
//
// Disconnects, unbinds, and removes all slot interfaces and removes the slot.
// All associated connections to remote peers are properly terminated.
//
// After unlinking the slot, it can be re-created again by re-using slot index.
func (s *Sender) Unlink(slot Slot) (err error) {
	logWrite(LogDebug,
		"entering Sender.Unlink(): sender=%p slot=%+v", s, slot,
	)
	defer func() {
		logWrite(LogDebug, "leaving Sender.Unlink(): sender=%p err=%#v", s, err)
	}()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.cPtr == nil {
		return errors.New("sender is closed")
	}

	var errCode C.int

	errCode = C.roc_sender_unlink(
		s.cPtr,
		(C.roc_slot)(slot))
	if errCode != 0 {
		return newNativeErr("roc_sender_unlink()", errCode)
	}

	return nil
}

// Encode samples to packets and transmit them to the receiver.
//
// Encodes samples to packets and enqueues them for transmission by the network worker
// thread of the context.
//
// If ClockSourceInternal is used, the function blocks until it's time to
// transmit the samples according to the configured sample rate. The function
// returns after encoding and enqueuing the packets, without waiting when the
// packets are actually transmitted.
//
// Until the sender is connected to at least one receiver, the stream is just
// dropped. If the sender is connected to multiple receivers, the stream is
// duplicated to each of them.
func (s *Sender) WriteFloats(frame []float32) (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.cPtr == nil {
		return errors.New("sender is closed")
	}

	if frame == nil {
		return errors.New("frame is nil")
	}

	if len(frame) == 0 {
		return nil
	}

	errCode := C.rocGoSenderWriteFloats(
		s.cPtr, (*C.float)(&frame[0]), (C.ulong)(len(frame)))
	if errCode != 0 {
		return newNativeErr("roc_sender_write()", errCode)
	}

	return nil
}

// Close the sender.
//
// Deinitializes and deallocates the sender, and detaches it from the context.
// The user should ensure that nobody uses the sender during and after this
// call. If this function fails, the sender is kept opened and attached to the
// context.
func (s *Sender) Close() (err error) {
	logWrite(LogDebug, "entering Sender.Close(): sender=%p", s)
	defer func() {
		logWrite(LogDebug, "leaving Sender.Close(): sender=%p err=%#v", s, err)
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cPtr != nil {
		errCode := C.roc_sender_close(s.cPtr)
		if errCode != 0 {
			return newNativeErr("roc_sender_close()", errCode)
		}

		s.cPtr = nil
	}

	return nil
}
