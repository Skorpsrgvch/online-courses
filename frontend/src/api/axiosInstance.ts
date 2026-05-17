import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios';
import DOMPurify from 'dompurify';
import { getAccessToken, authService } from './auth.service';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';


const commonConfig = {
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true, 
};

export const apiClient = axios.create(commonConfig);
export const rawClient = axios.create(commonConfig);

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

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getAccessToken();
    
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

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

apiClient.interceptors.response.use(
  (response) => {
    if (response.data) {
      response.data = sanitizeData(response.data);
    }
    return response;
  },
  async (error: AxiosError) => {

    const originalRequest = error.config as (InternalAxiosRequestConfig & { _retry?: boolean }) | undefined;
    const isAuthRequest = originalRequest?.url?.includes('/auth/login') || 
                          originalRequest?.url?.includes('/auth/register');

    if (error.response?.status === 401 && originalRequest && !originalRequest._retry && !isAuthRequest) {
      // Если нет access token в localStorage, пробуем обновить через куки
      if (!getAccessToken()) {
          if (isRefreshing) {
             return new Promise((resolve, reject) => {
                failedQueue.push({ resolve, reject });
             }).then(() => apiClient(originalRequest));
          }
          
          originalRequest._retry = true;
          isRefreshing = true;

          try {
            await authService.refreshToken();
            processQueue(null);
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
      } else {
          // Токен есть, но он протух или неверен - просто чистим и редиректим
          authService.clearTokens();
          window.location.href = '/login?reason=session_expired';
          return Promise.reject(error);
      }
    }
    const serverMessage = (error.response?.data as any)?.message || 
                          (error.response?.data as any)?.error || 
                          'Произошла ошибка соединения с сервером';
    
    // Если это не ошибка сети (например, 400, 409, 500), возвращаем конкретную ошибку
    if (error.response) {
      return Promise.reject(new Error(serverMessage));
    }

    // Если ошибки сети (нет ответа от сервера)
    return Promise.reject(new Error('Нет соединения с сервером'));
  }
);

export default apiClient;