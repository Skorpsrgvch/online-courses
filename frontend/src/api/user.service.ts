import apiClient from './axiosInstance';
import { User } from './types';

export interface UpdateProfileDto {
  name?: string;
  email?: string;
  avatar?: File | null;
}

export interface ChangePasswordDto {
  currentPassword: string;
  newPassword: string;
  confirmPassword: string;
}

export const userService = {
  /**
   * Получение данных текущего пользователя
   */
  getProfile: async (): Promise<User> => {
    const response = await apiClient.get<User>('/user/profile');
    return response.data;
  },

  /**
   * Обновление профиля (имя, email, аватар)
   */
  updateProfile: async (data: UpdateProfileDto): Promise<User> => {
    const formData = new FormData();
    if (data.name) formData.append('name', data.name);
    if (data.email) formData.append('email', data.email);
    if (data.avatar) formData.append('avatar', data.avatar);

    const response = await apiClient.put<User>('/user/profile', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },

  /**
   * Смена пароля
   */
  changePassword: async (data: ChangePasswordDto): Promise<void> => {
    if (data.newPassword !== data.confirmPassword) {
      throw new Error('Новые пароли не совпадают');
    }
    await apiClient.put('/user/change-password', {
      current_password: data.currentPassword,
      new_password: data.newPassword,
    });
  },

  /**
   * Удаление аккаунта (требует подтверждения на бэкенде)
   */
  deleteAccount: async (): Promise<void> => {
    await apiClient.delete('/user/account');
  },
  
  /**
   * Загрузка настроек приватности
   */
  getPrivacySettings: async (): Promise<Record<string, boolean>> => {
    const response = await apiClient.get<Record<string, boolean>>('/user/privacy-settings');
    return response.data;
  },

  /**
   * Обновление настроек приватности
   */
  updatePrivacySettings: async (settings: Record<string, boolean>): Promise<void> => {
    await apiClient.put('/user/privacy-settings', settings);
  }
};