import { useQuery } from '@tanstack/react-query';

import { get } from '@/lib/api';
import type { Category } from '@/types';

export function useCategories() {
  return useQuery<Category[]>({
    queryKey: ['categories'],
    queryFn: () => get<Category[]>('/categories'),
  });
}
