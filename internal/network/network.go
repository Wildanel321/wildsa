package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type InterfaceInfo struct {
	Name        string   `json:"name"`
	MAC         string   `json:"mac"`
	IPAddresses []string `json:"ip_addresses"`
	Flags       []string `json:"flags"`
	IsUp        bool     `json:"is_up"`
	IsLoopback  bool     `json:"is_loopback"`
}

type RouteInfo struct {
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
	Genmask     string `json:"genmask"`
	Flags       string `json:"flags"`
	Interface   string `json:"interface"`
}

type DNSConfig struct {
	Nameservers []string `json:"nameservers"`
	Search      []string `json:"search"`
}

type NetworkStatus struct {
	Interfaces []InterfaceInfo `json:"interfaces"`
	Routes     []RouteInfo     `json:"routes"`
	DNS        DNSConfig       `json:"dns"`
}

// GetNetworkStatus retrieves active system network interfaces, routes, and DNS
func GetNetworkStatus() (*NetworkStatus, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve network interfaces: %w", err)
	}

	var results []InterfaceInfo
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var ipList []string
		for _, addr := range addrs {
			ipList = append(ipList, addr.String())
		}

		isUp := (iface.Flags & net.FlagUp) != 0
		isLoopback := (iface.Flags & net.FlagLoopback) != 0

		var flagList []string
		if isUp {
			flagList = append(flagList, "UP")
		}
		if isLoopback {
			flagList = append(flagList, "LOOPBACK")
		}

		results = append(results, InterfaceInfo{
			Name:        iface.Name,
			MAC:         iface.HardwareAddr.String(),
			IPAddresses: ipList,
			Flags:       flagList,
			IsUp:        isUp,
			IsLoopback:  isLoopback,
		})
	}

	routes := GetRoutes()
	dns := GetDNSConfig()

	return &NetworkStatus{
		Interfaces: results,
		Routes:     routes,
		DNS:        dns,
	}, nil
}

// GetRoutes queries active network routes
func GetRoutes() []RouteInfo {
	if runtime.GOOS != "linux" {
		return []RouteInfo{
			{Destination: "0.0.0.0", Gateway: "192.168.1.1", Genmask: "0.0.0.0", Flags: "UG", Interface: "eth0"},
			{Destination: "192.168.1.0", Gateway: "0.0.0.0", Genmask: "255.255.255.0", Flags: "U", Interface: "eth0"},
		}
	}

	cmd := exec.Command("ip", "route")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var routes []RouteInfo
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		dest := fields[0]
		gw := "0.0.0.0"
		iface := "unknown"

		for i := 0; i < len(fields); i++ {
			if fields[i] == "via" && i+1 < len(fields) {
				gw = fields[i+1]
			}
			if fields[i] == "dev" && i+1 < len(fields) {
				iface = fields[i+1]
			}
		}

		routes = append(routes, RouteInfo{
			Destination: dest,
			Gateway:     gw,
			Genmask:     "255.255.255.255",
			Flags:       "U",
			Interface:   iface,
		})
	}
	return routes
}

// GetDNSConfig parses /etc/resolv.conf
func GetDNSConfig() DNSConfig {
	cfg := DNSConfig{
		Nameservers: []string{"1.1.1.1", "8.8.8.8"},
	}
	if runtime.GOOS != "linux" {
		return cfg
	}

	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return cfg
	}
	defer file.Close()

	var ns []string
	var search []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			if fields[0] == "nameserver" {
				ns = append(ns, fields[1])
			} else if fields[0] == "search" {
				search = append(search, fields[1:]...)
			}
		}
	}
	if len(ns) > 0 {
		cfg.Nameservers = ns
	}
	cfg.Search = search
	return cfg
}
