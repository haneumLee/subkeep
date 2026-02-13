import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { del, get, post, put } from '@/lib/api';
import type { CreateShareGroupRequest, ShareGroup } from '@/types';

export function useShareGroups() {
  return useQuery<ShareGroup[]>({
    queryKey: ['share-groups'],
    queryFn: () => get<ShareGroup[]>('/share-groups'),
  });
}

export function useShareGroup(id: string) {
  return useQuery<ShareGroup>({
    queryKey: ['share-groups', id],
    queryFn: () => get<ShareGroup>(`/share-groups/${id}`),
    enabled: !!id,
  });
}

export function useCreateShareGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateShareGroupRequest) => post<ShareGroup>('/share-groups', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['share-groups'] });
    },
  });
}

export function useUpdateShareGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateShareGroupRequest }) =>
      put<ShareGroup>(`/share-groups/${id}`, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['share-groups'] });
      queryClient.invalidateQueries({ queryKey: ['share-groups', variables.id] });
    },
  });
}

export function useDeleteShareGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => del(`/share-groups/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['share-groups'] });
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
    },
  });
}
