import React from 'react';

export const renderParsedText = (text: string) => {
  if (!text) return null;

  const parts = text.split('|||').map(i => i.trim()).filter(i => i.length > 0);
  if (parts.length === 0) return null;

  // Если всего одна часть, возвращаем простой текст
  if (parts.length === 1) {
    return <p className="text-gray-700 leading-relaxed whitespace-pre-line">{parts[0]}</p>;
  }

  return (
    <div className="space-y-4">
      {/* Вступление */}
      <p className="text-gray-700 leading-relaxed font-semibold">
        {parts[0]}
      </p>

      {/* Список пунктов */}
      <ul className="space-y-3 mt-4">
        {parts.slice(1).map((item, idx) => {
          // Проверяем, начинается ли строка с дефиса (с пробелом или без)
          const startsWithDash = item.startsWith('-');

          // Очищаем строку от первого символа, если это маркер списка (- или ✅)
          let cleanItem = item;
          if (startsWithDash) {
            // Удаляем '-' и возможные пробелы после него
            cleanItem = item.replace(/^-+\s*/, '').trim();
          } else {
            // На случай, если там уже есть эмодзи, тоже чистим (опционально)
            cleanItem = item.replace(/^[✅❌]\s*/, '').trim();
          }

          return (
            <li key={idx} className="flex items-start gap-3">
              {/* Иконка */}
              <span className={`flex-shrink-0 w-6 h-6 rounded-full flex items-center justify-center mt-0.5 ${startsWithDash
                  ? 'bg-white text-gray-700' // Белый фон, серая рамка для черточки
                  : 'bg-green-100 border-transparent text-green-600' // Зеленый фон для галочки
                }`}>
                {startsWithDash ? (
                  // Черточка (горизонтальная линия)
                  <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M20 12H4" />
                  </svg>
                ) : (
                  // Галочка
                  <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
                  </svg>
                )}
              </span>

              <span className="text-gray-700 leading-relaxed text-sm md:text-base">
                {cleanItem}
              </span>
            </li>
          );
        })}
      </ul>
    </div>
  );
};

export const getShortDescription = (text: string, maxLength = 120): string => {
  if (!text) return '';
  const firstPart = text.split('|||')[0].trim();
  if (firstPart.length <= maxLength) return firstPart;
  return firstPart.substring(0, maxLength).trim() + '...';
};