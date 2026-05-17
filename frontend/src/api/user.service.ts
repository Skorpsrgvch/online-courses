import apiClient from './axiosInstance';
import type { UserProfile, UserCoursesResponse } from './types';

export interface UpdateProfileDto {
  name?: string;
  email?: string;
}

export interface ChangePasswordDto {
  currentPassword: string;
  newPassword: string;
}

export const userService = {
  /**
   * Получение данных текущего пользователя
   */
  getProfile: async (): Promise<UserProfile> => {
    const response = await apiClient.get<UserProfile>('/user/profile');
    return response.data;
  },

  /**
   * Получение курсов пользователя с прогрессом
   */
  getCourses: async (): Promise<UserCoursesResponse> => {
    const response = await apiClient.get<UserCoursesResponse>('/user/courses');
    return response.data;
  },

  // НОВЫЙ МЕТОД
  updateProfile: async (data: UpdateProfileDto): Promise<{ message: string }> => {
    const response = await apiClient.put('/user/profile', data);
    return response.data;
  },
  
  // Запрос кода (вызывает ваш эндпоинт forgot-password)
  requestPasswordReset: async (email: string): Promise<void> => {
    await apiClient.post('/auth/forgot-password', { email });
  },

  // Подтверждение кода и смена пароля (вызывает reset-password)
  confirmPasswordReset: async (code: string, newPassword: string): Promise<void> => {
    // Обратите внимание: бэкенд может ожидать поле "code" или "token"
    await apiClient.post('/auth/reset-password', { 
      code: code, 
      new_password: newPassword 
    });
  },
};
