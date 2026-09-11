package ipc

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

type Action string

const (
	ActionSystemInfo       Action = "system.info"
	ActionHealthCheck      Action = "health.check"
	ActionServiceStatus    Action = "service.status"
	ActionServiceStart     Action = "service.start"
	ActionServiceStop      Action = "service.stop"
	ActionServiceRestart   Action = "service.restart"
	ActionServiceEnable    Action = "service.enable"
	ActionServiceDisable   Action = "service.disable"
	ActionPackageSearch    Action = "package.search"
	ActionPackageInstall   Action = "package.install"
	ActionPackageRemove    Action = "package.remove"
	ActionPackageUpdate    Action = "package.update"
	ActionFirewallStatus   Action = "firewall.status"
	ActionFirewallAllow    Action = "firewall.allow"
	ActionFirewallDeny     Action = "firewall.deny"
	ActionContainerList    Action = "container.list"
	ActionContainerStart   Action = "container.start"
	ActionContainerStop    Action = "container.stop"
	ActionContainerRestart Action = "container.restart"
	ActionContainerRemove  Action = "container.remove"
	ActionContainerImages  Action = "container.images"
	ActionSecurityStatus   Action = "security.status"
	ActionSecurityAudit    Action = "security.audit"
	ActionSecurityUpdates  Action = "security.updates"
	ActionProfileList      Action = "profile.list"
	ActionProfileShow      Action = "profile.show"
	ActionProfileApply     Action = "profile.apply"
)

type Request struct {
	ID        string                 `json:"id"`
	Action    Action                 `json:"action"`
	Params    map[string]interface{} `json:"params,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

type Response struct {
	ID        string      `json:"id"`
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// SendRequest sends a JSON IPC request over a socket connection and receives the structured response
func SendRequest(network, address string, req *Request, timeout time.Duration) (*Response, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to IPC socket (%s:%s): %w", network, address, err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(req); err != nil {
		return nil, fmt.Errorf("failed to encode IPC request: %w", err)
	}

	var resp Response
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&resp); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("agent closed connection prematurely")
		}
		return nil, fmt.Errorf("failed to decode IPC response: %w", err)
	}

	return &resp, nil
}
