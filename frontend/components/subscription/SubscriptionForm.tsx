'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { useCategories } from '@/lib/hooks/useCategories';
import {
  useCreateSubscription,
  useUpdateSubscription,
} from '@/lib/hooks/useSubscriptions';
import type { BillingCycle, CreateSubscriptionRequest, Subscription } from '@/types';

import SatisfactionInput from './SatisfactionInput';

const subscriptionSchema = z.object({
  serviceName: z
    .string()
    .min(1, '서비스명은 필수입니다')
    .max(50, '서비스명은 50자 이내여야 합니다'),
  categoryId: z.string().optional(),
  amount: z.number().min(0, '금액은 0 이상이어야 합니다'),
  billingCycle: z.enum(['weekly', 'monthly', 'yearly']),
  nextBillingDate: z.string().min(1, '다음 결제일은 필수입니다'),
  startDate: z.string().min(1, '시작일은 필수입니다'),
  autoRenew: z.boolean(),
  satisfactionScore: z.number().min(1).max(5).optional(),
  note: z.string().max(500, '메모는 500자 이내여야 합니다').optional(),
  serviceUrl: z.string().url('올바른 URL 형식이 아닙니다').optional().or(z.literal('')),
});

type SubscriptionFormData = z.infer<typeof subscriptionSchema>;

interface SubscriptionFormProps {
  isOpen: boolean;
  onClose: () => void;
  subscription?: Subscription;
}

