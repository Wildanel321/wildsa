package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/sawitos/sawit/internal/config"
	"github.com/sawitos/sawit/internal/containers"
	"github.com/sawitos/sawit/internal/firewall"
	"github.com/sawitos/sawit/internal/health"
	"github.com/sawitos/sawit/internal/ipc"
	"github.com/sawitos/sawit/internal/packages"
	"github.com/sawitos/sawit/internal/profiles"
	"github.com/sawitos/sawit/internal/security"
	"github.com/sawitos/sawit/internal/services"
	"github.com/sawitos/sawit/internal/system"
)

var Version = "0.1.0-dev"

func main() {
	socketPath := flag.String("socket", "/run/sawit/sawit-agent.sock", "Path to privileged IPC socket")
	tcpPort := flag.Int("tcp-fallback-port", 0, "Optional TCP fallback port for non-Unix development")
	versionFlag := flag.Bool("version", false, "Print sawit-agent version")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("sawit-agent version %s\n", Version)
		os.Exit(0)
	}

	var listener net.Listener
	var err error

	if *tcpPort > 0 {
		listener, err = net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *tcpPort))
		if err != nil {
			log.Fatalf("Failed to listen on TCP fallback port: %v", err)
		}
		log.Printf("SawitOS Agent (sawit-agent %s) listening on TCP fallback 127.0.0.1:%d", Version, *tcpPort)
	} else {
		// Clean up existing socket file if present
		os.Remove(*socketPath)
		listener, err = net.Listen("unix", *socketPath)
		if err != nil {
			log.Printf("Notice: Unix socket listen on %s failed (%v). Falling back to TCP 127.0.0.1:8081 for dev mode.", *socketPath, err)
			listener, err = net.Listen("tcp", "127.0.0.1:8081")
			if err != nil {
				log.Fatalf("Failed to start agent TCP fallback: %v", err)
			}
		} else {
			log.Printf("SawitOS Agent (sawit-agent %s) listening on Unix socket %s", Version, *socketPath)
		}
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var req ipc.Request
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&req); err != nil {
		if err != io.EOF {
			log.Printf("Error decoding agent request: %v", err)
		}
		return
	}

	resp := processRequest(&req)
	encoder := json.NewEncoder(conn)
	encoder.Encode(resp)
}

func processRequest(req *ipc.Request) *ipc.Response {
	resp := &ipc.Response{
		ID:        req.ID,
		Timestamp: time.Now(),
	}

	switch req.Action {
	case ipc.ActionSystemInfo:
		info, err := system.GetSystemInfo()
		if err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = info
		}

	case ipc.ActionHealthCheck:
		cfg := config.Default()
		report, err := health.Evaluate(cfg)
		if err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = report
		}

	case ipc.ActionServiceStatus:
		svcName, _ := req.Params["service"].(string)
		if svcName == "" {
			resp.Success = false
			resp.Error = "missing required parameter 'service'"
		} else {
			status, err := services.GetServiceStatus(svcName)
			if err != nil {
				resp.Success = false
				resp.Error = err.Error()
			} else {
				resp.Success = true
				resp.Data = status
			}
		}

	case ipc.ActionServiceStart:
		svcName, _ := req.Params["service"].(string)
		if err := services.StartService(svcName); err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = fmt.Sprintf("Service %s started successfully", svcName)
		}

	case ipc.ActionServiceStop:
		svcName, _ := req.Params["service"].(string)
		if err := services.StopService(svcName); err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = fmt.Sprintf("Service %s stopped successfully", svcName)
		}

	case ipc.ActionServiceRestart:
		svcName, _ := req.Params["service"].(string)
		if err := services.RestartService(svcName); err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = fmt.Sprintf("Service %s restarted successfully", svcName)
		}

	case ipc.ActionPackageSearch:
		query, _ := req.Params["query"].(string)
		list, err := packages.SearchPackages(query)
		if err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = list
		}

	case ipc.ActionPackageInstall:
		pkgName, _ := req.Params["package"].(string)
		if err := packages.InstallPackage(pkgName); err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = fmt.Sprintf("Package %s installed successfully", pkgName)
		}

	case ipc.ActionPackageRemove:
		pkgName, _ := req.Params["package"].(string)
		if err := packages.RemovePackage(pkgName); err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = fmt.Sprintf("Package %s removed successfully", pkgName)
		}

	case ipc.ActionFirewallStatus:
		fwStatus, err := firewall.GetStatus()
		if err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = fwStatus
		}

	case ipc.ActionFirewallAllow:
		portProto, _ := req.Params["rule"].(string)
		if err := firewall.AllowPort(portProto); err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = fmt.Sprintf("Firewall rule added: allow %s", portProto)
		}

	case ipc.ActionFirewallDeny:
		portProto, _ := req.Params["rule"].(string)
		if err := firewall.DenyPort(portProto); err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = fmt.Sprintf("Firewall rule added: deny %s", portProto)
		}

	case ipc.ActionContainerList:
		list, err := containers.ListContainers()
		if err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = list
		}

	case ipc.ActionContainerStart:
		id, _ := req.Params["container"].(string)
		if err := containers.StartContainer(id); err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = fmt.Sprintf("Container %s started", id)
		}

	case ipc.ActionContainerStop:
		id, _ := req.Params["container"].(string)
		if err := containers.StopContainer(id); err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = fmt.Sprintf("Container %s stopped", id)
		}

	case ipc.ActionContainerRestart:
		id, _ := req.Params["container"].(string)
		if err := containers.RestartContainer(id); err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = fmt.Sprintf("Container %s restarted", id)
		}

	case ipc.ActionContainerRemove:
		id, _ := req.Params["container"].(string)
		if err := containers.RemoveContainer(id); err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = fmt.Sprintf("Container %s removed", id)
		}

	case ipc.ActionContainerImages:
		imgs, err := containers.ListImages()
		if err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = imgs
		}

	case ipc.ActionSecurityStatus:
		resp.Success = true
		resp.Data = security.Audit()

	case ipc.ActionSecurityAudit:
		resp.Success = true
		resp.Data = security.RunAudit()

	case ipc.ActionSecurityUpdates:
		resp.Success = true
		resp.Data = security.GetSecurityUpdates()

	case ipc.ActionProfileList:
		resp.Success = true
		resp.Data = profiles.GetProfiles()

	case ipc.ActionProfileShow:
		name, _ := req.Params["name"].(string)
		p, err := profiles.GetProfile(name)
		if err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = p
		}

	case ipc.ActionProfileApply:
		name, _ := req.Params["name"].(string)
		res, err := profiles.ApplyProfile(name)
		if err != nil {
			resp.Success = false
			resp.Error = err.Error()
		} else {
			resp.Success = true
			resp.Data = res
		}

	default:
		resp.Success = false
		resp.Error = fmt.Sprintf("unsupported or unvalidated agent action %q", req.Action)
	}

	return resp
}
