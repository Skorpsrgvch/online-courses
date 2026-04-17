/**
 * Утилиты для работы с датами.
 */

/**
 * Форматирует дату в читаемый формат: "13 апреля 2026"
 */
export function formatDate(dateStr: string): string {
  const date = new Date(dateStr);
  if (isNaN(date.getTime())) return '—';

  return date.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
}

/**
 * Форматирует дату в компактный формат: "13.04.2026"
 */
export function formatDateShort(dateStr: string): string {
  const date = new Date(dateStr);
  if (isNaN(date.getTime())) return '—';

  return date.toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  });
}

/**
 * Возвращает относительное время: "2 дня назад", "только что"
 */
export function timeAgo(dateStr: string): string {
  const date = new Date(dateStr);
  if (isNaN(date.getTime())) return '—';

  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSec = Math.floor(diffMs / 1000);
  const diffMin = Math.floor(diffSec / 60);
  const diffHr = Math.floor(diffMin / 60);
  const diffDay = Math.floor(diffHr / 24);

  if (diffSec < 60) return 'только что';
  if (diffMin < 60) return `${diffMin} мин. назад`;
  if (diffHr < 24) return `${diffHr} ч. назад`;
  if (diffDay < 7) return `${diffDay} дн. назад`;

  return formatDateShort(dateStr);
}
