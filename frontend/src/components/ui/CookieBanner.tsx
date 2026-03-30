import React, { useState } from 'react';
import { useLocalStorage } from '../../hooks/useLocalStorage';

export const CookieBanner: React.FC = () => {
  const [cookiesAccepted, setCookiesAccepted] = useLocalStorage('cookies_accepted', false);
  const [preferences, setPreferences] = useState({ necessary: true, analytics: false, marketing: false });

  if (cookiesAccepted) return null;

  const handleAccept = (selected: typeof preferences) => {
    setCookiesAccepted(true);
    // Здесь можно инициировать загрузку скриптов аналитики
    console.log('Cookie preferences:', selected);
  };

  return (
    <div className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 p-4 shadow-lg z-40 transform transition-transform duration-300">
      <div className="max-w-7xl mx-auto flex flex-col md:flex-row items-center justify-between gap-4">
        <p className="text-sm text-gray-600">
          Мы используем cookie для улучшения работы сайта. Вы можете выбрать категории или принять все.
        </p>
        <div className="flex gap-2">
          <button 
            onClick={() => handleAccept({ necessary: true, analytics: false, marketing: false })}
            className="px-4 py-2 text-sm text-gray-600 border border-gray-300 rounded hover:bg-gray-50"
          >
            Только необходимые
          </button>
          <button 
            onClick={() => handleAccept(preferences)}
            className="px-4 py-2 text-sm bg-rose-500 text-white rounded hover:bg-rose-600"
          >
            Принять выбранные
          </button>
          <button 
            onClick={() => handleAccept({ necessary: true, analytics: true, marketing: true })}
            className="px-4 py-2 text-sm bg-gray-800 text-white rounded hover:bg-gray-900"
          >
            Принять все
          </button>
        </div>
      </div>
      {/* Чекбоксы настроек можно вынести в модальное окно при клике на "Настройки" */}
    </div>
  );
};