'use client';

import { useState, useEffect } from 'react';
import { getProfiles, applyProfile } from '@/lib/api';
import { ServerProfile } from '@/lib/types';
import { Server, Save, Layers, CheckCircle, Shield } from 'lucide-react';

export default function SettingsPage() {
  const [host, setHost] = useState('0.0.0.0');
  const [port, setPort] = useState(8080);
  const [logLevel, setLogLevel] = useState('info');
  const [saved, setSaved] = useState(false);
  const [profiles, setProfiles] = useState<ServerProfile[]>([]);
  const [selectedProfile, setSelectedProfile] = useState<string>('minimal');
  const [applyStatus, setApplyStatus] = useState<string | null>(null);
  const [applying, setApplying] = useState(false);

  useEffect(() => {
    getProfiles().then(setProfiles).catch(() => {});
  }, []);

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    setSaved(true);
    setTimeout(() => setSaved(false), 3000);
  };

  const handleApplyProfile = async (profileName: string) => {
    setApplying(true);
    try {
      const res = await applyProfile(profileName);
      setSelectedProfile(profileName);
      setApplyStatus(res.message);
      setTimeout(() => setApplyStatus(null), 4000);
    } catch (err: any) {
      alert(`Error applying profile: ${err.message}`);
    } finally {
      setApplying(false);
    }
  };

  return (
    <div className="space-y-6 max-w-5xl">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-white tracking-tight">SawitOS System Settings & Profiles</h1>
        <p className="text-sm text-slate-400 mt-1">Configure sawitd daemon parameters, server profile presets, and administrative policy</p>
      </div>

      {saved && (
        <div className="p-3.5 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 text-xs flex items-center gap-2">
          <CheckCircle className="w-4 h-4" />
          <span>Daemon settings updated successfully.</span>
        </div>
      )}

      {applyStatus && (
        <div className="p-3.5 rounded-xl bg-blue-500/10 border border-blue-500/30 text-blue-400 text-xs flex items-center gap-2">
          <CheckCircle className="w-4 h-4" />
          <span>{applyStatus}</span>
        </div>
      )}

      {/* Server Profiles Section */}
      <div className="glass-panel p-6 rounded-2xl border border-border space-y-4">
        <div className="flex items-center gap-3 pb-3 border-b border-border">
          <div className="p-2 rounded-xl bg-emerald-500/10 text-emerald-400">
            <Layers className="w-5 h-5" />
          </div>
          <div>
            <h2 className="text-base font-bold text-white">SawitOS Server Profile Presets</h2>
            <p className="text-xs text-slate-400">Select target deployment profile to apply packages, services, and firewall defaults</p>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {profiles.map((p) => {
            const isSelected = selectedProfile === p.name;
            return (
              <div
                key={p.name}
                className={`p-4 rounded-xl border transition-all flex flex-col justify-between space-y-3 ${
                  isSelected
                    ? 'bg-blue-600/10 border-blue-500/60 shadow-lg shadow-blue-500/10'
                    : 'bg-slate-900/50 border-border hover:border-slate-700'
                }`}
              >
                <div>
                  <div className="flex items-center justify-between">
                    <h3 className="text-sm font-bold text-white">{p.title}</h3>
                    <span className="font-mono text-[10px] text-slate-400 uppercase bg-slate-800 px-2 py-0.5 rounded">{p.name}</span>
                  </div>
                  <p className="text-xs text-slate-400 mt-2 line-clamp-2">{p.description}</p>
                </div>

                <div className="pt-2 border-t border-border/40 text-[11px] text-slate-400 space-y-1">
                  <div><span className="text-slate-300 font-semibold">Ports:</span> {p.firewall_rules.join(', ')}</div>
                  <div><span className="text-slate-300 font-semibold">Services:</span> {p.services.slice(0, 3).join(', ')}</div>
                </div>

                <button
                  onClick={() => handleApplyProfile(p.name)}
                  disabled={applying}
                  className={`w-full py-2 rounded-lg text-xs font-semibold transition-colors ${
                    isSelected
                      ? 'bg-blue-600 hover:bg-blue-500 text-white'
                      : 'bg-slate-800 hover:bg-slate-700 text-slate-200'
                  }`}
                >
                  {isSelected ? 'Active Profile' : 'Apply Profile'}
                </button>
              </div>
            );
          })}
        </div>
      </div>

      {/* Settings Form */}
      <form onSubmit={handleSave} className="glass-panel p-6 rounded-2xl border border-border space-y-6">
        <div className="flex items-center gap-3 pb-4 border-b border-border">
          <div className="p-2 rounded-xl bg-blue-500/10 text-blue-400">
            <Server className="w-5 h-5" />
          </div>
          <div>
            <h2 className="text-base font-bold text-white">sawitd Daemon Configuration</h2>
            <p className="text-xs text-slate-400">/etc/sawit/sawitd.yaml</p>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-xs font-semibold text-slate-300 uppercase mb-2">Bind Address</label>
            <input
              type="text"
              value={host}
              onChange={(e) => setHost(e.target.value)}
              className="w-full px-3.5 py-2.5 rounded-xl bg-slate-900/80 border border-border text-xs text-white focus:outline-none focus:border-blue-500"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-slate-300 uppercase mb-2">API Server Port</label>
            <input
              type="number"
              value={port}
              onChange={(e) => setPort(Number(e.target.value))}
              className="w-full px-3.5 py-2.5 rounded-xl bg-slate-900/80 border border-border text-xs text-white focus:outline-none focus:border-blue-500"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-slate-300 uppercase mb-2">Logging Verbosity</label>
            <select
              value={logLevel}
              onChange={(e) => setLogLevel(e.target.value)}
              className="w-full px-3.5 py-2.5 rounded-xl bg-slate-900/80 border border-border text-xs text-white focus:outline-none"
            >
              <option value="debug">DEBUG</option>
              <option value="info">INFO</option>
              <option value="warn">WARN</option>
              <option value="error">ERROR</option>
            </select>
          </div>
        </div>

        <div className="pt-4 border-t border-border flex justify-end">
          <button
            type="submit"
            className="px-5 py-2.5 rounded-xl bg-blue-600 hover:bg-blue-500 text-white text-xs font-semibold flex items-center gap-2 shadow-lg shadow-blue-500/20"
          >
            <Save className="w-4 h-4" />
            <span>Save Configuration</span>
          </button>
        </div>
      </form>
    </div>
  );
}

