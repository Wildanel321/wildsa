package main

import (
	"testing"
)

func TestSawitUpdateVersion(t *testing.T) {
	if Version == "" {
		t.Error("expected Version to be non-empty")
	}
}
