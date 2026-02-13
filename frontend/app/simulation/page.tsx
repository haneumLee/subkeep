'use client';

import { useState } from 'react';

import { CancelSimulation } from '@/components/simulation/CancelSimulation';
import { AddSimulation } from '@/components/simulation/AddSimulation';
import { cn, formatCurrency } from '@/lib/utils';
import { useQuery } from '@tanstack/react-query';
import { get } from '@/lib/api';
import type { DashboardSummary, Category } from '@/types';

type TabType = 'cancel' | 'add';

export default function SimulationPage() {
  const [activeTab, setActiveTab] = useState<TabType>('cancel');

  // 현재 월 총액 조회
  const { data: dashboardData, isLoading: dashboardLoading } = useQuery<DashboardSummary>({
    queryKey: ['dashboard'],
    queryFn: () => get<DashboardSummary>('/dashboard/summary'),
  });

  // 카테고리 목록 조회 (추가 시뮬레이션용)
  const { data: categoriesData } = useQuery<Category[]>({
    queryKey: ['categories'],
    queryFn: () => get<Category[]>('/categories'),
  });

  return (
    <div className="space-y-6">
      {/* 현재 월 총액 */}
      <div className="rounded-lg border border-gray-200 bg-white p-6">
        <div className="flex items-center justify-between">
          <div>
            <div className="text-sm text-gray-500">현재 월 총액</div>
            <div className="mt-1 text-3xl font-bold text-gray-900">
              {dashboardLoading ? (
                <div className="h-9 w-32 animate-pulse rounded bg-gray-200"></div>
              ) : (
                formatCurrency(dashboardData?.monthlyTotal || 0)
              )}
            </div>
          </div>
          <div className="rounded-lg bg-blue-50 p-4">
            <svg
              className="h-8 w-8 text-blue-600"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"
              />
            </svg>
          </div>
        </div>
      </div>

      {/* 탭 */}
      <div className="border-b border-gray-200">
        <nav className="-mb-px flex space-x-8">
          <button
            onClick={() => setActiveTab('cancel')}
            className={cn(
              'whitespace-nowrap border-b-2 px-1 py-4 text-sm font-medium transition-colors',
              activeTab === 'cancel'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
            )}
          >
            해지 시뮬레이션
          </button>
          <button
            onClick={() => setActiveTab('add')}
            className={cn(
              'whitespace-nowrap border-b-2 px-1 py-4 text-sm font-medium transition-colors',
              activeTab === 'add'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
            )}
          >
            추가 시뮬레이션
          </button>
        </nav>
      </div>

      {/* 탭 컨텐츠 */}
      <div className="py-6">
        {activeTab === 'cancel' && <CancelSimulation />}
        {activeTab === 'add' && <AddSimulation categories={categoriesData || []} />}
      </div>
    </div>
  );
}
