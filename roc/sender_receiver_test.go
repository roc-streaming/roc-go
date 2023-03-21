package roc

import (
	"sync"
	"testing"
)

type testEnv struct {
	Context  *Context
	Sender   *Sender
	Receiver *Receiver
}

func makeReceiverConfig() ReceiverConfig {
	return ReceiverConfig{
		FrameSampleRate:  44100,
		FrameChannels:    ChannelSetStereo,
		FrameEncoding:    FrameEncodingPcmFloat,
		ClockSource:      ClockInternal,
		ResamplerProfile: ResamplerProfileDisable,
	}
}

func makeSenderConfig() SenderConfig {
	return SenderConfig{
		FrameSampleRate:  44100,
		FrameChannels:    ChannelSetStereo,
		FrameEncoding:    FrameEncodingPcmFloat,
		ClockSource:      ClockInternal,
		ResamplerProfile: ResamplerProfileDisable,
		FecEncoding:      FecEncodingRs8m,
	}
}

func newTestEnv(t *testing.T) *testEnv {
	var (
		err error
		e   testEnv
	)
	// create context
	e.Context, err = OpenContext(ContextConfig{})
	if e.Context == nil || err != nil {
		t.Fatalf("Cannot create context: %v, error: %v", e.Context, err)
	}

	// create receiver
	e.Receiver, err = OpenReceiver(e.Context, makeReceiverConfig())
	if e.Receiver == nil || err != nil {
		t.Fatalf("Cannot create receiver: %v, error: %v", e.Receiver, err)
	}

	// bind receiver to two endpoints
	// 1)
	sourceEndpoint, err := ParseEndpoint("rtp+rs8m://127.0.0.1:0")
	if sourceEndpoint == nil || err != nil {
		t.Fatalf("Cannot create source address object, address object: %v, error: %v",
			sourceEndpoint, err)
	}
	err = e.Receiver.Bind(SlotDefault, InterfaceAudioSource, sourceEndpoint)
	if err != nil {
		t.Fatalf("Cannot bind receiver: %v", err)
	}

	// 2)
	repairEndpoint, err := ParseEndpoint("rs8m://127.0.0.1:0")
	if repairEndpoint == nil || err != nil {
		t.Fatalf("Cannot create repair address object, address object: %v, error: %v",
			repairEndpoint, err)
	}
	err = e.Receiver.Bind(SlotDefault, InterfaceAudioRepair, repairEndpoint)
	if err != nil {
		t.Fatalf("Cannot bind receiver: %v", err)
	}

	// create sender
	e.Sender, err = OpenSender(e.Context, makeSenderConfig())
	if e.Sender == nil || err != nil {
		t.Fatalf("Cannot create sender, sender: %v, error: %v", e.Sender, err)
	}

	// connect sender to receiver endpoints
	// 1)
	err = e.Sender.Connect(SlotDefault, InterfaceAudioSource, sourceEndpoint)
	if err != nil {
		t.Fatalf("Cannot connect sender to receiver: %v", err)
	}
	// 2)
	err = e.Sender.Connect(SlotDefault, InterfaceAudioRepair, repairEndpoint)
	if err != nil {
		t.Fatalf("Cannot connect sender to receiver: %v", err)
	}
	return &e
}

func TestSenderSetReuseaddr(t *testing.T) {
	e := newTestEnv(t)
	defer e.close(t)

	err := e.Sender.SetReuseaddr(SlotDefault, InterfaceAudioRepair, 1)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}

func (e *testEnv) close(t *testing.T) {
	err := e.Receiver.Close()
	if err != nil {
		t.Fail()
	}
	err = e.Sender.Close()
	if err != nil {
		t.Fail()
	}
	err = e.Context.Close() // remove after finalizers are done
	if err != nil {
		t.Fail()
	}
}

func TestSenderReceiver(t *testing.T) {
	e := newTestEnv(t)
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
