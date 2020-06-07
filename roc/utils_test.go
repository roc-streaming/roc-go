package roc

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func fail(expected interface{}, got interface{}, t *testing.T) {
	t.Errorf("Mismatch, expected: %v, got: %v", expected, got)
	t.FailNow()
}

func TestMain(m *testing.M) {
	// by default, disable logging; can be overridden by specific tests
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}
