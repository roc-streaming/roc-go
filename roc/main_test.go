package roc

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// by default, disable logging; can be overridden by specific tests
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}
