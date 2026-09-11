package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sawitos/sawit/internal/audit"
	"github.com/sawitos/sawit/internal/auth"
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
	"github.com/sawitos/sawit/internal/ws"
)

var Version = "0.1.0-dev"

func main() {
	configPath := flag.String("config", "/etc/sawit/sawitd.yaml", "Path to sawitd configuration file")
	versionFlag := flag.Bool("version", false, "Print sawitd version")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("sawitd version %s\n", Version)
		os.Exit(0)
	}

	cfg, err := config.LoadFromFile(*configPath)
	if err != nil {
		log.Printf("Notice: could not load config from %s (%v), using default settings", *configPath, err)
		cfg = config.Default()
	}

	// Initialize WebSocket Hub
	wsHub := ws.NewHub()
	go wsHub.Run()

	mux := http.NewServeMux()

	// Public Endpoints
	mux.HandleFunc("/api/v1/ping", handlePing)
	mux.HandleFunc("/api/v1/auth/login", handleAuthLogin)

	// Protected Endpoints
	mux.HandleFunc("/api/v1/auth/me", handleAuthMe)
	mux.HandleFunc("/api/v1/auth/logout", handleAuthLogout)
	mux.HandleFunc("/api/v1/system", handleSystem)
	mux.HandleFunc("/api/v1/info", handleInfo)
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) { handleHealth(w, r, cfg) })
	mux.HandleFunc("/api/v1/services", handleServices)
	mux.HandleFunc("/api/v1/services/start", handleServiceStart)
	mux.HandleFunc("/api/v1/services/stop", handleServiceStop)
	mux.HandleFunc("/api/v1/services/restart", handleServiceRestart)
	mux.HandleFunc("/api/v1/network", handleNetwork)
	mux.HandleFunc("/api/v1/storage", handleStorage)
	mux.HandleFunc("/api/v1/packages/search", handlePackageSearch)
	mux.HandleFunc("/api/v1/packages/install", handlePackageInstall)
	mux.HandleFunc("/api/v1/packages/remove", handlePackageRemove)
	mux.HandleFunc("/api/v1/packages/update", handlePackageUpdate)
	mux.HandleFunc("/api/v1/logs", handleLogs)
	mux.HandleFunc("/api/v1/firewall", handleFirewall)
	mux.HandleFunc("/api/v1/firewall/allow", handleFirewallAllow)
	mux.HandleFunc("/api/v1/firewall/deny", handleFirewallDeny)
	mux.HandleFunc("/api/v1/containers", handleContainerList)
	mux.HandleFunc("/api/v1/containers/start", handleContainerStart)
	mux.HandleFunc("/api/v1/containers/stop", handleContainerStop)
	mux.HandleFunc("/api/v1/containers/restart", handleContainerRestart)
	mux.HandleFunc("/api/v1/containers/remove", handleContainerRemove)
	mux.HandleFunc("/api/v1/containers/logs", handleContainerLogs)
	mux.HandleFunc("/api/v1/containers/images", handleContainerImages)
	mux.HandleFunc("/api/v1/containers/volumes", handleContainerVolumes)
	mux.HandleFunc("/api/v1/containers/networks", handleContainerNetworks)
	mux.HandleFunc("/api/v1/security/status", handleSecurityStatus)
	mux.HandleFunc("/api/v1/security/audit", handleSecurityAudit)
	mux.HandleFunc("/api/v1/security/updates", handleSecurityUpdates)
	mux.HandleFunc("/api/v1/profiles", handleProfilesList)
	mux.HandleFunc("/api/v1/profiles/apply", handleProfileApply)
	mux.HandleFunc("/api/v1/ws", func(w http.ResponseWriter, r *http.Request) { wsHub.ServeWebSocket(w, r) })
	mux.HandleFunc("/", handleWebOrStatic)

	// Wrap router with AuthMiddleware
	handler := auth.AuthMiddleware(mux)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("SawitOS Core Daemon (sawitd %s) starting on %s...", Version, addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("sawitd server failure: %v", err)
	}
}

func jsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{
		"status":  "pong",
		"service": "sawitd",
		"version": Version,
		"time":    time.Now().Format(time.RFC3339),
	})
}

func handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	user, err := auth.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		audit.Record(req.Username, r.RemoteAddr, "auth.login", "sawitd", "FAILED")
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	token, err := auth.GenerateToken(user.Username, user.Role)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate session token"})
		return
	}

	audit.Record(user.Username, r.RemoteAddr, "auth.login", "sawitd", "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func handleAuthMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	if !ok {
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{
		"username": claims.Username,
		"role":     claims.Role,
	})
}

func handleAuthLogout(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

func handleSystem(w http.ResponseWriter, r *http.Request) {
	info, err := system.GetSystemInfo()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	res, err := system.GetResourceStats()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	jsonResponse(w, http.StatusOK, system.SystemStatus{
		Info:      info,
		Resources: res,
		Daemon: system.DaemonStatus{
			Name:    "sawitd",
			State:   "running",
			Version: Version,
		},
	})
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	info, err := system.GetSystemInfo()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, info)
}

func handleHealth(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	report, err := health.Evaluate(cfg)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, report)
}

func handleServices(w http.ResponseWriter, r *http.Request) {
	list, err := services.ListServices()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, list)
}

func handleServiceStart(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Service string `json:"service"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Service == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload, required 'service' field"})
		return
	}
	if err := services.StartService(payload.Service); err != nil {
		audit.Record("admin", r.RemoteAddr, "service.start", payload.Service, "FAILED")
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "service.start", payload.Service, "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Service %s started", payload.Service)})
}

func handleServiceStop(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Service string `json:"service"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Service == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload, required 'service' field"})
		return
	}
	if err := services.StopService(payload.Service); err != nil {
		audit.Record("admin", r.RemoteAddr, "service.stop", payload.Service, "FAILED")
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "service.stop", payload.Service, "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Service %s stopped", payload.Service)})
}

func handleServiceRestart(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Service string `json:"service"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Service == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload, required 'service' field"})
		return
	}
	if err := services.RestartService(payload.Service); err != nil {
		audit.Record("admin", r.RemoteAddr, "service.restart", payload.Service, "FAILED")
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "service.restart", payload.Service, "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Service %s restarted", payload.Service)})
}

func handleNetwork(w http.ResponseWriter, r *http.Request) {
	netStatus, err := network.GetNetworkStatus()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, netStatus)
}

func handleStorage(w http.ResponseWriter, r *http.Request) {
	stStatus, err := storage.GetStorageStatus()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, stStatus)
}

func handlePackageSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "query parameter 'q' required"})
		return
	}
	list, err := packages.SearchPackages(query)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, list)
}

func handlePackageInstall(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Package string `json:"package"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Package == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload, required 'package' field"})
		return
	}
	if err := packages.InstallPackage(payload.Package); err != nil {
		audit.Record("admin", r.RemoteAddr, "package.install", payload.Package, "FAILED")
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "package.install", payload.Package, "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Package %s installed", payload.Package)})
}

func handlePackageRemove(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Package string `json:"package"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Package == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload, required 'package' field"})
		return
	}
	if err := packages.RemovePackage(payload.Package); err != nil {
		audit.Record("admin", r.RemoteAddr, "package.remove", payload.Package, "FAILED")
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "package.remove", payload.Package, "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Package %s removed", payload.Package)})
}

func handlePackageUpdate(w http.ResponseWriter, r *http.Request) {
	if err := packages.UpdatePackageIndex(); err != nil {
		audit.Record("admin", r.RemoteAddr, "package.update", "apt", "FAILED")
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "package.update", "apt", "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Package index updated"})
}

func handleLogs(w http.ResponseWriter, r *http.Request) {
	unit := r.URL.Query().Get("unit")
	priority := r.URL.Query().Get("priority")
	linesStr := r.URL.Query().Get("lines")
	lines := 50
	if linesStr != "" {
		if val, err := strconv.Atoi(linesStr); err == nil {
			lines = val
		}
	}

	opts := logs.LogOptions{
		Unit:     unit,
		Lines:    lines,
		Priority: priority,
	}

	entries, err := logs.ReadLogs(opts)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, entries)
}

func handleFirewall(w http.ResponseWriter, r *http.Request) {
	st, err := firewall.GetStatus()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, st)
}

func handleFirewallAllow(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Rule string `json:"rule"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Rule == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload, required 'rule' field (e.g. 80/tcp)"})
		return
	}
	if err := firewall.AllowPort(payload.Rule); err != nil {
		audit.Record("admin", r.RemoteAddr, "firewall.allow", payload.Rule, "FAILED")
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "firewall.allow", payload.Rule, "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Firewall rule added: allow %s", payload.Rule)})
}

