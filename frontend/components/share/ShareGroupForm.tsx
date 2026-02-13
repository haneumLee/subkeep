'use client';

import { useEffect } from 'react';
import { useForm, useFieldArray } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

import { Modal } from '@/components/ui/Modal';
import { Button } from '@/components/ui/Button';
import { useCreateShareGroup, useUpdateShareGroup } from '@/lib/hooks/useShareGroups';
import type { ShareGroup } from '@/types';

const memberSchema = z.object({
  nickname: z.string().min(1, '닉네임을 입력해주세요').max(50, '닉네임은 50자 이하여야 합니다'),
  role: z.string().max(30, '역할은 30자 이하여야 합니다').optional(),
});

const shareGroupSchema = z.object({
  name: z
    .string()
    .min(1, '그룹명을 입력해주세요')
    .max(100, '그룹명은 100자 이하여야 합니다'),
  description: z.string().max(500, '설명은 500자 이하여야 합니다').optional(),
  members: z
    .array(memberSchema)
    .min(2, '최소 2명의 멤버가 필요합니다')
    .max(20, '최대 20명까지 추가할 수 있습니다'),
});

type ShareGroupFormData = z.infer<typeof shareGroupSchema>;

interface ShareGroupFormProps {
  mode: 'create' | 'edit';
  group?: ShareGroup;
  isOpen: boolean;
  onClose: () => void;
}

export function ShareGroupForm({ mode, group, isOpen, onClose }: ShareGroupFormProps) {
  const createMutation = useCreateShareGroup();
  const updateMutation = useUpdateShareGroup();

  const {
    register,
    control,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ShareGroupFormData>({
    resolver: zodResolver(shareGroupSchema),
    defaultValues: {
      name: '',
      description: '',
      members: [
        { nickname: '', role: '' },
        { nickname: '', role: '' },
      ],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'members',
  });

  useEffect(() => {
    if (mode === 'edit' && group) {
      reset({
        name: group.name,
        description: group.description || '',
        members:
          group.members?.map((m) => ({
            nickname: m.nickname,
            role: m.role || '',
          })) || [],
      });
    }
  }, [mode, group, reset]);

  const onSubmit = async (data: ShareGroupFormData) => {
    try {
      const payload = {
        name: data.name,
        description: data.description || undefined,
        members: data.members.map((m) => ({
          nickname: m.nickname,
          role: m.role || undefined,
        })),
      };

      if (mode === 'create') {
        await createMutation.mutateAsync(payload);
      } else if (group) {
        await updateMutation.mutateAsync({ id: group.id, data: payload });
      }

      onClose();
      reset();
    } catch (error) {
      console.error('Failed to save share group:', error);
    }
  };

  const handleClose = () => {
    onClose();
    reset();
  };

  const isLoading = createMutation.isPending || updateMutation.isPending;

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      title={mode === 'create' ? '공유 그룹 생성' : '공유 그룹 수정'}
      showFooter={false}
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {/* 그룹명 */}
        <div>
          <label htmlFor="name" className="mb-1 block text-sm font-medium text-slate-700">
            그룹명 <span className="text-red-500">*</span>
          </label>
          <input
            {...register('name')}
            id="name"
            type="text"
            placeholder="예: 넷플릭스 가족 공유"
            className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
          {errors.name && <p className="mt-1 text-xs text-red-600">{errors.name.message}</p>}
        </div>

        {/* 설명 */}
        <div>
          <label
            htmlFor="description"
            className="mb-1 block text-sm font-medium text-slate-700"
          >
            설명
          </label>
          <textarea
            {...register('description')}
            id="description"
            rows={2}
            placeholder="그룹에 대한 설명을 입력하세요 (선택)"
            className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
          {errors.description && (
            <p className="mt-1 text-xs text-red-600">{errors.description.message}</p>
          )}
        </div>

        {/* 멤버 목록 */}
        <div>
          <div className="mb-2 flex items-center justify-between">
            <label className="block text-sm font-medium text-slate-700">
              멤버 <span className="text-red-500">*</span>
              <span className="ml-1 text-xs font-normal text-slate-500">(최소 2명)</span>
            </label>
            <Button
              type="button"
              variant="secondary"
              size="sm"
              onClick={() => append({ nickname: '', role: '' })}
              disabled={fields.length >= 20}
            >
              + 멤버 추가
            </Button>
          </div>

          <div className="space-y-2">
            {fields.map((field, index) => (
              <div key={field.id} className="flex gap-2">
                <div className="flex-1">
                  <input
                    {...register(`members.${index}.nickname`)}
                    type="text"
                    placeholder="닉네임"
                    className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                  {errors.members?.[index]?.nickname && (
                    <p className="mt-1 text-xs text-red-600">
                      {errors.members[index]?.nickname?.message}
                    </p>
                  )}
                </div>
                <div className="flex-1">
                  <input
                    {...register(`members.${index}.role`)}
                    type="text"
                    placeholder="역할 (선택)"
                    className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                  {errors.members?.[index]?.role && (
                    <p className="mt-1 text-xs text-red-600">
                      {errors.members[index]?.role?.message}
                    </p>
                  )}
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => remove(index)}
                  disabled={fields.length <= 2}
                  className="px-2"
                >
                  <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                    />
                  </svg>
                </Button>
              </div>
            ))}
          </div>
          {errors.members && (
            <p className="mt-1 text-xs text-red-600">{errors.members.message}</p>
          )}
        </div>

        {/* 버튼 */}
        <div className="flex justify-end gap-3 pt-2">
          <Button type="button" variant="secondary" onClick={handleClose} disabled={isLoading}>
            취소
          </Button>
          <Button type="submit" variant="primary" isLoading={isLoading}>
            {mode === 'create' ? '생성' : '수정'}
          </Button>
        </div>
      </form>
    </Modal>
  );
}
