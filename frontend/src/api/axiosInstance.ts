import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import DOMPurify from 'dompurify';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Критично для отправки httpOnly cookies
});

// Очистка входящих данных от потенциального XSS (защита медицинского контента)
const sanitizeData = (data: any): any => {
  if (typeof data === 'string') {
    return DOMPurify.sanitize(data);
  }
  if (Array.isArray(data)) {
    return data.map(sanitizeData);
  }
  if (data !== null && typeof data === 'object') {
    const cleanData: Record<string, any> = {};
    for (const key in data) {
      if (Object.prototype.hasOwnProperty.call(data, key)) {
        cleanData[key] = sanitizeData(data[key]);
      }
    }
    return cleanData;
  }
  return data;
};

// Интерцептор запроса
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Можно добавить логирование или дополнительную логику перед отправкой
    return config;
  },
  (error) => Promise.reject(error)
);

// Интерцептор ответа
apiClient.interceptors.response.use(
  (response) => {
    // Санификация данных, полученных от сервера, перед передачей в React
    if (response.data) {
      response.data = sanitizeData(response.data);
    }
    return response;
  },
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      // Токен истек или невалиден
      // В реальном приложении здесь можно попытаться вызвать /refresh
      // Если не удалось - редирект на логин
      window.location.href = '/login?reason=session_expired';
    }
    
    // Формируем понятную ошибку для UI
    const message = (error.response?.data as any)?.message || 'Произошла ошибка соединения с сервером';
    return Promise.reject(new Error(message));
  }
);

export default apiClient;