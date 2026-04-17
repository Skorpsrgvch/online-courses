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
  refresh_token: string;
  expires_in: number;
}

const ACCESS_TOKEN_KEY = 'access_token';
const REFRESH_TOKEN_KEY = 'refresh_token';

export function getAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY);
}

function storeTokens(access: string, refresh: string) {
  localStorage.setItem(ACCESS_TOKEN_KEY, access);
  localStorage.setItem(REFRESH_TOKEN_KEY, refresh);
}

function clearTokens() {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
}

export const authService = {
  login: async (data: LoginDto): Promise<AuthResponse & TokensResponse> => {
    const response = await apiClient.post<AuthResponse & TokensResponse>('/auth/login', data);
    const { access_token, refresh_token } = response.data;
    storeTokens(access_token, refresh_token);
    return response.data;
  },

  register: async (data: RegisterDto): Promise<void> => {
    await apiClient.post('/auth/register', data);
  },

  logout: async (): Promise<void> => {
    clearTokens();
  },

  getMe: async (): Promise<AuthResponse> => {
    // rawClient без интерцепторов — чтобы 401 не вызывал redirect
    const token = getAccessToken();
    const response = await rawClient.get<AuthResponse>('/auth/me', {
      headers: { Authorization: `Bearer ${token}` },
    });
    return response.data;
  },

  refreshToken: async (): Promise<TokensResponse> => {
    const refresh = getRefreshToken();
    if (!refresh) throw new Error('No refresh token');
    const response = await rawClient.post<TokensResponse>('/auth/refresh', null, {
      headers: { Authorization: `Bearer ${refresh}` },
    });
    const { access_token, refresh_token } = response.data;
    storeTokens(access_token, refresh_token);
    return response.data;
  },

  clearTokens,
};
