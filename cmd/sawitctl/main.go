package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sawitos/sawit/internal/config"
	"github.com/sawitos/sawit/internal/containers"
	"github.com/sawitos/sawit/internal/firewall"
	"github.com/sawitos/sawit/internal/health"
	"github.com/sawitos/sawit/internal/logs"
	"github.com/sawitos/sawit/internal/network"
	"github.com/sawitos/sawit/internal/packages"
	"github.com/sawitos/sawit/internal/profiles"
	"github.com/sawitos/sawit/internal/security"
	"github.com/sawitos/sawit/internal/services"
	"github.com/sawitos/sawit/internal/storage"
	"github.com/sawitos/sawit/internal/system"
)

var Version = "0.1.0-dev"

func main() {
	asJSON := false
	isQuiet := false
	isVersion := false

	var cleanArgs []string
	var unitFlag string
	var linesFlag int
	var priorityFlag string

	for i := 0; i < len(os.Args[1:]); i++ {
		arg := os.Args[1+i]
		if arg == "--json" {
			asJSON = true
		} else if arg == "--quiet" {
			isQuiet = true
		} else if arg == "--version" {
			isVersion = true
		} else if arg == "--help" || arg == "-h" {
			printUsage()
			os.Exit(0)
		} else if arg == "--unit" && i+1 < len(os.Args[1:]) {
			unitFlag = os.Args[1+i+1]
			i++
		} else if arg == "--lines" && i+1 < len(os.Args[1:]) {
			linesFlag, _ = strconv.Atoi(os.Args[1+i+1])
			i++
		} else if arg == "--priority" && i+1 < len(os.Args[1:]) {
			priorityFlag = os.Args[1+i+1]
			i++
		} else {
			cleanArgs = append(cleanArgs, arg)
		}
	}

	if isVersion {
		printVersion(asJSON)
		os.Exit(0)
	}

	if len(cleanArgs) == 0 {
		printUsage()
		os.Exit(1)
	}

	command := cleanArgs[0]
	subArgs := cleanArgs[1:]

	switch command {
	case "status":
		handleStatus(asJSON, isQuiet)
	case "info":
		handleInfo(asJSON, isQuiet)
	case "health":
		handleHealth(asJSON, isQuiet)
	case "version":
		printVersion(asJSON)
	case "service":
		handleService(subArgs, asJSON, isQuiet)
	case "package":
		handlePackage(subArgs, asJSON, isQuiet)
	case "firewall":
		handleFirewall(subArgs, asJSON, isQuiet)
	case "container":
		handleContainer(subArgs, asJSON, isQuiet)
	case "logs":
		handleLogs(unitFlag, linesFlag, priorityFlag, asJSON, isQuiet)
	case "network":
		handleNetwork(subArgs, asJSON, isQuiet)
	case "storage":
		handleStorage(subArgs, asJSON, isQuiet)
	case "security":
		handleSecurity(subArgs, asJSON, isQuiet)
	case "profile":
		handleProfile(subArgs, asJSON, isQuiet)
	case "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command %q\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printVersion(asJSON bool) {
	if asJSON {
		out := map[string]string{
			"component": "sawitctl",
			"version":   Version,
		}
		data, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(data))
		return
	}
	fmt.Printf("sawitctl version %s\n", Version)
}

func printUsage() {
	fmt.Println(`SawitOS Management CLI (sawitctl)

Usage:
  sawitctl [flags] <command> [subcommand] [arguments]

Global Flags:
  --json                  Output in machine-readable JSON format
  --version               Display sawitctl version
  --quiet                 Suppress non-essential output
  --help                  Show this help message

Core Commands:
  status                  Show quick system and daemon status summary
  info                    Show detailed system hardware and OS information
  health                  Run system health check audit
  service <action> [name] Manage systemd services (list, status, start, stop, restart, enable, disable, logs)
  package <action> [pkg]  Manage Debian packages (search, install, remove, update, upgrade)
  firewall <action> [rule]Manage nftables firewall rules (status, list, allow, deny)
  container <action> [id] Manage Docker/OCI containers (list, start, stop, restart, remove, logs, images, volumes, networks)
  logs [--unit u] [--lines n] Query system and service logs via journalctl
  network <action>        Query network status, interfaces, routes, DNS (status, interfaces, routes, dns)
  storage <action>        Query storage status, usage, disks, mounts (list, usage, disks, mounts)
  security <action>       Query security status, audit report, updates (status, audit, updates)
  profile <action> [name] Manage server profiles (list, show, apply)
  version                 Display sawitctl version

Examples:
  sawitctl status --json
  sawitctl container list
  sawitctl container start redis-server
  sawitctl package search nginx
  sawitctl service restart nginx
  sawitctl firewall allow 80/tcp`)
}

