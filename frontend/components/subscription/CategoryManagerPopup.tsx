'use client';

import { useState } from 'react';

import { Modal } from '@/components/ui/Modal';
import { CategoryForm } from '@/components/settings/CategoryForm';
import { useCategories, useDeleteCategory } from '@/lib/hooks/useCategories';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import type { Category } from '@/types';

interface CategoryManagerPopupProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function CategoryManagerPopup({ isOpen, onClose }: CategoryManagerPopupProps) {
  const { data: categories, isLoading } = useCategories();
  const deleteCategory = useDeleteCategory();

  const [editingCategory, setEditingCategory] = useState<Category | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const handleDeleteConfirm = () => {
    if (!deletingId) return;
    deleteCategory.mutate(deletingId, {
      onSuccess: () => setDeletingId(null),
    });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4">
      <div className="w-full max-w-md rounded-lg bg-white shadow-xl">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-gray-200 px-6 py-4">
          <h2 className="text-lg font-semibold text-gray-900">카테고리 관리</h2>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setShowCreateForm(true)}
              className="rounded-lg bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-700"
            >
              추가
            </button>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600"
              aria-label="닫기"
            >
              <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        {/* Category list */}
        <div className="max-h-80 overflow-y-auto px-6 py-4">
          {isLoading ? (
            <div className="flex justify-center py-8">
              <LoadingSpinner size="md" />
            </div>
          ) : categories && categories.length > 0 ? (
            <ul className="space-y-2">
              {categories.map((cat) => (
                <li
                  key={cat.id}
                  className="flex items-center justify-between rounded-lg border border-gray-200 p-3"
                >
                  <div className="flex items-center gap-3">
                    <span
                      className="inline-block h-4 w-4 rounded-full"
                      style={{ backgroundColor: cat.color || '#6366f1' }}
                    />
                    <span className="text-sm font-medium text-gray-900">{cat.name}</span>
                    {cat.isSystem && (
                      <span className="rounded bg-gray-100 px-1.5 py-0.5 text-xs text-gray-500">
                        시스템
                      </span>
                    )}
                  </div>
                  {!cat.isSystem && (
                    <div className="flex gap-2">
                      <button
                        onClick={() => setEditingCategory(cat)}
                        className="text-sm text-gray-500 hover:text-blue-600"
                      >
                        수정
                      </button>
                      <button
                        onClick={() => setDeletingId(cat.id)}
                        className="text-sm text-gray-500 hover:text-red-600"
                      >
                        삭제
                      </button>
                    </div>
                  )}
                </li>
              ))}
            </ul>
          ) : (
            <p className="py-8 text-center text-sm text-gray-500">
              등록된 카테고리가 없습니다
            </p>
          )}
        </div>

        {/* Footer */}
        <div className="border-t border-gray-200 px-6 py-4">
          <button
            onClick={onClose}
            className="w-full rounded-lg border-2 border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
          >
            닫기
          </button>
        </div>

        {/* Create Form */}
        {showCreateForm && (
          <CategoryForm onClose={() => setShowCreateForm(false)} />
        )}

        {/* Edit Form */}
        {editingCategory && (
          <CategoryForm
            category={editingCategory}
            onClose={() => setEditingCategory(null)}
          />
        )}

        {/* Delete Confirm */}
        <Modal
          isOpen={deletingId !== null}
          onClose={() => setDeletingId(null)}
          title="카테고리 삭제"
          confirmText="삭제"
          cancelText="취소"
          onConfirm={handleDeleteConfirm}
          confirmVariant="danger"
          isConfirmLoading={deleteCategory.isPending}
        >
          <p>이 카테고리를 삭제하시겠습니까? 해당 카테고리의 구독은 미분류로 변경됩니다.</p>
        </Modal>
      </div>
    </div>
  );
}
