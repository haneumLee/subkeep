import { getAccessToken, setTokens, clearTokens } from '@/lib/api';

jest.mock('axios', () => {
  const mockInstance = {
    get: jest.fn(),
    post: jest.fn(),
    put: jest.fn(),
    delete: jest.fn(),
    patch: jest.fn(),
    interceptors: {
      request: { use: jest.fn() },
      response: { use: jest.fn() },
    },
    defaults: { headers: { common: {} } },
  };

  return {
    __esModule: true,
    default: {
      create: jest.fn(() => mockInstance),
      post: jest.fn(),
    },
    // Expose the mock instance for test access
    __mockInstance: mockInstance,
  };
});

// eslint-disable-next-line @typescript-eslint/no-require-imports
const { __mockInstance: mockClient } = require('axios');
// eslint-disable-next-line @typescript-eslint/no-require-imports
const apiModule = require('@/lib/api');

describe('lib/api', () => {
  beforeEach(() => {
    localStorage.clear();
    mockClient.get.mockReset();
    mockClient.post.mockReset();
    mockClient.put.mockReset();
    mockClient.delete.mockReset();
    mockClient.patch.mockReset();
  });

  describe('getAccessToken', () => {
    it('should return null when no token is stored', () => {
      expect(getAccessToken()).toBeNull();
    });

    it('should return the stored access token', () => {
      localStorage.setItem('accessToken', 'test-access-token');
      expect(getAccessToken()).toBe('test-access-token');
    });
  });

  describe('setTokens', () => {
    it('should store both access and refresh tokens in localStorage', () => {
      setTokens('my-access', 'my-refresh');
      expect(localStorage.getItem('accessToken')).toBe('my-access');
      expect(localStorage.getItem('refreshToken')).toBe('my-refresh');
    });
  });

  describe('clearTokens', () => {
    it('should remove both tokens from localStorage', () => {
      localStorage.setItem('accessToken', 'token-a');
      localStorage.setItem('refreshToken', 'token-r');

      clearTokens();

      expect(localStorage.getItem('accessToken')).toBeNull();
      expect(localStorage.getItem('refreshToken')).toBeNull();
    });
  });

  describe('get', () => {
    it('should call apiClient.get and return response.data.data', async () => {
      const mockData = { id: '1', name: 'Netflix' };
      mockClient.get.mockResolvedValue({ data: { success: true, data: mockData } });

      const result = await apiModule.get('/subscriptions');

      expect(mockClient.get).toHaveBeenCalledWith('/subscriptions', { params: undefined });
      expect(result).toEqual(mockData);
    });

    it('should pass query params to apiClient.get', async () => {
      mockClient.get.mockResolvedValue({ data: { success: true, data: [] } });

      await apiModule.get('/subscriptions', { status: 'active' });

      expect(mockClient.get).toHaveBeenCalledWith('/subscriptions', {
        params: { status: 'active' },
      });
    });
  });

  describe('post', () => {
    it('should call apiClient.post and return response.data.data', async () => {
      const mockResponse = { id: '2', serviceName: 'Spotify' };
      mockClient.post.mockResolvedValue({
        data: { success: true, data: mockResponse },
      });

      const payload = { serviceName: 'Spotify', amount: 10900 };
      const result = await apiModule.post('/subscriptions', payload);

      expect(mockClient.post).toHaveBeenCalledWith('/subscriptions', payload);
      expect(result).toEqual(mockResponse);
    });
  });

  describe('put', () => {
    it('should call apiClient.put and return response.data.data', async () => {
      const mockResponse = { id: '1', serviceName: 'Updated' };
      mockClient.put.mockResolvedValue({
        data: { success: true, data: mockResponse },
      });

      const payload = { serviceName: 'Updated' };
      const result = await apiModule.put('/subscriptions/1', payload);

      expect(mockClient.put).toHaveBeenCalledWith('/subscriptions/1', payload);
      expect(result).toEqual(mockResponse);
    });
  });

  describe('del', () => {
    it('should call apiClient.delete and return response.data.data', async () => {
      mockClient.delete.mockResolvedValue({
        data: { success: true, data: undefined },
      });

      const result = await apiModule.del('/subscriptions/1');

      expect(mockClient.delete).toHaveBeenCalledWith('/subscriptions/1');
      expect(result).toBeUndefined();
    });
  });
});
