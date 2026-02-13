import { ReactNode } from 'react';

export default function ShareGroupsLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-slate-50">
      <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-slate-900">공유 그룹 관리</h1>
          <p className="mt-2 text-sm text-slate-600">
            구독을 함께 사용하는 그룹을 만들고 비용을 분담하세요
          </p>
        </div>
        {children}
      </div>
    </div>
  );
}
