package roc

/*
#cgo LDFLAGS: -lroc
#include <roc/receiver.h>
#include <roc/address.h>
#include <roc/sender.h>
#include <roc/log.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

func OpenReceiver(rocContext *Context, receiverConfig *ReceiverConfig) (*Receiver, error) {
	rocContextCPtr := (*C.roc_context)(unsafe.Pointer(rocContext))
	receiverConfigCPtr := (*C.roc_receiver_config)(unsafe.Pointer(receiverConfig))

	receiver := C.roc_receiver_open(rocContextCPtr, receiverConfigCPtr)
	if receiver == nil {
		return nil, ErrInvalidArguments
	}
	return (*Receiver)(receiver), nil
}

func (r *Receiver) Close() error {
	errCode := C.roc_receiver_close((*C.roc_receiver)(unsafe.Pointer(r)))
	if errCode == 0 {
		return nil
	}
	if errCode < 0 {
		return ErrInvalidArguments
	}
	return ErrInvalidApi
}
