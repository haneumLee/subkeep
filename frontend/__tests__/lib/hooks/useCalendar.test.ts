import {
  useMonthlyCalendar,
  useDayDetail,
  useUpcomingPayments,
} from '@/lib/hooks/useCalendar';

describe('lib/hooks/useCalendar', () => {
  describe('exports', () => {
    it('should export useMonthlyCalendar as a function', () => {
      expect(typeof useMonthlyCalendar).toBe('function');
    });

    it('should export useDayDetail as a function', () => {
      expect(typeof useDayDetail).toBe('function');
    });

    it('should export useUpcomingPayments as a function', () => {
      expect(typeof useUpcomingPayments).toBe('function');
    });
  });
});