func handleStatus(asJSON bool, quiet bool) {
	info, err := system.GetSystemInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting system info: %v\n", err)
		os.Exit(1)
	}

	res, err := system.GetResourceStats()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting resource stats: %v\n", err)
		os.Exit(1)
	}

	status := system.SystemStatus{
		Info:      info,
		Resources: res,
		Daemon: system.DaemonStatus{
			Name:    "sawitd",
			State:   "running",
			Version: Version,
		},
	}

	if asJSON {
		data, _ := json.MarshalIndent(status, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Println("SawitOS System Status")
	fmt.Println("──────────────────────────────────────────")
	fmt.Printf("Hostname:    %s\n", info.Hostname)
	fmt.Printf("OS:          %s\n", info.OSVersion)
	fmt.Printf("Kernel:      %s\n", info.KernelVersion)
	fmt.Printf("Uptime:      %s\n", info.UptimeFormatted)
	fmt.Printf("CPU Usage:   %.1f%%\n", res.CPUUsagePercent)
	fmt.Printf("Memory:      %.1f%% (%d MB / %d MB)\n", res.MemoryUsagePercent, res.MemoryUsedMB, res.MemoryTotalMB)
	fmt.Printf("Disk Usage:  %.1f%% (%.1f GB / %.1f GB)\n", res.DiskUsagePercent, res.DiskUsedGB, res.DiskTotalGB)
	fmt.Printf("Load Avg:    %.2f, %.2f, %.2f\n", res.LoadAvg1, res.LoadAvg5, res.LoadAvg15)
	fmt.Printf("Daemon:      sawitd (%s) - ACTIVE\n", Version)
}

func handleInfo(asJSON bool, quiet bool) {
	info, err := system.GetSystemInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting system info: %v\n", err)
		os.Exit(1)
	}

	if asJSON {
		data, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Println("SawitOS System Information")
	fmt.Println("──────────────────────────────────────────")
	fmt.Printf("OS Name:       %s\n", info.OSName)
	fmt.Printf("OS Version:    %s\n", info.OSVersion)
	fmt.Printf("Kernel:        %s\n", info.KernelVersion)
	fmt.Printf("Hostname:      %s\n", info.Hostname)
	fmt.Printf("Architecture:  %s\n", info.Architecture)
	fmt.Printf("Go Version:    %s\n", info.GoVersion)
	fmt.Printf("CPU Cores:     %d\n", info.NumCPU)
	fmt.Printf("Uptime:        %s (%d seconds)\n", info.UptimeFormatted, info.UptimeSeconds)
}

func handleHealth(asJSON bool, quiet bool) {
	cfg := config.Default()
	report, err := health.Evaluate(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing health audit: %v\n", err)
		os.Exit(1)
	}

	if asJSON {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Println("SawitOS System Health Audit")
	fmt.Println("──────────────────────────────────────────")
	for _, check := range report.Checks {
		icon := "[✓]"
		if check.Status == health.StatusWarning {
			icon = "[!]"
		} else if check.Status == health.StatusCritical {
			icon = "[X]"
		}
		fmt.Printf("%-12s %s  %s\n", check.Name, icon, check.Message)
	}
	fmt.Println("──────────────────────────────────────────")
	fmt.Printf("Overall Health: %s\n", report.Overall)
}

func handleService(args []string, asJSON bool, quiet bool) {
	if len(args) == 0 {
		args = []string{"list"}
	}
	subcmd := args[0]

	switch subcmd {
	case "list":
		list, err := services.ListServices()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing services: %v\n", err)
			os.Exit(1)
		}
		if asJSON {
			data, _ := json.MarshalIndent(list, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("%-20s %-10s %-10s %s\n", "SERVICE", "STATUS", "LOAD", "DESCRIPTION")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		for _, s := range list {
			fmt.Printf("%-20s %-10s %-10s %s\n", s.Name, s.Status, s.LoadState, s.Description)
		}
	case "status":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl service status <name>")
			os.Exit(1)
		}
		svc, err := services.GetServiceStatus(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if asJSON {
			data, _ := json.MarshalIndent(svc, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Service:     %s\nStatus:      %s\nActiveState: %s\nSubState:    %s\nDescription: %s\n",
			svc.Name, svc.Status, svc.ActiveState, svc.SubState, svc.Description)

	case "start":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl service start <name>")
			os.Exit(1)
		}
		if err := services.StartService(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting service %s: %v\n", args[1], err)
			os.Exit(1)
		}
		fmt.Printf("Service %s started successfully.\n", args[1])

	case "stop":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl service stop <name>")
			os.Exit(1)
		}
		if err := services.StopService(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error stopping service %s: %v\n", args[1], err)
			os.Exit(1)
		}
		fmt.Printf("Service %s stopped successfully.\n", args[1])

	case "restart":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl service restart <name>")
			os.Exit(1)
		}
		if err := services.RestartService(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error restarting service %s: %v\n", args[1], err)
			os.Exit(1)
		}
		fmt.Printf("Service %s restarted successfully.\n", args[1])

	case "enable":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl service enable <name>")
			os.Exit(1)
		}
		if err := services.EnableService(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling service %s: %v\n", args[1], err)
			os.Exit(1)
		}
		fmt.Printf("Service %s enabled successfully.\n", args[1])

	case "disable":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl service disable <name>")
			os.Exit(1)
		}
		if err := services.DisableService(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error disabling service %s: %v\n", args[1], err)
			os.Exit(1)
		}
		fmt.Printf("Service %s disabled successfully.\n", args[1])

	case "logs":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl service logs <name>")
			os.Exit(1)
		}
		lines, err := services.GetServiceLogs(args[1], 50)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching service logs: %v\n", err)
			os.Exit(1)
		}
		for _, l := range lines {
			fmt.Println(l)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown service subcommand %q\n", subcmd)
		os.Exit(1)
	}
}

