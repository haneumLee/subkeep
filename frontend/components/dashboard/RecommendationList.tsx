import Link from 'next/link';

import { formatCurrency, satisfactionStars } from '@/lib/utils';
import type { Recommendation } from '@/types';

interface RecommendationListProps {
  recommendations: Recommendation[];
}

export default function RecommendationList({ recommendations }: RecommendationListProps) {
  if (recommendations.length === 0) {
    return (
      <div className="rounded-lg border-2 border-gray-200 bg-white p-6">
        <h3 className="mb-4 text-lg font-semibold text-gray-900">해지 추천</h3>
        <div className="flex items-center justify-center py-8 text-gray-500">
          추천할 구독이 없습니다
        </div>
      </div>
    );
  }

  return (
    <div className="rounded-lg border-2 border-gray-200 bg-white p-6">
      <div className="mb-4 flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">해지 추천</h3>
        <span className="text-sm text-gray-500">{recommendations.length}건</span>
      </div>

      <div className="space-y-3">
        {recommendations.map((rec) => (
          <div
            key={rec.subscription.id}
            className="rounded-lg border border-gray-200 bg-gray-50 p-4 transition-all hover:border-red-300 hover:bg-red-50"
          >
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <h4 className="font-semibold text-gray-900">
                    {rec.subscription.serviceName}
                  </h4>
                  {rec.subscription.satisfactionScore !== null && (
                    <span className="text-sm text-yellow-600">
                      {satisfactionStars(rec.subscription.satisfactionScore)}
                    </span>
                  )}
                </div>
                <p className="mt-1 text-sm text-gray-600">{rec.reason}</p>
                <div className="mt-2 flex items-center gap-4 text-xs text-gray-500">
                  <span>월 {formatCurrency(rec.subscription.monthlyAmount)}</span>
                  {rec.subscription.category && (
                    <span className="rounded-full bg-gray-200 px-2 py-1">
                      {rec.subscription.category.name}
                    </span>
                  )}
                </div>
              </div>
              <div className="ml-4 text-right">
                <div className="text-sm font-semibold text-red-600">
                  절감 가능
                </div>
                <div className="text-lg font-bold text-red-700">
                  {formatCurrency(rec.potentialSaving)}
                </div>
                <div className="text-xs text-gray-500">/월</div>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="mt-4 text-center">
        <Link
          href="/subscriptions"
          className="text-sm font-medium text-blue-600 hover:text-blue-700"
        >
          전체 구독 관리하기 →
        </Link>
      </div>
    </div>
  );
}
