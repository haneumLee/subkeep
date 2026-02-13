'use client';

import { formatCurrency } from '@/lib/utils';
import type { ReportSummary } from '@/types';

interface ReportSummaryPanelProps {
  summary: ReportSummary;
}

export function ReportSummaryPanel({ summary }: ReportSummaryPanelProps) {
  return (
    <div className="rounded-xl bg-white p-6 shadow-sm">
      <h3 className="mb-4 text-lg font-semibold text-slate-900">구독 요약</h3>

      <div className="grid grid-cols-2 gap-4 lg:grid-cols-3">
        <div>
          <p className="text-sm text-slate-500">전체 구독</p>
          <p className="text-lg font-bold text-slate-900">{summary.totalSubscriptions}개</p>
        </div>
        <div>
          <p className="text-sm text-slate-500">활성</p>
          <p className="text-lg font-bold text-green-600">{summary.activeCount}개</p>
        </div>
        <div>
          <p className="text-sm text-slate-500">일시정지</p>
          <p className="text-lg font-bold text-amber-600">{summary.pausedCount}개</p>
        </div>
        <div>
          <p className="text-sm text-slate-500">최고가 구독</p>
          <p className="text-lg font-bold text-slate-900">
            {summary.mostExpensive || '-'}
          </p>
          {summary.mostExpensive && (
            <p className="text-sm text-slate-500">{formatCurrency(summary.mostExpensiveAmount)}</p>
          )}
        </div>
        <div>
          <p className="text-sm text-slate-500">평균 만족도</p>
          <p className="text-lg font-bold text-slate-900">
            {summary.averageSatisfaction > 0
              ? `${'★'.repeat(Math.round(summary.averageSatisfaction))}${'☆'.repeat(5 - Math.round(summary.averageSatisfaction))}`
              : '미평가'}
          </p>
        </div>
      </div>
    </div>
  );
}
