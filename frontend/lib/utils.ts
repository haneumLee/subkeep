import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

import type { BillingCycle } from '@/types';

/** Tailwind class merge helper */
export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}

/** 금액을 원화 형식으로 포맷 */
export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('ko-KR', {
    style: 'currency',
    currency: 'KRW',
    maximumFractionDigits: 0,
  }).format(amount);
}

/** 월 환산 금액 계산 (FRS F-03) */
export function calculateMonthlyAmount(amount: number, cycle: BillingCycle): number {
  switch (cycle) {
    case 'monthly':
      return amount;
    case 'yearly':
      return Math.round(amount / 12);
    case 'weekly':
      return Math.round((amount * 52) / 12);
    default:
      return amount;
  }
}

/** 연 환산 금액 계산 */
export function calculateAnnualAmount(amount: number, cycle: BillingCycle): number {
  return calculateMonthlyAmount(amount, cycle) * 12;
}

/** 날짜를 YYYY-MM-DD 형식으로 포맷 */
export function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('ko-KR', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
}

/** 결제 주기 한국어 표시 */
export function billingCycleLabel(cycle: BillingCycle): string {
  const labels: Record<BillingCycle, string> = {
    weekly: '주간',
    monthly: '월간',
    yearly: '연간',
  };
  return labels[cycle];
}

/** 만족도 점수를 별 문자열로 변환 */
export function satisfactionStars(score: number | null): string {
  if (score === null) return '미평가';
  return '★'.repeat(score) + '☆'.repeat(5 - score);
}
