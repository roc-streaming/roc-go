package roc

/*
#include <roc/sender.h>
#include <roc/config.h>

int rocGoSenderWriteFloats(roc_sender* sender, float* samples, unsigned long samples_size) {
    roc_frame frame = {(void*)samples, samples_size*sizeof(float)};
    return roc_sender_write(sender, &frame);
 }
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
)

// Sender peer.
//
// Sender gets an audio stream from the user, encodes it into network packets, and
// transmits them to a remote receiver.
//
// # Context
//
// Sender is automatically attached to a context when opened and detached from it when
// closed. The user should not close the context until the sender is closed.
//
// Sender work consists of two parts: stream encoding and packet transmission. The
// encoding part is performed in the sender itself, and the transmission part is
// performed in the context network worker threads.
//
// # Life cycle
//
// - A sender is created using OpenSender().
//
//   - Optionally, the sender parameters may be fine-tuned using Sender.Set().
//     functions.
//
//   - The sender either binds local endpoints using Sender.Bind(), allowing receivers
//     connecting to them, or itself connects to remote receiver endpoints using
//     Sender.Connect(). What approach to use is up to the user.
//
//   - The audio stream is iteratively written to the sender using Sender.Write*(). The
//     sender encodes the stream into packets and send to connected receiver(s).
//
// - The sender is destroyed using Sender.Close().
//
// The user is responsible for closing any opened sender before exiting the program.
//
// # Slots, interfaces, and endpoints
//
// Sender has one or multiple slots, which may be independently bound or connected.
// Slots may be used to connect sender to multiple receivers. Slots are numbered from
// zero and are created automatically. In simple cases just use SlotDefault.
//
// Each slot has its own set of interfaces, one per each type defined in Interface
// type. The interface defines the type of the communication with the remote peer
// and the set of the protocols supported by it.
//
// Supported actions with the interface:
//
//   - Call Sender.Bind() to bind the interface to a local Endpoint. In this
//     case the sender accepts connections from receivers and sends media stream to all
//     connected receivers.
//
//   - Call Sender.Connect() to connect the interface to a remote Endpoint.
//     In this case the sender initiates connection to the receiver and starts sending
//     media stream to it.
//
// Supported interface configurations:
//
//   - Connect InterfaceConsolidated to a remote endpoint (e.g. be an RTSP
//     client).
//
//   - Bind InterfaceConsolidated to a local endpoint (e.g. be an RTSP server).
//
//   - Connect InterfaceAudioSource, InterfaceAudioRepair (optionally,
//     for FEC), and InterfaceAudioControl (optionally, for control messages)
//     to remote endpoints (e.g. be an RTP/FECFRAME/RTCP sender).
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
// # Sample rate
//
// If the sample rate of the user frames and the sample rate of the network packets are
// different, the sender employs resampler to convert one rate to another.
//
// Resampling is a quite time-consuming operation. The user can choose between completely
// disabling resampling (and so use the same rate for frames and packets) or several
// resampler profiles providing different compromises between CPU consumption and quality.
//
// # Clock source
//
// Sender should encode samples at a constant rate that is configured when the sender
// is created. There are two ways to accomplish this:
//
//   - If the user enabled internal clock (ClockInternal), the sender employs a
//     CPU timer to block writes until it's time to encode the next bunch of samples
//     according to the configured sample rate.
//
//     This mode is useful when the user gets samples from a non-realtime source, e.g.
//     from an audio file.
//
//   - If the user enabled external clock (ClockExternal), the samples written to
//     the sender are encoded and sent immediately, and hence the user is responsible to
//     call write operation according to the sample rate.
//
//     This mode is useful when the user gets samples from a realtime source with its own
//     clock, e.g. from an audio device. Internal clock should not be used in this case
//     because the audio device and the CPU might have slightly different clocks, and the
//     difference will eventually lead to an underrun or an overrun.
//
// # Thread safety
//
// Can be used concurrently.
type Sender struct {
	mu   sync.RWMutex
	cPtr *C.roc_sender
}

// Open a new sender.
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

	cConfig := C.struct_roc_sender_config{
		frame_sample_rate:        (C.uint)(config.FrameSampleRate),
		frame_channels:           (C.roc_channel_set)(config.FrameChannels),
		frame_encoding:           (C.roc_frame_encoding)(config.FrameEncoding),
		packet_sample_rate:       (C.uint)(config.PacketSampleRate),
		packet_channels:          (C.roc_channel_set)(config.PacketChannels),
		packet_encoding:          (C.roc_packet_encoding)(config.PacketEncoding),
		packet_length:            (C.ulonglong)(config.PacketLength),
		packet_interleaving:      go2cBool(config.PacketInterleaving),
		clock_source:             (C.roc_clock_source)(config.ClockSource),
		resampler_backend:        (C.roc_resampler_backend)(config.ResamplerBackend),
		resampler_profile:        (C.roc_resampler_profile)(config.ResamplerProfile),
		fec_encoding:             (C.roc_fec_encoding)(config.FecEncoding),
		fec_block_source_packets: (C.uint)(config.FecBlockSourcePackets),
		fec_block_repair_packets: (C.uint)(config.FecBlockRepairPackets),
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

// Set sender interface outgoing address.
//
// Optional. Should be used only when connecting an interface to a remote endpoint.
//
// If set, explicitly defines the IP address of the OS network interface from which to
// send the outgoing packets. If not set, the outgoing interface is selected automatically
// by the OS, depending on the remote endpoint address.
//
// It is allowed to set outgoing address to `0.0.0.0` (for IPv4) or to `::` (for IPv6),
// to achieve the same behavior as if it wasn't set, i.e. to let the OS to select the
// outgoing interface automatically.
//
// By default, the outgoing address is not set.
//
// Each slot's interface can have only one outgoing address. The function should be called
// before calling Sender.Connect() for this slot and interface. It should not be
// called when calling Sender.Bind() for the interface.
//
// Automatically initializes slot with given index if it's used first time.
func (s *Sender) SetOutgoingAddress(slot Slot, iface Interface, ip string) (err error) {
	logWrite(LogDebug,
		"entering Sender.SetOutgoingAddress(): sender=%p slot=%v iface=%v ip=%v", s, slot, iface, ip,
	)
	defer func() {
		logWrite(LogDebug, "leaving Sender.SetOutgoingAddress(): sender=%p err=%#v", s, err)
	}()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.cPtr == nil {
		return errors.New("sender is closed")
	}

	cIP, err := go2cStr(ip)
	if err != nil {
		return fmt.Errorf("invalid ip: %w", err)
	}
	errCode := C.roc_sender_set_outgoing_address(
		s.cPtr,
		(C.roc_slot)(slot),
		(C.roc_interface)(iface),
		(*C.char)(&cIP[0]))
	if errCode != 0 {
		return newNativeErr("roc_sender_set_outgoing_address()", errCode)
	}

	return nil
}

// Set sender interface address reuse option.
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
func (s *Sender) SetReuseaddr(slot Slot, iface Interface, enabled bool) (err error) {
	logWrite(LogDebug,
		"entering Sender.SetReuseaddr(): sender=%p slot=%v iface=%v enabled=%v", s, slot, iface, enabled,
	)
	defer func() {
		logWrite(LogDebug, "leaving Sender.SetReuseaddr(): sender=%p err=%#v", s, err)
	}()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.cPtr == nil {
		return errors.New("sender is closed")
	}

	cEnabled := go2cBool(enabled)

	errCode := C.roc_sender_set_reuseaddr(
		s.cPtr,
		(C.roc_slot)(slot),
		(C.roc_interface)(iface),
		(C.int)(cEnabled),
	)

	if errCode != 0 {
		return newNativeErr("roc_sender_set_reuseaddr()", errCode)
	}

	return nil
}

// Connect the sender interface to a remote receiver endpoint.
//
// Checks that the endpoint is valid and supported by the interface, allocates
// a new outgoing port, and connects it to the remote endpoint.
//
// Each slot's interface can be bound or connected only once.
// May be called multiple times for different slots or interfaces.
//
// Automatically initializes slot with given index if it's used first time.
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

// Encode samples to packets and transmit them to the receiver.
//
// Encodes samples to packets and enqueues them for transmission by the network worker
// thread of the context.
//
// If ClockInternal is used, the function blocks until it's time to transmit the
// samples according to the configured sample rate. The function returns after encoding
// and enqueuing the packets, without waiting when the packets are actually transmitted.
//
// Until the sender is connected to at least one receiver, the stream is just dropped.
// If the sender is connected to multiple receivers, the stream is duplicated to
// each of them.
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
// Deinitializes and deallocates the sender, and detaches it from the context. The user
// should ensure that nobody uses the sender during and after this call. If this
// function fails, the sender is kept opened and attached to the context.
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
