'use client';

import { useState, useEffect } from 'react';
import { getServices, startService, stopService, restartService, getLogs } from '@/lib/api';
import { ServiceInfo, LogEntry } from '@/lib/types';
import Badge from '@/components/Badge';
import { Server, Play, Square, RotateCw, FileText, Search, RefreshCw, X } from 'lucide-react';

export default function ServicesPage() {
  const [services, setServices] = useState<ServiceInfo[]>([]);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(false);
  const [selectedLogsUnit, setSelectedLogsUnit] = useState<string | null>(null);
  const [unitLogs, setUnitLogs] = useState<LogEntry[]>([]);

  const fetchServicesData = async () => {
    setLoading(true);
    try {
      const data = await getServices();
      setServices(data);
    } catch (e) {
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchServicesData();
  }, []);

  const handleAction = async (action: 'start' | 'stop' | 'restart', name: string) => {
    try {
      if (action === 'start') await startService(name);
      if (action === 'stop') await stopService(name);
      if (action === 'restart') await restartService(name);
      await fetchServicesData();
    } catch (e) {
      alert(`Service action failed: ${e}`);
    }
  };

  const openLogsModal = async (unitName: string) => {
    setSelectedLogsUnit(unitName);
    try {
      const logs = await getLogs(unitName, 30);
      setUnitLogs(logs);
    } catch (e) {
      setUnitLogs([]);
    }
  };

  const filteredServices = services.filter(
    (s) =>
      s.name.toLowerCase().includes(search.toLowerCase()) ||
      s.description.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-white tracking-tight">Systemd Service Manager</h1>
          <p className="text-sm text-slate-400 mt-1">Control and inspect system services and daemons</p>
        </div>
        <button
          onClick={fetchServicesData}
          disabled={loading}
          className="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold flex items-center gap-2 border border-border transition-all"
        >
          <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
          <span>Refresh</span>
        </button>
      </div>

      {/* Search Bar */}
      <div className="relative max-w-md">
        <Search className="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
        <input
          type="text"
          placeholder="Filter services by name or description..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="w-full pl-10 pr-4 py-2.5 rounded-xl bg-slate-900/80 border border-border text-xs text-white placeholder:text-slate-500 focus:outline-none focus:border-blue-500 transition-all"
        />
      </div>

      {/* Service Table */}
      <div className="glass-panel rounded-2xl border border-border overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs text-slate-300">
            <thead className="bg-slate-900/80 border-b border-border text-slate-400 uppercase text-[11px] font-semibold tracking-wider">
              <tr>
                <th className="px-6 py-4">Service Name</th>
                <th className="px-6 py-4">Status</th>
                <th className="px-6 py-4">Load State</th>
                <th className="px-6 py-4">Description</th>
                <th className="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/60">
              {filteredServices.map((svc) => (
                <tr key={svc.name} className="hover:bg-slate-800/40 transition-colors">
                  <td className="px-6 py-4 font-semibold text-white flex items-center gap-2">
                    <Server className="w-4 h-4 text-blue-400 shrink-0" />
                    <span>{svc.name}</span>
                  </td>
                  <td className="px-6 py-4">
                    <Badge status={svc.status} />
                  </td>
                  <td className="px-6 py-4 font-mono text-slate-400">{svc.load_state}</td>
                  <td className="px-6 py-4 text-slate-400 max-w-xs truncate">{svc.description || 'System service'}</td>
                  <td className="px-6 py-4 text-right">
                    <div className="flex items-center justify-end gap-1.5">
                      <button
                        onClick={() => handleAction('start', svc.name)}
                        className="p-1.5 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 border border-emerald-500/20"
                        title="Start Service"
                      >
                        <Play className="w-3.5 h-3.5" />
                      </button>
                      <button
                        onClick={() => handleAction('stop', svc.name)}
                        className="p-1.5 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 border border-rose-500/20"
                        title="Stop Service"
                      >
                        <Square className="w-3.5 h-3.5" />
                      </button>
                      <button
                        onClick={() => handleAction('restart', svc.name)}
                        className="p-1.5 rounded-lg bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 border border-blue-500/20"
                        title="Restart Service"
                      >
                        <RotateCw className="w-3.5 h-3.5" />
                      </button>
                      <button
                        onClick={() => openLogsModal(svc.name)}
                        className="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 border border-border"
                        title="View Service Logs"
                      >
                        <FileText className="w-3.5 h-3.5" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Logs Modal */}
      {selectedLogsUnit && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="glass-panel w-full max-w-3xl rounded-2xl border border-border p-6 space-y-4 shadow-2xl">
            <div className="flex justify-between items-center pb-3 border-b border-border">
              <div className="flex items-center gap-2">
                <FileText className="w-5 h-5 text-blue-400" />
                <h3 className="text-base font-bold text-white">Systemd Logs: {selectedLogsUnit}</h3>
              </div>
              <button
                onClick={() => setSelectedLogsUnit(null)}
                className="p-1 rounded-lg hover:bg-slate-800 text-slate-400 hover:text-white"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="bg-slate-950 p-4 rounded-xl font-mono text-xs text-slate-300 max-h-96 overflow-y-auto space-y-1 border border-border/80">
              {unitLogs.length > 0 ? (
                unitLogs.map((log, idx) => (
                  <div key={idx} className="flex gap-2">
                    <span className="text-slate-500 shrink-0">{new Date(log.timestamp).toLocaleTimeString()}</span>
                    <span className="text-blue-400 font-semibold shrink-0">{log.process}:</span>
                    <span className="text-slate-300">{log.message}</span>
                  </div>
                ))
              ) : (
                <p className="text-slate-500">No logs found for unit {selectedLogsUnit}</p>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
