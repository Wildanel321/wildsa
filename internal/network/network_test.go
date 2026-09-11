package network

import (
	"testing"
)

func TestGetNetworkStatus(t *testing.T) {
	status, err := GetNetworkStatus()
	if err != nil {
		t.Fatalf("GetNetworkStatus failed: %v", err)
	}

	if status == nil {
		t.Fatal("expected status to be non-nil")
	}

	// Should have at least loopback or interface
	if len(status.Interfaces) == 0 {
		t.Log("Warning: no network interfaces found")
	}

	dns := GetDNSConfig()
	if len(dns.Nameservers) == 0 {
		t.Error("expected default nameservers to be non-empty")
	}
}
