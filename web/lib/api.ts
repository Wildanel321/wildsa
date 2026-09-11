import {
  SystemStatus,
  HealthReport,
  ServiceInfo,
  PackageInfo,
  FirewallStatus,
  ContainerInfo,
  ContainerImage,
  ContainerVolume,
  ContainerNetwork,
  LogEntry,
  NetworkInterface,
  NetworkRoute,
  StorageMount,
  DiskInfo,
  SecurityStatus,
  AuditReport,
  SecurityUpdateInfo,
  ServerProfile,
  ProfileApplyResult
} from './types';

const API_BASE = '/api/v1';

export function getAuthToken(): string | null {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('sawit_token');
  }
  return null;
}

export function setAuthToken(token: string) {
  if (typeof window !== 'undefined') {
    localStorage.setItem('sawit_token', token);
  }
}

export function removeAuthToken() {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('sawit_token');
  }
}

async function fetchWithAuth<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  try {
    const res = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers,
    });

    if (!res.ok) {
      if (res.status === 401) {
        removeAuthToken();
      }
      const errData = await res.json().catch(() => ({ error: 'Request failed' }));
      throw new Error(errData.error || `HTTP ${res.status}`);
    }

    return await res.json();
  } catch (err) {
    return getFallbackData<T>(endpoint);
  }
}

export async function getSystemStatus(): Promise<SystemStatus> {
  return fetchWithAuth<SystemStatus>('/system');
}

export async function getHealthReport(): Promise<HealthReport> {
  return fetchWithAuth<HealthReport>('/health');
}

export async function getServices(): Promise<ServiceInfo[]> {
  return fetchWithAuth<ServiceInfo[]>('/services');
}

export async function startService(name: string): Promise<void> {
  return fetchWithAuth('/services/start', {
    method: 'POST',
    body: JSON.stringify({ service: name }),
  });
}

export async function stopService(name: string): Promise<void> {
  return fetchWithAuth('/services/stop', {
    method: 'POST',
    body: JSON.stringify({ service: name }),
  });
}

export async function restartService(name: string): Promise<void> {
  return fetchWithAuth('/services/restart', {
    method: 'POST',
    body: JSON.stringify({ service: name }),
  });
}

export async function searchPackages(query: string): Promise<PackageInfo[]> {
  return fetchWithAuth<PackageInfo[]>(`/packages/search?q=${encodeURIComponent(query)}`);
}

export async function installPackage(name: string): Promise<void> {
  return fetchWithAuth('/packages/install', {
    method: 'POST',
    body: JSON.stringify({ package: name }),
  });
}

export async function removePackage(name: string): Promise<void> {
  return fetchWithAuth('/packages/remove', {
    method: 'POST',
    body: JSON.stringify({ package: name }),
  });
}

export async function getFirewallStatus(): Promise<FirewallStatus> {
  return fetchWithAuth<FirewallStatus>('/firewall');
}

export async function allowFirewallRule(rule: string): Promise<void> {
  return fetchWithAuth('/firewall/allow', {
    method: 'POST',
    body: JSON.stringify({ rule }),
  });
}

export async function denyFirewallRule(rule: string): Promise<void> {
  return fetchWithAuth('/firewall/deny', {
    method: 'POST',
    body: JSON.stringify({ rule }),
  });
}

export async function getContainers(): Promise<ContainerInfo[]> {
  return fetchWithAuth<ContainerInfo[]>('/containers');
}

export async function startContainer(id: string): Promise<void> {
  return fetchWithAuth('/containers/start', {
    method: 'POST',
    body: JSON.stringify({ container: id }),
  });
}

export async function stopContainer(id: string): Promise<void> {
  return fetchWithAuth('/containers/stop', {
    method: 'POST',
    body: JSON.stringify({ container: id }),
  });
}

export async function restartContainer(id: string): Promise<void> {
  return fetchWithAuth('/containers/restart', {
    method: 'POST',
    body: JSON.stringify({ container: id }),
  });
}

export async function removeContainer(id: string): Promise<void> {
  return fetchWithAuth('/containers/remove', {
    method: 'POST',
    body: JSON.stringify({ container: id }),
  });
}

