'use client';

import { CategoryPieChart } from '@/components/reports/CategoryPieChart';
import { CostSummaryCards } from '@/components/reports/CostSummaryCards';
import { MonthlyTrendChart } from '@/components/reports/MonthlyTrendChart';
import { ReportSummaryPanel } from '@/components/reports/ReportSummaryPanel';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import { useReportOverview } from '@/lib/hooks/useReports';

export default function ReportsPage() {
  const { data, isLoading } = useReportOverview();

  if (isLoading) {
    return (
      <div className="flex min-h-[400px] items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (!data) {
    return (
      <div className="flex min-h-[400px] items-center justify-center">
        <p className="text-slate-500">리포트 데이터를 불러올 수 없습니다.</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-900">리포트</h2>

      {/* Cost Summary */}
      <CostSummaryCards averageCost={data.averageCost} />

      {/* Charts Row */}
      <div className="grid gap-6 lg:grid-cols-2">
        <CategoryPieChart categories={data.categoryBreakdown} />
        <MonthlyTrendChart trends={data.monthlyTrend} />
      </div>

      {/* Report Summary */}
      <ReportSummaryPanel summary={data.summary} />
    </div>
  );
}
