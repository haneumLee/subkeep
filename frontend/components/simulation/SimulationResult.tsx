'use client';

import type { CategoryBreakdown, SimulationResult } from '@/types';
import { cn, formatCurrency } from '@/lib/utils';

interface SimulationResultProps {
  result: SimulationResult | null;
  loading?: boolean;
}

export function SimulationResult({ result, loading }: SimulationResultProps) {
  if (loading) {
    return (
      <div className="rounded-lg border border-gray-200 bg-white p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-6 w-32 rounded bg-gray-200"></div>
          <div className="h-10 w-full rounded bg-gray-200"></div>
          <div className="h-24 w-full rounded bg-gray-200"></div>
        </div>
      </div>
    );
  }

  if (!result) {
    return (
      <div className="rounded-lg border border-gray-200 bg-white p-6 text-center text-gray-500">
        시뮬레이션 결과가 여기에 표시됩니다
      </div>
    );
  }

  const isNegative = result.monthlyDifference < 0;
  const isPositive = result.monthlyDifference > 0;

  return (
    <div className="space-y-6 rounded-lg border border-gray-200 bg-white p-6">
      {/* 금액 변화 */}
      <div>
        <h3 className="mb-4 text-lg font-semibold text-gray-900">금액 변화</h3>
        <div className="flex items-center justify-between">
          <div className="text-center">
            <div className="text-sm text-gray-500">현재 월 총액</div>
            <div className="mt-1 text-2xl font-bold text-gray-900">
              {formatCurrency(result.currentMonthlyTotal)}
            </div>
          </div>

          <div className="flex-shrink-0 px-4">
            <svg
              className="h-8 w-8 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M13 7l5 5m0 0l-5 5m5-5H6"
              />
            </svg>
          </div>

          <div className="text-center">
            <div className="text-sm text-gray-500">변경 후 월 총액</div>
            <div className="mt-1 text-2xl font-bold text-gray-900">
              {formatCurrency(result.simulatedMonthlyTotal)}
            </div>
          </div>
        </div>

        <div className="mt-4 rounded-lg bg-gray-50 p-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <div className="text-sm text-gray-500">월 차액</div>
              <div
                className={cn('mt-1 text-xl font-bold', {
                  'text-green-600': isNegative, // 절감 (음수)
                  'text-red-600': isPositive, // 증가 (양수)
                  'text-gray-900': result.monthlyDifference === 0,
                })}
              >
                {isNegative && '-'}
                {formatCurrency(Math.abs(result.monthlyDifference))}
                {isNegative && ' 절감'}
                {isPositive && ' 증가'}
              </div>
            </div>
            <div>
              <div className="text-sm text-gray-500">연 차액</div>
              <div
                className={cn('mt-1 text-xl font-bold', {
                  'text-green-600': isNegative,
                  'text-red-600': isPositive,
                  'text-gray-900': result.annualDifference === 0,
                })}
              >
                {isNegative && '-'}
                {formatCurrency(Math.abs(result.annualDifference))}
                {isNegative && ' 절감'}
                {isPositive && ' 증가'}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 카테고리 비중 */}
      {result.categoryBreakdown.length > 0 && (
        <div>
          <h3 className="mb-4 text-lg font-semibold text-gray-900">카테고리 비중</h3>
          <div className="space-y-3">
            {result.categoryBreakdown.map((category: CategoryBreakdown) => (
              <div key={category.categoryId}>
                <div className="mb-1 flex items-center justify-between text-sm">
                  <span className="font-medium text-gray-700">{category.categoryName}</span>
                  <span className="text-gray-500">
                    {formatCurrency(category.amount)} ({category.percentage.toFixed(1)}%)
                  </span>
                </div>
                <div className="h-2 w-full overflow-hidden rounded-full bg-gray-200">
                  <div
                    className="h-full rounded-full transition-all duration-300"
                    style={{
                      width: `${category.percentage}%`,
                      backgroundColor: category.categoryColor || '#3b82f6',
                    }}
                  ></div>
                </div>
                <div className="mt-1 text-xs text-gray-500">{category.count}개 구독</div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
