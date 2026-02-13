import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

import { get, post, put, del } from '@/lib/api';
import type { Folder, CreateFolderRequest, UpdateFolderRequest } from '@/types';

export function useFolders() {
  return useQuery<Folder[]>({
    queryKey: ['folders'],
    queryFn: () => get<Folder[]>('/folders'),
  });
}

export function useCreateFolder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateFolderRequest) =>
      post<Folder>('/folders', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['folders'] });
    },
  });
}

export function useUpdateFolder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateFolderRequest }) =>
      put<Folder>(`/folders/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['folders'] });
    },
  });
}

export function useDeleteFolder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => del(`/folders/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['folders'] });
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
    },
  });
}
