import React, { useState } from 'react';
import { useLocalStorage } from '../../hooks/useLocalStorage';

export const CookieBanner: React.FC = () => {
  const [cookiesAccepted, setCookiesAccepted] = useLocalStorage('cookies_accepted', false);
  const [preferences] = useState({ necessary: true, analytics: false, marketing: false });

  if (cookiesAccepted) return null;

  const handleAccept = (_selected: typeof preferences) => {
    setCookiesAccepted(true);
  };

  return (
    <div className="fixed inset-x-0 bottom-0 z-40 max-h-[85vh] overflow-y-auto border-t border-gray-200 bg-white p-3 shadow-lg transition-transform duration-300 sm:p-4">
      <div className="mx-auto flex max-w-7xl flex-col gap-3 md:flex-row md:items-center md:justify-between md:gap-4">
        <p className="text-sm leading-relaxed text-gray-600 md:max-w-3xl">
          Мы используем cookie для улучшения работы сайта. Вы можете выбрать категории или принять все.{' '}
          <a href="/cookies" className="font-medium text-rose-600 hover:underline" style={{ textDecoration: 'none' }}>
            Подробнее о cookie
          </a>
        </p>

        <div className="grid grid-cols-1 gap-2 sm:grid-cols-3 md:flex md:shrink-0">
          <button
            onClick={() => handleAccept({ necessary: true, analytics: false, marketing: false })}
            className="min-h-10 rounded border border-gray-300 px-3 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50"
          >
            Только необходимые
          </button>
          <button
            onClick={() => handleAccept(preferences)}
            className="min-h-10 rounded bg-rose-500 px-3 py-2 text-sm font-medium text-white hover:bg-rose-600"
          >
            Принять выбранные
          </button>
          <button
            onClick={() => handleAccept({ necessary: true, analytics: true, marketing: true })}
            className="min-h-10 rounded bg-gray-800 px-3 py-2 text-sm font-medium text-white hover:bg-gray-900"
          >
            Принять все
          </button>
        </div>
      </div>
    </div>
  );
};
