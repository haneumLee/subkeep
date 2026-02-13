import { useQuery } from '@tanstack/react-query';

import { get } from '@/lib/api';
import type { DashboardSummary, Recommendation } from '@/types';

export function useDashboardSummary() {
  return useQuery<DashboardSummary>({
    queryKey: ['dashboard', 'summary'],
    queryFn: () => get<DashboardSummary>('/dashboard/summary'),
  });
}

export function useRecommendations() {
  return useQuery<Recommendation[]>({
    queryKey: ['dashboard', 'recommendations'],
    queryFn: () => get<Recommendation[]>('/dashboard/recommendations'),
  });
}
