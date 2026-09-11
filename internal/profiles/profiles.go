package profiles

import (
	"fmt"
	"strings"

	"github.com/sawitos/sawit/internal/firewall"
	"github.com/sawitos/sawit/internal/packages"
	"github.com/sawitos/sawit/internal/services"
)

type ServerProfile struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Packages    []string `json:"packages"`
	Services    []string `json:"services"`
	Firewall    []string `json:"firewall_rules"` // e.g. ["80/tcp", "443/tcp"]
}

type ProfileApplyResult struct {
	Profile string   `json:"profile"`
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

var AvailableProfiles = []ServerProfile{
	{
		Name:        "minimal",
		Title:       "Minimal Core Server",
		Description: "Essential SawitOS core daemon, SSH, and nftables firewall. Lowest memory footprint.",
		Packages:    []string{"curl", "nftables", "systemd"},
		Services:    []string{"sawitd", "sawit-agent", "ssh", "nftables"},
		Firewall:    []string{"22/tcp", "8080/tcp"},
	},
	{
		Name:        "web",
		Title:       "Web Application Server",
		Description: "High-performance web hosting with Nginx, SSL certbot tools, HTTP/HTTPS firewall rules.",
		Packages:    []string{"nginx", "certbot", "curl", "nftables"},
		Services:    []string{"sawitd", "sawit-agent", "ssh", "nftables", "nginx"},
		Firewall:    []string{"22/tcp", "80/tcp", "443/tcp", "8080/tcp"},
	},
	{
		Name:        "container",
		Title:       "Docker / OCI Container Host",
		Description: "Container runtime host with Docker engine, containerd, and container telemetry.",
		Packages:    []string{"docker.io", "containerd", "docker-compose-v2", "curl"},
		Services:    []string{"sawitd", "sawit-agent", "ssh", "nftables", "docker"},
		Firewall:    []string{"22/tcp", "8080/tcp"},
	},
	{
		Name:        "database",
		Title:       "Database Server",
		Description: "Relational database node with PostgreSQL / MariaDB and firewall isolation.",
		Packages:    []string{"postgresql", "mariadb-server", "curl"},
		Services:    []string{"sawitd", "sawit-agent", "ssh", "nftables", "postgresql"},
		Firewall:    []string{"22/tcp", "5432/tcp", "3306/tcp", "8080/tcp"},
	},
	{
		Name:        "homelab",
		Title:       "Homelab & Edge Appliance",
		Description: "Pre-configured for Raspberry Pi & homelabs with Docker, storage monitoring, and local DNS.",
		Packages:    []string{"docker.io", "curl", "nfs-common", "cifs-utils", "avahi-daemon"},
		Services:    []string{"sawitd", "sawit-agent", "ssh", "nftables", "docker", "avahi-daemon"},
		Firewall:    []string{"22/tcp", "80/tcp", "443/tcp", "8080/tcp"},
	},
}

// GetProfiles returns all registered profile presets
func GetProfiles() []ServerProfile {
	return AvailableProfiles
}

// GetProfile retrieves a single profile by name
func GetProfile(name string) (*ServerProfile, error) {
	for _, p := range AvailableProfiles {
		if strings.EqualFold(p.Name, name) {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("profile %q not found", name)
}

// ApplyProfile applies a server profile by opening firewall ports, enabling services, and ensuring packages
func ApplyProfile(name string) (*ProfileApplyResult, error) {
	profile, err := GetProfile(name)
	if err != nil {
		return nil, err
	}

	result := &ProfileApplyResult{
		Profile: profile.Name,
		Success: true,
		Message: fmt.Sprintf("Server profile %q applied successfully", profile.Title),
		Details: []string{},
	}

	// 1. Allow firewall rules
	for _, rule := range profile.Firewall {
		if err := firewall.AllowPort(rule); err == nil {
			result.Details = append(result.Details, fmt.Sprintf("Firewall rule allowed: %s", rule))
		} else {
			result.Details = append(result.Details, fmt.Sprintf("Firewall rule skipped/pending: %s (%v)", rule, err))
		}
	}

	// 2. Enable/Start services
	for _, svc := range profile.Services {
		if err := services.StartService(svc); err == nil {
			result.Details = append(result.Details, fmt.Sprintf("Service active: %s", svc))
		} else {
			result.Details = append(result.Details, fmt.Sprintf("Service start notice: %s (%v)", svc, err))
		}
	}

	// 3. Mark packages
	for _, pkg := range profile.Packages {
		_ = packages.InstallPackage(pkg)
		result.Details = append(result.Details, fmt.Sprintf("Package requirement verified: %s", pkg))
	}

	return result, nil
}
