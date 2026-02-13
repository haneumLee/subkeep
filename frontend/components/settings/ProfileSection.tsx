'use client';

import { useAuth } from '@/contexts/AuthContext';
import { formatDate } from '@/lib/utils';

export function ProfileSection() {
  const { user } = useAuth();

  if (!user) return null;

  const providerLabels: Record<string, string> = {
    google: 'Google',
    apple: 'Apple',
    naver: '네이버',
    kakao: '카카오',
  };

  return (
    <div className="rounded-xl bg-white p-6 shadow-sm">
      <h3 className="mb-4 text-lg font-semibold text-slate-900">프로필 정보</h3>

      <div className="space-y-4">
        <div className="flex items-center gap-4">
          {user.avatarUrl ? (
            <div className="h-16 w-16 overflow-hidden rounded-full">
              <img
                src={user.avatarUrl}
                alt={user.nickname || '프로필'}
                className="h-full w-full object-cover"
              />
            </div>
          ) : (
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-primary-100 text-2xl font-bold text-primary-600">
              {user.nickname?.[0]?.toUpperCase() || 'U'}
            </div>
          )}
          <div>
            <p className="text-lg font-medium text-slate-900">
              {user.nickname || '사용자'}
            </p>
            <p className="text-sm text-slate-500">{user.email || '이메일 없음'}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-3 border-t border-slate-100 pt-4 sm:grid-cols-2">
          <div>
            <p className="text-sm text-slate-500">로그인 방식</p>
            <p className="text-sm font-medium text-slate-900">
              {providerLabels[user.provider] || user.provider}
            </p>
          </div>
          <div>
            <p className="text-sm text-slate-500">가입일</p>
            <p className="text-sm font-medium text-slate-900">{formatDate(user.createdAt)}</p>
          </div>
          {user.lastLoginAt && (
            <div>
              <p className="text-sm text-slate-500">마지막 로그인</p>
              <p className="text-sm font-medium text-slate-900">{formatDate(user.lastLoginAt)}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
