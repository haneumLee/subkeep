import { formatCurrency } from '@/lib/utils';
import type { CategoryBreakdown } from '@/types';

interface CategoryChartProps {
  breakdown: CategoryBreakdown[];
}

const DEFAULT_COLORS = [
  '#3B82F6',
  '#8B5CF6',
  '#EC4899',
  '#F59E0B',
  '#10B981',
  '#6366F1',
  '#EF4444',
  '#14B8A6',
];

export default function CategoryChart({ breakdown }: CategoryChartProps) {
  if (breakdown.length === 0) {
    return (
      <div className="rounded-lg border-2 border-gray-200 bg-white p-6">
        <h3 className="mb-4 text-lg font-semibold text-gray-900">카테고리별 비중</h3>
        <div className="flex items-center justify-center py-12 text-gray-500">
          등록된 구독이 없습니다
        </div>
      </div>
    );
  }

  // Calculate conic-gradient stops
  let currentPercentage = 0;
  const gradientStops = breakdown
    .map((item, index) => {
      const color = item.categoryColor || DEFAULT_COLORS[index % DEFAULT_COLORS.length];
      const start = currentPercentage;
      const end = currentPercentage + item.percentage;
      currentPercentage = end;
      return `${color} ${start}% ${end}%`;
    })
    .join(', ');

  return (
    <div className="rounded-lg border-2 border-gray-200 bg-white p-6">
      <h3 className="mb-6 text-lg font-semibold text-gray-900">카테고리별 비중</h3>

      <div className="flex flex-col items-center gap-6 lg:flex-row">
        {/* Donut Chart */}
        <div className="flex-shrink-0">
          <div className="relative h-48 w-48">
            <div
              className="h-full w-full rounded-full"
              style={{
                background: `conic-gradient(${gradientStops})`,
              }}
            />
            <div className="absolute left-1/2 top-1/2 h-28 w-28 -translate-x-1/2 -translate-y-1/2 rounded-full bg-white" />
          </div>
        </div>

        {/* Legend */}
        <div className="flex-1 space-y-2">
          {breakdown.map((item, index) => {
            const color = item.categoryColor || DEFAULT_COLORS[index % DEFAULT_COLORS.length];
            return (
              <div key={item.categoryId} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div
                    className="h-3 w-3 rounded-full"
                    style={{ backgroundColor: color }}
                  />
                  <span className="text-sm font-medium text-gray-700">
                    {item.categoryName}
                  </span>
                  <span className="text-xs text-gray-500">({item.count}개)</span>
                </div>
                <div className="text-right">
                  <div className="text-sm font-semibold text-gray-900">
                    {formatCurrency(item.amount)}
                  </div>
                  <div className="text-xs text-gray-500">{item.percentage.toFixed(1)}%</div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
