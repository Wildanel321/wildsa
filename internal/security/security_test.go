package security_test

import (
	"testing"

	"github.com/sawitos/sawit/internal/security"
)

func TestAuditStatus(t *testing.T) {
	status := security.Audit()
	if !status.FirewallEnabled {
		t.Error("expected firewall enabled by default")
	}
}

func TestRunAuditReport(t *testing.T) {
	report := security.RunAudit()
	if report.Score <= 0 {
		t.Errorf("expected positive audit score, got %.1f", report.Score)
	}
	if len(report.CheckItems) == 0 {
		t.Error("expected non-empty security check items")
	}
}

func TestGetSecurityUpdates(t *testing.T) {
	up := security.GetSecurityUpdates()
	if up.PendingSecurityUpdates < 0 {
		t.Errorf("invalid pending updates count %d", up.PendingSecurityUpdates)
	}
}
