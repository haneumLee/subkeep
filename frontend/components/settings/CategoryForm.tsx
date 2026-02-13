'use client';

import { useState } from 'react';

import { Button } from '@/components/ui/Button';
import { Modal } from '@/components/ui/Modal';
import { useCreateCategory, useUpdateCategory } from '@/lib/hooks/useCategories';
import type { Category } from '@/types';

interface CategoryFormProps {
  category?: Category;
  onClose: () => void;
}

const PRESET_COLORS = [
  '#6366f1', '#8b5cf6', '#ec4899', '#ef4444',
  '#f97316', '#eab308', '#22c55e', '#14b8a6',
  '#06b6d4', '#3b82f6', '#6b7280', '#1e293b',
];

export function CategoryForm({ category, onClose }: CategoryFormProps) {
  const isEdit = !!category;
  const [name, setName] = useState(category?.name || '');
  const [color, setColor] = useState(category?.color || '#6366f1');
  const [icon, setIcon] = useState(category?.icon || '');

  const createCategory = useCreateCategory();
  const updateCategory = useUpdateCategory();

  const isLoading = createCategory.isPending || updateCategory.isPending;

  const handleSubmit = () => {
    const data = { name, color, icon: icon || undefined };

    if (isEdit && category) {
      updateCategory.mutate(
        { id: category.id, data },
        { onSuccess: onClose }
      );
    } else {
      createCategory.mutate(data, { onSuccess: onClose });
    }
  };

  return (
    <Modal
      isOpen={true}
      onClose={onClose}
      title={isEdit ? '카테고리 수정' : '카테고리 추가'}
      confirmText={isEdit ? '수정' : '추가'}
      cancelText="취소"
      onConfirm={handleSubmit}
      isConfirmLoading={isLoading}
    >
      <div className="space-y-4">
        {/* Name */}
        <div>
          <label htmlFor="cat-name" className="mb-1 block text-sm font-medium text-slate-700">
            이름
          </label>
          <input
            id="cat-name"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="카테고리 이름"
            className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
          />
        </div>

        {/* Color */}
        <div>
          <label className="mb-1 block text-sm font-medium text-slate-700">색상</label>
          <div className="flex flex-wrap gap-2">
            {PRESET_COLORS.map((c) => (
              <button
                key={c}
                type="button"
                onClick={() => setColor(c)}
                className={`h-8 w-8 rounded-full border-2 transition-transform ${
                  color === c ? 'scale-110 border-slate-900' : 'border-transparent hover:scale-105'
                }`}
                style={{ backgroundColor: c }}
                aria-label={`색상 ${c}`}
              />
            ))}
          </div>
        </div>

        {/* Icon */}
        <div>
          <label htmlFor="cat-icon" className="mb-1 block text-sm font-medium text-slate-700">
            아이콘 (선택)
          </label>
          <input
            id="cat-icon"
            type="text"
            value={icon}
            onChange={(e) => setIcon(e.target.value)}
            placeholder="이모지 또는 아이콘 이름"
            className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
          />
        </div>
      </div>
    </Modal>
  );
}
