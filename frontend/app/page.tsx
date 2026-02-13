'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/Button';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';

const OAUTH_PROVIDERS = [
  {
    id: 'google' as const,
    name: 'Google',
    color: 'bg-white text-slate-700 border border-slate-300 hover:bg-slate-50',
    icon: (
      <svg className="h-5 w-5" viewBox="0 0 24 24">
        <path
          fill="#4285F4"
          d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
        />
        <path
          fill="#34A853"
          d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
        />
        <path
          fill="#FBBC05"
          d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
        />
        <path
          fill="#EA4335"
          d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
        />
      </svg>
    ),
  },
  {
    id: 'apple' as const,
    name: 'Apple',
    color: 'bg-black text-white hover:bg-slate-800',
    icon: (
      <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
        <path d="M17.05 20.28c-.98.95-2.05.88-3.08.4-1.09-.5-2.08-.48-3.24 0-1.44.62-2.2.44-3.06-.4C2.79 15.25 3.51 7.59 9.05 7.31c1.35.07 2.29.74 3.08.8 1.18-.24 2.31-.93 3.57-.84 1.51.12 2.65.72 3.4 1.8-3.12 1.87-2.38 5.98.48 7.13-.57 1.5-1.31 2.99-2.54 4.09l.01-.01zM12.03 7.25c-.15-2.23 1.66-4.07 3.74-4.25.29 2.58-2.34 4.5-3.74 4.25z" />
      </svg>
    ),
  },
  {
    id: 'naver' as const,
    name: 'Naver',
    color: 'bg-[#03C75A] text-white hover:bg-[#02b350]',
    icon: (
      <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
        <path d="M16.273 12.845L7.376 0H0v24h7.726V11.156L16.624 24H24V0h-7.727v12.845z" />
      </svg>
    ),
  },
  {
    id: 'kakao' as const,
    name: 'Kakao',
    color: 'bg-[#FEE500] text-[#000000] hover:bg-[#f5dc00]',
    icon: (
      <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
        <path d="M12 3c5.799 0 10.5 3.664 10.5 8.185 0 4.52-4.701 8.184-10.5 8.184a13.5 13.5 0 0 1-1.727-.11l-4.408 2.883c-.501.265-.678.236-.472-.413l.892-3.678c-2.88-1.46-4.785-3.99-4.785-6.866C1.5 6.665 6.201 3 12 3zm5.907 8.06l1.47-1.424a.472.472 0 0 0-.656-.678l-1.928 1.866V9.282a.472.472 0 0 0-.944 0v2.557a.471.471 0 0 0 0 .222V13.5a.472.472 0 0 0 .944 0v-.602l.087-.085 2.096 2.84a.472.472 0 1 0 .757-.559l-1.826-2.476zm-5.39 1.556a.472.472 0 0 0 .472-.472v-4.31a.472.472 0 0 0-.944 0v4.31c0 .261.211.472.472.472zm-1.441-.472V9.282a.472.472 0 0 0-.472-.472H7.522a.472.472 0 1 0 0 .944h2.555v.984H7.522a.472.472 0 1 0 0 .944h2.555v.984a.472.472 0 0 0 .944 0zm9.618-2.032l-1.563 4.157a.472.472 0 0 0 .886.332l.328-.872h2.16l.328.872a.472.472 0 0 0 .886-.332l-1.563-4.157a.472.472 0 0 0-.886 0h-.576zm.944 2.703h-1.515l.758-2.014.757 2.014z" />
      </svg>
    ),
  },
];

export default function LoginPage() {
  const router = useRouter();
  const { isAuthenticated, isLoading } = useAuth();

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.push('/dashboard');
    }
  }, [isLoading, isAuthenticated, router]);

  const handleOAuthLogin = (provider: string) => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || '/api/v1';
    const redirectUri = `${window.location.origin}/auth/callback`;
    window.location.href = `${apiUrl}/auth/${provider}?redirect_uri=${encodeURIComponent(redirectUri)}`;
  };

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (isAuthenticated) {
    return null;
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100 px-4">
      <div className="w-full max-w-md">
        {/* Logo & Title */}
        <div className="mb-8 text-center">
          <h1 className="mb-2 text-4xl font-bold text-primary-600">Subkeep</h1>
          <p className="text-lg text-slate-600">
            구독 서비스를 한눈에 관리하세요
          </p>
          <p className="mt-2 text-sm text-slate-500">
            월간 구독 비용을 파악하고 최적화하여 불필요한 지출을 줄이세요
          </p>
        </div>

        {/* Login Card */}
        <div className="rounded-xl bg-white p-8 shadow-xl">
          <h2 className="mb-6 text-center text-xl font-semibold text-slate-800">
            소셜 로그인으로 시작하기
          </h2>

          <div className="space-y-3">
            {OAUTH_PROVIDERS.map((provider) => (
              <Button
                key={provider.id}
                onClick={() => handleOAuthLogin(provider.id)}
                className={`w-full ${provider.color}`}
                variant="secondary"
                size="lg"
              >
                <span className="mr-3">{provider.icon}</span>
                {provider.name}로 계속하기
              </Button>
            ))}
          </div>

          <p className="mt-6 text-center text-xs text-slate-500">
            로그인 시{' '}
            <a href="#" className="text-primary-600 hover:underline">
              서비스 이용약관
            </a>{' '}
            및{' '}
            <a href="#" className="text-primary-600 hover:underline">
              개인정보 처리방침
            </a>
            에 동의하게 됩니다
          </p>
        </div>

        {/* Features */}
        <div className="mt-8 grid grid-cols-3 gap-4 text-center">
          <div className="rounded-lg bg-white bg-opacity-70 p-4">
            <div className="mb-2 text-2xl">📊</div>
            <p className="text-xs font-medium text-slate-700">대시보드</p>
          </div>
          <div className="rounded-lg bg-white bg-opacity-70 p-4">
            <div className="mb-2 text-2xl">💰</div>
            <p className="text-xs font-medium text-slate-700">비용 최적화</p>
          </div>
          <div className="rounded-lg bg-white bg-opacity-70 p-4">
            <div className="mb-2 text-2xl">👥</div>
            <p className="text-xs font-medium text-slate-700">공유 관리</p>
          </div>
        </div>
      </div>
    </div>
  );
}
