import { render, screen, waitFor } from '@testing-library/react';

import { AuthProvider, useAuth } from '@/contexts/AuthContext';

jest.mock('@/lib/api', () => ({
  get: jest.fn(),
  getAccessToken: jest.fn(),
  clearTokens: jest.fn(),
}));

import { getAccessToken } from '@/lib/api';

const mockedGetAccessToken = getAccessToken as jest.MockedFunction<typeof getAccessToken>;

describe('contexts/AuthContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('AuthProvider', () => {
    it('should render children correctly', async () => {
      mockedGetAccessToken.mockReturnValue(null);

      render(
        <AuthProvider>
          <div data-testid="child">Hello</div>
        </AuthProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('child')).toBeInTheDocument();
      });

      expect(screen.getByTestId('child')).toHaveTextContent('Hello');
    });

    it('should have user as null when there is no token', async () => {
      mockedGetAccessToken.mockReturnValue(null);

      function TestConsumer() {
        const { user, isAuthenticated } = useAuth();
        return (
          <div>
            <span data-testid="user">{user === null ? 'null' : 'exists'}</span>
            <span data-testid="authenticated">{String(isAuthenticated)}</span>
          </div>
        );
      }

      render(
        <AuthProvider>
          <TestConsumer />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('user')).toHaveTextContent('null');
      });

      expect(screen.getByTestId('authenticated')).toHaveTextContent('false');
    });
  });

  describe('useAuth', () => {
    it('should throw an error when used outside of AuthProvider', () => {
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});

      function BadConsumer() {
        useAuth();
        return <div>should not render</div>;
      }

      expect(() => render(<BadConsumer />)).toThrow(
        'useAuth must be used within an AuthProvider'
      );

      consoleSpy.mockRestore();
    });
  });
});
