import { render, screen } from '@testing-library/react';

import { UpcomingPaymentsList } from '@/components/calendar/UpcomingPaymentsList';

describe('components/calendar/UpcomingPaymentsList', () => {
  it('should render nothing when payments list is empty', () => {
    const { container } = render(<UpcomingPaymentsList payments={[]} />);
    expect(container.firstChild).toBeNull();
  });

  it('should render upcoming payments', () => {
    const payments = [
      {
        date: '2025-05-20',
        daysUntil: 5,
        subscriptionId: 's1',
        serviceName: 'Netflix',
        amount: 15000,
        personalAmount: 15000,
        categoryName: '엔터테인먼트',
        categoryColor: '#ef4444',
      },
    ];

    render(<UpcomingPaymentsList payments={payments} />);
    expect(screen.getByText('다가오는 결제')).toBeInTheDocument();
    expect(screen.getByText('Netflix')).toBeInTheDocument();
    expect(screen.getByText('5일 후')).toBeInTheDocument();
  });

  it('should show "오늘" for daysUntil 0', () => {
    const payments = [
      {
        date: '2025-05-15',
        daysUntil: 0,
        subscriptionId: 's2',
        serviceName: 'Spotify',
        amount: 10900,
        personalAmount: 10900,
        categoryName: '음악',
        categoryColor: '#22c55e',
      },
    ];

    render(<UpcomingPaymentsList payments={payments} />);
    expect(screen.getByText('오늘')).toBeInTheDocument();
  });

  it('should show "내일" for daysUntil 1', () => {
    const payments = [
      {
        date: '2025-05-16',
        daysUntil: 1,
        subscriptionId: 's3',
        serviceName: 'YouTube',
        amount: 14900,
        personalAmount: 14900,
        categoryName: '엔터테인먼트',
        categoryColor: '#ef4444',
      },
    ];

    render(<UpcomingPaymentsList payments={payments} />);
    expect(screen.getByText('내일')).toBeInTheDocument();
  });
});
