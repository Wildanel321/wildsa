package logs

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Process   string    `json:"process"`
	PID       int       `json:"pid"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
}

type LogOptions struct {
	Unit     string `json:"unit"`
	Lines    int    `json:"lines"`
	Priority string `json:"priority"` // err, warning, info
}

// ReadLogs reads system logs using journalctl or fallback test logs
func ReadLogs(opts LogOptions) ([]LogEntry, error) {
	linesCount := 50
	if opts.Lines > 0 {
		linesCount = opts.Lines
	}

	if runtime.GOOS != "linux" {
		// Development mode log entry generator for non-Linux runtime test suites
		entries := []LogEntry{
			{
				Timestamp: time.Now().Add(-5 * time.Minute),
				Host:      "sawit-server",
				Process:   "sawitd",
				PID:       1024,
				Message:   "SawitOS Core Daemon started successfully on port 8080",
				Severity:  "info",
			},
			{
				Timestamp: time.Now().Add(-2 * time.Minute),
				Host:      "sawit-server",
				Process:   "sawit-agent",
				PID:       1025,
				Message:   "Listening on IPC socket /run/sawit/sawit-agent.sock",
				Severity:  "info",
			},
			{
				Timestamp: time.Now().Add(-1 * time.Minute),
				Host:      "sawit-server",
				Process:   "systemd",
				PID:       1,
				Message:   "Reached target Multi-User System Target",
				Severity:  "info",
			},
		}

		if opts.Unit != "" {
			var filtered []LogEntry
			for _, e := range entries {
				if strings.Contains(e.Process, opts.Unit) {
					filtered = append(filtered, e)
				}
			}
			return filtered, nil
		}
		return entries, nil
	}

	args := []string{"-n", strconv.Itoa(linesCount), "--no-pager", "--output=short-iso"}
	if opts.Unit != "" {
		args = append(args, "-u", opts.Unit)
	}
	if opts.Priority != "" {
		args = append(args, "-p", opts.Priority)
	}

	cmd := exec.Command("journalctl", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("journalctl failed: %w", err)
	}

	var results []LogEntry
	rawLines := strings.Split(string(output), "\n")
	for _, l := range rawLines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "--") {
			continue
		}
		fields := strings.Fields(l)
		if len(fields) < 4 {
			continue
		}

		tsStr := fields[0]
		host := fields[1]
		proc := fields[2]
		msg := strings.Join(fields[3:], " ")

		parsedTime, err := time.Parse(time.RFC3339, tsStr)
		if err != nil {
			parsedTime = time.Now()
		}

		results = append(results, LogEntry{
			Timestamp: parsedTime,
			Host:      host,
			Process:   proc,
			Message:   msg,
			Severity:  "info",
		})
	}
	return results, nil
}