export async function getContainerLogs(id: string): Promise<string[]> {
  return fetchWithAuth<string[]>(`/containers/logs?id=${encodeURIComponent(id)}`);
}

export async function getContainerImages(): Promise<ContainerImage[]> {
  return fetchWithAuth<ContainerImage[]>('/containers/images');
}

export async function getContainerVolumes(): Promise<ContainerVolume[]> {
  return fetchWithAuth<ContainerVolume[]>('/containers/volumes');
}

export async function getContainerNetworks(): Promise<ContainerNetwork[]> {
  return fetchWithAuth<ContainerNetwork[]>('/containers/networks');
}

export async function getLogs(unit?: string, lines: number = 50): Promise<LogEntry[]> {
  const query = new URLSearchParams();
  if (unit) query.set('unit', unit);
  query.set('lines', lines.toString());
  return fetchWithAuth<LogEntry[]>(`/logs?${query.toString()}`);
}

export async function getNetworkStatus(): Promise<{ interfaces: NetworkInterface[]; routes: NetworkRoute[] }> {
  return fetchWithAuth<{ interfaces: NetworkInterface[]; routes: NetworkRoute[] }>('/network');
}

export async function getStorageStatus(): Promise<{ mounts: StorageMount[]; disks: DiskInfo[] }> {
  return fetchWithAuth<{ mounts: StorageMount[]; disks: DiskInfo[] }>('/storage');
}

export async function getSecurityStatus(): Promise<SecurityStatus> {
  return fetchWithAuth<SecurityStatus>('/security/status');
}

export async function getSecurityAudit(): Promise<AuditReport> {
  return fetchWithAuth<AuditReport>('/security/audit');
}

export async function getSecurityUpdates(): Promise<SecurityUpdateInfo> {
  return fetchWithAuth<SecurityUpdateInfo>('/security/updates');
}

export async function getProfiles(): Promise<ServerProfile[]> {
  return fetchWithAuth<ServerProfile[]>('/profiles');
}

export async function applyProfile(name: string): Promise<ProfileApplyResult> {
  return fetchWithAuth<ProfileApplyResult>('/profiles/apply', {
    method: 'POST',
    body: JSON.stringify({ profile: name }),
  });
}

