import { LucideIcon } from 'lucide-react';

interface MetricCardProps {
  title: string;
  value: string | number;
  subtitle?: string;
  percent?: number;
  icon: LucideIcon;
  color?: 'blue' | 'emerald' | 'amber' | 'rose' | 'indigo';
}

export default function MetricCard({
  title,
  value,
  subtitle,
  percent,
  icon: Icon,
  color = 'blue',
}: MetricCardProps) {
  const colorMap = {
    blue: { text: 'text-blue-400', bg: 'bg-blue-500/10', bar: 'bg-blue-500' },
    emerald: { text: 'text-emerald-400', bg: 'bg-emerald-500/10', bar: 'bg-emerald-500' },
    amber: { text: 'text-amber-400', bg: 'bg-amber-500/10', bar: 'bg-amber-500' },
    rose: { text: 'text-rose-400', bg: 'bg-rose-500/10', bar: 'bg-rose-500' },
    indigo: { text: 'text-indigo-400', bg: 'bg-indigo-500/10', bar: 'bg-indigo-500' },
  };

  const activeColor = colorMap[color];

  return (
    <div className="glass-panel glass-panel-hover p-5 rounded-2xl border border-border">
      <div className="flex justify-between items-start mb-3">
        <div>
          <span className="text-xs font-semibold text-slate-400 tracking-wider uppercase">{title}</span>
          <h3 className="text-2xl font-bold text-white mt-1">{value}</h3>
        </div>
        <div className={`p-2.5 rounded-xl ${activeColor.bg} ${activeColor.text}`}>
          <Icon className="w-5 h-5" />
        </div>
      </div>

      {percent !== undefined && (
        <div className="mt-3">
          <div className="flex justify-between text-xs text-slate-400 mb-1">
            <span>Utilization</span>
            <span className="font-semibold text-slate-200">{percent.toFixed(1)}%</span>
          </div>
          <div className="w-full h-2 rounded-full bg-slate-800 overflow-hidden">
            <div
              className={`h-full rounded-full transition-all duration-500 ${activeColor.bar}`}
              style={{ width: `${Math.min(100, Math.max(0, percent))}%` }}
            ></div>
          </div>
        </div>
      )}

      {subtitle && <p className="text-xs text-slate-400 mt-2 font-mono">{subtitle}</p>}
    </div>
  );
}
