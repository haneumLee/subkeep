'use client';

import { formatCurrency } from '@/lib/utils';
import type { SplitType, SubscriptionShare } from '@/types';
import { cn } from '@/lib/utils';

interface ShareBadgeProps {
  share: SubscriptionShare;
  myShareAmount: number;
  className?: string;
}

export function ShareBadge({ share, myShareAmount, className }: ShareBadgeProps) {
  const splitTypeIcons: Record<SplitType, string> = {
    equal: '=',
    custom_amount: '₩',
    custom_ratio: '%',
  };

  const splitTypeLabels: Record<SplitType, string> = {
    equal: '균등',
    custom_amount: '금액',
    custom_ratio: '비율',
  };

  return (
    <div
      className={cn(
        'inline-flex items-center gap-2 rounded-lg bg-gradient-to-r from-blue-50 to-indigo-50 px-3 py-1.5 text-sm',
        className
      )}
    >
      <div className="flex items-center gap-1.5">
        <svg className="h-4 w-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
          />
        </svg>
        <span className="font-medium text-blue-700">
          {share.shareGroup?.name || '공유 그룹'}
        </span>
      </div>

      <div className="h-4 w-px bg-blue-200" />

      <div className="flex items-center gap-1">
        <span className="inline-flex h-5 w-5 items-center justify-center rounded bg-blue-100 text-xs font-bold text-blue-700">
          {splitTypeIcons[share.splitType]}
        </span>
        <span className="text-xs text-blue-600">{splitTypeLabels[share.splitType]}</span>
      </div>

      <div className="h-4 w-px bg-blue-200" />

      <div className="font-semibold text-blue-800">{formatCurrency(myShareAmount)}</div>
    </div>
  );
}
