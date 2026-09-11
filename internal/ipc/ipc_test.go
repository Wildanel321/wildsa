package ipc_test

import (
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/sawitos/sawit/internal/ipc"
)

func TestIPCCommunication(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on random TCP port: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()

	// Mock Server Goroutine
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		var req ipc.Request
		decoder := json.NewDecoder(conn)
		if err := decoder.Decode(&req); err != nil {
			return
		}

		resp := ipc.Response{
			ID:        req.ID,
			Success:   true,
			Data:      "test-ok",
			Timestamp: time.Now(),
		}
		json.NewEncoder(conn).Encode(resp)
	}()

	req := &ipc.Request{
		ID:        "req-1",
		Action:    ipc.ActionSystemInfo,
		Timestamp: time.Now(),
	}

	resp, err := ipc.SendRequest("tcp", addr, req, 2*time.Second)
	if err != nil {
		t.Fatalf("expected successful IPC request, got error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected resp.Success true, got false (error: %s)", resp.Error)
	}
	if resp.Data != "test-ok" {
		t.Errorf("expected resp.Data 'test-ok', got %v", resp.Data)
	}
}
