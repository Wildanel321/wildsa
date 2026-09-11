package main

import (
	"testing"
)

func TestSawitHealthVersion(t *testing.T) {
	if Version == "" {
		t.Error("expected Version to be non-empty")
	}
}
