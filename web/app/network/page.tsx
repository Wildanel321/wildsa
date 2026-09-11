'use client';

import { useState, useEffect } from 'react';
import { getNetworkStatus } from '@/lib/api';
import { NetworkInterface, NetworkRoute } from '@/lib/types';
import Badge from '@/components/Badge';
import { Network, ArrowUpRight, Globe, Server } from 'lucide-react';

export default function NetworkPage() {
  const [interfaces, setInterfaces] = useState<NetworkInterface[]>([]);
  const [routes, setRoutes] = useState<NetworkRoute[]>([]);

  useEffect(() => {
    getNetworkStatus().then((data) => {
      setInterfaces(data.interfaces || []);
      setRoutes(data.routes || []);
    });
  }, []);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-white tracking-tight">Network Interfaces & Routing</h1>
        <p className="text-sm text-slate-400 mt-1">Inspect network hardware interfaces, IP addresses, and kernel routing table</p>
      </div>

      {/* Interfaces Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
        {interfaces.map((iface) => (
          <div key={iface.name} className="glass-panel p-5 rounded-2xl border border-border space-y-3">
            <div className="flex justify-between items-center pb-2 border-b border-border">
              <div className="flex items-center gap-2">
                <Network className="w-5 h-5 text-blue-400" />
                <h3 className="text-base font-bold text-white">{iface.name}</h3>
              </div>
              <Badge status={iface.is_up ? 'HEALTHY' : 'stopped'} text={iface.is_up ? 'UP' : 'DOWN'} />
            </div>

            <div className="space-y-1.5 text-xs">
              <div className="flex justify-between py-1 border-b border-border/40">
                <span className="text-slate-400">Hardware MAC</span>
                <span className="font-mono text-slate-200">{iface.mac || 'N/A'}</span>
              </div>
              <div className="flex justify-between py-1 border-b border-border/40">
                <span className="text-slate-400">IP Addresses</span>
                <span className="font-mono text-emerald-400">{iface.ip_addresses.join(', ') || 'None'}</span>
              </div>
              <div className="flex justify-between py-1">
                <span className="text-slate-400">Flags</span>
                <span className="font-mono text-slate-400">{iface.flags.join(', ')}</span>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Routing Table */}
      <div className="glass-panel rounded-2xl border border-border overflow-hidden">
        <div className="px-6 py-4 border-b border-border flex items-center gap-2">
          <Globe className="w-4 h-4 text-blue-400" />
          <h3 className="text-sm font-bold text-white">Kernel IPv4 Routing Table</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs text-slate-300">
            <thead className="bg-slate-900/80 border-b border-border text-slate-400 uppercase text-[11px] font-semibold tracking-wider">
              <tr>
                <th className="px-6 py-4">Destination</th>
                <th className="px-6 py-4">Gateway</th>
                <th className="px-6 py-4">Genmask</th>
                <th className="px-6 py-4">Flags</th>
                <th className="px-6 py-4">Interface</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/60">
              {routes.map((r, idx) => (
                <tr key={idx} className="hover:bg-slate-800/40 transition-colors font-mono">
                  <td className="px-6 py-4 font-semibold text-white">{r.destination}</td>
                  <td className="px-6 py-4 text-emerald-400">{r.gateway}</td>
                  <td className="px-6 py-4 text-slate-400">{r.genmask}</td>
                  <td className="px-6 py-4 text-slate-400">{r.flags}</td>
                  <td className="px-6 py-4 text-blue-400">{r.interface}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
