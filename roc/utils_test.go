package roc

import "testing"

func fail(expected interface{}, got interface{}, t *testing.T) {
	t.Errorf("Mismatch, expected: %v, got: %v", expected, got)
	t.FailNow()
}
