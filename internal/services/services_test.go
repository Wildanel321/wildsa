package services

import (
	"testing"
)

func TestListServices(t *testing.T) {
	list, err := ListServices()
	if err != nil {
		t.Fatalf("ListServices failed: %v", err)
	}

	if len(list) == 0 {
		t.Error("expected at least mock or system services")
	}

	// Test GetServiceStatus for existing service or fallback
	svc, err := GetServiceStatus(list[0].Name)
	if err != nil {
		t.Fatalf("GetServiceStatus failed for %s: %v", list[0].Name, err)
	}

	if svc == nil {
		t.Fatal("expected service info to be non-nil")
	}
}

func TestServiceActionsEmptyName(t *testing.T) {
	if err := StartService(""); err == nil {
		t.Error("expected error when starting service with empty name")
	}

	if err := StopService(""); err == nil {
		t.Error("expected error when stopping service with empty name")
	}

	if err := RestartService(""); err == nil {
		t.Error("expected error when restarting service with empty name")
	}
}
