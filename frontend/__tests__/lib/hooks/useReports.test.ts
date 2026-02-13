import { useReportOverview } from '@/lib/hooks/useReports';

describe('lib/hooks/useReports', () => {
  describe('exports', () => {
    it('should export useReportOverview as a function', () => {
      expect(typeof useReportOverview).toBe('function');
    });
  });
});
