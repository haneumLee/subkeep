import { render, screen } from '@testing-library/react';

import { MonthlyTrendChart } from '@/components/reports/MonthlyTrendChart';

describe('components/reports/MonthlyTrendChart', () => {
  it('should show empty message when no trends', () => {
    render(<MonthlyTrendChart trends={[]} />);
    expect(screen.getByText('데이터가 없습니다.')).toBeInTheDocument();
  });

  it('should render month labels', () => {
    const trends = [
      { year: 2025, month: 3, amount: 50000, count: 3 },
      { year: 2025, month: 4, amount: 55000, count: 4 },
      { year: 2025, month: 5, amount: 48000, count: 3 },
    ];

    render(<MonthlyTrendChart trends={trends} />);
    expect(screen.getByText('3월')).toBeInTheDocument();
    expect(screen.getByText('4월')).toBeInTheDocument();
    expect(screen.getByText('5월')).toBeInTheDocument();
  });

  it('should render the title', () => {
    render(<MonthlyTrendChart trends={[{ year: 2025, month: 1, amount: 10000, count: 1 }]} />);
    expect(screen.getByText('월별 추이')).toBeInTheDocument();
  });
});
