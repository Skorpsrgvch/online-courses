import apiClient from './axiosInstance';
import type { Service } from './types';

// Функция маппинга ответа от сервера в формат приложения
const mapService = (data: any): Service => {
  return {
    id: data.id,
    title: data.title,
    description: data.description,
    price: data.price,
    duration_minutes: typeof data.duration_minutes === 'number' ? data.duration_minutes : 0,
  };
};

export const servicesService = {
  getAll: async (): Promise<Service[]> => {
    const response = await apiClient.get<Service[]>('/services');
    return response.data; 
  },

  getById: async (id: number): Promise<Service> => {
    const response = await apiClient.get<Service>(`/admin/services/${id}`);
    // Применяем маппинг для гарантии корректности типов
    return mapService(response.data);
  },

  create: async (data: Omit<Service, 'id'>): Promise<void> => {
    // Формируем полезную нагрузку, явно проверяя поля
    const payload = {
      title: data.title,
      description: data.description,
      price: data.price,
      duration_minutes: data.duration_minutes || 0, // Гарантируем, что отправляем число
    };
    
    // Логируем для отладки (удалите в продакшене)
    await apiClient.post('/admin/services', payload);
  },

  update: async (id: number, data: Partial<Service>): Promise<void> => {
    // Критически важно: убедитесь, что ключи соответствуют ожиданиям бэкенда
    const payload = {
      title: data.title,
      description: data.description,
      price: data.price,
      duration_minutes: data.duration_minutes !== undefined ? data.duration_minutes : 0,
    };

    // Логируем для отладки
    await apiClient.put(`/admin/services/${id}`, payload);
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/admin/services/${id}`);
  },

  reorder: async (serviceIds: number[]): Promise<void> => {
    await apiClient.put('/admin/services/reorder', { service_ids: serviceIds });
  },
};