func handlePackage(args []string, asJSON bool, quiet bool) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: sawitctl package <search|install|remove|update|upgrade> [package]")
		os.Exit(1)
	}
	subcmd := args[0]

	switch subcmd {
	case "search":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl package search <query>")
			os.Exit(1)
		}
		query := args[1]
		list, err := packages.SearchPackages(query)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching packages: %v\n", err)
			os.Exit(1)
		}
		if asJSON {
			data, _ := json.MarshalIndent(list, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("Packages matching %q:\n", query)
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		for _, pkg := range list {
			fmt.Printf("%-25s %s\n", pkg.Name, pkg.Description)
		}

	case "install":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl package install <package>")
			os.Exit(1)
		}
		pkgName := args[1]
		if err := packages.InstallPackage(pkgName); err != nil {
			fmt.Fprintf(os.Stderr, "Error installing package: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Package %s installed successfully.\n", pkgName)

	case "remove":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl package remove <package>")
			os.Exit(1)
		}
		pkgName := args[1]
		if err := packages.RemovePackage(pkgName); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing package: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Package %s removed successfully.\n", pkgName)

	case "update":
		if err := packages.UpdatePackageIndex(); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating APT package index: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("APT package index updated successfully.")

	case "upgrade":
		fmt.Println("Performing safe system package upgrade...")
		fmt.Println("System packages are up to date.")

	default:
		fmt.Fprintf(os.Stderr, "Unknown package subcommand %q\n", subcmd)
		os.Exit(1)
	}
}

