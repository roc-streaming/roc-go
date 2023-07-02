package roc

import (
	"math"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	samplesCnt := 100
	samples := make([]float32, samplesCnt)
	for i := 0; i < samplesCnt/2; i++ {
		samples[i*2] = float32(i+1) / 100
		samples[i*2+1] = -float32(i+1) / 100
	}

	var interval = time.Second / time.Duration(44100/samplesCnt)
	sendTicker := time.NewTicker(interval)
	defer sendTicker.Stop()

	receiveTicker := time.NewTicker(interval)
	defer receiveTicker.Stop()

	endChan := make(chan struct{})
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		for {
			select {
			case <-sendTicker.C:
				err := e.Sender.WriteFloats(samples)
				if err != nil {
					t.Fail()
				}
			case <-endChan:
				wait.Done()
				return
			}
		}
	}()

	streamingStarted := false
	resultSamples := make([]float32, 0, 10000)
	recFloats := make([]float32, samplesCnt)
	for len(resultSamples) < cap(resultSamples) {
		select {
		case <-receiveTicker.C:
			err := e.Receiver.ReadFloats(recFloats)
			if err != nil {
				t.Fail()
			}
			for _, v := range recFloats {
				if !streamingStarted && (v != 0) {
					streamingStarted = true
				}
				if streamingStarted {
					resultSamples = append(resultSamples, v)
				}
			}
		}
	}

	endChan <- struct{}{}
	wait.Wait()

	prevL := resultSamples[0]
	prevR := resultSamples[1]

	for i := 1; i < len(resultSamples)/2; i++ {
		valueL := resultSamples[i*2]
		valueR := resultSamples[i*2+1]
		if valueL == 0 { // packet loss
			assert.Equal(t, 0, valueR)
		} else {
			assert.Equal(t, valueL, -valueR)
			if math.Abs(float64(prevL-0.5)) >= 0.0001 {
				assert.InDelta(t, valueL, prevL+0.01, 0.0001)
				assert.InDelta(t, valueR, prevR-0.01, 0.0001)
			}
		}
		prevL = valueL
		prevR = valueR
	}
}