func handleFirewallDeny(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Rule string `json:"rule"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Rule == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload, required 'rule' field (e.g. 23/tcp)"})
		return
	}
	if err := firewall.DenyPort(payload.Rule); err != nil {
		audit.Record("admin", r.RemoteAddr, "firewall.deny", payload.Rule, "FAILED")
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "firewall.deny", payload.Rule, "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Firewall rule added: deny %s", payload.Rule)})
}

func handleContainerList(w http.ResponseWriter, r *http.Request) {
	list, err := containers.ListContainers()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, list)
}

func handleContainerStart(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Container string `json:"container"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Container == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload, required 'container' field"})
		return
	}
	if err := containers.StartContainer(payload.Container); err != nil {
		audit.Record("admin", r.RemoteAddr, "container.start", payload.Container, "FAILED")
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "container.start", payload.Container, "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Container %s started", payload.Container)})
}

func handleContainerStop(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Container string `json:"container"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Container == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload, required 'container' field"})
		return
	}
	if err := containers.StopContainer(payload.Container); err != nil {
		audit.Record("admin", r.RemoteAddr, "container.stop", payload.Container, "FAILED")
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "container.stop", payload.Container, "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Container %s stopped", payload.Container)})
}

func handleContainerRestart(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Container string `json:"container"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Container == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload, required 'container' field"})
		return
	}
	if err := containers.RestartContainer(payload.Container); err != nil {
		audit.Record("admin", r.RemoteAddr, "container.restart", payload.Container, "FAILED")
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "container.restart", payload.Container, "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Container %s restarted", payload.Container)})
}

func handleContainerRemove(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Container string `json:"container"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Container == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload, required 'container' field"})
		return
	}
	if err := containers.RemoveContainer(payload.Container); err != nil {
		audit.Record("admin", r.RemoteAddr, "container.remove", payload.Container, "FAILED")
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "container.remove", payload.Container, "SUCCESS")
	jsonResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Container %s removed", payload.Container)})
}

func handleContainerLogs(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "query parameter 'id' required"})
		return
	}
	lines, err := containers.GetContainerLogs(id, 50)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, lines)
}

func handleContainerImages(w http.ResponseWriter, r *http.Request) {
	imgs, err := containers.ListImages()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, imgs)
}

func handleContainerVolumes(w http.ResponseWriter, r *http.Request) {
	vols, err := containers.ListVolumes()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, vols)
}

func handleContainerNetworks(w http.ResponseWriter, r *http.Request) {
	nets, err := containers.ListNetworks()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, nets)
}

func handleSecurityStatus(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, security.Audit())
}

func handleSecurityAudit(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, security.RunAudit())
}

func handleSecurityUpdates(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, security.GetSecurityUpdates())
}

func handleProfilesList(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, profiles.GetProfiles())
}

func handleProfileApply(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Profile string `json:"profile"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Profile == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload, required 'profile' field"})
		return
	}
	res, err := profiles.ApplyProfile(payload.Profile)
	if err != nil {
		audit.Record("admin", r.RemoteAddr, "profile.apply", payload.Profile, "FAILED")
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	audit.Record("admin", r.RemoteAddr, "profile.apply", payload.Profile, "SUCCESS")
	jsonResponse(w, http.StatusOK, res)
}
