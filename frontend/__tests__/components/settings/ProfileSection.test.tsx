import { render, screen } from '@testing-library/react';

import { ProfileSection } from '@/components/settings/ProfileSection';

jest.mock('@/contexts/AuthContext', () => ({
  useAuth: jest.fn(),
}));

import { useAuth } from '@/contexts/AuthContext';
const mockedUseAuth = useAuth as jest.MockedFunction<typeof useAuth>;

describe('components/settings/ProfileSection', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render nothing when user is null', () => {
    mockedUseAuth.mockReturnValue({
      user: null,
      isAuthenticated: false,
      isLoading: false,
      logout: jest.fn(),
    });

    const { container } = render(<ProfileSection />);
    expect(container.firstChild).toBeNull();
  });

  it('should render user nickname and email', () => {
    mockedUseAuth.mockReturnValue({
      user: {
        id: 'u1',
        provider: 'google',
        providerUserId: 'g123',
        email: 'test@example.com',
        nickname: '테스트유저',
        avatarUrl: null,
        createdAt: '2025-01-01T00:00:00Z',
        updatedAt: '2025-01-01T00:00:00Z',
        lastLoginAt: '2025-05-01T00:00:00Z',
      },
      isAuthenticated: true,
      isLoading: false,
      logout: jest.fn(),
    });

    render(<ProfileSection />);
    expect(screen.getByText('테스트유저')).toBeInTheDocument();
    expect(screen.getByText('test@example.com')).toBeInTheDocument();
  });

  it('should render login provider label', () => {
    mockedUseAuth.mockReturnValue({
      user: {
        id: 'u1',
        provider: 'kakao',
        providerUserId: 'k123',
        email: null,
        nickname: '카카오유저',
        avatarUrl: null,
        createdAt: '2025-01-01T00:00:00Z',
        updatedAt: '2025-01-01T00:00:00Z',
        lastLoginAt: null,
      },
      isAuthenticated: true,
      isLoading: false,
      logout: jest.fn(),
    });

    render(<ProfileSection />);
    expect(screen.getByText('카카오')).toBeInTheDocument();
  });

  it('should render section title', () => {
    mockedUseAuth.mockReturnValue({
      user: {
        id: 'u1',
        provider: 'google',
        providerUserId: 'g123',
        email: 'test@example.com',
        nickname: '유저',
        avatarUrl: null,
        createdAt: '2025-01-01T00:00:00Z',
        updatedAt: '2025-01-01T00:00:00Z',
        lastLoginAt: null,
      },
      isAuthenticated: true,
      isLoading: false,
      logout: jest.fn(),
    });

    render(<ProfileSection />);
    expect(screen.getByText('프로필 정보')).toBeInTheDocument();
  });
});
