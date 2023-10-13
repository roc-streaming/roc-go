package roc

// Clock source for sender or receiver.
//
//go:generate stringer -type ClockSource -trimprefix Clock
type ClockSource int

const (
	// Sender or receiver is clocked by external user-defined clock.
	// Write and read operations are non-blocking. The user is responsible
	// to call them in time, according to the external clock.
	ClockExternal ClockSource = 0

	// Sender or receiver is clocked by an internal clock.
	// Write and read operations are blocking. They automatically wait until it's time
	// to process the next bunch of samples according to the configured sample rate.
	ClockInternal ClockSource = 1
)
