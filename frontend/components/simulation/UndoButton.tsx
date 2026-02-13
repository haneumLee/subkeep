'use client';

import { useEffect, useState } from 'react';

import { useUndoSimulation } from '@/lib/hooks/useSimulation';
import { cn } from '@/lib/utils';

interface UndoButtonProps {
  onUndo?: () => void;
  className?: string;
}

const UNDO_TIMEOUT_SECONDS = 30;

export function UndoButton({ onUndo, className }: UndoButtonProps) {
  const [remainingSeconds, setRemainingSeconds] = useState(UNDO_TIMEOUT_SECONDS);
  const [isVisible, setIsVisible] = useState(true);
  const undoMutation = useUndoSimulation();

  useEffect(() => {
    if (remainingSeconds <= 0) {
      setIsVisible(false);
      return;
    }

    const timer = setInterval(() => {
      setRemainingSeconds((prev) => prev - 1);
    }, 1000);

    return () => clearInterval(timer);
  }, [remainingSeconds]);

  const handleUndo = async () => {
    try {
      await undoMutation.mutateAsync();
      setIsVisible(false);
      onUndo?.();
    } catch (error) {
      console.error('Undo failed:', error);
    }
  };

  if (!isVisible) {
    return null;
  }

  return (
    <div
      className={cn(
        'fixed bottom-6 right-6 z-50 rounded-lg border border-orange-200 bg-orange-50 p-4 shadow-lg',
        className
      )}
    >
      <div className="flex items-center gap-4">
        <div className="flex-1">
          <div className="font-semibold text-orange-900">변경 사항이 적용되었습니다</div>
          <div className="text-sm text-orange-700">{remainingSeconds}초 내에 되돌릴 수 있습니다</div>
        </div>
        <button
          onClick={handleUndo}
          disabled={undoMutation.isPending}
          className={cn(
            'rounded-lg bg-orange-600 px-4 py-2 font-medium text-white transition-colors',
            'hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2',
            'disabled:cursor-not-allowed disabled:opacity-50'
          )}
        >
          {undoMutation.isPending ? '되돌리는 중...' : '되돌리기'}
        </button>
        <button
          onClick={() => setIsVisible(false)}
          className="text-orange-600 hover:text-orange-800"
          aria-label="닫기"
        >
          <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>

      {undoMutation.isError && (
        <div className="mt-2 text-sm text-red-600">
          되돌리기에 실패했습니다. 다시 시도해주세요.
        </div>
      )}
    </div>
  );
}
