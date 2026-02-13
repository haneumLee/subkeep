import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: '시뮬레이션 - Subkeep',
  description: '구독 해지 및 추가 시 예상 금액을 미리 확인하세요',
};

export default function SimulationLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-gray-50">
      <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">시뮬레이션</h1>
          <p className="mt-2 text-gray-600">
            구독을 해지하거나 추가했을 때의 금액 변화를 미리 확인하세요
          </p>
        </div>
        {children}
      </div>
    </div>
  );
}
