'use client';

import { formatCurrency } from '@/lib/utils';
import type { AverageCost } from '@/types';

interface CostSummaryCardsProps {
  averageCost: AverageCost;
}

export function CostSummaryCards({ averageCost }: CostSummaryCardsProps) {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
      <div className="rounded-xl bg-white p-4 shadow-sm">
        <p className="text-sm text-slate-500">주간 평균</p>
        <p className="text-xl font-bold text-slate-900">{formatCurrency(averageCost.weekly)}</p>
      </div>
      <div className="rounded-xl bg-white p-4 shadow-sm">
        <p className="text-sm text-slate-500">월간 평균</p>
        <p className="text-xl font-bold text-primary-600">{formatCurrency(averageCost.monthly)}</p>
      </div>
      <div className="rounded-xl bg-white p-4 shadow-sm">
        <p className="text-sm text-slate-500">연간 평균</p>
        <p className="text-xl font-bold text-slate-900">{formatCurrency(averageCost.annual)}</p>
      </div>
    </div>
  );
}
