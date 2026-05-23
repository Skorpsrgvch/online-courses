import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { servicesService } from '../../api/services.service';
import type { Service } from '../../api/types';

export const ServicesSection: React.FC = () => {
  const [services, setServices] = useState<Service[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null); // Добавлено состояние ошибки

  const loadServices = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await servicesService.getAll();
      const sortedData = [...data].sort((a, b) => a.id - b.id);
      setServices(sortedData.slice(0, 3));
    } catch (err: any) {
      console.error('Ошибка загрузки услуг:', err);
      setError('Не удалось загрузить каталог услуг');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadServices();
  }, []);

  const formatPrice = (price: number): string => {
    if (price === 0) return 'Бесплатно';
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RUB',
      maximumFractionDigits: 0,
    }).format(price).replace('₽', '₽');
  };

  const getShortDescription = (text: string): string => {
    if (!text) return 'Описание услуги скоро появится...';
    const firstPart = text.split('|||')[0];
    return firstPart.trim();
  };

  // Экран загрузки
  if (isLoading) {
    return (
      <section id="services" className="py-20 relative z-10">
        <div className="max-w-7xl mx-auto px-4 text-center">
          <div className="inline-block w-10 h-10 border-4 border-rose-200 border-t-rose-500 rounded-full animate-spin"></div>
          <p className="mt-4 text-gray-500">Загрузка услуг...</p>
        </div>
      </section>
    );
  }

  // Экран ошибки
  if (error) {
    return (
      <section id="services" className="py-20 relative z-10">
        <div className="max-w-7xl mx-auto px-4 text-center text-red-500 bg-red-50 p-6 rounded-xl">
          <p>{error}</p>
          <button
            onClick={loadServices}
            className="mt-4 text-sm underline hover:text-red-700 font-medium"
          >
            Попробовать снова
          </button>
        </div>
      </section>
    );
  }

  return (
    <section id="services" className="py-10 md:py-14 lg:py-16 relative z-10">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">

        {/* Заголовок секции */}
        <div className="text-center mb-14 max-w-3xl mx-auto">
          <h2 className="text-3xl md:text-4xl font-serif font-bold text-gray-900 mb-4">
            Услуги и программы
          </h2>
          <p className="text-base text-gray-600 leading-relaxed">
            Выберите подходящий формат работы для восстановления женского здоровья.
            От разовых консультаций до комплексных программ сопровождения.
          </p>
        </div>


        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
          {services.map((service) => (
            <div
              key={service.id}
              className="bg-white rounded-2xl shadow-lg overflow-hidden border border-gray-100 hover:shadow-2xl hover:shadow-rose-100/50 transition-all duration-300 transform hover:-translate-y-1 flex flex-col h-full p-8"
            >
              <div className="flex flex-col h-full gap-6">


                <div className="flex items-start gap-4 shrink-0">
                  {/* Иконка */}
                  <div className="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-rose-100 to-rose-200 rounded-xl flex items-center justify-center text-rose-600 shadow-sm group-hover:scale-110 transition-transform duration-300">
                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19.428 15.428a2 2 0 00-1.022-.547l-2.384-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
                    </svg>
                  </div>

                  {/* Заголовок */}
                  <h4 className="text-xl font-bold text-gray-900 font-serif leading-snug break-words  flex-grow">
                    {service.title}
                  </h4>
                </div>


                <div className="flex-grow flex flex-col gap-4">
                  <p className="text-gray-600 leading-relaxed text-sm md:text-base break-words  min-h-[4rem]">
                    {getShortDescription(service.description)}
                  </p>
                </div>

                {service.duration_minutes && service.duration_minutes > 0 && (
                  <div>
                    <span className="inline-flex items-center gap-1.5 px-3 py-1 bg-gray-50 rounded-full text-xs font-medium text-gray-500">
                      <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      {service.duration_minutes} мин
                    </span>
                  </div>
                )}

                <div className="mt-auto pt-6 border-t border-gray-100">

                  <div className="flex items-baseline justify-between mb-4">
                    <span className="text-sm text-gray-500 font-medium">Стоимость:</span>
                    <span className="text-xl font-bold text-gray-900">
                      {formatPrice(service.price)}
                    </span>
                  </div>

                  <Link
                    to="/services"
                    style={{ textDecoration: 'none' }}
                    className="block w-full py-2.5 px-4 bg-rose-500 text-white text-sm font-semibold rounded-2xl shadow-md hover:bg-rose-600 hover:shadow-lg transition-all duration-300 cursor-pointer text-center"
                  >
                    Подробнее об услуге
                  </Link>
                </div>

              </div>
            </div>
          ))}
        </div>

        {/* Кнопка "Все услуги" */}
        <div className="text-center mt-12">
          <Link
            to="/services"
            className="inline-flex items-center gap-2 px-8 py-3 bg-white border-2 border-rose-300 text-rose-600! font-semibold rounded-2xl! hover:border-rose-500 hover:bg-rose-50 hover:-translate-y-1 hover:shadow-lg hover:shadow-rose-100/50 transition-all duration-300 group"
            style={{ textDecoration: 'none' }}
          >
            Смотреть все услуги
            <svg className="w-5 h-5 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8l4 4m0 0l-4 4m4-4H3" />
            </svg>
          </Link>
        </div>
      </div>
    </section>
  );
};