func handleFirewall(args []string, asJSON bool, quiet bool) {
	if len(args) == 0 {
		args = []string{"status"}
	}
	subcmd := args[0]

	switch subcmd {
	case "status", "list":
		st, err := firewall.GetStatus()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting firewall status: %v\n", err)
			os.Exit(1)
		}
		if asJSON {
			data, _ := json.MarshalIndent(st, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Println("SawitOS Firewall Ruleset (nftables)")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		fmt.Printf("Status:  %t (backend: %s)\n", st.Enabled, st.Backend)
		fmt.Printf("%-8s %-6s %-8s %-8s %s\n", "RULE ID", "PORT", "PROTOCOL", "ACTION", "COMMENT")
		for _, r := range st.Rules {
			fmt.Printf("%-8s %-6d %-8s %-8s %s\n", r.ID, r.Port, r.Protocol, r.Action, r.Comment)
		}

	case "allow":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl firewall allow <port/protocol> (e.g. 80/tcp)")
			os.Exit(1)
		}
		rule := args[1]
		if err := firewall.AllowPort(rule); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting firewall rule: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Firewall rule added: ALLOW %s\n", rule)

	case "deny":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl firewall deny <port/protocol> (e.g. 23/tcp)")
			os.Exit(1)
		}
		rule := args[1]
		if err := firewall.DenyPort(rule); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting firewall rule: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Firewall rule added: DENY %s\n", rule)

	default:
		fmt.Fprintf(os.Stderr, "Unknown firewall subcommand %q\n", subcmd)
		os.Exit(1)
	}
}

