import { render, screen } from '@testing-library/react';

import { ReportSummaryPanel } from '@/components/reports/ReportSummaryPanel';

describe('components/reports/ReportSummaryPanel', () => {
  const summary = {
    totalSubscriptions: 5,
    activeCount: 3,
    pausedCount: 2,
    mostExpensive: 'Netflix',
    mostExpensiveAmount: 17000,
    averageSatisfaction: 4,
  };

  it('should render subscription counts', () => {
    render(<ReportSummaryPanel summary={summary} />);
    expect(screen.getByText('5개')).toBeInTheDocument();
    expect(screen.getByText('3개')).toBeInTheDocument();
    expect(screen.getByText('2개')).toBeInTheDocument();
  });

  it('should render most expensive subscription name', () => {
    render(<ReportSummaryPanel summary={summary} />);
    expect(screen.getByText('Netflix')).toBeInTheDocument();
  });

  it('should render section title', () => {
    render(<ReportSummaryPanel summary={summary} />);
    expect(screen.getByText('구독 요약')).toBeInTheDocument();
  });

  it('should show 미평가 when averageSatisfaction is 0', () => {
    const noSatisfaction = { ...summary, averageSatisfaction: 0 };
    render(<ReportSummaryPanel summary={noSatisfaction} />);
    expect(screen.getByText('미평가')).toBeInTheDocument();
  });
});
