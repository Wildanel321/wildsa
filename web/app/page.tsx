'use client';

import { useState, useEffect } from 'react';
import { getSystemStatus, getHealthReport, getServices } from '@/lib/api';
import { getWebSocketInstance } from '@/lib/websocket';
import { SystemStatus, HealthReport, ServiceInfo, ResourceStats } from '@/lib/types';
import MetricCard from '@/components/MetricCard';
import Badge from '@/components/Badge';
import { Cpu, HardDrive, Server, ShieldCheck, Activity, MemoryStick, Layers } from 'lucide-react';

export default function DashboardPage() {
  const [status, setStatus] = useState<SystemStatus | null>(null);
  const [health, setHealth] = useState<HealthReport | null>(null);
  const [services, setServices] = useState<ServiceInfo[]>([]);

  useEffect(() => {
    // Initial fetch
    getSystemStatus().then(setStatus).catch(() => {});
    getHealthReport().then(setHealth).catch(() => {});
    getServices().then(setServices).catch(() => {});

    // WebSocket real-time subscription
    const ws = getWebSocketInstance();
    const unsubscribe = ws.subscribe('metrics', (stats: ResourceStats) => {
      setStatus((prev) => (prev ? { ...prev, resources: stats } : prev));
    });

    return () => unsubscribe();
  }, []);

  const resources = status?.resources;
  const runningServicesCount = services.filter((s) => s.status === 'running').length;

  return (
    <div className="space-y-8">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl font-bold text-white tracking-tight">System Dashboard</h1>
        <p className="text-sm text-slate-400 mt-1">Real-time telemetric metrics and operational health summary</p>
      </div>

      {/* Top Metric Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <MetricCard
          title="CPU Usage"
          value={`${resources?.cpu_usage_percent.toFixed(1) || '0.0'}%`}
          subtitle={`${status?.info.num_cpu || 0} Cores (${status?.info.architecture || 'x86_64'})`}
          percent={resources?.cpu_usage_percent || 0}
          icon={Cpu}
          color="blue"
        />

        <MetricCard
          title="Memory Usage"
          value={`${resources?.memory_usage_percent.toFixed(1) || '0.0'}%`}
          subtitle={`${resources?.memory_used_mb || 0} MB / ${resources?.memory_total_mb || 0} MB`}
          percent={resources?.memory_usage_percent || 0}
          icon={MemoryStick}
          color="emerald"
        />

        <MetricCard
          title="Root Storage"
          value={`${resources?.disk_usage_percent.toFixed(1) || '0.0'}%`}
          subtitle={`${resources?.disk_used_gb.toFixed(1) || '0.0'} GB / ${resources?.disk_total_gb.toFixed(1) || '0.0'} GB`}
          percent={resources?.disk_usage_percent || 0}
          icon={HardDrive}
          color="indigo"
        />

        <MetricCard
          title="Active Services"
          value={`${runningServicesCount} / ${services.length}`}
          subtitle="Systemd Service Units"
          icon={Server}
          color="amber"
        />
      </div>

      {/* Main Grid Section */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column: Health Audit */}
        <div className="lg:col-span-2 glass-panel p-6 rounded-2xl border border-border space-y-5">
          <div className="flex justify-between items-center pb-4 border-b border-border">
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-xl bg-blue-500/10 text-blue-400">
                <ShieldCheck className="w-5 h-5" />
              </div>
              <div>
                <h2 className="text-base font-bold text-white">System Health Audit</h2>
                <p className="text-xs text-slate-400">Automated diagnostic evaluation (sawit-health)</p>
              </div>
            </div>
            {health && <Badge status={health.overall} />}
          </div>

          <div className="space-y-3">
            {health?.checks.map((check) => (
              <div
                key={check.name}
                className="flex items-center justify-between p-3 rounded-xl bg-slate-900/50 border border-border/60"
              >
                <div className="flex items-center gap-3">
                  <Badge status={check.status} text={check.name} />
                  <span className="text-xs text-slate-300 font-medium">{check.message}</span>
                </div>
                <span className="text-[11px] text-slate-500 font-mono">OK</span>
              </div>
            ))}
          </div>
        </div>

        {/* Right Column: System Specs & Info */}
        <div className="glass-panel p-6 rounded-2xl border border-border space-y-5">
          <div className="flex items-center gap-3 pb-4 border-b border-border">
            <div className="p-2 rounded-xl bg-indigo-500/10 text-indigo-400">
              <Layers className="w-5 h-5" />
            </div>
            <div>
              <h2 className="text-base font-bold text-white">System Specifications</h2>
              <p className="text-xs text-slate-400">SawitOS Core Layer</p>
            </div>
          </div>

          <div className="space-y-3 text-xs">
            <div className="flex justify-between py-1.5 border-b border-border/40">
              <span className="text-slate-400">OS Version</span>
              <span className="font-semibold text-slate-200">{status?.info.os_version}</span>
            </div>
            <div className="flex justify-between py-1.5 border-b border-border/40">
              <span className="text-slate-400">Kernel Version</span>
              <span className="font-mono text-slate-200">{status?.info.kernel_version}</span>
            </div>
            <div className="flex justify-between py-1.5 border-b border-border/40">
              <span className="text-slate-400">Architecture</span>
              <span className="font-mono text-slate-200">{status?.info.architecture}</span>
            </div>
            <div className="flex justify-between py-1.5 border-b border-border/40">
              <span className="text-slate-400">Go Compiler</span>
              <span className="font-mono text-slate-200">{status?.info.go_version}</span>
            </div>
            <div className="flex justify-between py-1.5 border-b border-border/40">
              <span className="text-slate-400">Daemon Version</span>
              <span className="font-mono text-emerald-400">sawitd v{status?.daemon.version}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
