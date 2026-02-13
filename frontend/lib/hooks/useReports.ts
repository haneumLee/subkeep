import { useQuery } from '@tanstack/react-query';

import { get } from '@/lib/api';
import type { ReportOverview } from '@/types';

export function useReportOverview() {
  return useQuery<ReportOverview>({
    queryKey: ['reports', 'overview'],
    queryFn: () => get<ReportOverview>('/reports/overview'),
  });
}
