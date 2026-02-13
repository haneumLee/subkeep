'use client';

import { useState } from 'react';

import { useDeleteSubscription } from '@/lib/hooks/useSubscriptions';
import { billingCycleLabel, formatCurrency, formatDate, satisfactionStars } from '@/lib/utils';
import type { Subscription } from '@/types';

import DeleteConfirmModal from './DeleteConfirmModal';

interface SubscriptionCardProps {
  subscription: Subscription;
  onEdit: (subscription: Subscription) => void;
}

export default function SubscriptionCard({ subscription, onEdit }: SubscriptionCardProps) {
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const deleteMutation = useDeleteSubscription();

  const handleDelete = () => {
    deleteMutation.mutate(subscription.id, {
      onSuccess: () => {
        setShowDeleteModal(false);
      },
    });
  };

  return (
    <>
      <div className="rounded-lg border-2 border-gray-200 bg-white p-6 transition-all hover:border-blue-300 hover:shadow-md">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2">
              <h3 className="text-lg font-semibold text-gray-900">
                {subscription.serviceName}
              </h3>
              {subscription.status === 'paused' && (
                <span className="rounded-full bg-gray-200 px-2 py-1 text-xs font-medium text-gray-600">
                  일시중지
                </span>
              )}
            </div>

            <div className="mt-2 flex items-baseline gap-2">
              <span className="text-2xl font-bold text-blue-600">
                {formatCurrency(subscription.monthlyAmount)}
              </span>
              <span className="text-sm text-gray-500">/ 월</span>
            </div>

            <div className="mt-3 space-y-1 text-sm text-gray-600">
              <div className="flex items-center gap-2">
                <span>결제주기:</span>
                <span className="font-medium">{billingCycleLabel(subscription.billingCycle)}</span>
              </div>
              <div className="flex items-center gap-2">
                <span>다음 결제일:</span>
                <span className="font-medium">{formatDate(subscription.nextBillingDate)}</span>
              </div>
              {subscription.satisfactionScore !== null && (
                <div className="flex items-center gap-2">
                  <span>만족도:</span>
                  <span className="text-yellow-500">
                    {satisfactionStars(subscription.satisfactionScore)}
                  </span>
                </div>
              )}
            </div>

            <div className="mt-3 flex flex-wrap gap-2">
              {subscription.category && (
                <span
                  className="rounded-full px-3 py-1 text-xs font-medium text-white"
                  style={{
                    backgroundColor: subscription.category.color || '#6B7280',
                  }}
                >
                  {subscription.category.name}
                </span>
              )}
              {!subscription.autoRenew && (
                <span className="rounded-full bg-yellow-100 px-3 py-1 text-xs font-medium text-yellow-800">
                  자동갱신 없음
                </span>
              )}
            </div>
          </div>

          <div className="ml-4 flex flex-col gap-2">
            <button
              onClick={() => onEdit(subscription)}
              className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
            >
              수정
            </button>
            <button
              onClick={() => setShowDeleteModal(true)}
              className="rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700"
            >
              삭제
            </button>
          </div>
        </div>

        {subscription.note && (
          <div className="mt-4 border-t border-gray-200 pt-3">
            <p className="text-sm text-gray-600">{subscription.note}</p>
          </div>
        )}
      </div>

      <DeleteConfirmModal
        subscription={subscription}
        isOpen={showDeleteModal}
        onClose={() => setShowDeleteModal(false)}
        onConfirm={handleDelete}
        isDeleting={deleteMutation.isPending}
      />
    </>
  );
}
