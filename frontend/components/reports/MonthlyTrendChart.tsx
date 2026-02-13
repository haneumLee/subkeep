'use client';

import { formatCurrency } from '@/lib/utils';
import type { MonthlyTrend } from '@/types';

interface MonthlyTrendChartProps {
  trends: MonthlyTrend[];
}

export function MonthlyTrendChart({ trends }: MonthlyTrendChartProps) {
  if (trends.length === 0) {
    return (
      <div className="rounded-xl bg-white p-6 shadow-sm">
        <h3 className="mb-4 text-lg font-semibold text-slate-900">월별 추이</h3>
        <p className="text-sm text-slate-500">데이터가 없습니다.</p>
      </div>
    );
  }

  const maxAmount = Math.max(...trends.map((t) => t.amount), 1);

  return (
    <div className="rounded-xl bg-white p-6 shadow-sm">
      <h3 className="mb-4 text-lg font-semibold text-slate-900">월별 추이</h3>

      <div className="flex items-end gap-2" style={{ height: '200px' }}>
        {trends.map((trend) => {
          const heightPercent = (trend.amount / maxAmount) * 100;
          return (
            <div
              key={`${trend.year}-${trend.month}`}
              className="flex flex-1 flex-col items-center justify-end"
            >
              <span className="mb-1 text-xs text-slate-500">
                {formatCurrency(trend.amount)}
              </span>
              <div
                className="w-full max-w-[40px] rounded-t-md bg-primary-500 transition-all"
                style={{ height: `${heightPercent}%`, minHeight: trend.amount > 0 ? '4px' : '0' }}
              />
              <span className="mt-1 text-xs text-slate-500">
                {trend.month}월
              </span>
            </div>
          );
        })}
      </div>
    </div>
  );
}
