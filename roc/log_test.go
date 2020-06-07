package roc

import (
	"io/ioutil"
	"log"
	"testing"
	"time"
)

const defaultLogLevel LogLevel = LogError

type testWriter struct {
	ch chan string
}

func makeTestWriter() testWriter {
	return testWriter{
		// capacity 1 is needed to ensure that at least one
		// message can be written without blocking
		ch: make(chan string, 1),
	}
}

func (tw testWriter) Write(buf []byte) (int, error) {
	select {
	case tw.ch <- string(buf):
	default:
		// drop message instead of blocking if the channel is full
	}
	return len(buf), nil
}

func (tw testWriter) wait() string {
	const waitTimeout = time.Minute

	select {
	case s := <-tw.ch:
		return s
	case <-time.After(waitTimeout):
		return ""
	}
}

func Test_roc_log_default(t *testing.T) {
	setupFuncs := []func(){
		func() {},
		func() { SetLoggerFunc(nil) },
		func() { SetLogger(nil) },
	}

	SetLogLevel(LogDebug)
	defer SetLogLevel(defaultLogLevel)

	for n, setupFn := range setupFuncs {
		tw := makeTestWriter()

		log.SetOutput(&tw)
		defer log.SetOutput(ioutil.Discard)

		setupFn()

		ctx, _ := OpenContext(&ContextConfig{})
		ctx.Close()

		if tw.wait() == "" {
			t.Fatalf("test %v: expected logs, didn't get them before timeout", n)
		}
	}
}

func Test_roc_log_func(t *testing.T) {
	SetLogLevel(LogDebug)
	defer SetLogLevel(defaultLogLevel)

	tw := makeTestWriter()
	SetLoggerFunc(func(_ LogLevel, component string, message string) {
		_, _ = tw.Write([]byte(component + ":" + message))
	})
	defer SetLoggerFunc(nil)

	ctx, _ := OpenContext(&ContextConfig{})
	ctx.Close()

	if tw.wait() == "" {
		t.Fatal("expected logs, didn't get them before timeout")
	}
}

func Test_roc_log_interface(t *testing.T) {
	SetLogLevel(LogDebug)
	defer SetLogLevel(defaultLogLevel)

	tw := makeTestWriter()
	logger := log.New(&tw, "", log.Lshortfile)

	SetLogger(logger)
	defer SetLogger(nil)

	ctx, _ := OpenContext(&ContextConfig{})
	ctx.Close()

	if tw.wait() == "" {
		t.Fatal("expected logs, didn't get them before timeout")
	}
}
