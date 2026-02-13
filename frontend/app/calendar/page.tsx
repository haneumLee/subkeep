'use client';

import { useState } from 'react';

import { CalendarGrid } from '@/components/calendar/CalendarGrid';
import { DayDetailModal } from '@/components/calendar/DayDetailModal';
import { UpcomingPaymentsList } from '@/components/calendar/UpcomingPaymentsList';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import { useMonthlyCalendar, useDayDetail, useUpcomingPayments } from '@/lib/hooks/useCalendar';
import { formatCurrency } from '@/lib/utils';

export default function CalendarPage() {
  const now = new Date();
  const [year, setYear] = useState(now.getFullYear());
  const [month, setMonth] = useState(now.getMonth() + 1);
  const [selectedDay, setSelectedDay] = useState(0);

  const { data: calendar, isLoading } = useMonthlyCalendar(year, month);
  const { data: dayDetail } = useDayDetail(year, month, selectedDay);
  const { data: upcoming } = useUpcomingPayments(30);

  const handlePrevMonth = () => {
    if (month === 1) {
      setYear(year - 1);
      setMonth(12);
    } else {
      setMonth(month - 1);
    }
    setSelectedDay(0);
  };

  const handleNextMonth = () => {
    if (month === 12) {
      setYear(year + 1);
      setMonth(1);
    } else {
      setMonth(month + 1);
    }
    setSelectedDay(0);
  };

  if (isLoading) {
    return (
      <div className="flex min-h-[400px] items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-900">결제일 캘린더</h2>

      {/* Summary Cards */}
      {calendar && (
        <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
          <div className="rounded-xl bg-white p-4 shadow-sm">
            <p className="text-sm text-slate-500">이번 달 총 결제</p>
            <p className="text-xl font-bold text-slate-900">{formatCurrency(calendar.totalAmount)}</p>
          </div>
          <div className="rounded-xl bg-white p-4 shadow-sm">
            <p className="text-sm text-slate-500">결제 건수</p>
            <p className="text-xl font-bold text-slate-900">{calendar.totalCount}건</p>
          </div>
          <div className="rounded-xl bg-white p-4 shadow-sm">
            <p className="text-sm text-slate-500">남은 결제</p>
            <p className="text-xl font-bold text-primary-600">{formatCurrency(calendar.remainingAmount)}</p>
          </div>
          <div className="rounded-xl bg-white p-4 shadow-sm">
            <p className="text-sm text-slate-500">남은 건수</p>
            <p className="text-xl font-bold text-primary-600">{calendar.remainingCount}건</p>
          </div>
        </div>
      )}

      {/* Calendar Grid */}
      <CalendarGrid
        year={year}
        month={month}
        days={calendar?.days || []}
        onPrevMonth={handlePrevMonth}
        onNextMonth={handleNextMonth}
        onDayClick={setSelectedDay}
        selectedDay={selectedDay}
      />

      {/* Upcoming Payments */}
      <UpcomingPaymentsList payments={upcoming || []} />

      {/* Day Detail Modal */}
      {dayDetail && (
        <DayDetailModal
          dayDetail={dayDetail}
          isOpen={selectedDay > 0}
          onClose={() => setSelectedDay(0)}
        />
      )}
    </div>
  );
}
