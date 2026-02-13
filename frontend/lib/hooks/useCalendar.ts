import { useQuery } from '@tanstack/react-query';

import { get } from '@/lib/api';
import type { DayDetail, MonthlyCalendar, UpcomingPayment } from '@/types';

export function useMonthlyCalendar(year: number, month: number) {
  return useQuery<MonthlyCalendar>({
    queryKey: ['calendar', 'monthly', year, month],
    queryFn: () => get<MonthlyCalendar>('/calendar/monthly', { year, month }),
  });
}

export function useDayDetail(year: number, month: number, day: number) {
  return useQuery<DayDetail>({
    queryKey: ['calendar', 'daily', year, month, day],
    queryFn: () => get<DayDetail>('/calendar/daily', { year, month, day }),
    enabled: day > 0,
  });
}

export function useUpcomingPayments(days: number = 30) {
  return useQuery<UpcomingPayment[]>({
    queryKey: ['calendar', 'upcoming', days],
    queryFn: () => get<UpcomingPayment[]>('/calendar/upcoming', { days }),
  });
}
