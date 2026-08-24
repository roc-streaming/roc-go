package roc

import (
	"io"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// by default, disable logging; can be overridden by specific tests
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}
