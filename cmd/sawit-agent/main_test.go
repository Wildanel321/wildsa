package main

import (
	"testing"

	"github.com/sawitos/sawit/internal/ipc"
)

func TestProcessRequestSystemInfo(t *testing.T) {
	req := &ipc.Request{
		ID:     "req-1",
		Action: ipc.ActionSystemInfo,
	}

	resp := processRequest(req)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	if !resp.Success {
		t.Errorf("expected request success, got error: %s", resp.Error)
	}
}

func TestProcessRequestUnknownAction(t *testing.T) {
	req := &ipc.Request{
		ID:     "req-2",
		Action: "unknown.action",
	}

	resp := processRequest(req)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	if resp.Success {
		t.Error("expected request failure for unknown action")
	}
}
