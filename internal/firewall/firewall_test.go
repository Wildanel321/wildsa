package firewall_test

import (
	"testing"

	"github.com/sawitos/sawit/internal/firewall"
)

func TestGetStatus(t *testing.T) {
	status, err := firewall.GetStatus()
	if err != nil {
		t.Fatalf("unexpected error getting firewall status: %v", err)
	}
	if status.Backend != "nftables" {
		t.Errorf("expected backend 'nftables', got %s", status.Backend)
	}
}

func TestAllowPortValidation(t *testing.T) {
	if err := firewall.AllowPort("8080/tcp"); err != nil {
		t.Fatalf("expected valid port 8080/tcp, got: %v", err)
	}

	invalid := []string{"invalid", "80", "abc/tcp", "80/http", "70000/tcp"}
	for _, inv := range invalid {
		if err := firewall.AllowPort(inv); err == nil {
			t.Errorf("expected error for invalid input %q, got nil", inv)
		}
	}
}
