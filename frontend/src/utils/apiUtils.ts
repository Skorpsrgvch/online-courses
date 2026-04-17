/**
 * Утилиты для работы с API.
 */

/**
 * Извлекает сообщение об ошибке из API-ответа.
 */
export function extractApiError(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }
  if (typeof error === 'object' && error !== null) {
    const obj = error as Record<string, unknown>;
    if (typeof obj.message === 'string') return obj.message;
  }
  return 'Произошла неизвестная ошибка';
}

/**
 * Проверяет, является ли статус-код ошибкой клиента (4xx).
 */
export function isClientError(status: number): boolean {
  return status >= 400 && status < 500;
}

/**
 * Проверяет, является ли статус-код ошибкой сервера (5xx).
 */
export function isServerError(status: number): boolean {
  return status >= 500 && status < 600;
}

/**
 * Форматирует число в строку с разделителями (1000 -> "1 000").
 */
export function formatNumber(num: number): string {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ');
}
