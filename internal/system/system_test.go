package system_test

import (
	"testing"

	"github.com/sawitos/sawit/internal/system"
)

func TestGetSystemInfo(t *testing.T) {
	info, err := system.GetSystemInfo()
	if err != nil {
		t.Fatalf("unexpected error getting system info: %v", err)
	}
	if info.Hostname == "" {
		t.Error("expected non-empty hostname")
	}
	if info.Architecture == "" {
		t.Error("expected non-empty architecture")
	}
}

func TestGetResourceStats(t *testing.T) {
	stats, err := system.GetResourceStats()
	if err != nil {
		t.Fatalf("unexpected error getting resource stats: %v", err)
	}
	if stats.MemoryTotalMB == 0 {
		t.Error("expected MemoryTotalMB > 0")
	}
	if stats.DiskTotalGB == 0 {
		t.Error("expected DiskTotalGB > 0")
	}
}

func TestFormatUptime(t *testing.T) {
	tests := []struct {
		sec      uint64
		expected string
	}{
		{45, "45s"},
		{125, "2m 5s"},
		{3665, "1h 1m 5s"},
		{90000, "1d 1h 0m"},
	}

	for _, tt := range tests {
		result := system.FormatUptime(tt.sec)
		if result != tt.expected {
			t.Errorf("FormatUptime(%d) = %q, expected %q", tt.sec, result, tt.expected)
		}
	}
}
