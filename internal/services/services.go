package services

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type ServiceInfo struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	ActiveState string `json:"active_state"`
	SubState    string `json:"sub_state"`
	LoadState   string `json:"load_state"`
	Description string `json:"description"`
}

// ListServices queries systemd unit files status
func ListServices() ([]ServiceInfo, error) {
	if runtime.GOOS != "linux" {
		// Mock services for non-Linux runtime development
		return []ServiceInfo{
			{Name: "sawitd", Status: "running", ActiveState: "active", SubState: "running", LoadState: "loaded", Description: "SawitOS Core Daemon"},
			{Name: "sawit-agent", Status: "running", ActiveState: "active", SubState: "running", LoadState: "loaded", Description: "SawitOS Privileged Agent"},
			{Name: "ssh", Status: "running", ActiveState: "active", SubState: "running", LoadState: "loaded", Description: "OpenBSD Secure Shell server"},
			{Name: "nftables", Status: "running", ActiveState: "active", SubState: "exited", LoadState: "loaded", Description: "nftables firewall"},
		}, nil
	}

	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--no-legend")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("systemctl command failed: %w", err)
	}

	var results []ServiceInfo
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		name := strings.TrimSuffix(fields[0], ".service")
		loadState := fields[1]
		activeState := fields[2]
		subState := fields[3]
		desc := ""
		if len(fields) >= 5 {
			desc = strings.Join(fields[4:], " ")
		}

		status := "stopped"
		if activeState == "active" && subState == "running" {
			status = "running"
		}

		results = append(results, ServiceInfo{
			Name:        name,
			Status:      status,
			ActiveState: activeState,
			SubState:    subState,
			LoadState:   loadState,
			Description: desc,
		})
	}
	return results, nil
}

// GetServiceStatus checks specific systemd unit state
func GetServiceStatus(name string) (*ServiceInfo, error) {
	list, err := ListServices()
	if err != nil {
		return nil, err
	}
	for _, svc := range list {
		if svc.Name == name || svc.Name == name+".service" {
			return &svc, nil
		}
	}
	return nil, fmt.Errorf("service %q not found", name)
}

// StartService executes systemctl start <name>
func StartService(name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if runtime.GOOS != "linux" {
		return nil
	}
	cmd := exec.Command("systemctl", "start", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start service %s: %s: %w", name, string(out), err)
	}
	return nil
}

// StopService executes systemctl stop <name>
func StopService(name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if runtime.GOOS != "linux" {
		return nil
	}
	cmd := exec.Command("systemctl", "stop", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop service %s: %s: %w", name, string(out), err)
	}
	return nil
}

// RestartService executes systemctl restart <name>
func RestartService(name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if runtime.GOOS != "linux" {
		return nil
	}
	cmd := exec.Command("systemctl", "restart", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to restart service %s: %s: %w", name, string(out), err)
	}
	return nil
}

// EnableService executes systemctl enable <name>
func EnableService(name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if runtime.GOOS != "linux" {
		return nil
	}
	cmd := exec.Command("systemctl", "enable", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to enable service %s: %s: %w", name, string(out), err)
	}
	return nil
}

// DisableService executes systemctl disable <name>
func DisableService(name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if runtime.GOOS != "linux" {
		return nil
	}
	cmd := exec.Command("systemctl", "disable", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to disable service %s: %s: %w", name, string(out), err)
	}
	return nil
}

// GetServiceLogs fetches recent journalctl lines for unit
func GetServiceLogs(name string, lines int) ([]string, error) {
	if lines <= 0 {
		lines = 50
	}
	if runtime.GOOS != "linux" {
		return []string{
			fmt.Sprintf("Systemd logs for unit %s.service", name),
			"Service active and running normally.",
		}, nil
	}
	cmd := exec.Command("journalctl", "-u", name, "-n", strconv.Itoa(lines), "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to read logs for service %s: %w", name, err)
	}
	return strings.Split(string(output), "\n"), nil
}
