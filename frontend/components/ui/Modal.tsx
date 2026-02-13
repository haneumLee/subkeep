'use client';

import { useEffect, type ReactNode } from 'react';

import { cn } from '@/lib/utils';

import { Button } from './Button';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: ReactNode;
  confirmText?: string;
  cancelText?: string;
  onConfirm?: () => void;
  isConfirmLoading?: boolean;
  confirmVariant?: 'primary' | 'danger';
  showFooter?: boolean;
}

export function Modal({
  isOpen,
  onClose,
  title,
  children,
  confirmText = '확인',
  cancelText = '취소',
  onConfirm,
  isConfirmLoading = false,
  confirmVariant = 'primary',
  showFooter = true,
}: ModalProps) {
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener('keydown', handleEscape);
      document.body.style.overflow = 'hidden';
    }

    return () => {
      document.removeEventListener('keydown', handleEscape);
      document.body.style.overflow = 'unset';
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Overlay */}
      <div
        className="absolute inset-0 bg-black bg-opacity-50 transition-opacity"
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Modal */}
      <div
        className={cn(
          'relative z-10 w-full max-w-md rounded-xl bg-white p-6 shadow-xl',
          'transform transition-all'
        )}
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-title"
      >
        {/* Header */}
        <div className="mb-4">
          <h3
            id="modal-title"
            className="text-lg font-semibold text-slate-900"
          >
            {title}
          </h3>
        </div>

        {/* Content */}
        <div className="mb-6 text-sm text-slate-600">{children}</div>

        {/* Footer */}
        {showFooter && (
          <div className="flex justify-end gap-3">
            <Button variant="secondary" onClick={onClose} disabled={isConfirmLoading}>
              {cancelText}
            </Button>
            {onConfirm && (
              <Button
                variant={confirmVariant}
                onClick={onConfirm}
                isLoading={isConfirmLoading}
              >
                {confirmText}
              </Button>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
