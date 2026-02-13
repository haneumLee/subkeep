'use client';

import { formatCurrency } from '@/lib/utils';
import type { ReportCategoryBreakdown } from '@/types';

interface CategoryPieChartProps {
  categories: ReportCategoryBreakdown[];
}

export function CategoryPieChart({ categories }: CategoryPieChartProps) {
  if (categories.length === 0) {
    return (
      <div className="rounded-xl bg-white p-6 shadow-sm">
        <h3 className="mb-4 text-lg font-semibold text-slate-900">카테고리별 지출</h3>
        <p className="text-sm text-slate-500">데이터가 없습니다.</p>
      </div>
    );
  }

  // Build conic-gradient stops
  let accumulated = 0;
  const gradientStops = categories.map((cat) => {
    const start = accumulated;
    accumulated += cat.percentage;
    return `${cat.color || '#6366f1'} ${start}% ${accumulated}%`;
  });

  const gradient = `conic-gradient(${gradientStops.join(', ')})`;

  return (
    <div className="rounded-xl bg-white p-6 shadow-sm">
      <h3 className="mb-4 text-lg font-semibold text-slate-900">카테고리별 지출</h3>

      <div className="flex flex-col items-center gap-6 sm:flex-row">
        {/* Donut Chart */}
        <div className="relative h-48 w-48 flex-shrink-0">
          <div
            className="h-full w-full rounded-full"
            style={{ background: gradient }}
          />
          <div className="absolute inset-0 m-auto h-24 w-24 rounded-full bg-white" />
        </div>

        {/* Legend */}
        <ul className="flex-1 space-y-2">
          {categories.map((cat) => (
            <li key={cat.categoryId} className="flex items-center justify-between text-sm">
              <div className="flex items-center gap-2">
                <span
                  className="inline-block h-3 w-3 rounded-full"
                  style={{ backgroundColor: cat.color || '#6366f1' }}
                />
                <span className="text-slate-700">{cat.categoryName}</span>
              </div>
              <div className="text-right">
                <span className="font-medium text-slate-900">
                  {formatCurrency(cat.monthlyAmount)}
                </span>
                <span className="ml-2 text-slate-400">({cat.percentage.toFixed(1)}%)</span>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
