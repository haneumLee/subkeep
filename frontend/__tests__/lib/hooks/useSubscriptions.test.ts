import {
  useSubscriptions,
  useSubscription,
  useCreateSubscription,
  useUpdateSubscription,
  useDeleteSubscription,
  useUpdateSatisfaction,
} from '@/lib/hooks/useSubscriptions';

describe('lib/hooks/useSubscriptions', () => {
  describe('exports', () => {
    it('should export useSubscriptions as a function', () => {
      expect(typeof useSubscriptions).toBe('function');
    });

    it('should export useSubscription as a function', () => {
      expect(typeof useSubscription).toBe('function');
    });

    it('should export useCreateSubscription as a function', () => {
      expect(typeof useCreateSubscription).toBe('function');
    });

    it('should export useUpdateSubscription as a function', () => {
      expect(typeof useUpdateSubscription).toBe('function');
    });

    it('should export useDeleteSubscription as a function', () => {
      expect(typeof useDeleteSubscription).toBe('function');
    });

    it('should export useUpdateSatisfaction as a function', () => {
      expect(typeof useUpdateSatisfaction).toBe('function');
    });
  });
});
