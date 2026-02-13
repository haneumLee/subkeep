import {
  useCategories,
  useCreateCategory,
  useUpdateCategory,
  useDeleteCategory,
} from '@/lib/hooks/useCategories';

describe('lib/hooks/useCategories', () => {
  describe('exports', () => {
    it('should export useCategories as a function', () => {
      expect(typeof useCategories).toBe('function');
    });

    it('should export useCreateCategory as a function', () => {
      expect(typeof useCreateCategory).toBe('function');
    });

    it('should export useUpdateCategory as a function', () => {
      expect(typeof useUpdateCategory).toBe('function');
    });

    it('should export useDeleteCategory as a function', () => {
      expect(typeof useDeleteCategory).toBe('function');
    });
  });
});
