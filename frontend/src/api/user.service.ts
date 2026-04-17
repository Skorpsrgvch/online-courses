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

  /**
   * Обновление профиля (пока заглушка — бэкенд не реализован)
   */
  updateProfile: async (_data: UpdateProfileDto): Promise<UserProfile> => {
    throw new Error('Обновление профиля пока не реализовано');
  },

  /**
   * Смена пароля (пока заглушка — бэкенд не реализован)
   */
  changePassword: async (_data: ChangePasswordDto): Promise<void> => {
    throw new Error('Смена пароля пока не реализована');
  },

  /**
   * Удаление аккаунта (пока заглушка — бэкенд не реализован)
   */
  deleteAccount: async (): Promise<void> => {
    throw new Error('Удаление аккаунта пока не реализовано');
  },
};
