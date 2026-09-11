export interface SystemInfo {
  os_name: string;
  os_version: string;
  kernel_version: string;
  hostname: string;
  architecture: string;
  go_version: string;
  num_cpu: number;
  uptime_seconds: number;
  uptime_formatted: string;
  timestamp: string;
}

export interface ResourceStats {
  cpu_usage_percent: number;
  memory_total_mb: number;
  memory_used_mb: number;
  memory_free_mb: number;
  memory_usage_percent: number;
  disk_total_gb: number;
  disk_used_gb: number;
  disk_free_gb: number;
  disk_usage_percent: number;
  load_avg_1: number;
  load_avg_5: number;
  load_avg_15: number;
  timestamp: string;
}

export interface SystemStatus {
  info: SystemInfo;
  resources: ResourceStats;
  daemon: {
    name: string;
    state: string;
    version: string;
  };
}

export interface HealthCheck {
  name: string;
  status: 'HEALTHY' | 'WARNING' | 'CRITICAL';
  message: string;
}

export interface HealthReport {
  overall: 'HEALTHY' | 'WARNING' | 'CRITICAL';
  checks: HealthCheck[];
  timestamp: string;
}

export interface ServiceInfo {
  name: string;
  status: 'running' | 'stopped';
  active_state: string;
  sub_state: string;
  load_state: string;
  description: string;
}

export interface PackageInfo {
  name: string;
  version: string;
  description: string;
  status: 'installed' | 'available';
}

export interface FirewallRule {
  id: string;
  port: number;
  protocol: 'tcp' | 'udp';
  action: 'allow' | 'deny';
  comment: string;
}

export interface FirewallStatus {
  enabled: boolean;
  backend: string;
  rules: FirewallRule[];
}

export interface ContainerInfo {
  id: string;
  names: string;
  image: string;
  status: string;
  state: 'running' | 'exited' | 'paused';
  ports: string;
  created: string;
  cpu_percent: number;
  memory_mb: number;
}

export interface ContainerImage {
  id: string;
  repository: string;
  tag: string;
  size: string;
  created: string;
}

export interface ContainerVolume {
  name: string;
  driver: string;
  scope: string;
  mount_point: string;
}

export interface ContainerNetwork {
  name: string;
  id: string;
  driver: string;
  scope: string;
}

export interface LogEntry {
  timestamp: string;
  host: string;
  process: string;
  pid?: number;
  message: string;
  severity: string;
}

export interface NetworkInterface {
  name: string;
  mac: string;
  ip_addresses: string[];
  flags: string[];
  is_up: boolean;
  is_loopback: boolean;
}

export interface NetworkRoute {
  destination: string;
  gateway: string;
  genmask: string;
  flags: string;
  interface: string;
}

export interface StorageMount {
  device: string;
  mount_point: string;
  fs_type: string;
  total_gb: number;
  used_gb: number;
  free_gb: number;
  usage_pct: number;
}

export interface DiskInfo {
  name: string;
  type: string;
  size_gb: number;
  model: string;
  filesystem: string;
  mount_point: string;
}

export interface AuditCheckItem {
  id: string;
  category: string;
  title: string;
  passed: boolean;
  severity: 'HIGH' | 'MEDIUM' | 'LOW';
  details: string;
}

export interface AuditReport {
  timestamp: string;
  passed: number;
  failed: number;
  score: number;
  check_items: AuditCheckItem[];
}

export interface SecurityUpdateInfo {
  pending_security_updates: number;
  packages: string[];
  last_checked: string;
}

export interface SecurityStatus {
  firewall_enabled: boolean;
  firewall_backend: string;
  root_ssh_disabled: boolean;
  ssh_password_auth: boolean;
  pending_sec_updates: number;
  warnings: string[];
}

export interface ServerProfile {
  name: string;
  title: string;
  description: string;
  packages: string[];
  services: string[];
  firewall_rules: string[];
}

export interface ProfileApplyResult {
  profile: string;
  success: boolean;
  message: string;
  details: string[];
}


