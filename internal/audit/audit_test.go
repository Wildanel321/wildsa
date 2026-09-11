package audit_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sawitos/sawit/internal/audit"
)

func TestAuditRecordFile(t *testing.T) {
	tmpDir := t.TempDir()
	auditFile := filepath.Join(tmpDir, "audit.log")

	audit.SetAuditLogPath(auditFile)
	audit.Record("admin", "192.168.1.100", "service.restart", "nginx", "SUCCESS")

	content, err := os.ReadFile(auditFile)
	if err != nil {
		t.Fatalf("failed to read audit log file: %v", err)
	}

	strContent := string(content)
	if !strings.Contains(strContent, `"user":"admin"`) {
		t.Errorf("expected audit record with user 'admin', got: %s", strContent)
	}
	if !strings.Contains(strContent, `"action":"service.restart"`) {
		t.Errorf("expected action 'service.restart', got: %s", strContent)
	}
}
