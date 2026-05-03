import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios';
import DOMPurify from 'dompurify';
import { getAccessToken, getRefreshToken, authService } from './auth.service';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Клиент без интерцепторов — для refresh-запросов
export const rawClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value: unknown) => void;
  reject: (error: unknown) => void;
}> = [];

function processQueue(error: unknown | null) {
  failedQueue.forEach((promise) => {
    if (error) promise.reject(error);
    else promise.resolve(null);
  });
  failedQueue = [];
}

// Интерцептор запроса — добавляем JWT
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getAccessToken();

     console.log('[Axios] Token present:', !!token);
     
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Очистка входящих данных от XSS
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

// Интерцептор ответа — автоматический refresh при 401
apiClient.interceptors.response.use(
  (response) => {
    if (response.data) {
      response.data = sanitizeData(response.data);
    }
    return response;
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as (InternalAxiosRequestConfig & { _retry?: boolean }) | undefined;

    // Если 401 и ещё не пробовали refresh
    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      // Если нет refresh token — сразу редирект
      if (!getRefreshToken()) {
        authService.clearTokens();
        window.location.href = '/login?reason=session_expired';
        return Promise.reject(error);
      }

      if (isRefreshing) {
        // Ждём пока другой запрос обновит токен
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(() => apiClient(originalRequest));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        await authService.refreshToken();
        processQueue(null);
        // Повторяем оригинальный запрос с новым токеном
        const newToken = getAccessToken();
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`;
        }
        return apiClient(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError);
        authService.clearTokens();
        window.location.href = '/login?reason=session_expired';
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    // Для не-401 ошибок — просто прокидываем
    const message = (error.response?.data as any)?.message || 'Произошла ошибка соединения с сервером';
    return Promise.reject(new Error(message));
  }
);

export default apiClient;