export default function SubscriptionForm({ isOpen, onClose, subscription }: SubscriptionFormProps) {
  const [showDuplicateWarning, setShowDuplicateWarning] = useState(false);
  const [showHighAmountWarning, setShowHighAmountWarning] = useState(false);
  const [satisfactionScore, setSatisfactionScore] = useState(3);

  const { data: categories } = useCategories();
  const createMutation = useCreateSubscription();
  const updateMutation = useUpdateSubscription();

  const isEditing = !!subscription;

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    watch,
    setValue,
  } = useForm<SubscriptionFormData>({
    resolver: zodResolver(subscriptionSchema),
    defaultValues: {
      serviceName: '',
      categoryId: '',
      amount: 0,
      billingCycle: 'monthly',
      nextBillingDate: new Date().toISOString().split('T')[0],
      startDate: new Date().toISOString().split('T')[0],
      autoRenew: true,
      note: '',
      serviceUrl: '',
    },
  });

  const watchedAmount = watch('amount');

  useEffect(() => {
    if (subscription) {
      reset({
        serviceName: subscription.serviceName,
        categoryId: subscription.categoryId || '',
        amount: subscription.amount,
        billingCycle: subscription.billingCycle,
        nextBillingDate: subscription.nextBillingDate.split('T')[0],
        startDate: subscription.startDate.split('T')[0],
        autoRenew: subscription.autoRenew,
        note: subscription.note || '',
        serviceUrl: subscription.serviceUrl || '',
      });
      setSatisfactionScore(subscription.satisfactionScore || 3);
    } else {
      reset({
        serviceName: '',
        categoryId: '',
        amount: 0,
        billingCycle: 'monthly',
        nextBillingDate: new Date().toISOString().split('T')[0],
        startDate: new Date().toISOString().split('T')[0],
        autoRenew: true,
        note: '',
        serviceUrl: '',
      });
      setSatisfactionScore(3);
    }
  }, [subscription, reset]);

  const onSubmit = (data: SubscriptionFormData) => {
    // Check for high amount
    if (data.amount > 1000000 && !showHighAmountWarning) {
      setShowHighAmountWarning(true);
      return;
    }

    const requestData: CreateSubscriptionRequest = {
      ...data,
      categoryId: data.categoryId || undefined,
      satisfactionScore,
      note: data.note || undefined,
      serviceUrl: data.serviceUrl || undefined,
    };

    if (isEditing && subscription) {
      updateMutation.mutate(
        { id: subscription.id, data: requestData },
        {
          onSuccess: () => {
            onClose();
            reset();
          },
        }
      );
    } else {
      createMutation.mutate(requestData, {
        onSuccess: () => {
          onClose();
          reset();
        },
        onError: (error: Error) => {
          if (error.message.includes('duplicate') || error.message.includes('중복')) {
            setShowDuplicateWarning(true);
          }
        },
      });
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4">
      <div className="max-h-[90vh] w-full max-w-2xl overflow-y-auto rounded-lg bg-white p-6 shadow-xl">
        <h2 className="mb-6 text-2xl font-bold text-gray-900">
          {isEditing ? '구독 수정' : '구독 추가'}
        </h2>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {/* Service Name */}
          <div>
            <label className="block text-sm font-medium text-gray-700">
              서비스명 <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              {...register('serviceName')}
              className="mt-1 w-full rounded-lg border-2 border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
              placeholder="Netflix, Spotify 등"
            />
            {errors.serviceName && (
              <p className="mt-1 text-sm text-red-600">{errors.serviceName.message}</p>
            )}
          </div>

          {/* Category */}
          <div>
            <label className="block text-sm font-medium text-gray-700">카테고리</label>
            <select
              {...register('categoryId')}
              className="mt-1 w-full rounded-lg border-2 border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
            >
              <option value="">선택 안 함</option>
              {categories?.map((category) => (
                <option key={category.id} value={category.id}>
                  {category.name}
                </option>
              ))}
            </select>
          </div>

          {/* Amount & Billing Cycle */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">
                결제 금액 <span className="text-red-500">*</span>
              </label>
              <input
                type="number"
                {...register('amount', { valueAsNumber: true })}
                className="mt-1 w-full rounded-lg border-2 border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
                placeholder="0"
                min="0"
              />
              {errors.amount && (
                <p className="mt-1 text-sm text-red-600">{errors.amount.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700">
                결제 주기 <span className="text-red-500">*</span>
              </label>
              <select
                {...register('billingCycle')}
                className="mt-1 w-full rounded-lg border-2 border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
              >
                <option value="weekly">주간</option>
                <option value="monthly">월간</option>
                <option value="yearly">연간</option>
              </select>
            </div>
          </div>

          {/* Dates */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">
                시작일 <span className="text-red-500">*</span>
              </label>
              <input
                type="date"
                {...register('startDate')}
                className="mt-1 w-full rounded-lg border-2 border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
              />
              {errors.startDate && (
                <p className="mt-1 text-sm text-red-600">{errors.startDate.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700">
                다음 결제일 <span className="text-red-500">*</span>
              </label>
              <input
                type="date"
                {...register('nextBillingDate')}
                className="mt-1 w-full rounded-lg border-2 border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
              />
              {errors.nextBillingDate && (
                <p className="mt-1 text-sm text-red-600">{errors.nextBillingDate.message}</p>
              )}
            </div>
          </div>

          {/* Auto Renew */}
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              {...register('autoRenew')}
              id="autoRenew"
              className="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <label htmlFor="autoRenew" className="text-sm font-medium text-gray-700">
              자동 갱신
            </label>
          </div>

          {/* Satisfaction Score */}
          <div>
            <label className="block text-sm font-medium text-gray-700">만족도</label>
            <div className="mt-1">
              <SatisfactionInput
                value={satisfactionScore}
                onChange={setSatisfactionScore}
              />
            </div>
          </div>

          {/* Service URL */}
          <div>
            <label className="block text-sm font-medium text-gray-700">서비스 URL</label>
            <input
              type="url"
              {...register('serviceUrl')}
              className="mt-1 w-full rounded-lg border-2 border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
              placeholder="https://example.com"
            />
            {errors.serviceUrl && (
              <p className="mt-1 text-sm text-red-600">{errors.serviceUrl.message}</p>
            )}
          </div>

          {/* Note */}
          <div>
            <label className="block text-sm font-medium text-gray-700">메모</label>
            <textarea
              {...register('note')}
              className="mt-1 w-full rounded-lg border-2 border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
              rows={3}
              placeholder="구독 관련 메모 (최대 500자)"
              maxLength={500}
            />
            {errors.note && (
              <p className="mt-1 text-sm text-red-600">{errors.note.message}</p>
            )}
          </div>

          {/* Action Buttons */}
          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={() => {
                onClose();
                reset();
              }}
              className="flex-1 rounded-lg border-2 border-gray-300 px-4 py-2 font-medium text-gray-700 transition-colors hover:bg-gray-50"
            >
              취소
            </button>
            <button
              type="submit"
              disabled={createMutation.isPending || updateMutation.isPending}
              className="flex-1 rounded-lg bg-blue-600 px-4 py-2 font-medium text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {createMutation.isPending || updateMutation.isPending
                ? '저장 중...'
                : isEditing
                  ? '수정'
                  : '추가'}
            </button>
          </div>
        </form>
      </div>

      {/* Duplicate Warning Modal */}
      {showDuplicateWarning && (
        <div className="fixed inset-0 z-[60] flex items-center justify-center bg-black bg-opacity-50 p-4">
          <div className="w-full max-w-md rounded-lg bg-white p-6 shadow-xl">
            <h3 className="text-lg font-semibold text-gray-900">중복된 서비스명</h3>
            <p className="mt-2 text-sm text-gray-600">
              이미 동일한 서비스명의 구독이 존재합니다. 계속 진행하시겠습니까?
            </p>
            <div className="mt-4 flex gap-3">
              <button
                onClick={() => setShowDuplicateWarning(false)}
                className="flex-1 rounded-lg border-2 border-gray-300 px-4 py-2 font-medium text-gray-700 transition-colors hover:bg-gray-50"
              >
                취소
              </button>
              <button
                onClick={() => {
                  setShowDuplicateWarning(false);
                  onClose();
                }}
                className="flex-1 rounded-lg bg-blue-600 px-4 py-2 font-medium text-white transition-colors hover:bg-blue-700"
              >
                계속 진행
              </button>
            </div>
          </div>
        </div>
      )}

      {/* High Amount Warning Modal */}
      {showHighAmountWarning && (
        <div className="fixed inset-0 z-[60] flex items-center justify-center bg-black bg-opacity-50 p-4">
          <div className="w-full max-w-md rounded-lg bg-white p-6 shadow-xl">
            <h3 className="text-lg font-semibold text-gray-900">높은 금액 감지</h3>
            <p className="mt-2 text-sm text-gray-600">
              입력하신 금액이 100만원을 초과합니다. 금액이 정확한지 확인해주세요.
            </p>
            <div className="mt-4 flex gap-3">
              <button
                onClick={() => setShowHighAmountWarning(false)}
                className="flex-1 rounded-lg border-2 border-gray-300 px-4 py-2 font-medium text-gray-700 transition-colors hover:bg-gray-50"
              >
                취소
              </button>
              <button
                onClick={() => {
                  setShowHighAmountWarning(false);
                  handleSubmit(onSubmit)();
                }}
                className="flex-1 rounded-lg bg-blue-600 px-4 py-2 font-medium text-white transition-colors hover:bg-blue-700"
              >
                계속 진행
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
