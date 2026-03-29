import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';

interface CookiePreferences {
  necessary: boolean;
  analytics: boolean;
  marketing: boolean;
}

const CookieBanner: React.FC = () => {
  const [isVisible, setIsVisible] = useState(false);
  const [preferences, setPreferences] = useState<CookiePreferences>({
    necessary: true,
    analytics: false,
    marketing: false
  });

  useEffect(() => {
    const savedPreferences = localStorage.getItem('cookiePreferences');
    if (savedPreferences) {
      try {
        const parsed = JSON.parse(savedPreferences);
        setPreferences(parsed);
      } catch (e) {
        localStorage.removeItem('cookiePreferences');
      }
    }
    
    // Показываем баннер только если нет сохраненных предпочтений
    if (!savedPreferences) {
      setIsVisible(true);
    }
  }, []);

  const handleAcceptAll = () => {
    const newPreferences = {
      necessary: true,
      analytics: true,
      marketing: true
    };
    
    savePreferences(newPreferences);
    setIsVisible(false);
    
    // Инициализация аналитики и маркетинговых инструментов
    if (newPreferences.analytics) initAnalytics();
    if (newPreferences.marketing) initMarketing();
  };

  const handleSavePreferences = () => {
    savePreferences(preferences);
    setIsVisible(false);
    
    // Инициализация в зависимости от выбранных опций
    if (preferences.analytics) initAnalytics();
    if (preferences.marketing) initMarketing();
  };

  const savePreferences = (prefs: CookiePreferences) => {
    localStorage.setItem('cookiePreferences', JSON.stringify(prefs));
  };

  const initAnalytics = () => {
    // Инициализация аналитики (Yandex.Metrika, Google Analytics)
    console.log('Analytics initialized');
  };

  const initMarketing = () => {
    // Инициализация маркетинговых инструментов
    console.log('Marketing tools initialized');
  };

  if (!isVisible) return null;

  return (
    <motion.div
      initial={{ opacity: 0, y: 100 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: 100 }}
      transition={{ duration: 0.3 }}
      className="fixed bottom-0 left-0 right-0 bg-white shadow-2xl m-4 rounded-2xl border border-gray-100 z-50"
    >
      <div className="max-w-4xl mx-auto">
        <div className="p-6">
          <h3 className="text-xl font-bold text-gray-800 mb-2">
            Мы заботимся о вашей конфиденциальности
          </h3>
          <p className="text-gray-600 mb-4">
            Наш сайт использует файлы cookie для улучшения работы и персонализации контента. 
            Вы можете настроить свои предпочтения в любое время.
          </p>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
            <div className="border border-gray-200 rounded-xl p-4">
              <h4 className="font-semibold text-gray-800 mb-2 flex items-center">
                <svg className="w-5 h-5 mr-2 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                </h4>
                Обязательные
              </h4>
              <p className="text-gray-600 text-sm">Необходимы для работы сайта</p>
              <div className="mt-2 bg-green-100 text-green-800 text-xs px-2 py-1 rounded inline-block">
                Всегда активны
              </div>
            </div>
            
            <div className="border border-gray-200 rounded-xl p-4">
              <h4 className="font-semibold text-gray-800 mb-2">Аналитика</h4>
              <p className="text-gray-600 text-sm">Помогают улучшать сайт</p>
              <label className="inline-flex items-center mt-2">
                <input 
                  type="checkbox" 
                  className="rounded text-pink-500"
                  checked={preferences.analytics}
                  onChange={(e) => setPreferences(prev => ({ ...prev, analytics: e.target.checked }))}
                />
                <span className="ml-2 text-sm">Разрешить</span>
              </label>
            </div>
            
            <div className="border border-gray-200 rounded-xl p-4">
              <h4 className="font-semibold text-gray-800 mb-2">Маркетинг</h4>
              <p className="text-gray-600 text-sm">Для персонализированной рекламы</p>
              <label className="inline-flex items-center mt-2">
                <input 
                  type="checkbox" 
                  className="rounded text-pink-500"
                  checked={preferences.marketing}
                  onChange={(e) => setPreferences(prev => ({ ...prev, marketing: e.target.checked }))}
                />
                <span className="ml-2 text-sm">Разрешить</span>
              </label>
            </div>
          </div>
          
          <div className="flex flex-col sm:flex-row justify-end gap-3">
            <button 
              onClick={() => setIsVisible(false)}
              className="px-4 py-2 text-gray-600 hover:text-gray-800"
            >
              Отклонить
            </button>
            <button 
              onClick={handleAcceptAll}
              className="px-4 py-2 bg-gray-100 text-gray-800 rounded-lg hover:bg-gray-200 transition-colors"
            >
              Принять все
            </button>
            <button 
              onClick={handleSavePreferences}
              className="px-4 py-2 bg-pink-500 text-white rounded-lg hover:bg-pink-600 transition-colors"
            >
              Сохранить настройки
            </button>
          </div>
          
          <div className="mt-4 text-center text-sm text-gray-500">
            <button className="hover:underline">
              Подробнее о файлах cookie
            </button>
          </div>
        </div>
      </div>
    </motion.div>
  );
};

export default CookieBanner;