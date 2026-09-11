'use client';

import { useState } from 'react';
import { searchPackages, installPackage, removePackage } from '@/lib/api';
import { PackageInfo } from '@/lib/types';
import { Package, Search, Download, Trash2, CheckCircle, AlertCircle } from 'lucide-react';

export default function PackagesPage() {
  const [query, setQuery] = useState('nginx');
  const [packages, setPackages] = useState<PackageInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query) return;
    setLoading(true);
    setMessage('');
    try {
      const data = await searchPackages(query);
      setPackages(data);
    } catch (err: any) {
      setMessage(`Search failed: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleInstall = async (name: string) => {
    try {
      await installPackage(name);
      setMessage(`Package ${name} installed successfully.`);
    } catch (err: any) {
      setMessage(`Install failed: ${err.message}`);
    }
  };

  const handleRemove = async (name: string) => {
    try {
      await removePackage(name);
      setMessage(`Package ${name} removed successfully.`);
    } catch (err: any) {
      setMessage(`Remove failed: ${err.message}`);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-white tracking-tight">Debian Package Manager (APT)</h1>
        <p className="text-sm text-slate-400 mt-1">Search, install, and manage system packages from Debian APT repositories</p>
      </div>

      {/* Search Bar */}
      <form onSubmit={handleSearch} className="flex gap-3 max-w-xl">
        <div className="relative flex-1">
          <Search className="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search APT package repository (e.g. nginx, curl, python3)..."
            className="w-full pl-10 pr-4 py-2.5 rounded-xl bg-slate-900/80 border border-border text-xs text-white placeholder:text-slate-500 focus:outline-none focus:border-blue-500 transition-all"
          />
        </div>
        <button
          type="submit"
          disabled={loading}
          className="px-5 py-2.5 rounded-xl bg-blue-600 hover:bg-blue-500 text-white text-xs font-semibold shadow-lg shadow-blue-500/20 transition-all disabled:opacity-50"
        >
          {loading ? 'Searching...' : 'Search APT'}
        </button>
      </form>

      {message && (
        <div className="p-3.5 rounded-xl bg-slate-900 border border-blue-500/30 text-blue-400 text-xs flex items-center gap-2">
          <CheckCircle className="w-4 h-4 shrink-0" />
          <span>{message}</span>
        </div>
      )}

      {/* Packages Table */}
      <div className="glass-panel rounded-2xl border border-border overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs text-slate-300">
            <thead className="bg-slate-900/80 border-b border-border text-slate-400 uppercase text-[11px] font-semibold tracking-wider">
              <tr>
                <th className="px-6 py-4">Package Name</th>
                <th className="px-6 py-4">Description</th>
                <th className="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/60">
              {packages.length > 0 ? (
                packages.map((pkg) => (
                  <tr key={pkg.name} className="hover:bg-slate-800/40 transition-colors">
                    <td className="px-6 py-4 font-semibold text-white flex items-center gap-2">
                      <Package className="w-4 h-4 text-indigo-400 shrink-0" />
                      <span>{pkg.name}</span>
                    </td>
                    <td className="px-6 py-4 text-slate-400 max-w-md truncate">{pkg.description || 'Debian package'}</td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end gap-2">
                        <button
                          onClick={() => handleInstall(pkg.name)}
                          className="px-3 py-1.5 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 border border-emerald-500/20 flex items-center gap-1.5"
                        >
                          <Download className="w-3.5 h-3.5" />
                          <span>Install</span>
                        </button>
                        <button
                          onClick={() => handleRemove(pkg.name)}
                          className="px-3 py-1.5 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 border border-rose-500/20 flex items-center gap-1.5"
                        >
                          <Trash2 className="w-3.5 h-3.5" />
                          <span>Remove</span>
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={3} className="px-6 py-8 text-center text-slate-500">
                    No packages loaded. Search above to query APT repositories.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
