'use client';

import { useEffect, useState } from 'react';

import type { AddSimulationRequest, BillingCycle, Category, SimulationResult } from '@/types';
import { useAddSimulation } from '@/lib/hooks/useSimulation';
import { formatCurrency, billingCycleLabel, cn } from '@/lib/utils';
import { SimulationResult as SimulationResultComponent } from './SimulationResult';

interface AddSimulationProps {
  categories?: Category[];
}

const BILLING_CYCLES: BillingCycle[] = ['weekly', 'monthly', 'yearly'];

export function AddSimulation({ categories = [] }: AddSimulationProps) {
  const addSimulation = useAddSimulation();

  const [formData, setFormData] = useState<AddSimulationRequest>({
    serviceName: '',
    amount: 0,
    billingCycle: 'monthly',
    categoryId: undefined,
  });
  const [simulationResult, setSimulationResult] = useState<SimulationResult | null>(null);
  const [isFormValid, setIsFormValid] = useState(false);

  // 폼 유효성 검사
  useEffect(() => {
    const isValid = formData.serviceName.trim().length > 0 && formData.amount > 0;
    setIsFormValid(isValid);
  }, [formData]);

  // 입력 변경 시 실시간 시뮬레이션
  useEffect(() => {
    if (!isFormValid) {
      setSimulationResult(null);
      return;
    }

    const runSimulation = async () => {
      try {
        const result = await addSimulation.mutateAsync(formData);
        setSimulationResult(result);
      } catch (error) {
        console.error('Simulation failed:', error);
      }
    };

    // 디바운싱
    const timer = setTimeout(runSimulation, 300);
    return () => clearTimeout(timer);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [formData, isFormValid]);

  const handleInputChange = (field: keyof AddSimulationRequest, value: string | number) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  const handleReset = () => {
    setFormData({
      serviceName: '',
      amount: 0,
      billingCycle: 'monthly',
      categoryId: undefined,
    });
    setSimulationResult(null);
  };

  return (
    <div className="space-y-6">
      {/* 입력 폼 */}
      <div className="rounded-lg border border-gray-200 bg-white p-6">
        <h3 className="mb-4 font-semibold text-gray-900">가상 구독 정보 입력</h3>
        <div className="space-y-4">
          {/* 서비스명 */}
          <div>
            <label htmlFor="serviceName" className="block text-sm font-medium text-gray-700">
              서비스명 <span className="text-red-500">*</span>
            </label>
            <input
              id="serviceName"
              type="text"
              value={formData.serviceName}
              onChange={(e) => handleInputChange('serviceName', e.target.value)}
              placeholder="예: Netflix"
              maxLength={50}
              className={cn(
                'mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2',
                'focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500'
              )}
            />
            <div className="mt-1 text-xs text-gray-500">
              {formData.serviceName.length}/50자
            </div>
          </div>

          {/* 금액 */}
          <div>
            <label htmlFor="amount" className="block text-sm font-medium text-gray-700">
              금액 <span className="text-red-500">*</span>
            </label>
            <div className="relative mt-1">
              <input
                id="amount"
                type="number"
                min="0"
                max="1000000"
                step="100"
                value={formData.amount || ''}
                onChange={(e) => handleInputChange('amount', Number(e.target.value))}
                placeholder="0"
                className={cn(
                  'block w-full rounded-lg border border-gray-300 px-3 py-2 pr-12',
                  'focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500'
                )}
              />
              <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3">
                <span className="text-gray-500">원</span>
              </div>
            </div>
            {formData.amount > 1000000 && (
              <div className="mt-1 text-xs text-orange-600">
                100만원을 초과하는 금액입니다. 확인해주세요.
              </div>
            )}
            {formData.amount > 0 && (
              <div className="mt-1 text-xs text-gray-500">
                {formatCurrency(formData.amount)}
              </div>
            )}
          </div>

          {/* 결제 주기 */}
          <div>
            <label htmlFor="billingCycle" className="block text-sm font-medium text-gray-700">
              결제 주기 <span className="text-red-500">*</span>
            </label>
            <div className="mt-2 grid grid-cols-3 gap-2">
              {BILLING_CYCLES.map((cycle) => (
                <button
                  key={cycle}
                  type="button"
                  onClick={() => handleInputChange('billingCycle', cycle)}
                  className={cn(
                    'rounded-lg border px-4 py-2 text-sm font-medium transition-colors',
                    formData.billingCycle === cycle
                      ? 'border-blue-600 bg-blue-50 text-blue-600'
                      : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50'
                  )}
                >
                  {billingCycleLabel(cycle)}
                </button>
              ))}
            </div>
          </div>

          {/* 카테고리 */}
          <div>
            <label htmlFor="categoryId" className="block text-sm font-medium text-gray-700">
              카테고리 (선택)
            </label>
            <select
              id="categoryId"
              value={formData.categoryId || ''}
              onChange={(e) => {
                const value = e.target.value;
                setFormData((prev) => ({
                  ...prev,
                  categoryId: value ? value : undefined,
                }));
              }}
              className={cn(
                'mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2',
                'focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500'
              )}
            >
              <option value="">카테고리 없음</option>
              {categories.map((category) => (
                <option key={category.id} value={category.id}>
                  {category.name}
                </option>
              ))}
            </select>
          </div>

          {/* 버튼 */}
          <div className="flex justify-end gap-3 pt-2">
            <button
              type="button"
              onClick={handleReset}
              className="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
            >
              초기화
            </button>
          </div>
        </div>
      </div>

      {/* 시뮬레이션 결과 */}
      <SimulationResultComponent
        result={simulationResult}
        loading={addSimulation.isPending}
      />

      {addSimulation.isError && (
        <div className="rounded-lg bg-red-50 p-4 text-sm text-red-600">
          시뮬레이션 실패. 입력값을 확인하고 다시 시도해주세요.
        </div>
      )}

      {/* 안내 메시지 */}
      <div className="rounded-lg bg-blue-50 p-4 text-sm text-blue-700">
        <div className="font-medium">안내</div>
        <ul className="mt-2 list-disc space-y-1 pl-5">
          <li>입력한 정보는 실제 구독에 저장되지 않습니다.</li>
          <li>현재 구독에 추가했을 때의 금액 변화를 미리 확인할 수 있습니다.</li>
          <li>실제 구독을 추가하려면 구독 관리 페이지를 이용하세요.</li>
        </ul>
      </div>
    </div>
  );
}
