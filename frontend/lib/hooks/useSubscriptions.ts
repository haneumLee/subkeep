import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { del, get, getPaginated, patch, post, put } from '@/lib/api';
import type {
  CreateSubscriptionRequest,
  PaginatedResponse,
  Subscription,
  UpdateSubscriptionRequest,
} from '@/types';

interface SubscriptionFilters {
  status?: string;
  categoryId?: string;
  folderId?: string;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  page?: number;
  perPage?: number;
}

export function useSubscriptions(filters?: SubscriptionFilters) {
  return useQuery<PaginatedResponse<Subscription>>({
    queryKey: ['subscriptions', filters],
    queryFn: () => getPaginated<Subscription>('/subscriptions', filters as Record<string, unknown>),
  });
}

export function useSubscription(id: string) {
  return useQuery<Subscription>({
    queryKey: ['subscriptions', id],
    queryFn: () => get<Subscription>(`/subscriptions/${id}`),
    enabled: !!id,
  });
}

export function useCreateSubscription() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateSubscriptionRequest) => post<Subscription>('/subscriptions', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}

export function useUpdateSubscription() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateSubscriptionRequest }) =>
      put<Subscription>(`/subscriptions/${id}`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['subscriptions', variables.id] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}

export function useDeleteSubscription() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => del(`/subscriptions/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}

export function useUpdateSatisfaction() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, score }: { id: string; score: number }) =>
      patch<Subscription>(`/subscriptions/${id}/satisfaction`, { score }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['subscriptions', variables.id] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}
