import { render, screen } from '@testing-library/react';

import { CostSummaryCards } from '@/components/reports/CostSummaryCards';

describe('components/reports/CostSummaryCards', () => {
  it('should render all three average cost labels', () => {
    const averageCost = {
      weekly: 12500,
      monthly: 50000,
      annual: 600000,
    };

    render(<CostSummaryCards averageCost={averageCost} />);
    expect(screen.getByText('주간 평균')).toBeInTheDocument();
    expect(screen.getByText('월간 평균')).toBeInTheDocument();
    expect(screen.getByText('연간 평균')).toBeInTheDocument();
  });
});
