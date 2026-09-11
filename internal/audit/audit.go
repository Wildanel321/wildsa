package audit

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user"`
	SourceIP  string    `json:"source_ip"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	Status    string    `json:"status"` // SUCCESS, FAILED
}

var (
	auditMutex sync.Mutex
	logPath    = "/var/log/sawit/audit.log"
)

// SetAuditLogPath configures target audit log output path
func SetAuditLogPath(path string) {
	auditMutex.Lock()
	defer auditMutex.Unlock()
	logPath = path
}

// Record logs an administrative security audit entry to the audit log file
func Record(user, sourceIP, action, target, status string) {
	entry := AuditEntry{
		Timestamp: time.Now(),
		User:      user,
		SourceIP:  sourceIP,
		Action:    action,
		Target:    target,
		Status:    status,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return
	}

	auditMutex.Lock()
	defer auditMutex.Unlock()

	// Append to audit log file or print to stdout if file cannot be opened
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("[AUDIT] %s", string(data))
		return
	}
	defer f.Close()

	fmt.Fprintln(f, string(data))
}
