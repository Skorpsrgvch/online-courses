import apiClient, { rawClient } from './axiosInstance';
import type { AuthResponse } from './types';

export interface LoginDto {
  email: string;
  password: string;
}

export interface RegisterDto {
  email: string;
  password: string;
  full_name: string;
}

export interface TokensResponse {
  access_token: string;
  expires_in?: number;
}

const ACCESS_TOKEN_KEY = 'access_token';

export function getAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}


function storeTokens(access: string) {
  localStorage.setItem(ACCESS_TOKEN_KEY, access);

}

function clearTokens() {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
}

export const authService = {
  login: async (data: LoginDto): Promise<AuthResponse & TokensResponse> => {
    const response = await apiClient.post<AuthResponse & TokensResponse>('/auth/login', data);
    const { access_token } = response.data;
    
    if (!access_token) {
        throw new Error('No access token in response');
    }

    storeTokens(access_token);
    return response.data;
  },

  register: async (data: RegisterDto): Promise<void> => {
    await apiClient.post('/auth/register', data);
  },

  logout: async (): Promise<void> => {
    try {
      // Отправляем запрос на бэкенд, чтобы он удалил токен из БД и очистил куку
      await apiClient.post('/auth/logout'); 
    } catch (error) {
      console.error('Logout error:', error);
      // Игнорируем ошибку сети, все равно чистим локальное состояние
    } finally {
      // Всегда очищаем локальные токены и перенаправляем
      clearTokens();
      window.location.href = '/'; // Или '/login'
    }
  },

  getMe: async (): Promise<AuthResponse> => {
    const token = getAccessToken();
    if (!token) throw new Error('No access token');
    
    const response = await rawClient.get<AuthResponse>('/auth/me', {
      headers: { Authorization: `Bearer ${token}` },
    });
    return response.data;
  },

  refreshToken: async (): Promise<TokensResponse> => {

    const response = await rawClient.post<TokensResponse>('/auth/refresh', null, {

    });
    
    const { access_token } = response.data;
    if (!access_token) {
        throw new Error('No access token in refresh response');
    }

    storeTokens(access_token);
    return response.data;
  },

  clearTokens,
};