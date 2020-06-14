package roc

import (
	"sync"
	"testing"
)

type testEnv struct {
	Sender   *Sender
	Receiver *Receiver
	Context  *Context
}

func makeReceiverConfig() *ReceiverConfig {
	return &ReceiverConfig{
		FrameSampleRate:  44100,
		FrameChannels:    ChannelSetStereo,
		FrameEncoding:    FrameEncodingPcmFloat,
		AutomaticTiming:  1,
		ResamplerProfile: ResamplerDisable,
	}
}

func makeSenderConfig() *SenderConfig {
	return &SenderConfig{
		FrameSampleRate:  44100,
		FrameChannels:    ChannelSetStereo,
		FrameEncoding:    FrameEncodingPcmFloat,
		AutomaticTiming:  1,
		ResamplerProfile: ResamplerDisable,
		FecCode:          FecRs8m,
	}
}

func newTestEnv(t *testing.T) *testEnv {
	var (
		err error
		e   testEnv
	)
	// create context
	e.Context, err = OpenContext(&ContextConfig{})
	if e.Context == nil || err != nil {
		t.Errorf("Cannot create context: %v, error: %v", e.Context, err)
	}

	// create receiver
	e.Receiver, err = OpenReceiver(e.Context, makeReceiverConfig())
	if e.Receiver == nil || err != nil {
		t.Errorf("Cannot create receiver: %v, error: %v", e.Receiver, err)
	}

	// bind receiver to two ports
	// 1)
	sourceAddr, err := NewAddress(AfAuto, "127.0.0.1", 0)
	if sourceAddr == nil || err != nil {
		t.Errorf("Cannot create source address object, address object: %v, error: %v",
			sourceAddr, err)
	}
	err = e.Receiver.Bind(PortAudioSource, ProtoRtpRs8mSource, sourceAddr)
	if err != nil {
		t.Errorf("Cannot bind receiver: %v", err)
	}

	// 2)
	repairAddr, err := NewAddress(AfAuto, "127.0.0.1", 0)
	if repairAddr == nil || err != nil {
		t.Errorf("Cannot create repair address object, address object: %v, error: %v",
			repairAddr, err)
	}
	err = e.Receiver.Bind(PortAudioRepair, ProtoRs8mRepair, repairAddr)
	if err != nil {
		t.Errorf("Cannot bind receiver: %v", err)
	}

	// create sender
	e.Sender, err = OpenSender(e.Context, makeSenderConfig())
	if e.Sender == nil || err != nil {
		t.Errorf("Cannot create sender, sender: %v, error: %v", e.Sender, err)
	}

	// bind sender to a port
	senderAddr, err := NewAddress(AfAuto, "127.0.0.1", 0)
	if senderAddr == nil || err != nil {
		t.Errorf("Cannot create sender address object, address object: %v, error: %v",
			senderAddr, err)
	}
	err = e.Sender.Bind(senderAddr)
	if err != nil {
		t.Errorf("Cannot bind sender: %v", err)
	}

	// connect sender to receiver ports
	// 1)
	err = e.Sender.Connect(PortAudioSource, ProtoRtpRs8mSource, sourceAddr)
	if err != nil {
		t.Errorf("Cannot connect sender to receiver: %v", err)
	}
	// 2)
	err = e.Sender.Connect(PortAudioRepair, ProtoRs8mRepair, repairAddr)
	if err != nil {
		t.Errorf("Cannot connect sender to receiver: %v", err)
	}
	return &e
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

func Test_roc_sender_write_receiver_read(t *testing.T) {
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