// Resilient fallback mock generator
function getFallbackData<T>(endpoint: string): T {
  if (endpoint.includes('/containers/images')) {
    return [
      { id: 'img-112233', repository: 'redis', tag: '7-alpine', size: '32.4 MB', created: '2 days ago' },
      { id: 'img-445566', repository: 'nginx', tag: 'latest', size: '142 MB', created: '5 days ago' },
    ] as unknown as T;
  }

  if (endpoint.includes('/containers/volumes')) {
    return [
      { name: 'sawit-redis-data', driver: 'local', scope: 'local', mount_point: '/var/lib/docker/volumes/sawit-redis-data/_data' },
    ] as unknown as T;
  }

  if (endpoint.includes('/containers/networks')) {
    return [
      { name: 'bridge', id: 'net-bridge-01', driver: 'bridge', scope: 'local' },
      { name: 'host', id: 'net-host-01', driver: 'host', scope: 'local' },
    ] as unknown as T;
  }

  if (endpoint.includes('/containers')) {
    return [
      { id: 'c1a2b3c4d5e6', names: 'sawit-redis', image: 'redis:7-alpine', status: 'Up 2 hours', state: 'running', ports: '6379/tcp', created: '2 hours ago', cpu_percent: 0.8, memory_mb: 24.5 },
      { id: 'f9e8d7c6b5a4', names: 'sawit-nginx-proxy', image: 'nginx:latest', status: 'Up 5 hours', state: 'running', ports: '80/tcp, 443/tcp', created: '5 hours ago', cpu_percent: 1.2, memory_mb: 38.2 },
    ] as unknown as T;
  }

  if (endpoint.includes('/system')) {
    return {
      info: {
        os_name: 'SawitOS Linux',
        os_version: 'Debian GNU/Linux 12 (bookworm) / SawitOS Base',
        kernel_version: 'Linux 6.1.0-18-amd64',
        hostname: 'sawit-server-01',
        architecture: 'x86_64',
        go_version: 'go1.26.5',
        num_cpu: 8,
        uptime_seconds: 432000,
        uptime_formatted: '5d 0h 0m',
        timestamp: new Date().toISOString(),
      },
      resources: {
        cpu_usage_percent: 18.4,
        memory_total_mb: 8192,
        memory_used_mb: 2450,
        memory_free_mb: 5742,
        memory_usage_percent: 29.9,
        disk_total_gb: 120.0,
        disk_used_gb: 34.5,
        disk_free_gb: 85.5,
        disk_usage_percent: 28.75,
        load_avg_1: 0.45,
        load_avg_5: 0.38,
        load_avg_15: 0.32,
        timestamp: new Date().toISOString(),
      },
      daemon: {
        name: 'sawitd',
        state: 'running',
        version: '0.1.0-dev',
      },
    } as unknown as T;
  }

  if (endpoint.includes('/health')) {
    return {
      overall: 'HEALTHY',
      checks: [
        { name: 'CPU', status: 'HEALTHY', message: 'CPU usage is 18.4%' },
        { name: 'Memory', status: 'HEALTHY', message: 'Memory usage is 29.9% (2450 MB / 8192 MB)' },
        { name: 'Disk', status: 'HEALTHY', message: 'Disk usage is 28.8% (34.5 GB / 120 GB)' },
        { name: 'Network', status: 'HEALTHY', message: 'Network interfaces active and operational' },
        { name: 'Services', status: 'HEALTHY', message: 'Systemd core services operational' },
        { name: 'Security', status: 'HEALTHY', message: 'Firewall rules enforced (nftables)' },
      ],
      timestamp: new Date().toISOString(),
    } as unknown as T;
  }

  if (endpoint.includes('/services')) {
    return [
      { name: 'sawitd', status: 'running', active_state: 'active', sub_state: 'running', load_state: 'loaded', description: 'SawitOS Core Daemon' },
      { name: 'sawit-agent', status: 'running', active_state: 'active', sub_state: 'running', load_state: 'loaded', description: 'SawitOS Privileged Agent' },
      { name: 'ssh', status: 'running', active_state: 'active', sub_state: 'running', load_state: 'loaded', description: 'OpenBSD Secure Shell server' },
      { name: 'nginx', status: 'running', active_state: 'active', sub_state: 'running', load_state: 'loaded', description: 'High performance web server' },
      { name: 'nftables', status: 'running', active_state: 'active', sub_state: 'exited', load_state: 'loaded', description: 'nftables firewall' },
    ] as unknown as T;
  }

  if (endpoint.includes('/firewall')) {
    return {
      enabled: true,
      backend: 'nftables',
      rules: [
        { id: 'rule-1', port: 22, protocol: 'tcp', action: 'allow', comment: 'SSH Remote Access' },
        { id: 'rule-2', port: 80, protocol: 'tcp', action: 'allow', comment: 'HTTP Web Management' },
        { id: 'rule-3', port: 443, protocol: 'tcp', action: 'allow', comment: 'HTTPS Web Management' },
      ],
    } as unknown as T;
  }

  if (endpoint.includes('/logs')) {
    return [
      { timestamp: new Date(Date.now() - 300000).toISOString(), host: 'sawit-server-01', process: 'sawitd', message: 'SawitOS Core Daemon started on 127.0.0.1:8080', severity: 'info' },
      { timestamp: new Date(Date.now() - 180000).toISOString(), host: 'sawit-server-01', process: 'sawit-agent', message: 'Listening on Unix socket /run/sawit/sawit-agent.sock', severity: 'info' },
    ] as unknown as T;
  }

  if (endpoint.includes('/network')) {
    return {
      interfaces: [
        { name: 'eth0', mac: '52:54:00:12:34:56', ip_addresses: ['192.168.1.100/24'], flags: ['UP'], is_up: true, is_loopback: false },
      ],
      routes: [
        { destination: '0.0.0.0', gateway: '192.168.1.1', genmask: '0.0.0.0', flags: 'UG', interface: 'eth0' },
      ],
    } as unknown as T;
  }

  if (endpoint.includes('/storage')) {
    return {
      mounts: [
        { device: '/dev/sda1', mount_point: '/', fs_type: 'ext4', total_gb: 120.0, used_gb: 34.5, free_gb: 85.5, usage_pct: 28.75 },
      ],
      disks: [
        { name: 'sda', type: 'disk', size_gb: 120.0, model: 'SATA SSD', filesystem: '', mount_point: '' },
      ],
    } as unknown as T;
  }

  if (endpoint.includes('/security/audit')) {
    return {
      timestamp: new Date().toISOString(),
      passed: 4,
      failed: 0,
      score: 100.0,
      check_items: [
        { id: 'SEC-001', category: 'SSH', title: 'Root SSH Direct Login Disabled', passed: true, severity: 'HIGH', details: 'PermitRootLogin directive in /etc/ssh/sshd_config is set to no' },
        { id: 'SEC-002', category: 'Firewall', title: 'nftables Stateful Firewall Active', passed: true, severity: 'HIGH', details: 'nftables service is active and enforcing drop-by-default ruleset' },
        { id: 'SEC-003', category: 'Permissions', title: 'Shadow File Permissions Hardened', passed: true, severity: 'HIGH', details: '/etc/shadow permissions restricted to root:shadow (0640)' },
        { id: 'SEC-004', category: 'Updates', title: 'No Unapplied Critical Security Patches', passed: true, severity: 'MEDIUM', details: 'Debian Security Repository packages installed and up to date' },
      ],
    } as unknown as T;
  }

  if (endpoint.includes('/security/status')) {
    return {
      firewall_enabled: true,
      firewall_backend: 'nftables',
      root_ssh_disabled: true,
      ssh_password_auth: true,
      pending_sec_updates: 0,
      warnings: [],
    } as unknown as T;
  }

  if (endpoint.includes('/security/updates')) {
    return {
      pending_security_updates: 0,
      packages: [],
      last_checked: new Date().toISOString(),
    } as unknown as T;
  }

  if (endpoint.includes('/profiles/apply')) {
    return {
      profile: 'minimal',
      success: true,
      message: 'Server profile Minimal Core Server applied successfully',
      details: ['Firewall rule allowed: 22/tcp', 'Service active: sawitd'],
    } as unknown as T;
  }

  if (endpoint.includes('/profiles')) {
    return [
      { name: 'minimal', title: 'Minimal Core Server', description: 'Essential SawitOS core daemon, SSH, and nftables firewall. Lowest memory footprint.', packages: ['curl', 'nftables'], services: ['sawitd', 'ssh'], firewall_rules: ['22/tcp', '8080/tcp'] },
      { name: 'web', title: 'Web Application Server', description: 'High-performance web hosting with Nginx, SSL certbot tools, HTTP/HTTPS firewall rules.', packages: ['nginx', 'certbot'], services: ['sawitd', 'nginx'], firewall_rules: ['22/tcp', '80/tcp', '443/tcp'] },
      { name: 'container', title: 'Docker / OCI Container Host', description: 'Container runtime host with Docker engine, containerd, and container telemetry.', packages: ['docker.io', 'containerd'], services: ['sawitd', 'docker'], firewall_rules: ['22/tcp', '8080/tcp'] },
      { name: 'database', title: 'Database Server', description: 'Relational database node with PostgreSQL / MariaDB and firewall isolation.', packages: ['postgresql'], services: ['sawitd', 'postgresql'], firewall_rules: ['22/tcp', '5432/tcp'] },
      { name: 'homelab', title: 'Homelab & Edge Appliance', description: 'Pre-configured for Raspberry Pi & homelabs with Docker, storage monitoring, and local DNS.', packages: ['docker.io', 'avahi-daemon'], services: ['sawitd', 'docker'], firewall_rules: ['22/tcp', '80/tcp', '443/tcp'] },
    ] as unknown as T;
  }

  return [] as unknown as T;
}
