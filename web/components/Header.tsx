'use client';

import { useState, useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { removeAuthToken, getSystemStatus } from '@/lib/api';
import { SystemStatus } from '@/lib/types';
import { Shield, Clock, Terminal as TerminalIcon, LogOut, Activity } from 'lucide-react';

export default function Header() {
  const router = useRouter();
  const pathname = usePathname();
  const [status, setStatus] = useState<SystemStatus | null>(null);

  useEffect(() => {
    if (pathname === '/login') return;
    getSystemStatus().then(setStatus).catch(() => {});
  }, [pathname]);

  if (pathname === '/login') return null;

  const handleLogout = () => {
    removeAuthToken();
    router.push('/login');
  };

  return (
    <header className="glass-panel border-b border-border px-6 py-3 sticky top-0 z-20 flex items-center justify-between">
      {/* System Status Indicators */}
      <div className="flex items-center gap-6">
        <div className="flex items-center gap-2">
          <div className="w-2.5 h-2.5 rounded-full bg-emerald-400 animate-ping"></div>
          <span className="text-sm font-semibold text-white tracking-wide">
            {status?.info.hostname || 'sawit-server'}
          </span>
        </div>

        <div className="hidden md:flex items-center gap-4 text-xs text-slate-400 font-mono border-l border-border pl-6">
          <div className="flex items-center gap-1.5">
            <Clock className="w-3.5 h-3.5 text-blue-400" />
            <span>Uptime: {status?.info.uptime_formatted || '0m'}</span>
          </div>
          <div className="flex items-center gap-1.5">
            <Activity className="w-3.5 h-3.5 text-emerald-400" />
            <span>Load: {status?.resources.load_avg_1.toFixed(2) || '0.00'}</span>
          </div>
          <div className="flex items-center gap-1.5">
            <Shield className="w-3.5 h-3.5 text-indigo-400" />
            <span>Firewall: Active</span>
          </div>
        </div>
      </div>

      {/* User Actions */}
      <div className="flex items-center gap-3">
        <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-slate-800/60 border border-border text-xs text-slate-300">
          <div className="w-2 h-2 rounded-full bg-emerald-400"></div>
          <span className="font-medium">admin</span>
        </div>

        <button
          onClick={handleLogout}
          className="p-2 rounded-lg bg-slate-800/40 hover:bg-rose-500/20 text-slate-400 hover:text-rose-400 border border-border transition-all duration-150"
          title="Logout"
        >
          <LogOut className="w-4 h-4" />
        </button>
      </div>
    </header>
  );
}
