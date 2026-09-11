package logs_test

import (
	"testing"

	"github.com/sawitos/sawit/internal/logs"
)

func TestReadLogs(t *testing.T) {
	opts := logs.LogOptions{
		Lines: 10,
	}
	entries, err := logs.ReadLogs(opts)
	if err != nil {
		t.Fatalf("unexpected error reading logs: %v", err)
	}
	if len(entries) == 0 {
		t.Error("expected non-empty log entries")
	}
}

func TestReadLogsFiltered(t *testing.T) {
	opts := logs.LogOptions{
		Unit:  "sawitd",
		Lines: 5,
	}
	entries, err := logs.ReadLogs(opts)
	if err != nil {
		t.Fatalf("unexpected error reading filtered logs: %v", err)
	}
	if len(entries) == 0 {
		t.Error("expected log entries matching unit 'sawitd'")
	}
}
