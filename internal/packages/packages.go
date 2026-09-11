package packages

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type PackageInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Status      string `json:"status"` // installed, available
}

type PackageUpdateStatus struct {
	UpdatesAvailable int           `json:"updates_available"`
	Packages         []PackageInfo `json:"packages"`
}

// SearchPackages searches Debian APT package repository
func SearchPackages(query string) ([]PackageInfo, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	if runtime.GOOS != "linux" {
		// Non-Linux development mode fallback data
		return []PackageInfo{
			{Name: query, Version: "1.24.0-1", Description: "High performance HTTP and reverse proxy server", Status: "available"},
			{Name: query + "-common", Version: "1.24.0-1", Description: "Common files for " + query, Status: "available"},
		}, nil
	}

	cmd := exec.Command("apt-cache", "search", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("apt-cache search failed: %w", err)
	}

	var results []PackageInfo
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " - ", 2)
		name := parts[0]
		desc := ""
		if len(parts) > 1 {
			desc = parts[1]
		}
		results = append(results, PackageInfo{
			Name:        name,
			Version:     "available",
			Description: desc,
			Status:      "available",
		})
	}
	return results, nil
}

// UpdatePackageIndex refreshes APT repository lists
func UpdatePackageIndex() error {
	if runtime.GOOS != "linux" {
		return nil
	}
	cmd := exec.Command("apt-get", "update", "-y")
	cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("apt-get update failed: %s: %w", string(out), err)
	}
	return nil
}

// InstallPackage installs specified package via APT
func InstallPackage(pkgName string) error {
	if pkgName == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	if runtime.GOOS != "linux" {
		return nil
	}
	cmd := exec.Command("apt-get", "install", "-y", pkgName)
	cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("apt-get install %s failed: %s: %w", pkgName, string(out), err)
	}
	return nil
}

// RemovePackage removes specified package via APT
func RemovePackage(pkgName string) error {
	if pkgName == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	if runtime.GOOS != "linux" {
		return nil
	}
	cmd := exec.Command("apt-get", "remove", "-y", pkgName)
	cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("apt-get remove %s failed: %s: %w", pkgName, string(out), err)
	}
	return nil
}
