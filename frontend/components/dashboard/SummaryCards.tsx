import { formatCurrency } from '@/lib/utils';
import type { DashboardSummary } from '@/types';

interface SummaryCardsProps {
  summary: DashboardSummary;
}

export default function SummaryCards({ summary }: SummaryCardsProps) {
  const cards = [
    {
      title: '월 총액',
      value: formatCurrency(summary.monthlyTotal),
      description: '이번 달 구독료',
      color: 'bg-blue-50 border-blue-200',
      textColor: 'text-blue-900',
    },
    {
      title: '연 총액',
      value: formatCurrency(summary.annualTotal),
      description: '연간 예상 비용',
      color: 'bg-purple-50 border-purple-200',
      textColor: 'text-purple-900',
    },
    {
      title: '활성 구독',
      value: `${summary.activeCount}개`,
      description: '현재 이용 중',
      color: 'bg-green-50 border-green-200',
      textColor: 'text-green-900',
    },
    {
      title: '일시중지',
      value: `${summary.pausedCount}개`,
      description: '중지된 구독',
      color: 'bg-gray-50 border-gray-200',
      textColor: 'text-gray-900',
    },
  ];

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {cards.map((card) => (
        <div
          key={card.title}
          className={`rounded-lg border-2 ${card.color} p-6 transition-all hover:shadow-md`}
        >
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">{card.title}</p>
              <p className={`mt-2 text-3xl font-bold ${card.textColor}`}>{card.value}</p>
              <p className="mt-1 text-xs text-gray-500">{card.description}</p>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
