'use client';

import { useState } from 'react';

import SubscriptionCard from '@/components/subscription/SubscriptionCard';
import SubscriptionForm from '@/components/subscription/SubscriptionForm';
import FolderPanel from '@/components/subscription/FolderPanel';
import FolderManagerPopup from '@/components/subscription/FolderManagerPopup';
import CategoryManagerPopup from '@/components/subscription/CategoryManagerPopup';
import { useCategories } from '@/lib/hooks/useCategories';
import { useSubscriptions } from '@/lib/hooks/useSubscriptions';
import type { Subscription } from '@/types';

export default function SubscriptionsPage() {
  const [showForm, setShowForm] = useState(false);
  const [editingSubscription, setEditingSubscription] = useState<Subscription | undefined>();
  const [showFolderManager, setShowFolderManager] = useState(false);
  const [showCategoryManager, setShowCategoryManager] = useState(false);

  // Filters
  const [statusFilter, setStatusFilter] = useState<string>('active');
  const [categoryFilter, setCategoryFilter] = useState<string>('');
  const [folderFilter, setFolderFilter] = useState<string>('');
  const [sortBy, setSortBy] = useState<string>('nextBillingDate');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [currentPage, setCurrentPage] = useState(1);

  const { data: categories } = useCategories();
  const {
    data: subscriptionsData,
    isLoading,
    error,
  } = useSubscriptions({
    status: statusFilter || undefined,
    categoryId: categoryFilter || undefined,
    folderId: folderFilter || undefined,
    sortBy,
    sortOrder,
    page: currentPage,
    perPage: 20,
  });

  const handleEdit = (subscription: Subscription) => {
    setEditingSubscription(subscription);
    setShowForm(true);
  };

  const handleCloseForm = () => {
    setShowForm(false);
    setEditingSubscription(undefined);
  };

  const handleFilterChange = () => {
    setCurrentPage(1);
  };

  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-red-600">오류가 발생했습니다</h2>
          <p className="mt-2 text-gray-600">구독 목록을 불러올 수 없습니다.</p>
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

  return (
    <div className="flex gap-6">
      {/* Folder Panel - Left Sidebar */}
      <FolderPanel
        selectedFolderId={folderFilter}
        onSelectFolder={(id) => {
          setFolderFilter(id);
          setCurrentPage(1);
        }}
      />

      {/* Main Content */}
      <div className="flex-1 space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-slate-900">구독 관리</h2>
            <p className="mt-1 text-sm text-slate-600">
              전체 {subscriptionsData?.meta.total || 0}개의 구독
            </p>
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setShowFolderManager(true)}
              className="rounded-lg border-2 border-gray-300 px-4 py-3 font-semibold text-gray-700 transition-colors hover:bg-gray-50"
            >
              📁 폴더 관리
            </button>
            <button
              onClick={() => setShowCategoryManager(true)}
              className="rounded-lg border-2 border-gray-300 px-4 py-3 font-semibold text-gray-700 transition-colors hover:bg-gray-50"
            >
              🏷️ 카테고리 관리
            </button>
            <button
              onClick={() => setShowForm(true)}
              className="rounded-lg bg-blue-600 px-6 py-3 font-semibold text-white transition-colors hover:bg-blue-700"
            >
              + 구독 추가
            </button>
          </div>
        </div>

        {/* Filters */}
        <div className="mb-6 rounded-lg border-2 border-gray-200 bg-white p-4">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {/* Status Filter */}
            <div>
              <label className="block text-sm font-medium text-gray-700">상태</label>
              <select
                value={statusFilter}
                onChange={(e) => {
                  setStatusFilter(e.target.value);
                  handleFilterChange();
                }}
                className="mt-1 w-full rounded-lg border-2 border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
              >
                <option value="">전체</option>
                <option value="active">활성</option>
                <option value="paused">일시중지</option>
                <option value="cancelled">해지</option>
              </select>
            </div>

            {/* Category Filter */}
            <div>
              <label className="block text-sm font-medium text-gray-700">카테고리</label>
              <select
                value={categoryFilter}
                onChange={(e) => {
                  setCategoryFilter(e.target.value);
                  handleFilterChange();
                }}
                className="mt-1 w-full rounded-lg border-2 border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
              >
                <option value="">전체</option>
                {categories?.map((category) => (
                  <option key={category.id} value={category.id}>
                    {category.name}
                  </option>
                ))}
              </select>
            </div>

            {/* Sort By */}
            <div>
              <label className="block text-sm font-medium text-gray-700">정렬 기준</label>
              <select
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value)}
                className="mt-1 w-full rounded-lg border-2 border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
              >
                <option value="nextBillingDate">결제일순</option>
                <option value="amount">금액순</option>
                <option value="satisfactionScore">만족도순</option>
                <option value="serviceName">이름순</option>
              </select>
            </div>

            {/* Sort Order */}
            <div>
              <label className="block text-sm font-medium text-gray-700">정렬 순서</label>
              <select
                value={sortOrder}
                onChange={(e) => setSortOrder(e.target.value as 'asc' | 'desc')}
                className="mt-1 w-full rounded-lg border-2 border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
              >
                <option value="asc">오름차순</option>
                <option value="desc">내림차순</option>
              </select>
            </div>
          </div>
        </div>

        {/* Subscriptions List */}
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="text-center">
              <div className="h-12 w-12 animate-spin rounded-full border-4 border-blue-600 border-t-transparent"></div>
              <p className="mt-4 text-gray-600">로딩 중...</p>
            </div>
          </div>
        ) : subscriptionsData?.data && subscriptionsData.data.length > 0 ? (
          <>
            <div className="space-y-4">
              {subscriptionsData.data.map((subscription) => (
                <SubscriptionCard
                  key={subscription.id}
                  subscription={subscription}
                  onEdit={handleEdit}
                />
              ))}
            </div>

            {/* Pagination */}
            {subscriptionsData.meta.totalPages > 1 && (
              <div className="mt-6 flex items-center justify-center gap-2">
                <button
                  onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                  className="rounded-lg border-2 border-gray-300 px-4 py-2 font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  이전
                </button>

                <div className="flex gap-1">
                  {Array.from({ length: subscriptionsData.meta.totalPages }, (_, i) => i + 1)
                    .filter((page) => {
                      const distance = Math.abs(page - currentPage);
                      return distance === 0 || distance === 1 || page === 1 || page === subscriptionsData.meta.totalPages;
                    })
                    .map((page, index, array) => {
                      if (index > 0 && array[index - 1] !== page - 1) {
                        return (
                          <span key={`ellipsis-${page}`} className="px-2 py-2 text-gray-500">
                            ...
                          </span>
                        );
                      }
                      return (
                        <button
                          key={page}
                          onClick={() => setCurrentPage(page)}
                          className={`rounded-lg px-4 py-2 font-medium transition-colors ${
                            currentPage === page
                              ? 'bg-blue-600 text-white'
                              : 'border-2 border-gray-300 text-gray-700 hover:bg-gray-50'
                          }`}
                        >
                          {page}
                        </button>
                      );
                    })}
                </div>

                <button
                  onClick={() => setCurrentPage((p) => Math.min(subscriptionsData.meta.totalPages, p + 1))}
                  disabled={currentPage === subscriptionsData.meta.totalPages}
                  className="rounded-lg border-2 border-gray-300 px-4 py-2 font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  다음
                </button>
              </div>
            )}
          </>
        ) : (
          <div className="rounded-lg border-2 border-gray-200 bg-white p-12 text-center">
            <p className="text-gray-500">구독이 없습니다</p>
            <button
              onClick={() => setShowForm(true)}
              className="mt-4 text-sm font-medium text-blue-600 hover:text-blue-700"
            >
              첫 구독 추가하기
            </button>
          </div>
        )}

      {/* Subscription Form Modal */}
      <SubscriptionForm
        isOpen={showForm}
        onClose={handleCloseForm}
        subscription={editingSubscription}
      />

      {/* Folder Manager Popup */}
      <FolderManagerPopup
        isOpen={showFolderManager}
        onClose={() => setShowFolderManager(false)}
      />

      {/* Category Manager Popup */}
      <CategoryManagerPopup
        isOpen={showCategoryManager}
        onClose={() => setShowCategoryManager(false)}
      />
      </div>
    </div>
  );
}
