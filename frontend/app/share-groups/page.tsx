'use client';

import { useState } from 'react';

import { Button } from '@/components/ui/Button';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import { ShareGroupCard } from '@/components/share/ShareGroupCard';
import { ShareGroupForm } from '@/components/share/ShareGroupForm';
import { useShareGroups } from '@/lib/hooks/useShareGroups';

export default function ShareGroupsPage() {
  const { data: shareGroups, isLoading, error } = useShareGroups();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  if (isLoading) {
    return (
      <div className="flex min-h-[400px] items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-lg bg-red-50 p-6 text-center">
        <svg
          className="mx-auto h-12 w-12 text-red-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
        </svg>
        <h3 className="mt-4 text-lg font-semibold text-red-900">데이터를 불러올 수 없습니다</h3>
        <p className="mt-2 text-sm text-red-700">잠시 후 다시 시도해주세요</p>
      </div>
    );
  }

  const hasGroups = shareGroups && shareGroups.length > 0;

  return (
    <>
      <div className="mb-6 flex items-center justify-between">
        <div className="text-sm text-slate-600">
          {hasGroups ? `총 ${shareGroups.length}개의 그룹` : '아직 생성된 그룹이 없습니다'}
        </div>
        <Button variant="primary" onClick={() => setIsCreateModalOpen(true)}>
          <svg
            className="mr-2 h-5 w-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 4v16m8-8H4"
            />
          </svg>
          그룹 추가
        </Button>
      </div>

      {hasGroups ? (
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {shareGroups.map((group) => (
            <ShareGroupCard key={group.id} group={group} />
          ))}
        </div>
      ) : (
        <div className="rounded-lg border-2 border-dashed border-slate-300 bg-white p-12 text-center">
          <svg
            className="mx-auto h-16 w-16 text-slate-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
            />
          </svg>
          <h3 className="mt-4 text-lg font-semibold text-slate-900">공유 그룹이 없습니다</h3>
          <p className="mt-2 text-sm text-slate-600">
            가족, 친구와 함께 구독을 공유하고 비용을 절약하세요
          </p>
          <Button
            variant="primary"
            onClick={() => setIsCreateModalOpen(true)}
            className="mt-6"
          >
            첫 번째 그룹 만들기
          </Button>
        </div>
      )}

      <ShareGroupForm
        mode="create"
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
      />
    </>
  );
}
