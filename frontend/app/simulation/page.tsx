'use client';

import { useCallback, useEffect, useState } from 'react';

import { SimulationResult as SimulationResultComponent } from '@/components/simulation/SimulationResult';
import { UndoButton } from '@/components/simulation/UndoButton';
import { useCombinedSimulation, useApplySimulation } from '@/lib/hooks/useSimulation';
import { useSubscriptions } from '@/lib/hooks/useSubscriptions';
import { useCategories } from '@/lib/hooks/useCategories';
import { formatCurrency, billingCycleLabel, satisfactionStars, cn } from '@/lib/utils';
import type {
  BillingCycle,
  CombinedSimulationItem,
  SimulationResult,
  Subscription,
} from '@/types';

const BILLING_CYCLES: BillingCycle[] = ['weekly', 'monthly', 'yearly'];

export default function SimulationPage() {
  const { data: subscriptionsData, isLoading: subscriptionsLoading } = useSubscriptions({
    status: 'active',
  });
  const { data: categories } = useCategories();
  const combinedSimulation = useCombinedSimulation();
  const applySimulation = useApplySimulation();

  // Cancel state
  const [cancelIds, setCancelIds] = useState<string[]>([]);

  // Add state - support multiple virtual items
  const [addItems, setAddItems] = useState<CombinedSimulationItem[]>([]);
  const [addForm, setAddForm] = useState<CombinedSimulationItem>({
    serviceName: '',
    amount: 0,
    billingCycle: 'monthly',
  });

  // Result & UI state
  const [simulationResult, setSimulationResult] = useState<SimulationResult | null>(null);
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [showUndoButton, setShowUndoButton] = useState(false);

  // Sorted subscriptions by satisfaction (low first)
  const sortedSubscriptions = [...(subscriptionsData?.data || [])].sort(
    (a: Subscription, b: Subscription) => {
      const scoreA = a.satisfactionScore ?? 6;
      const scoreB = b.satisfactionScore ?? 6;
      return scoreA - scoreB;
    }
  );

  // Run combined simulation whenever cancel/add selections change
  const runSimulation = useCallback(async () => {
    if (cancelIds.length === 0 && addItems.length === 0) {
      setSimulationResult(null);
      return;
    }
    try {
      const result = await combinedSimulation.mutateAsync({
        cancelSubscriptionIds: cancelIds,
        addItems,
      });
      setSimulationResult(result);
    } catch (error) {
      console.error('Simulation failed:', error);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [cancelIds, addItems]);

  useEffect(() => {
    const timer = setTimeout(runSimulation, 200);
    return () => clearTimeout(timer);
  }, [runSimulation]);

  // Cancel handlers
  const handleToggleCancel = (id: string) => {
    setCancelIds((prev) =>
      prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]
    );
  };

  const handleSelectAllCancel = () => {
    if (cancelIds.length === sortedSubscriptions.length) {
      setCancelIds([]);
    } else {
      setCancelIds(sortedSubscriptions.map((sub: Subscription) => sub.id));
    }
  };

  // Add handlers
  const handleAddItem = () => {
    if (!addForm.serviceName.trim() || addForm.amount <= 0) return;
    setAddItems((prev) => [...prev, { ...addForm }]);
    setAddForm({ serviceName: '', amount: 0, billingCycle: 'monthly' });
  };

  const handleRemoveItem = (index: number) => {
    setAddItems((prev) => prev.filter((_, i) => i !== index));
  };

  // Apply (cancel only - add items are virtual)
  const handleApply = () => {
    if (cancelIds.length === 0) return;
    setShowConfirmModal(true);
  };

  const handleConfirmApply = async () => {
    try {
      await applySimulation.mutateAsync({
        action: 'cancel',
        subscriptionIds: cancelIds,
      });
      setShowConfirmModal(false);
      setCancelIds([]);
      setSimulationResult(null);
      setShowUndoButton(true);
    } catch (error) {
      console.error('Apply simulation failed:', error);
    }
  };

  // Reset all
  const handleReset = () => {
    setCancelIds([]);
    setAddItems([]);
    setSimulationResult(null);
  };

  return (
    <div className="space-y-6">
      {/* 시뮬레이션 결과 - 상단 */}
      <SimulationResultComponent
        result={simulationResult}
        loading={combinedSimulation.isPending}
      />

      {/* 해지 + 추가 패널 - 좌우 배치 */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* 좌측: 해지 패널 */}
        <div className="rounded-lg border border-gray-200 bg-white">
          <div className="border-b border-gray-200 p-4">
            <div className="flex items-center justify-between">
              <h3 className="font-semibold text-gray-900">
                해지 선택
                {cancelIds.length > 0 && (
                  <span className="ml-2 text-sm font-normal text-red-600">
                    {cancelIds.length}개 선택
                  </span>
                )}
              </h3>
              <button
                onClick={handleSelectAllCancel}
                className="text-sm text-blue-600 hover:text-blue-700"
              >
                {cancelIds.length === sortedSubscriptions.length ? '전체 해제' : '전체 선택'}
              </button>
            </div>
            <p className="mt-1 text-xs text-gray-500">
              만족도가 낮은 순으로 정렬
            </p>
          </div>

          <div className="max-h-[480px] divide-y divide-gray-200 overflow-y-auto">
            {subscriptionsLoading ? (
              <div className="animate-pulse space-y-3 p-4">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-14 rounded bg-gray-100" />
                ))}
              </div>
            ) : sortedSubscriptions.length === 0 ? (
              <div className="p-8 text-center text-sm text-gray-500">
                활성 구독이 없습니다
              </div>
            ) : (
              sortedSubscriptions.map((subscription: Subscription) => (
                <label
                  key={subscription.id}
                  className={cn(
                    'flex cursor-pointer items-center gap-3 px-4 py-3 transition-colors hover:bg-gray-50',
                    { 'bg-red-50': cancelIds.includes(subscription.id) }
                  )}
                >
                  <input
                    type="checkbox"
                    checked={cancelIds.includes(subscription.id)}
                    onChange={() => handleToggleCancel(subscription.id)}
                    className="h-4 w-4 rounded border-gray-300 text-red-600 focus:ring-red-500"
                  />
                  <div className="min-w-0 flex-1">
                    <div className="font-medium text-gray-900">{subscription.serviceName}</div>
                    <div className="flex items-center gap-1.5 text-xs text-gray-500">
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
                    <div className="text-[10px] text-gray-400">만족도</div>
                    <div className="text-sm text-yellow-500">
                      {satisfactionStars(subscription.satisfactionScore)}
                    </div>
                  </div>
                </label>
              ))
            )}
          </div>
        </div>

        {/* 우측: 추가 패널 */}
        <div className="rounded-lg border border-gray-200 bg-white">
          <div className="border-b border-gray-200 p-4">
            <h3 className="font-semibold text-gray-900">
              추가 시뮬레이션
              {addItems.length > 0 && (
                <span className="ml-2 text-sm font-normal text-blue-600">
                  {addItems.length}개 추가
                </span>
              )}
            </h3>
            <p className="mt-1 text-xs text-gray-500">
              가상 구독을 추가하여 비용 변화를 확인
            </p>
          </div>

          <div className="p-4">
            {/* 입력 폼 */}
            <div className="space-y-3">
              <div>
                <input
                  type="text"
                  value={addForm.serviceName}
                  onChange={(e) =>
                    setAddForm((prev) => ({ ...prev, serviceName: e.target.value }))
                  }
                  placeholder="서비스명 (예: Netflix)"
                  maxLength={50}
                  className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
                />
              </div>

              <div className="grid grid-cols-2 gap-2">
                <div className="relative">
                  <input
                    type="number"
                    min={0}
                    max={1000000}
                    step={100}
                    value={addForm.amount || ''}
                    onChange={(e) =>
                      setAddForm((prev) => ({
                        ...prev,
                        amount: Number(e.target.value),
                      }))
                    }
                    placeholder="금액"
                    className="w-full rounded-lg border border-gray-300 px-3 py-2 pr-8 text-sm focus:border-blue-500 focus:outline-none"
                  />
                  <span className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-xs text-gray-400">
                    원
                  </span>
                </div>

                <select
                  value={addForm.billingCycle}
                  onChange={(e) =>
                    setAddForm((prev) => ({
                      ...prev,
                      billingCycle: e.target.value as BillingCycle,
                    }))
                  }
                  className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
                >
                  {BILLING_CYCLES.map((cycle) => (
                    <option key={cycle} value={cycle}>
                      {billingCycleLabel(cycle)}
                    </option>
                  ))}
                </select>
              </div>

              <div className="flex gap-2">
                <select
                  value={addForm.categoryId || ''}
                  onChange={(e) =>
                    setAddForm((prev) => ({
                      ...prev,
                      categoryId: e.target.value || undefined,
                    }))
                  }
                  className="flex-1 rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
                >
                  <option value="">카테고리 없음</option>
                  {categories?.map((cat) => (
                    <option key={cat.id} value={cat.id}>
                      {cat.name}
                    </option>
                  ))}
                </select>

                <button
                  onClick={handleAddItem}
                  disabled={!addForm.serviceName.trim() || addForm.amount <= 0}
                  className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  추가
                </button>
              </div>
            </div>

            {/* 추가된 가상 구독 목록 */}
            {addItems.length > 0 && (
              <div className="mt-4 max-h-[280px] divide-y divide-gray-100 overflow-y-auto rounded-lg border border-gray-200">
                {addItems.map((item, index) => (
                  <div
                    key={index}
                    className="flex items-center justify-between px-3 py-2.5"
                  >
                    <div className="min-w-0 flex-1">
                      <div className="text-sm font-medium text-gray-900">
                        {item.serviceName}
                      </div>
                      <div className="text-xs text-gray-500">
                        {formatCurrency(item.amount)} · {billingCycleLabel(item.billingCycle)}
                      </div>
                    </div>
                    <button
                      onClick={() => handleRemoveItem(index)}
                      className="ml-2 text-gray-400 hover:text-red-500"
                    >
                      <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </div>
                ))}
              </div>
            )}

            {/* 안내 */}
            <div className="mt-4 rounded-lg bg-gray-50 p-3 text-xs text-gray-500">
              추가 항목은 시뮬레이션 전용이며 실제 구독에 저장되지 않습니다.
            </div>
          </div>
        </div>
      </div>

      {/* 하단 버튼 */}
      <div className="flex items-center justify-between">
        <button
          onClick={handleReset}
          disabled={cancelIds.length === 0 && addItems.length === 0}
          className="rounded-lg border border-gray-300 px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
        >
          초기화
        </button>
        <button
          onClick={handleApply}
          disabled={cancelIds.length === 0 || applySimulation.isPending}
          className={cn(
            'rounded-lg bg-red-500 px-6 py-2.5 text-sm font-medium text-white transition-colors',
            'hover:bg-red-600 disabled:cursor-not-allowed disabled:opacity-50'
          )}
        >
          {applySimulation.isPending
            ? '적용 중...'
            : `해지 적용 (${cancelIds.length}개)`}
        </button>
      </div>

      {/* 확인 모달 */}
      {showConfirmModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
          <div className="max-w-md rounded-lg bg-white p-6 shadow-xl">
            <h3 className="text-lg font-semibold text-gray-900">구독 해지 확인</h3>
            <p className="mt-2 text-gray-600">
              선택한 {cancelIds.length}개의 구독을 해지하시겠습니까?
              <br />
              이 작업은 30초 이내에 되돌릴 수 있습니다.
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
