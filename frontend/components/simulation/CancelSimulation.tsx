'use client';

import { useEffect, useState } from 'react';

import type { SimulationResult, Subscription } from '@/types';
import { useCancelSimulation, useApplySimulation } from '@/lib/hooks/useSimulation';
import { useSubscriptions } from '@/lib/hooks/useSubscriptions';
import { formatCurrency, billingCycleLabel, satisfactionStars, cn } from '@/lib/utils';
import { SimulationResult as SimulationResultComponent } from './SimulationResult';
import { UndoButton } from './UndoButton';

export function CancelSimulation() {
  const { data: subscriptionsData, isLoading: subscriptionsLoading } = useSubscriptions({
    status: 'active',
  });
  const cancelSimulation = useCancelSimulation();
  const applySimulation = useApplySimulation();

  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [simulationResult, setSimulationResult] = useState<SimulationResult | null>(null);
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [showUndoButton, setShowUndoButton] = useState(false);

  // 만족도 낮은 순으로 정렬
  const sortedSubscriptions =
    subscriptionsData?.data.sort((a: Subscription, b: Subscription) => {
      const scoreA = a.satisfactionScore ?? 6; // null은 맨 뒤로
      const scoreB = b.satisfactionScore ?? 6;
      return scoreA - scoreB;
    }) || [];

  // 선택 변경 시 실시간 시뮬레이션
  useEffect(() => {
    if (selectedIds.length === 0) {
      setSimulationResult(null);
      return;
    }

    const runSimulation = async () => {
      try {
        const result = await cancelSimulation.mutateAsync({ subscriptionIds: selectedIds });
        setSimulationResult(result);
      } catch (error) {
        console.error('Simulation failed:', error);
      }
    };

    runSimulation();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedIds]);

  const handleToggleSubscription = (id: string) => {
    setSelectedIds((prev) => (prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]));
  };

  const handleSelectAll = () => {
    if (selectedIds.length === sortedSubscriptions.length) {
      setSelectedIds([]);
    } else {
      setSelectedIds(sortedSubscriptions.map((sub: Subscription) => sub.id));
    }
  };

  const handleApply = () => {
    if (selectedIds.length === 0) return;
    setShowConfirmModal(true);
  };

  const handleConfirmApply = async () => {
    try {
      await applySimulation.mutateAsync({
        action: 'cancel',
        subscriptionIds: selectedIds,
      });
      setShowConfirmModal(false);
      setSelectedIds([]);
      setSimulationResult(null);
      setShowUndoButton(true);
    } catch (error) {
      console.error('Apply simulation failed:', error);
    }
  };

  if (subscriptionsLoading) {
    return (
      <div className="animate-pulse space-y-4">
        <div className="h-12 rounded bg-gray-200"></div>
        <div className="h-64 rounded bg-gray-200"></div>
      </div>
    );
  }

  if (!sortedSubscriptions.length) {
    return (
      <div className="rounded-lg border border-gray-200 bg-white p-8 text-center">
        <div className="text-gray-500">활성 구독이 없습니다</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* 구독 목록 */}
      <div className="rounded-lg border border-gray-200 bg-white">
        <div className="border-b border-gray-200 p-4">
          <div className="flex items-center justify-between">
            <h3 className="font-semibold text-gray-900">해지할 구독 선택</h3>
            <button
              onClick={handleSelectAll}
              className="text-sm text-blue-600 hover:text-blue-700"
            >
              {selectedIds.length === sortedSubscriptions.length ? '전체 해제' : '전체 선택'}
            </button>
          </div>
          <div className="mt-2 text-sm text-gray-500">
            만족도가 낮은 순으로 정렬되어 있습니다
          </div>
        </div>

        <div className="divide-y divide-gray-200">
          {sortedSubscriptions.map((subscription: Subscription) => (
            <label
              key={subscription.id}
              className={cn(
                'flex cursor-pointer items-center gap-4 p-4 transition-colors hover:bg-gray-50',
                {
                  'bg-blue-50': selectedIds.includes(subscription.id),
                }
              )}
            >
              <input
                type="checkbox"
                checked={selectedIds.includes(subscription.id)}
                onChange={() => handleToggleSubscription(subscription.id)}
                className="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
              <div className="flex-1">
                <div className="flex items-center justify-between">
                  <div>
                    <div className="font-medium text-gray-900">{subscription.serviceName}</div>
                    <div className="mt-1 flex items-center gap-2 text-sm text-gray-500">
                      <span>{formatCurrency(subscription.amount)}</span>
                      <span>·</span>
                      <span>{billingCycleLabel(subscription.billingCycle)}</span>
                      {subscription.category && (
                        <>
                          <span>·</span>
                          <span>{subscription.category.name}</span>
                        </>
                      )}
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-sm text-gray-500">만족도</div>
                    <div className="mt-1 text-yellow-500">
                      {satisfactionStars(subscription.satisfactionScore)}
                    </div>
                  </div>
                </div>
              </div>
            </label>
          ))}
        </div>
      </div>

      {/* 시뮬레이션 결과 */}
      <SimulationResultComponent
        result={simulationResult}
        loading={cancelSimulation.isPending}
      />

      {/* 적용 버튼 */}
      <div className="flex justify-end gap-3">
        {cancelSimulation.isError && (
          <div className="flex-1 text-sm text-red-600">
            시뮬레이션 실패. 다시 시도해주세요.
          </div>
        )}
        <button
          onClick={handleApply}
          disabled={selectedIds.length === 0 || applySimulation.isPending}
          className={cn(
            'rounded-lg bg-red-600 px-6 py-3 font-medium text-white transition-colors',
            'hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2',
            'disabled:cursor-not-allowed disabled:opacity-50'
          )}
        >
          {applySimulation.isPending ? '적용 중...' : '해지 적용'}
        </button>
      </div>

      {/* 확인 모달 */}
      {showConfirmModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
          <div className="max-w-md rounded-lg bg-white p-6 shadow-xl">
            <h3 className="text-lg font-semibold text-gray-900">구독 해지 확인</h3>
            <p className="mt-2 text-gray-600">
              선택한 {selectedIds.length}개의 구독을 해지하시겠습니까?
              <br />
              이 작업은 되돌릴 수 있습니다.
            </p>
            {applySimulation.isError && (
              <div className="mt-4 rounded-lg bg-red-50 p-3 text-sm text-red-600">
                적용에 실패했습니다. 다시 시도해주세요.
              </div>
            )}
            <div className="mt-6 flex justify-end gap-3">
              <button
                onClick={() => setShowConfirmModal(false)}
                disabled={applySimulation.isPending}
                className="rounded-lg border border-gray-300 px-4 py-2 font-medium text-gray-700 hover:bg-gray-50"
              >
                취소
              </button>
              <button
                onClick={handleConfirmApply}
                disabled={applySimulation.isPending}
                className={cn(
                  'rounded-lg bg-red-600 px-4 py-2 font-medium text-white hover:bg-red-700',
                  'disabled:cursor-not-allowed disabled:opacity-50'
                )}
              >
                {applySimulation.isPending ? '적용 중...' : '확인'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* 되돌리기 버튼 */}
      {showUndoButton && <UndoButton onUndo={() => setShowUndoButton(false)} />}
    </div>
  );
}
