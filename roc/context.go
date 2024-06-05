package roc

/*
#include <roc/context.h>
*/
import "C"

import (
	"errors"
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
func OpenContext(config ContextConfig) (ctx *Context, err error) {
	logWrite(LogDebug, "entering OpenContext(): config=%+v", config)
	defer func() {
		logWrite(LogDebug, "leaving OpenContext(): context=%p err=%#v", ctx, err)
	}()

	checkVersionFn()

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

	ctx = &Context{
		cPtr: cCtx,
	}

	return ctx, nil
}

// Register custom encoding.
//
// Registers encoding with given encodingID. Registered encodings complement
// built-in encodings defined by \ref roc_packet_encoding enum. Whenever you need to
// specify packet encoding, you can use both built-in and registered encodings.
//
// On sender, you should register custom encoding and set to PacketEncoding field
// of SenderConfig, if you need to force specific encoding of packets, but
// built-in set of encodings is not enough.
//
// On receiver, you should register custom encoding with same id and specification,
// if you did so on sender, and you're not using any signaling protocol (like RTSP)
// that is capable of automatic exchange of encoding information.
//
// In case of RTP, encoding id is mapped directly to payload type field (PT).
func (c *Context) RegisterEncoding(encodingID int, encoding MediaEncoding) (err error) {
	logWrite(LogDebug,
		"entering Context.RegisterEncoding(): context=%p id=%+v encoding=%+v",
		c, encodingID, encoding,
	)
	defer func() {
		logWrite(LogDebug, "leaving Context.RegisterEncoding(): context=%p err=%#v", c, err)
	}()

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.cPtr == nil {
		return errors.New("context is closed")
	}

	cEncoding := C.struct_roc_media_encoding{
		rate:     C.uint(encoding.Rate),
		format:   C.roc_format(encoding.Format),
		channels: C.roc_channel_layout(encoding.Channels),
		tracks:   C.uint(encoding.Tracks),
	}

	var errCode C.int

	errCode = C.roc_context_register_encoding(
		c.cPtr,
		C.int(encodingID),
		&cEncoding)
	if errCode != 0 {
		return newNativeErr("roc_context_register_encoding()", errCode)
	}

	return nil
}

// Close the context.
// Stops any started background threads, deinitializes and deallocates the context.
// The user should ensure that nobody uses the context during and after this call.
// If this function fails, the context is kept opened.
func (c *Context) Close() (err error) {
	logWrite(LogDebug, "entering Context.Close(): context=%p", c)
	defer func() {
		logWrite(LogDebug, "leaving Context.Close(): context=%p err=%#v", c, err)
	}()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cPtr != nil {
		errCode := C.roc_context_close(c.cPtr)
		if errCode != 0 {
			return newNativeErr("roc_context_close()", errCode)
		}

		c.cPtr = nil
	}

	return nil
}
