import apiClient from './axiosInstance';
import { User, AuthResponse } from './types';

export interface LoginDto {
  email: string;
  password: string;
}

export interface RegisterDto {
  email: string;
  password: string;
  name: string;
  agreeToTerms: boolean; // Обязательно для ФЗ-152
}

export const authService = {
  login: async (data: LoginDto): Promise<User> => {
    const response = await apiClient.post<AuthResponse>('/auth/login', data);
    return response.data.user;
  },

  register: async (data: RegisterDto): Promise<User> => {
    if (!data.agreeToTerms) {
      throw new Error('Необходимо согласие на обработку персональных данных');
    }
    const response = await apiClient.post<AuthResponse>('/auth/register', data);
    return response.data.user;
  },

  logout: async (): Promise<void> => {
    await apiClient.post('/auth/logout');
    // Сервер очистит cookie, клиенту остается только очистить локальное состояние
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await apiClient.get<AuthResponse>('/auth/me');
    return response.data.user;
  },

  refreshToken: async (): Promise<void> => {
    await apiClient.post('/auth/refresh');
  }
};