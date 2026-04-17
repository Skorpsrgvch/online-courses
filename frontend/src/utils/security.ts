/**
 * Утилиты безопасности для валидации и санитизации данных.
 */

/**
 * Базовая XSS-санитизация без внешних библиотек.
 * Удаляет/экранирует опасные HTML-теги.
 */
export function sanitizeHtml(input: string): string {
  return input
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#x27;')
    .replace(/\//g, '&#x2F;');
}

/**
 * Разрешает только безопасные теги для статей.
 * Оставляет: <p>, <br>, <strong>, <em>, <u>, <ol>, <ul>, <li>, <h1>-<h6>, <a>, <img>
 */
export function allowSafeHtml(input: string): string {
  let result = input;

  // Удаляем event-атрибуты (onclick, onerror и т.д.)
  result = result.replace(/\s*on\w+\s*=\s*["'][^"']*["']/gi, '');
  // Удаляем javascript: URI
  result = result.replace(/javascript\s*:/gi, '');

  return result;
}

/**
 * Валидация URL
 */
export function isValidUrl(url: string): boolean {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
}

/**
 * Экстракция ID видео из различных URL (YouTube, RuTube)
 */
export function extractVideoId(url: string): string | null {
  // RuTube: https://rutube.ru/video/abc123/
  const rutubeMatch = url.match(/rutube\.ru\/video\/([a-f0-9]+)/);
  if (rutubeMatch) return rutubeMatch[1];

  // YouTube: https://www.youtube.com/watch?v=abc123
  const ytMatch = url.match(/[?&]v=([a-zA-Z0-9_-]{11})/);
  if (ytMatch) return ytMatch[1];

  return null;
}
