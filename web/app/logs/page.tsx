'use client';

import { useState, useEffect } from 'react';
import { getLogs } from '@/lib/api';
import { LogEntry } from '@/lib/types';
import { FileText, Search, RefreshCw, Filter } from 'lucide-react';

export default function LogsPage() {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [unit, setUnit] = useState('');
  const [lines, setLines] = useState(50);
  const [filterText, setFilterText] = useState('');
  const [loading, setLoading] = useState(false);

  const fetchLogs = async () => {
    setLoading(true);
    try {
      const data = await getLogs(unit || undefined, lines);
      setLogs(data);
    } catch (e) {
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLogs();
  }, [unit, lines]);

  const filteredLogs = logs.filter((log) =>
    log.message.toLowerCase().includes(filterText.toLowerCase()) ||
    log.process.toLowerCase().includes(filterText.toLowerCase())
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-white tracking-tight">System Logs (journalctl)</h1>
          <p className="text-sm text-slate-400 mt-1">Live journalctl log aggregation, filtering, and inspection</p>
        </div>
        <button
          onClick={fetchLogs}
          disabled={loading}
          className="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold flex items-center gap-2 border border-border transition-all"
        >
          <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
          <span>Refresh</span>
        </button>
      </div>

      {/* Filter Controls */}
      <div className="flex flex-wrap gap-3 items-center">
        <div className="relative flex-1 min-w-[200px]">
          <Search className="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
          <input
            type="text"
            placeholder="Search log contents..."
            value={filterText}
            onChange={(e) => setFilterText(e.target.value)}
            className="w-full pl-10 pr-4 py-2 rounded-xl bg-slate-900/80 border border-border text-xs text-white placeholder:text-slate-500 focus:outline-none focus:border-blue-500"
          />
        </div>

        <input
          type="text"
          placeholder="Filter by Unit (e.g. sawitd)"
          value={unit}
          onChange={(e) => setUnit(e.target.value)}
          className="px-3.5 py-2 rounded-xl bg-slate-900/80 border border-border text-xs text-white placeholder:text-slate-500 focus:outline-none focus:border-blue-500"
        />

        <select
          value={lines}
          onChange={(e) => setLines(Number(e.target.value))}
          className="px-3 py-2 rounded-xl bg-slate-900/80 border border-border text-xs text-white focus:outline-none"
        >
          <option value={25}>25 Lines</option>
          <option value={50}>50 Lines</option>
          <option value={100}>100 Lines</option>
          <option value={200}>200 Lines</option>
        </select>
      </div>

      {/* Console Output Panel */}
      <div className="glass-panel rounded-2xl border border-border overflow-hidden">
        <div className="px-6 py-3 bg-slate-900/90 border-b border-border flex justify-between items-center text-xs text-slate-400 font-mono">
          <span>CONSOLE LOG STREAM</span>
          <span>{filteredLogs.length} ENTRIES</span>
        </div>
        <div className="p-4 bg-slate-950 font-mono text-xs text-slate-300 max-h-[550px] overflow-y-auto space-y-1.5 border-t border-border/40">
          {filteredLogs.length > 0 ? (
            filteredLogs.map((log, idx) => (
              <div key={idx} className="flex items-start gap-3 hover:bg-slate-900/40 p-1 rounded transition-colors">
                <span className="text-slate-500 shrink-0">{new Date(log.timestamp).toLocaleTimeString()}</span>
                <span className="text-blue-400 font-semibold shrink-0">{log.host}</span>
                <span className="text-emerald-400 font-semibold shrink-0">{log.process}:</span>
                <span className="text-slate-300 break-all">{log.message}</span>
              </div>
            ))
          ) : (
            <p className="text-slate-500 text-center py-8">No log lines matched the criteria.</p>
          )}
        </div>
      </div>
    </div>
  );
}
