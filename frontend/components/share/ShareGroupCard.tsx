'use client';

import { useState } from 'react';

import type { ShareGroup } from '@/types';
import { Button } from '@/components/ui/Button';

import { ShareGroupForm } from './ShareGroupForm';
import { DeleteShareGroupModal } from './DeleteShareGroupModal';

interface ShareGroupCardProps {
  group: ShareGroup;
}

export function ShareGroupCard({ group }: ShareGroupCardProps) {
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

  const memberCount = group.members?.length || 0;
  const memberNames = group.members?.map((m) => m.nickname).join(', ') || '멤버 없음';

  return (
    <>
      <div className="rounded-lg border border-slate-200 bg-white p-6 shadow-sm transition-shadow hover:shadow-md">
        <div className="mb-4 flex items-start justify-between">
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-slate-900">{group.name}</h3>
            {group.description && (
              <p className="mt-1 text-sm text-slate-600">{group.description}</p>
            )}
          </div>
        </div>

        <div className="mb-4 space-y-2">
          <div className="flex items-center text-sm text-slate-600">
            <svg
              className="mr-2 h-4 w-4"
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
            <span className="font-medium">{memberCount}명</span>
          </div>

          {memberCount > 0 && (
            <div className="text-sm text-slate-500">
              <span className="font-medium">멤버:</span> {memberNames}
            </div>
          )}
        </div>

        <div className="flex gap-2">
          <Button
            variant="secondary"
            size="sm"
            onClick={() => setIsEditModalOpen(true)}
            className="flex-1"
          >
            수정
          </Button>
          <Button
            variant="danger"
            size="sm"
            onClick={() => setIsDeleteModalOpen(true)}
            className="flex-1"
          >
            삭제
          </Button>
        </div>
      </div>

      {isEditModalOpen && (
        <ShareGroupForm
          mode="edit"
          group={group}
          isOpen={isEditModalOpen}
          onClose={() => setIsEditModalOpen(false)}
        />
      )}

      {isDeleteModalOpen && (
        <DeleteShareGroupModal
          group={group}
          isOpen={isDeleteModalOpen}
          onClose={() => setIsDeleteModalOpen(false)}
        />
      )}
    </>
  );
}
