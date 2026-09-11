package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func handleWebOrStatic(w http.ResponseWriter, r *http.Request) {
	// If path starts with /api/, return 404 JSON for unknown API routes
	if strings.HasPrefix(r.URL.Path, "/api/") {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": "API route not found"})
		return
	}

	// Try serving static files from /usr/share/sawit/web or ./web/out
	staticDirs := []string{"/usr/share/sawit/web", "./web/out", "./web"}
	for _, dir := range staticDirs {
		filePath := filepath.Join(dir, filepath.Clean(r.URL.Path))
		if stat, err := os.Stat(filePath); err == nil && !stat.IsDir() {
			http.ServeFile(w, r, filePath)
			return
		}
	}

	// Serve embedded SawitOS Web UI Dashboard
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(embeddedWebUIHTML))
}

const embeddedWebUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SawitOS — Server Management Dashboard</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=Fira+Code:wght@400;500&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-dark: #0b0f19;
            --card-bg: #111827;
            --card-border: #1f2937;
            --accent: #10b981;
            --accent-hover: #059669;
            --text-main: #f9fafb;
            --text-muted: #9ca3af;
            --danger: #ef4444;
            --warning: #f59e0b;
        }

        * { box-sizing: border-box; margin: 0; padding: 0; font-family: 'Inter', sans-serif; }
        body { background-color: var(--bg-dark); color: var(--text-main); min-height: 100vh; display: flex; flex-direction: column; }

        /* Header */
        header { background: rgba(17, 24, 39, 0.8); backdrop-filter: blur(12px); border-bottom: 1px solid var(--card-border); padding: 1rem 2rem; display: flex; justify-content: space-between; align-items: center; position: sticky; top: 0; z-index: 100; }
        .logo { display: flex; align-items: center; gap: 0.75rem; font-weight: 700; font-size: 1.25rem; color: #fff; }
        .logo-icon { background: linear-gradient(135deg, #10b981, #047857); padding: 0.4rem 0.6rem; border-radius: 8px; font-size: 1.2rem; }
        .header-right { display: flex; align-items: center; gap: 1rem; }
        .badge { background: rgba(16, 185, 129, 0.1); color: var(--accent); border: 1px solid rgba(16, 185, 129, 0.2); padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.8rem; font-weight: 600; display: flex; align-items: center; gap: 0.4rem; }
        .dot { width: 8px; height: 8px; background: var(--accent); border-radius: 50%; display: inline-block; animation: pulse 2s infinite; }

        @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.4; } }

        /* Container */
        .container { max-width: 1280px; width: 100%; margin: 0 auto; padding: 2rem; flex: 1; }

        /* Login Modal */
        #login-section { max-width: 400px; margin: 4rem auto; background: var(--card-bg); border: 1px solid var(--card-border); border-radius: 12px; padding: 2.5rem; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5); }
        .form-group { margin-bottom: 1.25rem; }
        .form-group label { display: block; font-size: 0.85rem; color: var(--text-muted); margin-bottom: 0.5rem; font-weight: 500; }
        .form-control { width: 100%; background: #1f2937; border: 1px solid #374151; color: #fff; padding: 0.75rem 1rem; border-radius: 8px; font-size: 0.95rem; outline: none; transition: border 0.2s; }
        .form-control:focus { border-color: var(--accent); }
        .btn { width: 100%; background: var(--accent); color: #fff; border: none; padding: 0.75rem 1rem; border-radius: 8px; font-weight: 600; font-size: 1rem; cursor: pointer; transition: background 0.2s; }
        .btn:hover { background: var(--accent-hover); }

        /* Dashboard View */
        #dashboard-section { display: none; }
        .grid-4 { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 1.5rem; margin-bottom: 2rem; }
        .card { background: var(--card-bg); border: 1px solid var(--card-border); border-radius: 12px; padding: 1.5rem; }
        .card-title { font-size: 0.85rem; color: var(--text-muted); font-weight: 500; margin-bottom: 0.5rem; }
        .card-value { font-size: 1.8rem; font-weight: 700; color: #fff; }
        .progress-bg { background: #1f2937; height: 6px; border-radius: 9999px; overflow: hidden; margin-top: 1rem; }
        .progress-fill { background: var(--accent); height: 100%; width: 0%; transition: width 0.5s ease; }

        /* Sections */
        .section-title { font-size: 1.1rem; font-weight: 600; margin-bottom: 1rem; display: flex; justify-content: space-between; align-items: center; }
        table { width: 100%; border-collapse: collapse; text-align: left; }
        th, td { padding: 0.75rem 1rem; border-bottom: 1px solid var(--card-border); font-size: 0.9rem; }
        th { color: var(--text-muted); font-weight: 500; background: rgba(31, 41, 55, 0.4); }
        .status-running { color: var(--accent); font-weight: 600; }
        .status-stopped { color: var(--text-muted); }
        .btn-sm { padding: 0.35rem 0.75rem; font-size: 0.8rem; border-radius: 6px; cursor: pointer; border: 1px solid var(--card-border); background: #1f2937; color: #fff; }
        .btn-sm:hover { background: #374151; }

        .hidden { display: none !important; }
        .error-msg { background: rgba(239, 68, 68, 0.1); border: 1px solid rgba(239, 68, 68, 0.3); color: var(--danger); padding: 0.75rem; border-radius: 8px; font-size: 0.85rem; margin-bottom: 1rem; }
    </style>
</head>
<body>

    <header>
        <div class="logo">
            <span class="logo-icon">🌴</span>
            <span>SawitOS</span>
        </div>
        <div class="header-right">
            <div class="badge"><span class="dot"></span> <span id="header-status">Daemon Online</span></div>
            <button id="logout-btn" class="btn-sm hidden" onclick="logout()">Logout</button>
        </div>
    </header>

    <div class="container">
        <!-- LOGIN MODAL -->
        <div id="login-section">
            <h2 style="margin-bottom: 0.5rem; font-weight: 700;">Sign in to SawitOS</h2>
            <p style="color: var(--text-muted); font-size: 0.85rem; margin-bottom: 1.5rem;">Control panel authentication</p>
            <div id="login-error" class="error-msg hidden"></div>
            <form onsubmit="handleLogin(event)">
                <div class="form-group">
                    <label>Username</label>
                    <input type="text" id="username" class="form-control" value="admin" required>
                </div>
                <div class="form-group">
                    <label>Password</label>
                    <input type="password" id="password" class="form-control" value="admin123" required>
                </div>
                <button type="submit" class="btn">Login to Dashboard</button>
            </form>
        </div>

        <!-- DASHBOARD VIEW -->
        <div id="dashboard-section">
            <div class="grid-4">
                <div class="card">
                    <div class="card-title">CPU Usage</div>
                    <div class="card-value" id="val-cpu">0.0%</div>
                    <div class="progress-bg"><div class="progress-fill" id="bar-cpu"></div></div>
                </div>
                <div class="card">
                    <div class="card-title">Memory Usage</div>
                    <div class="card-value" id="val-ram">0 / 0 MB</div>
                    <div class="progress-bg"><div class="progress-fill" id="bar-ram"></div></div>
                </div>
                <div class="card">
                    <div class="card-title">Disk Storage</div>
                    <div class="card-value" id="val-disk">0.0 / 0.0 GB</div>
                    <div class="progress-bg"><div class="progress-fill" id="bar-disk"></div></div>
                </div>
                <div class="card">
                    <div class="card-title">Uptime & Host</div>
                    <div class="card-value" id="val-uptime" style="font-size: 1.3rem;">-</div>
                    <div style="font-size: 0.8rem; color: var(--text-muted); margin-top: 0.5rem;" id="val-host">Raspberry Pi</div>
                </div>
            </div>

            <!-- SYSTEM METADATA & SERVICES -->
            <div class="card" style="margin-bottom: 2rem;">
                <div class="section-title">
                    <span>System Services Status</span>
                    <button class="btn-sm" onclick="fetchServices()">Refresh Services</button>
                </div>
                <table>
                    <thead>
                        <tr>
                            <th>Service Name</th>
                            <th>Status</th>
                            <th>Sub State</th>
                            <th>Description</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody id="services-tbody">
                        <tr><td colspan="5" style="color: var(--text-muted);">Loading system services...</td></tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>

    <script>
        var jwtToken = localStorage.getItem('sawit_token') || '';

        function checkAuth() {
            if (jwtToken) {
                document.getElementById('login-section').classList.add('hidden');
                document.getElementById('dashboard-section').classList.remove('hidden');
                document.getElementById('logout-btn').classList.remove('hidden');
                startPolling();
            } else {
                document.getElementById('login-section').classList.remove('hidden');
                document.getElementById('dashboard-section').classList.add('hidden');
                document.getElementById('logout-btn').classList.add('hidden');
            }
        }

        async function handleLogin(e) {
            e.preventDefault();
            var errDiv = document.getElementById('login-error');
            errDiv.classList.add('hidden');

            var user = document.getElementById('username').value;
            var pass = document.getElementById('password').value;

            try {
                var res = await fetch('/api/v1/auth/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username: user, password: pass })
                });

                var data = await res.json();
                if (!res.ok) throw new Error(data.error || 'Login failed');

                jwtToken = data.token;
                localStorage.setItem('sawit_token', jwtToken);
                checkAuth();
            } catch (err) {
                errDiv.textContent = err.message;
                errDiv.classList.remove('hidden');
            }
        }

        function logout() {
            jwtToken = '';
            localStorage.removeItem('sawit_token');
            checkAuth();
        }

        async function fetchSystemStats() {
            if (!jwtToken) return;
            try {
                var res = await fetch('/api/v1/system', {
                    headers: { 'Authorization': 'Bearer ' + jwtToken }
                });
                if (res.status === 401) { logout(); return; }
                var data = await res.json();

                var resData = data.resources || {};
                var infoData = data.info || {};

                // CPU
                var cpuPct = (resData.cpu_usage_percent || 0).toFixed(1);
                document.getElementById('val-cpu').textContent = cpuPct + '%';
                document.getElementById('bar-cpu').style.width = cpuPct + '%';

                // RAM
                var usedRam = resData.memory_used_mb || 0;
                var totalRam = resData.memory_total_mb || 0;
                var ramPct = (resData.memory_usage_percent || 0).toFixed(1);
                document.getElementById('val-ram').textContent = usedRam + ' / ' + totalRam + ' MB';
                document.getElementById('bar-ram').style.width = ramPct + '%';

                // Disk
                var usedDisk = (resData.disk_used_gb || 0).toFixed(1);
                var totalDisk = (resData.disk_total_gb || 0).toFixed(1);
                var diskPct = (resData.disk_usage_percent || 0).toFixed(1);
                document.getElementById('val-disk').textContent = usedDisk + ' / ' + totalDisk + ' GB';
                document.getElementById('bar-disk').style.width = diskPct + '%';

                // Uptime & Host
                document.getElementById('val-uptime').textContent = infoData.uptime_formatted || '-';
                document.getElementById('val-host').textContent = (infoData.hostname || 'Raspberry Pi') + ' (' + (infoData.architecture || 'arm64') + ' / ' + (infoData.os_name || 'linux') + ')';

            } catch (err) {
                console.error('Failed to fetch system stats:', err);
            }
        }

        async function fetchServices() {
            if (!jwtToken) return;
            try {
                var res = await fetch('/api/v1/services', {
                    headers: { 'Authorization': 'Bearer ' + jwtToken }
                });
                if (!res.ok) return;
                var list = await res.json();

                var tbody = document.getElementById('services-tbody');
                tbody.innerHTML = '';

                list.forEach(function(svc) {
                    var tr = document.createElement('tr');
                    var isRunning = svc.status === 'running';
                    var actionName = isRunning ? 'stop' : 'start';
                    var actionLabel = isRunning ? 'Stop' : 'Start';
                    var statusClass = isRunning ? 'status-running' : 'status-stopped';
                    var desc = svc.description || '-';

                    tr.innerHTML = '<td style="font-family: \'Fira Code\', monospace; font-weight: 500;">' + svc.name + '</td>' +
                        '<td class="' + statusClass + '">' + svc.status + '</td>' +
                        '<td style="color: var(--text-muted);">' + svc.sub_state + '</td>' +
                        '<td>' + desc + '</td>' +
                        '<td><button class="btn-sm" onclick="toggleService(\'' + svc.name + '\', \'' + actionName + '\')">' + actionLabel + '</button></td>';
                    tbody.appendChild(tr);
                });
            } catch (err) {
                console.error('Failed to fetch services:', err);
            }
        }

        async function toggleService(name, action) {
            if (!jwtToken) return;
            try {
                await fetch('/api/v1/services/' + action, {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + jwtToken 
                    },
                    body: JSON.stringify({ service: name })
                });
                fetchServices();
            } catch (err) {
                alert('Failed to ' + action + ' service ' + name);
            }
        }

        var pollInterval = null;
        function startPolling() {
            fetchSystemStats();
            fetchServices();
            if (pollInterval) clearInterval(pollInterval);
            pollInterval = setInterval(fetchSystemStats, 3000);
        }

        // Initialize
        checkAuth();
    </script>
</body>
</html>
`