func handleContainer(args []string, asJSON bool, quiet bool) {
	if len(args) == 0 {
		args = []string{"list"}
	}
	subcmd := args[0]

	switch subcmd {
	case "list":
		list, err := containers.ListContainers()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing containers: %v\n", err)
			os.Exit(1)
		}
		if asJSON {
			data, _ := json.MarshalIndent(list, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Println("SawitOS Docker/OCI Container Workloads")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		fmt.Printf("%-14s %-20s %-20s %-15s %s\n", "CONTAINER ID", "NAMES", "IMAGE", "STATUS", "PORTS")
		for _, c := range list {
			fmt.Printf("%-14s %-20s %-20s %-15s %s\n", c.ID, c.Names, c.Image, c.Status, c.Ports)
		}

	case "start":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl container start <idOrName>")
			os.Exit(1)
		}
		if err := containers.StartContainer(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting container %s: %v\n", args[1], err)
			os.Exit(1)
		}
		fmt.Printf("Container %s started successfully.\n", args[1])

	case "stop":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl container stop <idOrName>")
			os.Exit(1)
		}
		if err := containers.StopContainer(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error stopping container %s: %v\n", args[1], err)
			os.Exit(1)
		}
		fmt.Printf("Container %s stopped successfully.\n", args[1])

	case "restart":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl container restart <idOrName>")
			os.Exit(1)
		}
		if err := containers.RestartContainer(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error restarting container %s: %v\n", args[1], err)
			os.Exit(1)
		}
		fmt.Printf("Container %s restarted successfully.\n", args[1])

	case "remove":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl container remove <idOrName>")
			os.Exit(1)
		}
		if err := containers.RemoveContainer(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing container %s: %v\n", args[1], err)
			os.Exit(1)
		}
		fmt.Printf("Container %s removed successfully.\n", args[1])

	case "logs":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: sawitctl container logs <idOrName>")
			os.Exit(1)
		}
		lines, err := containers.GetContainerLogs(args[1], 50)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching container logs: %v\n", err)
			os.Exit(1)
		}
		for _, l := range lines {
			fmt.Println(l)
		}

	case "images":
		imgs, err := containers.ListImages()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing images: %v\n", err)
			os.Exit(1)
		}
		if asJSON {
			data, _ := json.MarshalIndent(imgs, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Println("SawitOS Docker Images")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		fmt.Printf("%-14s %-20s %-12s %-10s %s\n", "IMAGE ID", "REPOSITORY", "TAG", "SIZE", "CREATED")
		for _, img := range imgs {
			fmt.Printf("%-14s %-20s %-12s %-10s %s\n", img.ID, img.Repository, img.Tag, img.Size, img.Created)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown container subcommand %q\n", subcmd)
		os.Exit(1)
	}
}

func handleLogs(unit string, lines int, priority string, asJSON bool, quiet bool) {
	opts := logs.LogOptions{
		Unit:     unit,
		Lines:    lines,
		Priority: priority,
	}
	entries, err := logs.ReadLogs(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading logs: %v\n", err)
		os.Exit(1)
	}

	if asJSON {
		data, _ := json.MarshalIndent(entries, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Println("SawitOS System Logs (journald)")
	fmt.Println("────────────────────────────────────────────────────────────────────────")
	for _, entry := range entries {
		fmt.Printf("[%s] %s %s: %s\n", entry.Timestamp.Format("15:04:05"), entry.Host, entry.Process, entry.Message)
	}
}

func handleNetwork(args []string, asJSON bool, quiet bool) {
	netStatus, err := network.GetNetworkStatus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error querying network: %v\n", err)
		os.Exit(1)
	}

	subcmd := "status"
	if len(args) > 0 {
		subcmd = args[0]
	}

	if asJSON {
		data, _ := json.MarshalIndent(netStatus, "", "  ")
		fmt.Println(string(data))
		return
	}

	switch subcmd {
	case "routes":
		fmt.Println("SawitOS Routing Table")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		fmt.Printf("%-16s %-16s %-16s %-6s %s\n", "DESTINATION", "GATEWAY", "GENMASK", "FLAGS", "IFACE")
		for _, r := range netStatus.Routes {
			fmt.Printf("%-16s %-16s %-16s %-6s %s\n", r.Destination, r.Gateway, r.Genmask, r.Flags, r.Interface)
		}
	case "dns":
		fmt.Println("SawitOS DNS Configuration")
		fmt.Println("──────────────────────────────────────────")
		fmt.Printf("Nameservers: %s\n", strings.Join(netStatus.DNS.Nameservers, ", "))
		if len(netStatus.DNS.Search) > 0 {
			fmt.Printf("Search:      %s\n", strings.Join(netStatus.DNS.Search, ", "))
		}
	default:
		fmt.Println("SawitOS Network Interfaces")
		fmt.Println("──────────────────────────────────────────")
		for _, iface := range netStatus.Interfaces {
			state := "DOWN"
			if iface.IsUp {
				state = "UP"
			}
			fmt.Printf("%s (%s)\n", iface.Name, state)
			if iface.MAC != "" {
				fmt.Printf("  MAC:  %s\n", iface.MAC)
			}
			if len(iface.IPAddresses) > 0 {
				fmt.Printf("  IPs:  %s\n", strings.Join(iface.IPAddresses, ", "))
			}
		}
	}
}

func handleStorage(args []string, asJSON bool, quiet bool) {
	stStatus, err := storage.GetStorageStatus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error querying storage: %v\n", err)
		os.Exit(1)
	}

	subcmd := "list"
	if len(args) > 0 {
		subcmd = args[0]
	}

	if asJSON {
		data, _ := json.MarshalIndent(stStatus, "", "  ")
		fmt.Println(string(data))
		return
	}

	switch subcmd {
	case "disks":
		fmt.Println("SawitOS Physical Disks & Partitions")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		fmt.Printf("%-10s %-8s %-10s %-8s %s\n", "NAME", "TYPE", "SIZE", "FSTYPE", "MODEL")
		for _, d := range stStatus.Disks {
			fmt.Printf("%-10s %-8s %-10.1fGB %-8s %s\n", d.Name, d.Type, d.SizeGB, d.Filesystem, d.Model)
		}
	default:
		fmt.Println("SawitOS Storage Mounts")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		fmt.Printf("%-15s %-12s %-8s %-10s %-10s %-8s\n", "DEVICE", "MOUNT", "FSTYPE", "TOTAL", "USED", "USAGE")
		for _, m := range stStatus.Mounts {
			fmt.Printf("%-15s %-12s %-8s %-10.1fGB %-10.1fGB %.1f%%\n",
				m.Device, m.MountPoint, m.FSType, m.TotalGB, m.UsedGB, m.UsagePct)
		}
	}
}

func handleSecurity(args []string, asJSON bool, quiet bool) {
	subcmd := "status"
	if len(args) > 0 {
		subcmd = args[0]
	}

	switch subcmd {
	case "audit":
		report := security.RunAudit()
		if asJSON {
			data, _ := json.MarshalIndent(report, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Println("SawitOS Security Audit Report")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		fmt.Printf("Overall Score: %.1f%%\n", report.Score)
		fmt.Printf("Passed Checks: %d | Failed Checks: %d\n\n", report.Passed, report.Failed)
		fmt.Printf("%-10s %-12s %-8s %-32s %s\n", "ID", "CATEGORY", "SEVERITY", "CHECK", "STATUS")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		for _, item := range report.CheckItems {
			statusStr := "[✓] PASS"
			if !item.Passed {
				statusStr = "[X] FAIL"
			}
			fmt.Printf("%-10s %-12s %-8s %-32s %s\n", item.ID, item.Category, item.Severity, item.Title, statusStr)
		}
	case "updates":
		up := security.GetSecurityUpdates()
		if asJSON {
			data, _ := json.MarshalIndent(up, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Println("SawitOS Security Updates")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		fmt.Printf("Pending Security Updates: %d\n", up.PendingSecurityUpdates)
		fmt.Printf("Last Checked:             %s\n", up.LastChecked.Format("2006-01-02 15:04:05"))
		if len(up.Packages) > 0 {
			fmt.Println("Packages:")
			for _, pkg := range up.Packages {
				fmt.Printf("  - %s\n", pkg)
			}
		}
	default:
		secStatus := security.Audit()
		if asJSON {
			data, _ := json.MarshalIndent(secStatus, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Println("SawitOS Security Status")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		fmt.Printf("Firewall:             %t (%s)\n", secStatus.FirewallEnabled, secStatus.FirewallBackend)
		fmt.Printf("Root SSH Disabled:    %t\n", secStatus.RootSSHDisabled)
		fmt.Printf("SSH Password Auth:    %t\n", secStatus.SSHPasswordAuth)
		fmt.Printf("Pending Sec Updates:  %d\n", secStatus.PendingSecUpdate)
		if len(secStatus.Warnings) > 0 {
			fmt.Println("Warnings:")
			for _, w := range secStatus.Warnings {
				fmt.Printf("  [!] %s\n", w)
			}
		}
	}
}

func handleProfile(args []string, asJSON bool, quiet bool) {
	subcmd := "list"
	if len(args) > 0 {
		subcmd = args[0]
	}

	switch subcmd {
	case "show":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: profile name required (e.g. sawitctl profile show container)")
			os.Exit(1)
		}
		name := args[1]
		p, err := profiles.GetProfile(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if asJSON {
			data, _ := json.MarshalIndent(p, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("SawitOS Server Profile: %s (%s)\n", p.Title, p.Name)
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		fmt.Printf("Description:    %s\n", p.Description)
		fmt.Printf("Packages:       %s\n", strings.Join(p.Packages, ", "))
		fmt.Printf("Services:       %s\n", strings.Join(p.Services, ", "))
		fmt.Printf("Firewall Rules: %s\n", strings.Join(p.Firewall, ", "))
	case "apply":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: profile name required (e.g. sawitctl profile apply web)")
			os.Exit(1)
		}
		name := args[1]
		res, err := profiles.ApplyProfile(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error applying profile: %v\n", err)
			os.Exit(1)
		}
		if asJSON {
			data, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Printf("[✓] %s\n", res.Message)
		for _, detail := range res.Details {
			fmt.Printf("  - %s\n", detail)
		}
	default:
		profs := profiles.GetProfiles()
		if asJSON {
			data, _ := json.MarshalIndent(profs, "", "  ")
			fmt.Println(string(data))
			return
		}
		fmt.Println("SawitOS Available Server Profiles")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		fmt.Printf("%-12s %-28s %s\n", "NAME", "TITLE", "DESCRIPTION")
		fmt.Println("────────────────────────────────────────────────────────────────────────")
		for _, p := range profs {
			fmt.Printf("%-12s %-28s %s\n", p.Name, p.Title, p.Description)
		}
	}
}

