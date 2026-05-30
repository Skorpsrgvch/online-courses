import React, { useState, useEffect } from 'react';

interface CookiePreferences {
  necessary: boolean;
  analytics: boolean;
  marketing: boolean;
}

const DEFAULT_PREFERENCES: CookiePreferences = {
  necessary: true,
  analytics: false,
  marketing: false,
};

const CookieBanner: React.FC = () => {
  const [isVisible, setIsVisible] = useState(false);
  const [preferences, setPreferences] = useState<CookiePreferences>(DEFAULT_PREFERENCES);
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    // Проверка наличия window для безопасности (SSR compatibility)
    if (typeof window === 'undefined') return;

    try {
      const savedPreferences = localStorage.getItem('cookiePreferences');
      
      if (savedPreferences) {
        const parsed = JSON.parse(savedPreferences) as CookiePreferences;
        // Валидация распарсенных данных
        if (parsed && typeof parsed.necessary === 'boolean') {
          setPreferences(parsed);
          // Если предпочтения уже сохранены, баннер не показываем
          setIsLoaded(true);
          return;
        } else {
          // Если данные повреждены, очищаем их
          localStorage.removeItem('cookiePreferences');
        }
      }
      
      // Если сохраненных настроек нет, показываем баннер
      setIsVisible(true);
    } catch (error) {
      console.error('Ошибка чтения настроек cookie:', error);
      localStorage.removeItem('cookiePreferences');
      setIsVisible(true);
    } finally {
      setIsLoaded(true);
    }
  }, []);

  const handleAcceptAll = () => {
    const newPreferences: CookiePreferences = {
      necessary: true,
      analytics: true,
      marketing: true,
    };
    
    saveAndClose(newPreferences);
  };

  const handleSavePreferences = () => {
    saveAndClose(preferences);
  };

  const handleRejectAll = () => {
    const newPreferences: CookiePreferences = {
      necessary: true,
      analytics: false,
      marketing: false,
    };
    saveAndClose(newPreferences);
  };

  const saveAndClose = (prefs: CookiePreferences) => {
    if (typeof window === 'undefined') return;
    
    localStorage.setItem('cookiePreferences', JSON.stringify(prefs));
    setIsVisible(false);
    
    // Инициализация скриптов только при изменении статуса с false на true
    if (prefs.analytics) initAnalytics();
    if (prefs.marketing) initMarketing();
  };

  const initAnalytics = () => {
    // Здесь код инициализации Яндекс.Метрики или Google Analytics
    // Пример: (window as any).ym?.(XXXXXX, 'hit', window.location.href);
  };

  const initMarketing = () => {
    // Здесь код инициализации пикселей VK, Facebook и т.д.
  };

  if (!isLoaded || !isVisible) return null;

  return (
    <div 
      role="alertdialog" 
      aria-labelledby="cookie-title" 
      aria-describedby="cookie-desc"
      className="fixed bottom-0 left-0 right-0 z-50 bg-white/95 backdrop-blur-sm shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.1)] border-t border-rose-100 animate-slide-up"
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div className="flex flex-col lg:flex-row gap-6 items-start lg:items-center justify-between">
          
          {/* Текстовый блок */}
          <div className="flex-1">
            <h3 id="cookie-title" className="text-lg font-bold text-gray-900 mb-2 flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-rose-500"></span>
              Мы заботимся о вашей конфиденциальности
            </h3>
            <p id="cookie-desc" className="text-sm text-gray-600 leading-relaxed max-w-3xl">
              Наш сайт использует файлы cookie для обеспечения работы сайта, аналитики и персонализации контента. 
              Вы можете выбрать, какие категории файлов cookie вы разрешаете использовать. 
              Подробнее читайте в нашей <a href="/privacy" className="text-rose-600 hover:underline font-medium">Политике конфиденциальности</a>.
            </p>
          </div>

          {/* Блок управления */}
          <div className="w-full lg:w-auto flex flex-col sm:flex-row gap-3 shrink-0">
            <button 
              onClick={handleRejectAll}
              className="px-5 py-2.5 text-sm font-medium text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors focus:ring-2 focus:ring-gray-300"
            >
              Отклонить все
            </button>
            
            <button 
              onClick={handleSavePreferences}
              className="px-5 py-2.5 text-sm font-medium text-white bg-rose-500 rounded-lg hover:bg-rose-600 shadow-md hover:shadow-lg transition-all focus:ring-2 focus:ring-rose-300"
            >
              Сохранить выбор
            </button>
            
            <button 
              onClick={handleAcceptAll}
              className="px-5 py-2.5 text-sm font-medium text-gray-800 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors focus:ring-2 focus:ring-gray-200"
            >
              Принять все
            </button>
          </div>
        </div>

        {/* Детальные настройки (аккордеон или просто чекбоксы для наглядности) */}
        <div className="mt-6 pt-6 border-t border-gray-100 grid grid-cols-1 sm:grid-cols-3 gap-4 text-sm">
          <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg opacity-75 cursor-not-allowed">
            <input type="checkbox" checked disabled className="rounded text-rose-500 focus:ring-rose-500" />
            <div>
              <span className="block font-semibold text-gray-800">Обязательные</span>
              <span className="text-xs text-gray-500">Необходимы для работы сайта</span>
            </div>
          </div>
          
          <label className="flex items-center gap-3 p-3 hover:bg-gray-50 rounded-lg cursor-pointer transition-colors">
            <input 
              type="checkbox" 
              className="rounded text-rose-500 focus:ring-rose-500 w-4 h-4"
              checked={preferences.analytics}
              onChange={(e) => setPreferences(prev => ({ ...prev, analytics: e.target.checked }))}
            />
            <div>
              <span className="block font-semibold text-gray-800">Аналитика</span>
              <span className="text-xs text-gray-500">Помогают улучшать сервис</span>
            </div>
          </label>
          
          <label className="flex items-center gap-3 p-3 hover:bg-gray-50 rounded-lg cursor-pointer transition-colors">
            <input 
              type="checkbox" 
              className="rounded text-rose-500 focus:ring-rose-500 w-4 h-4"
              checked={preferences.marketing}
              onChange={(e) => setPreferences(prev => ({ ...prev, marketing: e.target.checked }))}
            />
            <div>
              <span className="block font-semibold text-gray-800">Маркетинг</span>
              <span className="text-xs text-gray-500">Персонализированная реклама</span>
            </div>
          </label>
        </div>
      </div>
      
      {/* Глобальные стили для анимации, если они не настроены в tailwind.config.js */}
      <style>{`
        @keyframes slide-up {
          from { transform: translateY(100%); opacity: 0; }
          to { transform: translateY(0); opacity: 1; }
        }
        .animate-slide-up {
          animation: slide-up 0.4s cubic-bezier(0.16, 1, 0.3, 1) forwards;
        }
      `}</style>
    </div>
  );
};

export default CookieBanner;
