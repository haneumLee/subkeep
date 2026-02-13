import { render, screen } from '@testing-library/react';

import { CalendarGrid } from '@/components/calendar/CalendarGrid';

describe('components/calendar/CalendarGrid', () => {
  const defaultProps = {
    year: 2025,
    month: 5,
    days: [],
    onPrevMonth: jest.fn(),
    onNextMonth: jest.fn(),
    onDayClick: jest.fn(),
    selectedDay: 0,
  };

  it('should render the month and year header', () => {
    render(<CalendarGrid {...defaultProps} />);
    expect(screen.getByText('2025년 5월')).toBeInTheDocument();
  });

  it('should render weekday labels', () => {
    render(<CalendarGrid {...defaultProps} />);
    expect(screen.getByText('월')).toBeInTheDocument();
    expect(screen.getByText('수')).toBeInTheDocument();
    expect(screen.getByText('금')).toBeInTheDocument();
  });

  it('should render day numbers', () => {
    render(<CalendarGrid {...defaultProps} />);
    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('15')).toBeInTheDocument();
    expect(screen.getByText('31')).toBeInTheDocument();
  });

  it('should show payment amount on days with subscriptions', () => {
    const days = [
      {
        date: '2025-05-15',
        totalAmount: 15000,
        subscriptions: [
          {
            subscriptionId: 's1',
            serviceName: 'Netflix',
            amount: 15000,
            monthlyAmount: 15000,
            personalAmount: 15000,
            billingCycle: 'monthly',
            categoryName: '엔터테인먼트',
            categoryColor: '#ef4444',
            autoRenew: true,
          },
        ],
      },
    ];

    render(<CalendarGrid {...defaultProps} days={days} />);
    expect(screen.getByText('₩15,000')).toBeInTheDocument();
  });
});
