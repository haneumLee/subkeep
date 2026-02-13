'use client';

import { Modal } from '@/components/ui/Modal';
import { useDeleteShareGroup } from '@/lib/hooks/useShareGroups';
import type { ShareGroup } from '@/types';

interface DeleteShareGroupModalProps {
  group: ShareGroup;
  isOpen: boolean;
  onClose: () => void;
}

export function DeleteShareGroupModal({ group, isOpen, onClose }: DeleteShareGroupModalProps) {
  const deleteMutation = useDeleteShareGroup();

  const handleConfirm = async () => {
    try {
      await deleteMutation.mutateAsync(group.id);
      onClose();
    } catch (error) {
      console.error('Failed to delete share group:', error);
    }
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="공유 그룹 삭제"
      confirmText="삭제"
      cancelText="취소"
      onConfirm={handleConfirm}
      isConfirmLoading={deleteMutation.isPending}
      confirmVariant="danger"
    >
      <div className="space-y-3">
        <p className="text-slate-700">
          <span className="font-semibold">{group.name}</span> 그룹을 삭제하시겠습니까?
        </p>
        <div className="rounded-md bg-amber-50 p-3 text-sm text-amber-800">
          <div className="flex">
            <svg
              className="mr-2 mt-0.5 h-5 w-5 flex-shrink-0"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                clipRule="evenodd"
              />
            </svg>
            <div>
              <p className="font-medium">경고</p>
              <p className="mt-1">
                이 그룹과 연결된 구독들의 공유 설정이 모두 해제됩니다. 이 작업은 되돌릴 수 없습니다.
              </p>
            </div>
          </div>
        </div>
      </div>
    </Modal>
  );
}
