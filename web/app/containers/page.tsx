'use client';

import { useState, useEffect } from 'react';
import {
  getContainers,
  startContainer,
  stopContainer,
  restartContainer,
  removeContainer,
  getContainerLogs,
  getContainerImages,
  getContainerVolumes,
  getContainerNetworks
} from '@/lib/api';
import { ContainerInfo, ContainerImage, ContainerVolume, ContainerNetwork } from '@/lib/types';
import Badge from '@/components/Badge';
import { Box, Play, Square, RotateCw, Trash2, FileText, RefreshCw, Layers, HardDrive, Network, X } from 'lucide-react';

export default function ContainersPage() {
  const [activeTab, setActiveTab] = useState<'containers' | 'images' | 'volumes' | 'networks'>('containers');
  const [containers, setContainers] = useState<ContainerInfo[]>([]);
  const [images, setImages] = useState<ContainerImage[]>([]);
  const [volumes, setVolumes] = useState<ContainerVolume[]>([]);
  const [networks, setNetworks] = useState<ContainerNetwork[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedLogsContainer, setSelectedLogsContainer] = useState<string | null>(null);
  const [containerLogs, setContainerLogs] = useState<string[]>([]);

  const fetchContainerData = async () => {
    setLoading(true);
    try {
      if (activeTab === 'containers') setContainers(await getContainers());
      if (activeTab === 'images') setImages(await getContainerImages());
      if (activeTab === 'volumes') setVolumes(await getContainerVolumes());
      if (activeTab === 'networks') setNetworks(await getContainerNetworks());
    } catch (e) {
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchContainerData();
  }, [activeTab]);

  const handleAction = async (action: 'start' | 'stop' | 'restart' | 'remove', id: string) => {
    try {
      if (action === 'start') await startContainer(id);
      if (action === 'stop') await stopContainer(id);
      if (action === 'restart') await restartContainer(id);
      if (action === 'remove') await removeContainer(id);
      await fetchContainerData();
    } catch (e) {
      alert(`Container action failed: ${e}`);
    }
  };

  const openLogsModal = async (id: string) => {
    setSelectedLogsContainer(id);
    try {
      const logs = await getContainerLogs(id);
      setContainerLogs(logs);
    } catch (e) {
      setContainerLogs([]);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-white tracking-tight">Docker / OCI Container Workloads</h1>
          <p className="text-sm text-slate-400 mt-1">Manage containerized workloads, images, storage volumes, and virtual networks</p>
        </div>
        <button
          onClick={fetchContainerData}
          disabled={loading}
          className="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold flex items-center gap-2 border border-border transition-all"
        >
          <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
          <span>Refresh</span>
        </button>
      </div>

      {/* Tabs Header */}
      <div className="flex gap-2 border-b border-border/60 pb-3">
        {[
          { id: 'containers', label: 'Containers', icon: Box },
          { id: 'images', label: 'Images', icon: Layers },
          { id: 'volumes', label: 'Volumes', icon: HardDrive },
          { id: 'networks', label: 'Networks', icon: Network },
        ].map((tab) => {
          const Icon = tab.icon;
          const isActive = activeTab === tab.id;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as any)}
              className={`px-4 py-2 rounded-xl text-xs font-semibold flex items-center gap-2 transition-all ${
                isActive
                  ? 'bg-blue-600/20 text-blue-400 border border-blue-500/30 shadow-md shadow-blue-500/10'
                  : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/40'
              }`}
            >
              <Icon className="w-4 h-4" />
              <span>{tab.label}</span>
            </button>
          );
        })}
      </div>

      {/* Tab 1: Containers */}
      {activeTab === 'containers' && (
        <div className="glass-panel rounded-2xl border border-border overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs text-slate-300">
              <thead className="bg-slate-900/80 border-b border-border text-slate-400 uppercase text-[11px] font-semibold tracking-wider">
                <tr>
                  <th className="px-6 py-4">Container</th>
                  <th className="px-6 py-4">Image</th>
                  <th className="px-6 py-4">Status</th>
                  <th className="px-6 py-4">Ports</th>
                  <th className="px-6 py-4 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/60">
                {containers.map((c) => (
                  <tr key={c.id} className="hover:bg-slate-800/40 transition-colors">
                    <td className="px-6 py-4">
                      <div className="font-semibold text-white flex items-center gap-2">
                        <Box className="w-4 h-4 text-blue-400 shrink-0" />
                        <span>{c.names}</span>
                      </div>
                      <span className="text-[11px] font-mono text-slate-500">{c.id.substring(0, 12)}</span>
                    </td>
                    <td className="px-6 py-4 font-mono text-indigo-400">{c.image}</td>
                    <td className="px-6 py-4">
                      <Badge status={c.state === 'running' ? 'HEALTHY' : 'stopped'} text={c.status} />
                    </td>
                    <td className="px-6 py-4 font-mono text-slate-400">{c.ports || 'None'}</td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end gap-1.5">
                        <button
                          onClick={() => handleAction('start', c.id)}
                          className="p-1.5 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 border border-emerald-500/20"
                          title="Start Container"
                        >
                          <Play className="w-3.5 h-3.5" />
                        </button>
                        <button
                          onClick={() => handleAction('stop', c.id)}
                          className="p-1.5 rounded-lg bg-amber-500/10 hover:bg-amber-500/20 text-amber-400 border border-amber-500/20"
                          title="Stop Container"
                        >
                          <Square className="w-3.5 h-3.5" />
                        </button>
                        <button
                          onClick={() => handleAction('restart', c.id)}
                          className="p-1.5 rounded-lg bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 border border-blue-500/20"
                          title="Restart Container"
                        >
                          <RotateCw className="w-3.5 h-3.5" />
                        </button>
                        <button
                          onClick={() => openLogsModal(c.id)}
                          className="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 border border-border"
                          title="Logs"
                        >
                          <FileText className="w-3.5 h-3.5" />
                        </button>
                        <button
                          onClick={() => handleAction('remove', c.id)}
                          className="p-1.5 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 border border-rose-500/20"
                          title="Remove Container"
                        >
                          <Trash2 className="w-3.5 h-3.5" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Tab 2: Images */}
      {activeTab === 'images' && (
        <div className="glass-panel rounded-2xl border border-border overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs text-slate-300">
              <thead className="bg-slate-900/80 border-b border-border text-slate-400 uppercase text-[11px] font-semibold tracking-wider">
                <tr>
                  <th className="px-6 py-4">Repository</th>
                  <th className="px-6 py-4">Tag</th>
                  <th className="px-6 py-4">Image ID</th>
                  <th className="px-6 py-4">Size</th>
                  <th className="px-6 py-4">Created</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/60">
                {images.map((img) => (
                  <tr key={img.id} className="hover:bg-slate-800/40 transition-colors">
                    <td className="px-6 py-4 font-semibold text-white">{img.repository}</td>
                    <td className="px-6 py-4 font-mono text-emerald-400">{img.tag}</td>
                    <td className="px-6 py-4 font-mono text-slate-400">{img.id.substring(0, 12)}</td>
                    <td className="px-6 py-4 text-slate-300 font-semibold">{img.size}</td>
                    <td className="px-6 py-4 text-slate-400">{img.created}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Tab 3: Volumes */}
      {activeTab === 'volumes' && (
        <div className="glass-panel rounded-2xl border border-border overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs text-slate-300">
              <thead className="bg-slate-900/80 border-b border-border text-slate-400 uppercase text-[11px] font-semibold tracking-wider">
                <tr>
                  <th className="px-6 py-4">Volume Name</th>
                  <th className="px-6 py-4">Driver</th>
                  <th className="px-6 py-4">Mount Point</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/60">
                {volumes.map((vol) => (
                  <tr key={vol.name} className="hover:bg-slate-800/40 transition-colors">
                    <td className="px-6 py-4 font-semibold text-white">{vol.name}</td>
                    <td className="px-6 py-4 font-mono text-blue-400">{vol.driver}</td>
                    <td className="px-6 py-4 font-mono text-slate-400">{vol.mount_point}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Tab 4: Networks */}
      {activeTab === 'networks' && (
        <div className="glass-panel rounded-2xl border border-border overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs text-slate-300">
              <thead className="bg-slate-900/80 border-b border-border text-slate-400 uppercase text-[11px] font-semibold tracking-wider">
                <tr>
                  <th className="px-6 py-4">Network Name</th>
                  <th className="px-6 py-4">Driver</th>
                  <th className="px-6 py-4">Scope</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/60">
                {networks.map((net) => (
                  <tr key={net.id} className="hover:bg-slate-800/40 transition-colors">
                    <td className="px-6 py-4 font-semibold text-white">{net.name}</td>
                    <td className="px-6 py-4 font-mono text-emerald-400">{net.driver}</td>
                    <td className="px-6 py-4 font-mono text-slate-400">{net.scope}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Logs Modal */}
      {selectedLogsContainer && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="glass-panel w-full max-w-3xl rounded-2xl border border-border p-6 space-y-4 shadow-2xl">
            <div className="flex justify-between items-center pb-3 border-b border-border">
              <div className="flex items-center gap-2">
                <FileText className="w-5 h-5 text-blue-400" />
                <h3 className="text-base font-bold text-white">Container Logs: {selectedLogsContainer}</h3>
              </div>
              <button
                onClick={() => setSelectedLogsContainer(null)}
                className="p-1 rounded-lg hover:bg-slate-800 text-slate-400 hover:text-white"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="bg-slate-950 p-4 rounded-xl font-mono text-xs text-slate-300 max-h-96 overflow-y-auto space-y-1 border border-border/80">
              {containerLogs.map((l, idx) => (
                <div key={idx}>{l}</div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
