import { render, screen } from '@testing-library/react';

import { CategoryManager } from '@/components/settings/CategoryManager';

jest.mock('@/lib/hooks/useCategories', () => ({
  useCategories: jest.fn(),
  useCreateCategory: jest.fn(() => ({ mutate: jest.fn(), isPending: false })),
  useUpdateCategory: jest.fn(() => ({ mutate: jest.fn(), isPending: false })),
  useDeleteCategory: jest.fn(() => ({ mutate: jest.fn(), isPending: false })),
}));

import { useCategories } from '@/lib/hooks/useCategories';
const mockedUseCategories = useCategories as jest.MockedFunction<typeof useCategories>;

describe('components/settings/CategoryManager', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should show loading spinner when loading', () => {
    mockedUseCategories.mockReturnValue({
      data: undefined,
      isLoading: true,
    } as ReturnType<typeof useCategories>);

    render(<CategoryManager />);
    // LoadingSpinner renders a div with role=status
    const spinner = document.querySelector('[class*="animate-spin"]');
    expect(spinner).toBeTruthy();
  });

  it('should render category list', () => {
    mockedUseCategories.mockReturnValue({
      data: [
        {
          id: 'c1',
          userId: null,
          name: '엔터테인먼트',
          color: '#ef4444',
          icon: null,
          sortOrder: 0,
          isSystem: true,
          createdAt: '2025-01-01T00:00:00Z',
          updatedAt: '2025-01-01T00:00:00Z',
        },
        {
          id: 'c2',
          userId: 'u1',
          name: '개인',
          color: '#3b82f6',
          icon: null,
          sortOrder: 1,
          isSystem: false,
          createdAt: '2025-01-01T00:00:00Z',
          updatedAt: '2025-01-01T00:00:00Z',
        },
      ],
      isLoading: false,
    } as ReturnType<typeof useCategories>);

    render(<CategoryManager />);
    expect(screen.getByText('엔터테인먼트')).toBeInTheDocument();
    expect(screen.getByText('개인')).toBeInTheDocument();
    expect(screen.getByText('시스템')).toBeInTheDocument();
  });

  it('should show section title', () => {
    mockedUseCategories.mockReturnValue({
      data: [],
      isLoading: false,
    } as ReturnType<typeof useCategories>);

    render(<CategoryManager />);
    expect(screen.getByText('카테고리 관리')).toBeInTheDocument();
  });

  it('should show empty message when no categories', () => {
    mockedUseCategories.mockReturnValue({
      data: [],
      isLoading: false,
    } as ReturnType<typeof useCategories>);

    render(<CategoryManager />);
    expect(screen.getByText('등록된 카테고리가 없습니다.')).toBeInTheDocument();
  });
});
