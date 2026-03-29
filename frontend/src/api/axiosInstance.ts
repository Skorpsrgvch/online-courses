import axios, { AxiosInstance } from 'axios';
import { clearAuthTokens } from '@/utils/security';

const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.PROD 
    ? 'https://api.healthplatform.com/api' 
    : 'http://localhost:8080/api',
  withCredentials: true, // Критично для httpOnly cookies
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

// Перехватчик запросов для добавления заголовков
apiClient.interceptors.request.use(
  (config) => {
    // Дополнительные заголовки безопасности
    config.headers['X-Requested-With'] = 'XMLHttpRequest';
    
    // Для development можно добавить заголовок, но в production используется только httpOnly
    if (import.meta.env.DEV && config.url?.includes('/api')) {
      config.headers['Authorization'] = `Bearer ${document.cookie.replace(/(?:(?:^|.*;\s*)token\s*=\s*([^;]*).*$)|^.*$/, "$1")}`;
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Перехватчик ответов для обработки ошибок и безопасности
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    // Защита от XSS: очистка данных перед обработкой
    if (error.response?.data) {
      error.response.data = sanitizeResponse(error.response.data);
    }
    
    // Обработка 401 ошибки (неавторизованный доступ)
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        // Попытка обновить токен (если есть refresh эндпоинт)
        await apiClient.post('/auth/refresh');
        return apiClient(originalRequest);
      } catch (refreshError) {
        // Если обновление токена не удалось - выход из системы
        clearAuthTokens();
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }
    
    // Обработка 403 ошибки (доступ запрещен)
    if (error.response?.status === 403) {
      window.location.href = '/';
    }
    
    return Promise.reject(error);
  }
);

// Функция очистки ответа от потенциально опасного контента
function sanitizeResponse(data: any): any {
  if (typeof data !== 'object' || data === null) return data;
  
  const cleanData = Array.isArray(data) ? [] : {};
  
  for (const key in data) {
    if (Object.prototype.hasOwnProperty.call(data, key)) {
      const value = data[key];
      
      if (typeof value === 'string') {
        // Очистка строк от XSS
        cleanData[key] = DOMPurify.sanitize(value, { 
          ALLOWED_TAGS: ['b', 'i', 'em', 'strong', 'p', 'br', 'ul', 'ol', 'li'],
          FORBID_ATTR: ['style', 'on*']
        });
      } else if (typeof value === 'object' && value !== null) {
        // Рекурсивная очистка вложенных объектов
        cleanData[key] = sanitizeResponse(value);
      } else {
        cleanData[key] = value;
      }
    }
  }
  
  return cleanData;
}

export default apiClient;