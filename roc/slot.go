package roc

// Network slot.
//
// A peer (sender or receiver) may have multiple slots, which may be independently
// bound or connected. You can use multiple slots on sender to connect it to multiple
// receiver addresses, and you can use multiple slots on receiver to bind it to
// multiple receiver address.
//
// Slots are numbered from zero and are created implicitly. Just specify slot index
// when binding or connecting endpoint, and slot will be automatically created if it
// was not created yet.
//
// In simple cases, just use SlotDefault.
//
// Each slot has its own set of interfaces, dedicated to different kinds of endpoints.
// See Interface for details.
type Slot int

const (
	// Alias for the slot with index zero.
	SlotDefault Slot = 0
)
