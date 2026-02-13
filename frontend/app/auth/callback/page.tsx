'use client';

import { Suspense, useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

import { post, setTokens } from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import type { LoginResponse } from '@/types';

export default function AuthCallbackPage() {
  return (
    <Suspense fallback={
      <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100">
        <div className="rounded-xl bg-white p-8 shadow-xl">
          <div className="text-center">
            <LoadingSpinner size="lg" className="mx-auto mb-4" />
            <p className="text-sm text-slate-600">로딩 중...</p>
          </div>
        </div>
      </div>
    }>
      <AuthCallbackContent />
    </Suspense>
  );
}

function AuthCallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { refreshUser } = useAuth();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleCallback = async () => {
      const code = searchParams.get('code');
      const provider = searchParams.get('provider') || searchParams.get('state');
      const errorParam = searchParams.get('error');

      // Dev login bypass: tokens are passed directly via query params.
      const devToken = searchParams.get('dev_token');
      const devRefresh = searchParams.get('dev_refresh');
      if (devToken && devRefresh) {
        setTokens(devToken, devRefresh);
        await refreshUser();
        router.push('/dashboard');
        return;
      }

      if (errorParam) {
        setError(`인증 실패: ${errorParam}`);
        setTimeout(() => {
          router.push('/');
        }, 3000);
        return;
      }

      if (!code || !provider) {
        setError('잘못된 인증 요청입니다');
        setTimeout(() => {
          router.push('/');
        }, 3000);
        return;
      }

      try {
        const redirectUri = `${window.location.origin}/auth/callback`;
        const response = await post<LoginResponse>(
          `/auth/${provider}/callback`,
          { code, redirectUri }
        );

        setTokens(response.tokens.accessToken, response.tokens.refreshToken);

        await refreshUser();

        router.push('/dashboard');
      } catch (err) {
        console.error('OAuth callback error:', err);
        setError('로그인 처리 중 오류가 발생했습니다');
        setTimeout(() => {
          router.push('/');
        }, 3000);
      }
    };

    handleCallback();
  }, [searchParams, router, refreshUser]);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100">
      <div className="rounded-xl bg-white p-8 shadow-xl">
        {error ? (
          <div className="text-center">
            <svg
              className="mx-auto mb-4 h-12 w-12 text-red-500"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <h2 className="mb-2 text-lg font-semibold text-slate-800">
              로그인 실패
            </h2>
            <p className="text-sm text-slate-600">{error}</p>
            <p className="mt-4 text-xs text-slate-500">
              잠시 후 로그인 페이지로 이동합니다...
            </p>
          </div>
        ) : (
          <div className="text-center">
            <LoadingSpinner size="lg" className="mx-auto mb-4" />
            <h2 className="mb-2 text-lg font-semibold text-slate-800">
              로그인 처리 중
            </h2>
            <p className="text-sm text-slate-600">잠시만 기다려주세요...</p>
          </div>
        )}
      </div>
    </div>
  );
}
