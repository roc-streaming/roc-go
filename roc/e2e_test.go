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
	ReceiverConfig ReceiverConfig
	SenderConfig   SenderConfig
	Context        *Context
	Receiver       *Receiver
	Sender         *Sender
}

type e2eParams struct {
	clockSource ClockSource
	fecEncoding FecEncoding
	sourceURI   string
	repairURI   string
}

func newE2E(t *testing.T, params e2eParams) *e2e {
	var (
		err error
		e   e2e
	)

	// create context
	e.Context, err = OpenContext(makeContextConfig())
	require.NoError(t, err)
	require.NotNil(t, e.Context)

	// create receiver
	e.ReceiverConfig = makeReceiverConfig()
	e.ReceiverConfig.ClockSource = params.clockSource
	e.Receiver, err = OpenReceiver(e.Context, e.ReceiverConfig)
	require.NoError(t, err)
	require.NotNil(t, e.Receiver)

	// create sender
	e.SenderConfig = makeSenderConfig()
	e.SenderConfig.ClockSource = params.clockSource
	e.SenderConfig.FecEncoding = params.fecEncoding
	e.Sender, err = OpenSender(e.Context, e.SenderConfig)
	require.NoError(t, err)
	require.NotNil(t, e.Sender)

	// create source endpoint
	sourceEndpoint, err := ParseEndpoint(params.sourceURI)
	require.NoError(t, err)
	require.NotNil(t, sourceEndpoint)

	// bind receiver to source endpoint
	err = e.Receiver.Bind(SlotDefault, InterfaceAudioSource, sourceEndpoint)
	require.NoError(t, err)
	require.NotEmpty(t, sourceEndpoint.Port)

	var repairEndpoint *Endpoint
	if params.repairURI != "" {
		// create repair endpoint
		repairEndpoint, err = ParseEndpoint(params.repairURI)
		require.NoError(t, err)
		require.NotNil(t, repairEndpoint)

		// bind receiver to repair endpoint
		err = e.Receiver.Bind(SlotDefault, InterfaceAudioRepair, repairEndpoint)
		require.NoError(t, err)
		require.NotEmpty(t, repairEndpoint.Port)
	}

	// connect sender to receiver source endpoint
	err = e.Sender.Connect(SlotDefault, InterfaceAudioSource, sourceEndpoint)
	require.NoError(t, err)

	if params.repairURI != "" {
		// connect sender to receiver repair endpoint
		err = e.Sender.Connect(SlotDefault, InterfaceAudioRepair, repairEndpoint)
		require.NoError(t, err)
	}

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
	tests := []struct {
		name   string
		params e2eParams
	}{
		{
			name: "default",
			params: e2eParams{
				clockSource: ClockSourceExternal,
				fecEncoding: FecEncodingDisable,
				sourceURI:   "rtp://127.0.0.1:0"},
		},
		{
			name: "fec",
			params: e2eParams{
				clockSource: ClockSourceExternal,
				fecEncoding: FecEncodingRs8m,
				sourceURI:   "rtp+rs8m://127.0.0.1:0",
				repairURI:   "rs8m://127.0.0.1:0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newE2E(t, tt.params)
			defer e.close(t)

			samplesCnt := 100
			testSamples := generateTestSamples(samplesCnt)

			senderInterval := time.Second /
				time.Duration(int(e.SenderConfig.FrameEncoding.Rate)/samplesCnt)
			sendTicker := time.NewTicker(senderInterval)
			defer sendTicker.Stop()

			receiverInterval := time.Second /
				time.Duration(int(e.ReceiverConfig.FrameEncoding.Rate)/samplesCnt)
			receiveTicker := time.NewTicker(receiverInterval)
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

			validationState := validationState{}
			samples := make([]float32, samplesCnt)
			for validationState.nonZeroSamplesCount < 10000 {
				select {
				case <-receiveTicker.C:
					err := e.Receiver.ReadFloats(samples)
					if err != nil {
						t.Fail()
					}
					validateSamples(t, &validationState, samples, samplesCnt)
				}
			}

			endChan <- struct{}{}
			wait.Wait()
		})
	}

}

func TestEnd2End_Blocking(t *testing.T) {
	tests := []struct {
		name   string
		params e2eParams
	}{
		{
			name: "default",
			params: e2eParams{
				clockSource: ClockSourceInternal,
				fecEncoding: FecEncodingDisable,
				sourceURI:   "rtp://127.0.0.1:0"},
		},
		{
			name: "fec",
			params: e2eParams{
				clockSource: ClockSourceInternal,
				fecEncoding: FecEncodingRs8m,
				sourceURI:   "rtp+rs8m://127.0.0.1:0",
				repairURI:   "rs8m://127.0.0.1:0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newE2E(t, tt.params)
			defer e.close(t)

			samplesCnt := 100
			testSamples := generateTestSamples(samplesCnt)

			endChan := make(chan struct{})
			var wait sync.WaitGroup
			wait.Add(1)
			go func() {
				for {
					select {
					case <-endChan:
						wait.Done()
						return
					default:
						err := e.Sender.WriteFloats(testSamples)
						if err != nil {
							t.Fail()
						}
					}
				}
			}()

			validationState := validationState{}
			samples := make([]float32, samplesCnt)
			for validationState.nonZeroSamplesCount < 10000 {
				err := e.Receiver.ReadFloats(samples)
				if err != nil {
					t.Fail()
				}
				validateSamples(t, &validationState, samples, samplesCnt)
			}

			endChan <- struct{}{}
			wait.Wait()
		})
	}

}

func generateTestSamples(samplesCnt int) []float32 {
	testSamples := make([]float32, samplesCnt)
	for i := 0; i < samplesCnt/NumChannels; i++ {
		testSamples[i*NumChannels] = float32(i+1) / 100
		testSamples[i*NumChannels+1] = -float32(i+1) / 100
	}
	return testSamples
}

type validationState struct {
	nonZeroSamplesCount int
	prevL               int
	prevR               int
}

func validateSamples(t *testing.T, state *validationState, samples []float32, samplesCnt int) {
	samplesStr := samplesToString(samples)
	for i := 0; i < len(samples); i += NumChannels {
		valueL := int(math.Round(float64(samples[i] * 100)))
		valueR := int(math.Round(float64(samples[i+1] * 100)))

		if valueL == 0 { // packet loss or streaming not started yet
			assert.Equal(t, 0, valueR)
		} else {
			state.nonZeroSamplesCount++
			assert.Equal(t, valueL, -valueR)
			if state.prevL != 0 {
				require.Equal(t, valueL, state.prevL%(samplesCnt/NumChannels)+1,
					"prevL: %d, valueL: %d, index: %d, samples: %s", state.prevL, valueL, i, samplesStr)
				require.Equal(t, valueR, state.prevR%(samplesCnt/NumChannels)-1,
					"prevR: %d, valueR: %d, index: %d, samples: %s", state.prevR, valueR, i+1, samplesStr)
			}
		}
		state.prevL = valueL
		state.prevR = valueR
	}
}

func samplesToString(samples []float32) string {
	strValues := make([]string, len(samples))
	for i, sample := range samples {
		strValues[i] = fmt.Sprintf("%d=%.2f", i, sample)
	}
	return "[" + strings.Join(strValues, ", ") + "]"
}
