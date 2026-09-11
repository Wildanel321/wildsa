interface BadgeProps {
  status: 'HEALTHY' | 'WARNING' | 'CRITICAL' | 'running' | 'stopped' | 'active' | 'allow' | 'deny';
  text?: string;
}

export default function Badge({ status, text }: BadgeProps) {
  const label = text || status.toUpperCase();

  const variantMap: Record<string, string> = {
    HEALTHY: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30',
    running: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30',
    active: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30',
    allow: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30',
    WARNING: 'bg-amber-500/15 text-amber-400 border-amber-500/30',
    CRITICAL: 'bg-rose-500/15 text-rose-400 border-rose-500/30',
    stopped: 'bg-slate-800 text-slate-400 border-slate-700',
    deny: 'bg-rose-500/15 text-rose-400 border-rose-500/30',
  };

  const styleClass = variantMap[status] || 'bg-slate-800 text-slate-300 border-slate-700';

  return (
    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold border ${styleClass}`}>
      <span className="w-1.5 h-1.5 rounded-full bg-current"></span>
      {label}
    </span>
  );
}
