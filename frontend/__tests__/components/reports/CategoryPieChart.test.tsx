import { render, screen } from '@testing-library/react';

import { CategoryPieChart } from '@/components/reports/CategoryPieChart';

describe('components/reports/CategoryPieChart', () => {
  it('should show empty message when no categories', () => {
    render(<CategoryPieChart categories={[]} />);
    expect(screen.getByText('데이터가 없습니다.')).toBeInTheDocument();
  });

  it('should render category names and amounts', () => {
    const categories = [
      {
        categoryId: 'c1',
        categoryName: '엔터테인먼트',
        color: '#ef4444',
        monthlyAmount: 30000,
        percentage: 60,
        count: 2,
      },
      {
        categoryId: 'c2',
        categoryName: '생산성',
        color: '#3b82f6',
        monthlyAmount: 20000,
        percentage: 40,
        count: 1,
      },
    ];

    render(<CategoryPieChart categories={categories} />);
    expect(screen.getByText('엔터테인먼트')).toBeInTheDocument();
    expect(screen.getByText('생산성')).toBeInTheDocument();
    expect(screen.getByText('카테고리별 지출')).toBeInTheDocument();
  });
});
