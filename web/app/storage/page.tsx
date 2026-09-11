'use client';

import { useState, useEffect } from 'react';
import { getStorageStatus } from '@/lib/api';
import { StorageMount, DiskInfo } from '@/lib/types';
import { HardDrive, Disc, Database } from 'lucide-react';

export default function StoragePage() {
  const [mounts, setMounts] = useState<StorageMount[]>([]);
  const [disks, setDisks] = useState<DiskInfo[]>([]);

  useEffect(() => {
    getStorageStatus().then((data) => {
      setMounts(data.mounts || []);
      setDisks(data.disks || []);
    });
  }, []);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-white tracking-tight">Storage & Partition Layout</h1>
        <p className="text-sm text-slate-400 mt-1">Filesystem mounts, disk utilization, and block storage devices</p>
      </div>

      {/* Storage Mount Cards */}
      <div className="space-y-4">
        <h2 className="text-base font-bold text-white flex items-center gap-2">
          <HardDrive className="w-5 h-5 text-indigo-400" />
          <span>Active Filesystem Mounts</span>
        </h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
          {mounts.map((m) => (
            <div key={m.mount_point} className="glass-panel p-5 rounded-2xl border border-border space-y-3">
              <div className="flex justify-between items-center pb-2 border-b border-border">
                <div className="font-semibold text-white flex items-center gap-2">
                  <Database className="w-4 h-4 text-blue-400" />
                  <span>{m.mount_point}</span>
                </div>
                <span className="text-xs font-mono text-slate-400 uppercase">{m.fs_type}</span>
              </div>

              <div>
                <div className="flex justify-between text-xs text-slate-400 mb-1">
                  <span>
                    Used: {m.used_gb.toFixed(1)} GB / {m.total_gb.toFixed(1)} GB
                  </span>
                  <span className="font-semibold text-slate-200">{m.usage_pct.toFixed(1)}%</span>
                </div>
                <div className="w-full h-2.5 rounded-full bg-slate-800 overflow-hidden">
                  <div
                    className="h-full rounded-full bg-indigo-500 transition-all duration-500"
                    style={{ width: `${m.usage_pct}%` }}
                  ></div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Disks Table */}
      <div className="glass-panel rounded-2xl border border-border overflow-hidden">
        <div className="px-6 py-4 border-b border-border flex items-center gap-2">
          <Disc className="w-4 h-4 text-blue-400" />
          <h3 className="text-sm font-bold text-white">Physical Disk Layout (lsblk)</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs text-slate-300">
            <thead className="bg-slate-900/80 border-b border-border text-slate-400 uppercase text-[11px] font-semibold tracking-wider">
              <tr>
                <th className="px-6 py-4">Name</th>
                <th className="px-6 py-4">Type</th>
                <th className="px-6 py-4">Size</th>
                <th className="px-6 py-4">Filesystem</th>
                <th className="px-6 py-4">Model</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/60">
              {disks.map((d) => (
                <tr key={d.name} className="hover:bg-slate-800/40 transition-colors">
                  <td className="px-6 py-4 font-mono font-semibold text-white">{d.name}</td>
                  <td className="px-6 py-4 font-mono uppercase text-blue-400">{d.type}</td>
                  <td className="px-6 py-4 font-semibold text-slate-200">{d.size_gb.toFixed(1)} GB</td>
                  <td className="px-6 py-4 text-slate-400">{d.filesystem || '-'}</td>
                  <td className="px-6 py-4 text-slate-400">{d.model || 'Block Device'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
