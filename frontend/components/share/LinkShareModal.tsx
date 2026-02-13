'use client';

import { useEffect, useMemo, useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

import { Modal } from '@/components/ui/Modal';
import { Button } from '@/components/ui/Button';
import { useLinkShare, useSubscriptionShare, useUpdateShare } from '@/lib/hooks/useSubscriptionShare';
import { useShareGroups } from '@/lib/hooks/useShareGroups';
import { formatCurrency } from '@/lib/utils';
import type { SplitType, Subscription } from '@/types';

const linkShareSchema = z.object({
  shareGroupId: z.string().min(1, '공유 그룹을 선택해주세요'),
  splitType: z.enum(['equal', 'custom_amount', 'custom_ratio'], {
    required_error: '분담 방식을 선택해주세요',
  }),
  myShareAmount: z.number().optional(),
  myShareRatio: z.number().optional(),
});

type LinkShareFormData = z.infer<typeof linkShareSchema>;

interface LinkShareModalProps {
  subscription: Subscription;
  isOpen: boolean;
  onClose: () => void;
}

export function LinkShareModal({ subscription, isOpen, onClose }: LinkShareModalProps) {
  const { data: shareGroups } = useShareGroups();
  const { data: existingShare } = useSubscriptionShare(subscription.id);
  const linkMutation = useLinkShare();
  const updateMutation = useUpdateShare();

  const [totalMembers, setTotalMembers] = useState(2);

  const isEditMode = !!existingShare;

  const {
    register,
    watch,
    setValue,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<LinkShareFormData>({
    resolver: zodResolver(linkShareSchema),
    defaultValues: {
      shareGroupId: '',
      splitType: 'equal',
      myShareAmount: undefined,
      myShareRatio: undefined,
    },
  });

  const selectedGroupId = watch('shareGroupId');
  const splitType = watch('splitType');
  const myShareAmount = watch('myShareAmount');
  const myShareRatio = watch('myShareRatio');

  const selectedGroup = useMemo(
    () => shareGroups?.find((g) => g.id === selectedGroupId),
    [shareGroups, selectedGroupId]
  );

  useEffect(() => {
    if (selectedGroup?.members) {
      setTotalMembers(selectedGroup.members.length);
    }
  }, [selectedGroup]);

  useEffect(() => {
    if (isEditMode && existingShare) {
      reset({
        shareGroupId: existingShare.shareGroupId,
        splitType: existingShare.splitType,
        myShareAmount: existingShare.myShareAmount || undefined,
        myShareRatio: existingShare.myShareRatio || undefined,
      });
      setTotalMembers(existingShare.totalMembersSnapshot);
    }
  }, [isEditMode, existingShare, reset]);

  const calculatedShare = useMemo(() => {
    if (totalMembers < 1) {
      return { myShare: 0, error: '분담 인원은 1명 이상이어야 합니다' };
    }

    switch (splitType) {
      case 'equal':
        return {
          myShare: Math.round(subscription.amount / totalMembers),
          error: null,
        };

      case 'custom_amount':
        if (!myShareAmount) {
          return { myShare: 0, error: null };
        }
        if (myShareAmount < 0) {
          return { myShare: 0, error: '분담금은 0원 이상이어야 합니다' };
        }
        if (myShareAmount > subscription.amount * 1.05) {
          return {
            myShare: myShareAmount,
            error: '분담금이 구독 금액의 105%를 초과합니다',
          };
        }
        return { myShare: myShareAmount, error: null };

      case 'custom_ratio':
        if (!myShareRatio) {
          return { myShare: 0, error: null };
        }
        if (myShareRatio < 0 || myShareRatio > 100) {
          return { myShare: 0, error: '비율은 0~100% 사이여야 합니다' };
        }
        return {
          myShare: Math.round((subscription.amount * myShareRatio) / 100),
          error: null,
        };

      default:
        return { myShare: 0, error: null };
    }
  }, [splitType, subscription.amount, totalMembers, myShareAmount, myShareRatio]);

  const onSubmit = async (data: LinkShareFormData) => {
    if (calculatedShare.error) {
      return;
    }

    try {
      const payload = {
        shareGroupId: data.shareGroupId,
        splitType: data.splitType,
        totalMembersSnapshot: totalMembers,
        myShareAmount: data.splitType === 'custom_amount' ? data.myShareAmount : undefined,
        myShareRatio: data.splitType === 'custom_ratio' ? data.myShareRatio : undefined,
      };

      if (isEditMode && existingShare) {
        await updateMutation.mutateAsync({ id: existingShare.id, data: payload });
      } else {
        await linkMutation.mutateAsync({ subscriptionId: subscription.id, data: payload });
      }

      onClose();
      reset();
    } catch (error) {
      console.error('Failed to link share:', error);
    }
  };

  const handleClose = () => {
    onClose();
    reset();
  };

  const isLoading = linkMutation.isPending || updateMutation.isPending;

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      title={isEditMode ? '공유 설정 수정' : '공유 연결'}
      showFooter={false}
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {/* 구독 정보 */}
        <div className="rounded-lg bg-slate-50 p-3">
          <div className="text-sm font-medium text-slate-900">{subscription.serviceName}</div>
          <div className="mt-1 text-lg font-semibold text-primary-600">
            {formatCurrency(subscription.amount)} / {subscription.billingCycle === 'monthly' ? '월' : subscription.billingCycle === 'yearly' ? '년' : '주'}
          </div>
        </div>

        {/* 공유 그룹 선택 */}
        <div>
          <label htmlFor="shareGroupId" className="mb-1 block text-sm font-medium text-slate-700">
            공유 그룹 <span className="text-red-500">*</span>
          </label>
          <select
            {...register('shareGroupId')}
            id="shareGroupId"
            disabled={isEditMode}
            className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500 disabled:bg-slate-100 disabled:text-slate-500"
          >
            <option value="">선택하세요</option>
            {shareGroups?.map((group) => (
              <option key={group.id} value={group.id}>
                {group.name} ({group.members?.length || 0}명)
              </option>
            ))}
          </select>
          {errors.shareGroupId && (
            <p className="mt-1 text-xs text-red-600">{errors.shareGroupId.message}</p>
          )}
        </div>

        {/* 분담 인원 */}
        {selectedGroup && (
          <div>
            <label htmlFor="totalMembers" className="mb-1 block text-sm font-medium text-slate-700">
              분담 인원
            </label>
            <input
              type="number"
              id="totalMembers"
              value={totalMembers}
              onChange={(e) => setTotalMembers(Number(e.target.value))}
              min={1}
              max={selectedGroup.members?.length || 1}
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
            <p className="mt-1 text-xs text-slate-500">
              그룹 전체 인원: {selectedGroup.members?.length || 0}명
            </p>
          </div>
        )}

        {/* 분담 방식 */}
        <div>
          <label className="mb-2 block text-sm font-medium text-slate-700">
            분담 방식 <span className="text-red-500">*</span>
          </label>
          <div className="space-y-2">
            <label className="flex cursor-pointer items-center rounded-lg border border-slate-300 p-3 hover:bg-slate-50">
              <input
                {...register('splitType')}
                type="radio"
                value="equal"
                className="mr-3 h-4 w-4 text-primary-600 focus:ring-primary-500"
              />
              <div className="flex-1">
                <div className="text-sm font-medium text-slate-900">균등 분담</div>
                <div className="text-xs text-slate-500">총액을 인원수로 균등하게 나눕니다</div>
              </div>
            </label>

            <label className="flex cursor-pointer items-center rounded-lg border border-slate-300 p-3 hover:bg-slate-50">
              <input
                {...register('splitType')}
                type="radio"
                value="custom_amount"
                className="mr-3 h-4 w-4 text-primary-600 focus:ring-primary-500"
              />
              <div className="flex-1">
                <div className="text-sm font-medium text-slate-900">커스텀 (금액)</div>
                <div className="text-xs text-slate-500">내 분담금을 직접 입력합니다</div>
              </div>
            </label>

            <label className="flex cursor-pointer items-center rounded-lg border border-slate-300 p-3 hover:bg-slate-50">
              <input
                {...register('splitType')}
                type="radio"
                value="custom_ratio"
                className="mr-3 h-4 w-4 text-primary-600 focus:ring-primary-500"
              />
              <div className="flex-1">
                <div className="text-sm font-medium text-slate-900">커스텀 (비율)</div>
                <div className="text-xs text-slate-500">내 분담 비율(%)을 입력합니다</div>
              </div>
            </label>
          </div>
          {errors.splitType && (
            <p className="mt-1 text-xs text-red-600">{errors.splitType.message}</p>
          )}
        </div>

        {/* 커스텀 금액 입력 */}
        {splitType === 'custom_amount' && (
          <div>
            <label htmlFor="myShareAmount" className="mb-1 block text-sm font-medium text-slate-700">
              내 분담금 <span className="text-red-500">*</span>
            </label>
            <input
              {...register('myShareAmount', { valueAsNumber: true })}
              type="number"
              id="myShareAmount"
              placeholder="0"
              min={0}
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
            {errors.myShareAmount && (
              <p className="mt-1 text-xs text-red-600">{errors.myShareAmount.message}</p>
            )}
          </div>
        )}

        {/* 커스텀 비율 입력 */}
        {splitType === 'custom_ratio' && (
          <div>
            <label htmlFor="myShareRatio" className="mb-1 block text-sm font-medium text-slate-700">
              내 분담 비율 (%) <span className="text-red-500">*</span>
            </label>
            <input
              {...register('myShareRatio', { valueAsNumber: true })}
              type="number"
              id="myShareRatio"
              placeholder="0"
              min={0}
              max={100}
              step={0.1}
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
            {errors.myShareRatio && (
              <p className="mt-1 text-xs text-red-600">{errors.myShareRatio.message}</p>
            )}
          </div>
        )}

        {/* 계산 결과 미리보기 */}
        <div className="rounded-lg bg-primary-50 p-4">
          <div className="mb-2 text-sm font-medium text-slate-700">미리보기</div>
          {calculatedShare.error ? (
            <div className="flex items-center text-sm text-red-600">
              <svg className="mr-2 h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  clipRule="evenodd"
                />
              </svg>
              {calculatedShare.error}
            </div>
          ) : (
            <div className="space-y-1">
              <div className="flex justify-between text-sm">
                <span className="text-slate-600">총 구독 금액</span>
                <span className="font-semibold text-slate-900">
                  {formatCurrency(subscription.amount)}
                </span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-slate-600">분담 인원</span>
                <span className="font-semibold text-slate-900">{totalMembers}명</span>
              </div>
              <div className="mt-2 flex justify-between border-t border-primary-200 pt-2 text-base">
                <span className="font-medium text-slate-900">내 부담액</span>
                <span className="font-bold text-primary-600">
                  {formatCurrency(calculatedShare.myShare)}
                </span>
              </div>
            </div>
          )}
        </div>

        {/* 버튼 */}
        <div className="flex justify-end gap-3 pt-2">
          <Button type="button" variant="secondary" onClick={handleClose} disabled={isLoading}>
            취소
          </Button>
          <Button
            type="submit"
            variant="primary"
            isLoading={isLoading}
            disabled={!!calculatedShare.error}
          >
            {isEditMode ? '수정' : '연결'}
          </Button>
        </div>
      </form>
    </Modal>
  );
}
