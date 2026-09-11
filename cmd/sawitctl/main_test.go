package main_test

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSawitctlBuildAndVersion(t *testing.T) {
	tmpDir := t.TempDir()
	binName := "sawitctl"
	if runtime.GOOS == "windows" {
		binName = "sawitctl.exe"
	}
	binaryPath := filepath.Join(tmpDir, binName)

	// Build sawitctl binary
	cmdBuild := exec.Command("go", "build", "-o", binaryPath, "./cmd/sawitctl")
	cmdBuild.Dir = "../.."
	if out, err := cmdBuild.CombinedOutput(); err != nil {
		t.Fatalf("failed to build sawitctl: %v\nOutput:\n%s", err, string(out))
	}

	// Run sawitctl --version
	cmdVersion := exec.Command(binaryPath, "--version")
	out, err := cmdVersion.CombinedOutput()
	if err != nil {
		t.Fatalf("sawitctl --version failed: %v", err)
	}
	if !strings.Contains(string(out), "sawitctl version") {
		t.Errorf("expected version output, got: %s", string(out))
	}

	// Run sawitctl status --json
	cmdStatusJSON := exec.Command(binaryPath, "status", "--json")
	outJSON, err := cmdStatusJSON.CombinedOutput()
	if err != nil {
		t.Fatalf("sawitctl status --json failed: %v", err)
	}
	if !strings.Contains(string(outJSON), `"os_name"`) {
		t.Errorf("expected json status output with 'os_name', got: %s", string(outJSON))
	}

	// Run sawitctl health --json
	cmdHealthJSON := exec.Command(binaryPath, "health", "--json")
	outHealthJSON, err := cmdHealthJSON.CombinedOutput()
	if err != nil {
		t.Fatalf("sawitctl health --json failed: %v", err)
	}
	if !strings.Contains(string(outHealthJSON), `"overall"`) {
		t.Errorf("expected json health output with 'overall', got: %s", string(outHealthJSON))
	}
}
