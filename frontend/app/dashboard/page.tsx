'use client';

import Link from 'next/link';
import { useState } from 'react';

import CategoryChart from '@/components/dashboard/CategoryChart';
import RecommendationList from '@/components/dashboard/RecommendationList';
import SummaryCards from '@/components/dashboard/SummaryCards';
import SubscriptionCard from '@/components/subscription/SubscriptionCard';
import SubscriptionForm from '@/components/subscription/SubscriptionForm';
import { useDashboardSummary, useRecommendations } from '@/lib/hooks/useDashboard';
import { useSubscriptions } from '@/lib/hooks/useSubscriptions';
import type { Subscription } from '@/types';

export default function DashboardPage() {
  const [showForm, setShowForm] = useState(false);
  const [editingSubscription, setEditingSubscription] = useState<Subscription | undefined>();

  const { data: summary, isLoading: summaryLoading, error: summaryError } = useDashboardSummary();
  const { data: recommendations, isLoading: recsLoading } = useRecommendations();
  const {
    data: subscriptionsData,
    isLoading: subsLoading,
    error: subsError,
  } = useSubscriptions({
    status: 'active',
    perPage: 5,
    sortBy: 'nextBillingDate',
    sortOrder: 'asc',
  });

  const handleEdit = (subscription: Subscription) => {
    setEditingSubscription(subscription);
    setShowForm(true);
  };

  const handleCloseForm = () => {
    setShowForm(false);
    setEditingSubscription(undefined);
  };

  if (summaryError || subsError) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-red-600">오류가 발생했습니다</h2>
          <p className="mt-2 text-gray-600">대시보드를 불러올 수 없습니다.</p>
          <button
            onClick={() => window.location.reload()}
            className="mt-4 rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
          >
            다시 시도
          </button>
        </div>
      </div>
    );
  }

  if (summaryLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <div className="h-12 w-12 animate-spin rounded-full border-4 border-blue-600 border-t-transparent"></div>
          <p className="mt-4 text-gray-600">로딩 중...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 p-4 sm:p-6 lg:p-8">
      <div className="mx-auto max-w-7xl">
        {/* Header */}
        <div className="mb-6 flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">대시보드</h1>
            <p className="mt-1 text-gray-600">구독 서비스를 한눈에 확인하세요</p>
          </div>
          <button
            onClick={() => setShowForm(true)}
            className="rounded-lg bg-blue-600 px-6 py-3 font-semibold text-white transition-colors hover:bg-blue-700"
          >
            + 구독 추가
          </button>
        </div>

        {/* Summary Cards */}
        {summary && (
          <div className="mb-6">
            <SummaryCards summary={summary} />
          </div>
        )}

        {/* Charts and Recommendations Grid */}
        <div className="mb-6 grid gap-6 lg:grid-cols-2">
          {/* Category Chart */}
          {summary && (
            <div>
              <CategoryChart breakdown={summary.categoryBreakdown} />
            </div>
          )}

          {/* Recommendations */}
          {!recsLoading && recommendations && (
            <div>
              <RecommendationList recommendations={recommendations} />
            </div>
          )}
        </div>

        {/* Recent Subscriptions */}
        <div className="rounded-lg border-2 border-gray-200 bg-white p-6">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-xl font-semibold text-gray-900">최근 구독</h2>
            <Link
              href="/subscriptions"
              className="text-sm font-medium text-blue-600 hover:text-blue-700"
            >
              전체 보기 →
            </Link>
          </div>

          {subsLoading ? (
            <div className="py-8 text-center text-gray-500">로딩 중...</div>
          ) : subscriptionsData?.data && subscriptionsData.data.length > 0 ? (
            <div className="space-y-4">
              {subscriptionsData.data.map((subscription) => (
                <SubscriptionCard
                  key={subscription.id}
                  subscription={subscription}
                  onEdit={handleEdit}
                />
              ))}
            </div>
          ) : (
            <div className="py-8 text-center">
              <p className="text-gray-500">등록된 구독이 없습니다</p>
              <button
                onClick={() => setShowForm(true)}
                className="mt-4 text-sm font-medium text-blue-600 hover:text-blue-700"
              >
                첫 구독 추가하기
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Subscription Form Modal */}
      <SubscriptionForm
        isOpen={showForm}
        onClose={handleCloseForm}
        subscription={editingSubscription}
      />
    </div>
  );
}
