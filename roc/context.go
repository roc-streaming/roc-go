package roc

/*
#include <roc/context.h>
*/
import "C"

import (
	"fmt"
	"sync"
)

// Shared context.
//
// Context contains memory pools and network worker threads, shared among objects attached
// to the context. It is allowed both to create a separate context for every object, or
// to create a single context shared between multiple objects.
//
// # Life cycle
//
// A context is created using OpenContext() and destroyed using Context.Close().
//
// Objects can be attached and detached to an opened context at any moment from any
// thread. However, the user should ensure that the context is not closed until there
// are no objects attached to the context.
//
// The user is responsible for closing any opened context before exiting the program.
//
// # Thread safety
//
// Can be used concurrently.
//
// # See also
//
// See also Sender, Receiver.
type Context struct {
	mu   sync.RWMutex
	cPtr *C.roc_context
}

// Open a new context.
// Allocates and initializes a new context. May start some background threads.
// User is responsible to call Context.Close to free context resources.
func OpenContext(config ContextConfig) (*Context, error) {
	logWrite(LogDebug, fmt.Sprintf("%v", logFormat{op: "roc_context_open_begin()",
		params: map[string]interface{}{
			"config":  config,
			"version": Version(),
		}}))

	cConfig := C.struct_roc_context_config{
		max_packet_size: C.uint(config.MaxPacketSize),
		max_frame_size:  C.uint(config.MaxFrameSize),
	}

	var cCtx *C.roc_context
	errCode := C.roc_context_open(&cConfig, &cCtx)
	if errCode != 0 {
		return nil, newNativeErr("roc_context_open()", errCode)
	}
	if cCtx == nil {
		panic("roc_context_open() returned nil")
	}

	ctx := &Context{
		cPtr: cCtx,
	}

	logWrite(LogDebug, fmt.Sprintf("%v", logFormat{op: "roc_context_open_end()",
		params: map[string]interface{}{
			"config": config,
		}}))

	return ctx, nil
}

// Close the context.
// Stops any started background threads, deinitializes and deallocates the context.
// The user should ensure that nobody uses the context during and after this call.
// If this function fails, the context is kept opened.
func (c *Context) Close() error {
	logWrite(LogDebug, fmt.Sprintf("%v", logFormat{op: "roc_context_close_begin()"}))

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cPtr != nil {
		errCode := C.roc_context_close(c.cPtr)
		if errCode != 0 {
			return newNativeErr("roc_context_close()", errCode)
		}

		c.cPtr = nil
	}

	logWrite(LogDebug, fmt.Sprintf("%v", logFormat{op: "roc_context_close_end()"}))
	return nil
}
