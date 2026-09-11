'use client';

import { useState, useEffect } from 'react';
import { getFirewallStatus, allowFirewallRule, denyFirewallRule, getSecurityAudit, getSecurityUpdates } from '@/lib/api';
import { FirewallStatus, AuditReport, SecurityUpdateInfo } from '@/lib/types';
import Badge from '@/components/Badge';
import { Shield, Plus, Lock, CheckCircle, AlertTriangle, RefreshCw, FileText } from 'lucide-react';

export default function SecurityPage() {
  const [firewall, setFirewall] = useState<FirewallStatus | null>(null);
  const [audit, setAudit] = useState<AuditReport | null>(null);
  const [updates, setUpdates] = useState<SecurityUpdateInfo | null>(null);
  const [newRule, setNewRule] = useState('');
  const [ruleType, setRuleType] = useState<'allow' | 'deny'>('allow');
  const [loading, setLoading] = useState(false);

  const fetchSecurityData = async () => {
    try {
      const [fwData, auditData, updatesData] = await Promise.all([
        getFirewallStatus(),
        getSecurityAudit(),
        getSecurityUpdates(),
      ]);
      setFirewall(fwData);
      setAudit(auditData);
      setUpdates(updatesData);
    } catch (e) {}
  };

  useEffect(() => {
    fetchSecurityData();
  }, []);

  const handleAddRule = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newRule) return;
    setLoading(true);
    try {
      if (ruleType === 'allow') await allowFirewallRule(newRule);
      if (ruleType === 'deny') await denyFirewallRule(newRule);
      setNewRule('');
      await fetchSecurityData();
    } catch (err: any) {
      alert(`Firewall error: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white tracking-tight">Security & Audit Center</h1>
          <p className="text-sm text-slate-400 mt-1">Automated posture checks, nftables rule filtering & SSH hardening</p>
        </div>
        <button
          onClick={fetchSecurityData}
          className="px-3.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs text-white flex items-center gap-2 border border-border w-fit"
        >
          <RefreshCw className="w-3.5 h-3.5" />
          <span>Run Audit Now</span>
        </button>
      </div>

      {/* Security Posture Overview Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="glass-panel p-5 rounded-2xl border border-border flex items-center justify-between">
          <div>
            <p className="text-xs font-semibold text-slate-400">Security Score</p>
            <h2 className="text-2xl font-bold text-emerald-400 mt-1">{audit?.score ? `${audit.score.toFixed(0)}%` : '100%'}</h2>
            <p className="text-[11px] text-slate-400 mt-1">{audit?.passed || 4} passed, {audit?.failed || 0} failed checks</p>
          </div>
          <div className="p-3 rounded-2xl bg-emerald-500/10 text-emerald-400">
            <Shield className="w-7 h-7" />
          </div>
        </div>

        <div className="glass-panel p-5 rounded-2xl border border-border flex items-center justify-between">
          <div>
            <p className="text-xs font-semibold text-slate-400">nftables Stateful Firewall</p>
            <h2 className="text-xl font-bold text-white mt-1">{firewall?.enabled ? 'Active & Enforcing' : 'Inactive'}</h2>
            <p className="text-[11px] text-slate-400 mt-1">Backend: {firewall?.backend || 'nftables'}</p>
          </div>
          <div className="p-3 rounded-2xl bg-blue-500/10 text-blue-400">
            <Lock className="w-7 h-7" />
          </div>
        </div>

        <div className="glass-panel p-5 rounded-2xl border border-border flex items-center justify-between">
          <div>
            <p className="text-xs font-semibold text-slate-400">Pending Security Updates</p>
            <h2 className="text-2xl font-bold text-white mt-1">{updates?.pending_security_updates || 0}</h2>
            <p className="text-[11px] text-slate-400 mt-1">Debian Security Repository</p>
          </div>
          <div className="p-3 rounded-2xl bg-indigo-500/10 text-indigo-400">
            <CheckCircle className="w-7 h-7" />
          </div>
        </div>
      </div>

      {/* Audit Checklist Table */}
      <div className="glass-panel rounded-2xl border border-border overflow-hidden">
        <div className="px-6 py-4 border-b border-border flex items-center justify-between">
          <div className="flex items-center gap-2">
            <FileText className="w-4 h-4 text-emerald-400" />
            <h3 className="text-sm font-bold text-white">Automated Security Audit Findings</h3>
          </div>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs text-slate-300">
            <thead className="bg-slate-900/80 border-b border-border text-slate-400 uppercase text-[11px] font-semibold tracking-wider">
              <tr>
                <th className="px-6 py-4">ID</th>
                <th className="px-6 py-4">Category</th>
                <th className="px-6 py-4">Title & Details</th>
                <th className="px-6 py-4">Severity</th>
                <th className="px-6 py-4">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/60">
              {audit?.check_items.map((item) => (
                <tr key={item.id} className="hover:bg-slate-800/40 transition-colors">
                  <td className="px-6 py-4 font-mono text-slate-400">{item.id}</td>
                  <td className="px-6 py-4 font-semibold text-white">{item.category}</td>
                  <td className="px-6 py-4">
                    <p className="font-semibold text-white">{item.title}</p>
                    <p className="text-[11px] text-slate-400 mt-0.5">{item.details}</p>
                  </td>
                  <td className="px-6 py-4">
                    <span className={`px-2 py-0.5 rounded text-[10px] font-bold ${
                      item.severity === 'HIGH' ? 'bg-red-500/20 text-red-400' : 'bg-amber-500/20 text-amber-400'
                    }`}>
                      {item.severity}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <Badge status={item.passed ? 'HEALTHY' : 'CRITICAL'} text={item.passed ? 'PASS' : 'FAIL'} />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Firewall Rules Section */}
      <div className="glass-panel p-6 rounded-2xl border border-border space-y-4">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-3 border-b border-border">
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-xl bg-blue-500/10 text-blue-400">
              <Shield className="w-5 h-5" />
            </div>
            <div>
              <h2 className="text-base font-bold text-white">Active Firewall Rules (nftables)</h2>
              <p className="text-xs text-slate-400">Add or manage ingress port rules</p>
            </div>
          </div>

          <form onSubmit={handleAddRule} className="flex gap-2 w-full sm:w-auto">
            <input
              type="text"
              value={newRule}
              onChange={(e) => setNewRule(e.target.value)}
              placeholder="Rule (e.g. 8080/tcp)"
              className="flex-1 px-3.5 py-2 rounded-xl bg-slate-900/80 border border-border text-xs text-white placeholder:text-slate-500 focus:outline-none focus:border-blue-500"
            />
            <select
              value={ruleType}
              onChange={(e: any) => setRuleType(e.target.value)}
              className="px-3 py-2 rounded-xl bg-slate-900/80 border border-border text-xs text-white focus:outline-none"
            >
              <option value="allow">ALLOW</option>
              <option value="deny">DENY</option>
            </select>
            <button
              type="submit"
              disabled={loading}
              className="px-4 py-2 rounded-xl bg-blue-600 hover:bg-blue-500 text-white text-xs font-semibold flex items-center gap-1.5 shadow-lg shadow-blue-500/20"
            >
              <Plus className="w-4 h-4" />
              <span>Add</span>
            </button>
          </form>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs text-slate-300">
            <thead className="bg-slate-900/80 border-b border-border text-slate-400 uppercase text-[11px] font-semibold tracking-wider">
              <tr>
                <th className="px-6 py-4">Rule ID</th>
                <th className="px-6 py-4">Port</th>
                <th className="px-6 py-4">Protocol</th>
                <th className="px-6 py-4">Action</th>
                <th className="px-6 py-4">Comment</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/60">
              {firewall?.rules.map((rule) => (
                <tr key={rule.id} className="hover:bg-slate-800/40 transition-colors">
                  <td className="px-6 py-4 font-mono text-slate-400">{rule.id}</td>
                  <td className="px-6 py-4 font-semibold text-white">{rule.port}</td>
                  <td className="px-6 py-4 uppercase font-mono text-blue-400">{rule.protocol}</td>
                  <td className="px-6 py-4">
                    <Badge status={rule.action} />
                  </td>
                  <td className="px-6 py-4 text-slate-400">{rule.comment || 'Firewall rule'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

