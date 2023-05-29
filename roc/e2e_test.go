package roc

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type e2e struct {
	Context  *Context
	Receiver *Receiver
	Sender   *Sender
}

func newE2E(t *testing.T) *e2e {
	var (
		err error
		e   e2e
	)

	// create context
	e.Context, err = OpenContext(makeContextConfig())
	require.NoError(t, err)
	require.NotNil(t, e.Context)

	// create receiver
	e.Receiver, err = OpenReceiver(e.Context, makeReceiverConfig())
	require.NoError(t, err)
	require.NotNil(t, e.Receiver)

	// create sender
	e.Sender, err = OpenSender(e.Context, makeSenderConfig())
	require.NoError(t, err)
	require.NotNil(t, e.Sender)

	// create source endpoint
	sourceEndpoint, err := ParseEndpoint("rtp+rs8m://127.0.0.1:0")
	require.NoError(t, err)
	require.NotNil(t, sourceEndpoint)

	// bind receiver to source endpoint
	err = e.Receiver.Bind(SlotDefault, InterfaceAudioSource, sourceEndpoint)
	require.NoError(t, err)
	require.NotEmpty(t, sourceEndpoint.Port)

	// create repair endpoint
	repairEndpoint, err := ParseEndpoint("rs8m://127.0.0.1:0")
	require.NoError(t, err)
	require.NotNil(t, repairEndpoint)

	// bind receiver to repair endpoint
	err = e.Receiver.Bind(SlotDefault, InterfaceAudioRepair, repairEndpoint)
	require.NoError(t, err)
	require.NotEmpty(t, repairEndpoint.Port)

	// connect sender to receiver source endpoint
	err = e.Sender.Connect(SlotDefault, InterfaceAudioSource, sourceEndpoint)
	require.NoError(t, err)

	// connect sender to receiver repair endpoint
	err = e.Sender.Connect(SlotDefault, InterfaceAudioRepair, repairEndpoint)
	require.NoError(t, err)

	return &e
}

func (e *e2e) close(t *testing.T) {
	var err error

	err = e.Receiver.Close()
	require.NoError(t, err)

	err = e.Sender.Close()
	require.NoError(t, err)

	err = e.Context.Close()
	require.NoError(t, err)
}

func TestEnd2End_Default(t *testing.T) {
	e := newE2E(t)
	defer e.close(t)

	samplesCnt := 2
	samples := make([]float32, samplesCnt)
	for i := 0; i < samplesCnt; i++ {
		samples[i] = float32(i + 1)
	}

	endChan := make(chan struct{})
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		for {
			err := e.Sender.WriteFloats(samples)
			if err != nil {
				t.Fail()
			}
			select {
			case <-endChan:
				wait.Done()
				return
			default:
			}
		}
	}()

	for {
		recFloats := make([]float32, samplesCnt)
		err := e.Receiver.ReadFloats(recFloats)
		if err != nil {
			t.Fail()
		}

		nonZeroSamplesCnt := 0
		for i := 0; i < samplesCnt; i++ {
			if recFloats[i] != 0 {
				nonZeroSamplesCnt++
			}
		}
		if nonZeroSamplesCnt == samplesCnt {
			break
		}
	}

	endChan <- struct{}{}
	wait.Wait()
}
