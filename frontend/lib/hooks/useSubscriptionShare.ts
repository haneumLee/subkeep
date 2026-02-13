import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { del, get, post, put } from '@/lib/api';
import type { LinkShareRequest, SubscriptionShare } from '@/types';

export function useLinkShare() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ subscriptionId, data }: { subscriptionId: string; data: LinkShareRequest }) =>
      post<SubscriptionShare>(`/subscriptions/${subscriptionId}/share`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['subscriptions', variables.subscriptionId] });
      queryClient.invalidateQueries({ queryKey: ['subscription-share', variables.subscriptionId] });
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}

export function useSubscriptionShare(subscriptionId: string) {
  return useQuery<SubscriptionShare>({
    queryKey: ['subscription-share', subscriptionId],
    queryFn: () => get<SubscriptionShare>(`/subscriptions/${subscriptionId}/share`),
    enabled: !!subscriptionId,
  });
}

export function useUpdateShare() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<LinkShareRequest> }) =>
      put<SubscriptionShare>(`/subscription-shares/${id}`, data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['subscription-share', data.subscriptionId] });
      queryClient.invalidateQueries({ queryKey: ['subscriptions', data.subscriptionId] });
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}

export function useUnlinkShare() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, subscriptionId }: { id: string; subscriptionId: string }) =>
      del(`/subscription-shares/${id}`),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['subscription-share', variables.subscriptionId] });
      queryClient.invalidateQueries({ queryKey: ['subscriptions', variables.subscriptionId] });
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}
