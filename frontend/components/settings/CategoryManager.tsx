'use client';

import { useState } from 'react';

import { Button } from '@/components/ui/Button';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import { Modal } from '@/components/ui/Modal';
import { CategoryForm } from '@/components/settings/CategoryForm';
import { useCategories, useDeleteCategory } from '@/lib/hooks/useCategories';
import type { Category, CreateCategoryRequest } from '@/types';

export function CategoryManager() {
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

  if (isLoading) {
    return (
      <div className="flex min-h-[200px] items-center justify-center">
        <LoadingSpinner size="md" />
      </div>
    );
  }

  return (
    <div className="rounded-xl bg-white p-6 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <h3 className="text-lg font-semibold text-slate-900">카테고리 관리</h3>
        <Button size="sm" onClick={() => setShowCreateForm(true)}>
          추가
        </Button>
      </div>

      {/* Category List */}
      <ul className="space-y-2">
        {categories?.map((cat) => (
          <li
            key={cat.id}
            className="flex items-center justify-between rounded-lg border border-slate-100 p-3"
          >
            <div className="flex items-center gap-3">
              <span
                className="inline-block h-4 w-4 rounded-full"
                style={{ backgroundColor: cat.color || '#6366f1' }}
              />
              <span className="text-sm font-medium text-slate-900">{cat.name}</span>
              {cat.isSystem && (
                <span className="rounded bg-slate-100 px-1.5 py-0.5 text-xs text-slate-500">
                  시스템
                </span>
              )}
            </div>
            {!cat.isSystem && (
              <div className="flex gap-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setEditingCategory(cat)}
                >
                  수정
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setDeletingId(cat.id)}
                >
                  삭제
                </Button>
              </div>
            )}
          </li>
        ))}
      </ul>

      {categories?.length === 0 && (
        <p className="py-8 text-center text-sm text-slate-500">
          등록된 카테고리가 없습니다.
        </p>
      )}

      {/* Create Form */}
      {showCreateForm && (
        <CategoryForm
          onClose={() => setShowCreateForm(false)}
        />
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
  );
}
