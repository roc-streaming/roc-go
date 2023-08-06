package roc

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const NumChannels = 2

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
	testSamples := make([]float32, samplesCnt)
	for i := 0; i < samplesCnt/NumChannels; i++ {
		testSamples[i*NumChannels] = float32(i+1) / 100
		testSamples[i*NumChannels+1] = -float32(i+1) / 100
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
				err := e.Sender.WriteFloats(testSamples)
				if err != nil {
					t.Fail()
				}
			case <-endChan:
				wait.Done()
				return
			}
		}
	}()

	nonZeroSamplesCount := 0
	prevL := 0
	prevR := 0
	samples := make([]float32, samplesCnt)
	for nonZeroSamplesCount < 10000 {
		select {
		case <-receiveTicker.C:
			err := e.Receiver.ReadFloats(samples)
			if err != nil {
				t.Fail()
			}
			samplesStr := samplesToString(samples)
			for i := 0; i < len(samples); i += NumChannels {
				valueL := int(math.Round(float64(samples[i] * 100)))
				valueR := int(math.Round(float64(samples[i+1] * 100)))

				if valueL == 0 { // packet loss or streaming not started yet
					assert.Equal(t, 0, valueR)
				} else {
					nonZeroSamplesCount++
					assert.Equal(t, valueL, -valueR)
					if prevL != 0 {
						require.Equal(t, valueL, prevL%(samplesCnt/NumChannels)+1,
							"prevL: %d, valueL: %d, index: %d, samples: %s", prevL, valueL, i, samplesStr)
						require.Equal(t, valueR, prevR%(samplesCnt/NumChannels)-1,
							"prevR: %d, valueR: %d, index: %d, samples: %s", prevR, valueR, i+1, samplesStr)
					}
				}
				prevL = valueL
				prevR = valueR
			}
		}
	}

	endChan <- struct{}{}
	wait.Wait()

}

func samplesToString(samples []float32) string {
	strValues := make([]string, len(samples))
	for i, sample := range samples {
		strValues[i] = fmt.Sprintf("%d=%.2f", i, sample)
	}
	return "[" + strings.Join(strValues, ", ") + "]"
}
