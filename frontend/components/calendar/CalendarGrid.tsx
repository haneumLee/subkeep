'use client';

import type { CalendarDay } from '@/types';
import { formatCurrency } from '@/lib/utils';

interface CalendarGridProps {
  year: number;
  month: number;
  days: CalendarDay[];
  onPrevMonth: () => void;
  onNextMonth: () => void;
  onDayClick: (day: number) => void;
  selectedDay: number;
}

const WEEKDAY_LABELS = ['일', '월', '화', '수', '목', '금', '토'];

export function CalendarGrid({
  year,
  month,
  days,
  onPrevMonth,
  onNextMonth,
  onDayClick,
  selectedDay,
}: CalendarGridProps) {
  const firstDayOfMonth = new Date(year, month - 1, 1).getDay();
  const daysInMonth = new Date(year, month, 0).getDate();
  const dayMap = new Map(days.map((d) => [new Date(d.date).getDate(), d]));

  const cells: (number | null)[] = [];
  for (let i = 0; i < firstDayOfMonth; i++) cells.push(null);
  for (let d = 1; d <= daysInMonth; d++) cells.push(d);

  return (
    <div className="rounded-xl bg-white p-4 shadow-sm lg:p-6">
      {/* Header */}
      <div className="mb-4 flex items-center justify-between">
        <button
          onClick={onPrevMonth}
          className="rounded-lg p-2 text-slate-600 hover:bg-slate-100"
          aria-label="이전 달"
        >
          <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h3 className="text-lg font-semibold text-slate-900">
          {year}년 {month}월
        </h3>
        <button
          onClick={onNextMonth}
          className="rounded-lg p-2 text-slate-600 hover:bg-slate-100"
          aria-label="다음 달"
        >
          <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>

      {/* Weekday Labels */}
      <div className="mb-2 grid grid-cols-7 gap-1 text-center text-sm font-medium text-slate-500">
        {WEEKDAY_LABELS.map((label) => (
          <div key={label} className="py-2">
            {label}
          </div>
        ))}
      </div>

      {/* Day Cells */}
      <div className="grid grid-cols-7 gap-1">
        {cells.map((day, index) => {
          if (day === null) {
            return <div key={`empty-${index}`} className="min-h-[80px]" />;
          }

          const dayData = dayMap.get(day);
          const isSelected = day === selectedDay;
          const hasPayments = dayData && dayData.subscriptions.length > 0;

          return (
            <button
              key={day}
              onClick={() => hasPayments && onDayClick(day)}
              className={`min-h-[80px] rounded-lg p-1 text-left transition-colors ${
                isSelected
                  ? 'bg-primary-50 ring-2 ring-primary-500'
                  : hasPayments
                    ? 'hover:bg-slate-50 cursor-pointer'
                    : 'cursor-default'
              }`}
            >
              <span
                className={`text-sm font-medium ${
                  isSelected ? 'text-primary-700' : 'text-slate-700'
                }`}
              >
                {day}
              </span>
              {hasPayments && (
                <div className="mt-1">
                  <p className="text-xs font-semibold text-primary-600">
                    {formatCurrency(dayData.totalAmount)}
                  </p>
                  <div className="mt-0.5 flex flex-wrap gap-0.5">
                    {dayData.subscriptions.slice(0, 3).map((sub) => (
                      <span
                        key={sub.subscriptionId}
                        className="inline-block h-1.5 w-1.5 rounded-full"
                        style={{ backgroundColor: sub.categoryColor || '#6366f1' }}
                        title={sub.serviceName}
                      />
                    ))}
                    {dayData.subscriptions.length > 3 && (
                      <span className="text-[10px] text-slate-400">
                        +{dayData.subscriptions.length - 3}
                      </span>
                    )}
                  </div>
                </div>
              )}
            </button>
          );
        })}
      </div>
    </div>
  );
}
