'use client';

import { useState } from 'react';

import { useDeleteSubscription, useUpdateSubscription } from '@/lib/hooks/useSubscriptions';
import { billingCycleLabel, formatCurrency, formatDate, satisfactionStars } from '@/lib/utils';
import type { Subscription } from '@/types';

import DeleteConfirmModal from './DeleteConfirmModal';

interface SubscriptionCardProps {
  subscription: Subscription;
  onEdit: (subscription: Subscription) => void;
}

export default function SubscriptionCard({ subscription, onEdit }: SubscriptionCardProps) {
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [showCancelModal, setShowCancelModal] = useState(false);
  const deleteMutation = useDeleteSubscription();
  const updateMutation = useUpdateSubscription();

  const handleDelete = () => {
    deleteMutation.mutate(subscription.id, {
      onSuccess: () => {
        setShowDeleteModal(false);
      },
    });
  };

  const handleToggleStatus = () => {
    const newStatus = subscription.status === 'active' ? 'paused' : 'active';
    updateMutation.mutate({
      id: subscription.id,
      data: { status: newStatus },
    });
  };

  const handleCancel = () => {
    updateMutation.mutate(
      { id: subscription.id, data: { status: 'cancelled' } },
      {
        onSuccess: () => setShowCancelModal(false),
      }
    );
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
                <span className="rounded-full bg-yellow-100 px-2 py-1 text-xs font-medium text-yellow-700">
                  일시중지
                </span>
              )}
              {subscription.status === 'cancelled' && (
                <span className="rounded-full bg-red-100 px-2 py-1 text-xs font-medium text-red-700">
                  해지
                </span>
              )}
              {subscription.isTrial && (
                <span className="rounded-full bg-purple-100 px-2 py-1 text-xs font-medium text-purple-700">
                  체험
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
              {subscription.folder && (
                <span className="rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-700">
                  📁 {subscription.folder.name}
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
            {subscription.status !== 'cancelled' && (
              <button
                onClick={handleToggleStatus}
                disabled={updateMutation.isPending}
                className={`rounded-lg px-4 py-2 text-sm font-medium text-white transition-colors disabled:opacity-50 ${
                  subscription.status === 'active'
                    ? 'bg-yellow-500 hover:bg-yellow-600'
                    : 'bg-green-500 hover:bg-green-600'
                }`}
              >
                {subscription.status === 'active' ? '일시중지' : '활성'}
              </button>
            )}
            {subscription.status !== 'cancelled' && (
              <button
                onClick={() => setShowCancelModal(true)}
                className="rounded-lg bg-orange-500 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-orange-600"
              >
                해지
              </button>
            )}
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

      {/* Cancel Confirm Modal */}
      {showCancelModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4">
          <div className="w-full max-w-md rounded-lg bg-white p-6 shadow-xl">
            <h3 className="text-lg font-semibold text-gray-900">구독 해지</h3>
            <p className="mt-2 text-sm text-gray-600">
              &quot;{subscription.serviceName}&quot; 구독을 해지하시겠습니까?
              해지된 구독은 목록에서 숨겨지며, 해지 상태 필터로 확인할 수 있습니다.
            </p>
            <div className="mt-4 flex gap-3">
              <button
                onClick={() => setShowCancelModal(false)}
                className="flex-1 rounded-lg border-2 border-gray-300 px-4 py-2 font-medium text-gray-700 hover:bg-gray-50"
              >
                취소
              </button>
              <button
                onClick={handleCancel}
                disabled={updateMutation.isPending}
                className="flex-1 rounded-lg bg-orange-500 px-4 py-2 font-medium text-white hover:bg-orange-600 disabled:opacity-50"
              >
                {updateMutation.isPending ? '처리 중...' : '해지'}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
