'use client';

import { formatCurrency } from '@/lib/utils';
import type { UpcomingPayment } from '@/types';

interface UpcomingPaymentsListProps {
  payments: UpcomingPayment[];
}

export function UpcomingPaymentsList({ payments }: UpcomingPaymentsListProps) {
  if (payments.length === 0) return null;

  return (
    <div className="rounded-xl bg-white p-4 shadow-sm lg:p-6">
      <h3 className="mb-4 text-lg font-semibold text-slate-900">다가오는 결제</h3>
      <ul className="space-y-3">
        {payments.map((payment) => (
          <li
            key={`${payment.subscriptionId}-${payment.date}`}
            className="flex items-center justify-between rounded-lg border border-slate-100 p-3"
          >
            <div className="flex items-center gap-3">
              <span
                className="inline-block h-3 w-3 rounded-full"
                style={{ backgroundColor: payment.categoryColor || '#6366f1' }}
              />
              <div>
                <p className="text-sm font-medium text-slate-900">{payment.serviceName}</p>
                <p className="text-xs text-slate-500">{payment.categoryName}</p>
              </div>
            </div>
            <div className="text-right">
              <p className="text-sm font-semibold text-slate-900">
                {formatCurrency(payment.amount)}
              </p>
              <p className={`text-xs ${payment.daysUntil <= 3 ? 'text-red-500 font-medium' : 'text-slate-500'}`}>
                {payment.daysUntil === 0
                  ? '오늘'
                  : payment.daysUntil === 1
                    ? '내일'
                    : `${payment.daysUntil}일 후`}
              </p>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
