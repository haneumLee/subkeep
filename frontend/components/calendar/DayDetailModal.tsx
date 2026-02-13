'use client';

import { Modal } from '@/components/ui/Modal';
import { formatCurrency } from '@/lib/utils';
import type { DayDetail } from '@/types';

interface DayDetailModalProps {
  dayDetail: DayDetail;
  isOpen: boolean;
  onClose: () => void;
}

function formatDetailDate(dateString: string): string {
  const date = new Date(dateString);
  return `${date.getMonth() + 1}월 ${date.getDate()}일 결제 상세`;
}

export function DayDetailModal({ dayDetail, isOpen, onClose }: DayDetailModalProps) {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={formatDetailDate(dayDetail.date)}
      showFooter={false}
    >
      <div className="space-y-3">
        <div className="flex items-center justify-between border-b border-slate-100 pb-2">
          <span className="text-sm text-slate-500">총 결제 금액</span>
          <span className="text-lg font-bold text-slate-900">
            {formatCurrency(dayDetail.totalAmount)}
          </span>
        </div>

        <ul className="space-y-2">
          {dayDetail.subscriptions.map((sub) => (
            <li
              key={sub.subscriptionId}
              className="flex items-center justify-between rounded-lg bg-slate-50 p-3"
            >
              <div className="flex items-center gap-2">
                <span
                  className="inline-block h-3 w-3 rounded-full"
                  style={{ backgroundColor: sub.categoryColor || '#6366f1' }}
                />
                <div>
                  <p className="text-sm font-medium text-slate-900">{sub.serviceName}</p>
                  <p className="text-xs text-slate-500">{sub.categoryName}</p>
                </div>
              </div>
              <div className="text-right">
                <p className="text-sm font-semibold text-slate-900">
                  {formatCurrency(sub.amount)}
                </p>
                {sub.personalAmount !== sub.amount && (
                  <p className="text-xs text-primary-600">
                    내 부담: {formatCurrency(sub.personalAmount)}
                  </p>
                )}
              </div>
            </li>
          ))}
        </ul>
      </div>
    </Modal>
  );
}
