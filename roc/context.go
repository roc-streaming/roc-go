package roc

/*
#include <roc/receiver.h>
*/
import "C"

func OpenContext(contextConfig *ContextConfig) (*Context, error) {
	var cCtxConfig C.struct_roc_context_config
	cCtxConfig.max_packet_size = C.uint(contextConfig.MaxPacketSize)
	cCtxConfig.max_frame_size = C.uint(contextConfig.MaxFrameSize)

	c := C.roc_context_open(&cCtxConfig)
	if c == nil {
		return nil, ErrInvalidArguments
	}
	return (*Context)(c), nil
}

func (c *Context) Close() error {
	errCode := C.roc_context_close((*C.roc_context)(c))

	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return ErrInvalidArguments
	}
	return ErrInvalidApi
}
