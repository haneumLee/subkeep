import {
  formatCurrency,
  calculateMonthlyAmount,
  calculateAnnualAmount,
  billingCycleLabel,
} from '@/lib/utils';

import type { BillingCycle } from '@/types';

describe('lib/utils', () => {
  describe('formatCurrency', () => {
    it('should format 1000 as Korean Won currency', () => {
      expect(formatCurrency(1000)).toBe('\u20A91,000');
    });

    it('should format 0 as Korean Won currency', () => {
      expect(formatCurrency(0)).toBe('\u20A90');
    });

    it('should format large amounts correctly', () => {
      expect(formatCurrency(1000000)).toBe('\u20A91,000,000');
    });

    it('should round decimal amounts (no fraction digits)', () => {
      expect(formatCurrency(1234.56)).toBe('\u20A91,235');
    });

    it('should format negative amounts', () => {
      const result = formatCurrency(-5000);
      expect(result).toContain('5,000');
    });
  });

  describe('calculateMonthlyAmount', () => {
    it('should return the same amount for monthly cycle', () => {
      expect(calculateMonthlyAmount(10000, 'monthly')).toBe(10000);
    });

    it('should divide by 12 and round for yearly cycle', () => {
      expect(calculateMonthlyAmount(12000, 'yearly')).toBe(1000);
      expect(calculateMonthlyAmount(10000, 'yearly')).toBe(833);
      expect(calculateMonthlyAmount(15000, 'yearly')).toBe(1250);
    });

    it('should multiply by 52 and divide by 12 and round for weekly cycle', () => {
      expect(calculateMonthlyAmount(1000, 'weekly')).toBe(Math.round((1000 * 52) / 12));
      expect(calculateMonthlyAmount(2000, 'weekly')).toBe(Math.round((2000 * 52) / 12));
    });

    it('should return the amount for unknown cycle (default case)', () => {
      expect(calculateMonthlyAmount(5000, 'daily' as BillingCycle)).toBe(5000);
    });

    it('should handle 0 amount', () => {
      expect(calculateMonthlyAmount(0, 'monthly')).toBe(0);
      expect(calculateMonthlyAmount(0, 'yearly')).toBe(0);
      expect(calculateMonthlyAmount(0, 'weekly')).toBe(0);
    });
  });

  describe('calculateAnnualAmount', () => {
    it('should return amount * 12 for monthly cycle', () => {
      expect(calculateAnnualAmount(10000, 'monthly')).toBe(120000);
    });

    it('should return the yearly amount (monthly * 12) for yearly cycle', () => {
      expect(calculateAnnualAmount(12000, 'yearly')).toBe(12000);
      expect(calculateAnnualAmount(10000, 'yearly')).toBe(833 * 12);
    });

    it('should calculate annual from weekly cycle', () => {
      const monthlyFromWeekly = Math.round((1000 * 52) / 12);
      expect(calculateAnnualAmount(1000, 'weekly')).toBe(monthlyFromWeekly * 12);
    });

    it('should handle 0 amount', () => {
      expect(calculateAnnualAmount(0, 'monthly')).toBe(0);
      expect(calculateAnnualAmount(0, 'yearly')).toBe(0);
      expect(calculateAnnualAmount(0, 'weekly')).toBe(0);
    });
  });

  describe('billingCycleLabel', () => {
    it('should return "월간" for monthly', () => {
      expect(billingCycleLabel('monthly')).toBe('월간');
    });

    it('should return "연간" for yearly', () => {
      expect(billingCycleLabel('yearly')).toBe('연간');
    });

    it('should return "주간" for weekly', () => {
      expect(billingCycleLabel('weekly')).toBe('주간');
    });
  });
});
