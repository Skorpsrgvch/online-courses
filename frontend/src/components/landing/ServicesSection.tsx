import React from 'react';
import { Link } from 'react-router-dom';

interface Service {
  id: number;
  title: string;
  description: string;
  price: string;
  icon: React.ReactNode;
  duration?: string;
}

export const ServicesSection: React.FC = () => {
  const services: Service[] = [
    {
      id: 1,
      title: 'Индивидуальная консультация',
      description: 'Подробный разбор вашей ситуации, диагностика состояния тазового дна и составление персонального плана восстановления.',
      price: '3 500 ₽',
      duration: '60 мин',
      icon: (
        <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
      ),
    },
    {
      id: 2,
      title: 'Диагностика тазового дна',
      description: 'Комплексное обследование мышц тазового дна на специальном оборудовании для выявления проблем и дисфункций.',
      price: '2 800 ₽',
      duration: '45 мин',
      icon: (
        <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
        </svg>
      ),
    },
    {
      id: 3,
      title: 'Персональная реабилитация',
      description: 'Курс индивидуальных занятий по восстановлению после родов или операций. Подбор упражнений под ваш ритм жизни.',
      price: 'от 2 000 ₽',
      duration: 'за занятие',
      icon: (
        <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
        </svg>
      ),
    },
    {
      id: 4,
      title: 'Подготовка к родам',
      description: 'Онлайн-курс для будущих мам: дыхательные техники, работа с мышцами, психологическая подготовка и план родов.',
      price: '4 200 ₽',
      duration: 'курс',
      icon: (
        <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      ),
    },
    {
      id: 5,
      title: 'Ведение в послеродовом периоде',
      description: 'Комплексная программа восстановления "под ключ". Контроль специалиста, коррекция упражнений и поддержка 24/7.',
      price: '12 000 ₽',
      duration: 'месяц',
      icon: (
        <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
        </svg>
      ),
    },
    {
      id: 6,
      title: 'Онлайн-сопровождение',
      description: 'Удаленная работа со специалистом через видеосвязь. Проверка техники выполнения упражнений и ведение дневника.',
      price: '5 900 ₽',
      duration: 'месяц',
      icon: (
        <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
      ),
    },
  ];

  return (
    <section id="services" className="py-20 bg-gray-50/50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Заголовок секции */}
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl md:text-4xl font-serif font-bold text-gray-900 mb-4">
            Услуги и программы
          </h2>
          <p className="text-lg text-gray-600 leading-relaxed">
            Выберите подходящий формат работы для восстановления женского здоровья. 
            От разовых консультаций до комплексных программ сопровождения.
          </p>
        </div>

        {/* Сетка услуг */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {services.map((service) => (
            <div
              key={service.id}
              className="group bg-white rounded-3xl p-8 shadow-sm hover:shadow-xl transition-all duration-300 border border-gray-100 hover:border-rose-100 flex flex-col h-full transform hover:-translate-y-1"
            >
              {/* Иконка */}
              <div className="w-16 h-16 bg-linear-to-br from-rose-100 to-rose-200 rounded-2xl flex items-center justify-center text-rose-600 mb-6 group-hover:scale-110 transition-transform duration-300 shadow-inner">
                {service.icon}
              </div>

              {/* Контент */}
              <div className="flex-grow mb-6">
                <div className="flex justify-between items-start mb-3">
                  <h3 className="text-xl font-bold text-gray-900 font-serif group-hover:text-rose-600 transition-colors">
                    {service.title}
                  </h3>
                </div>
                <p className="text-gray-600 leading-relaxed mb-4 text-sm md:text-base">
                  {service.description}
                </p>
                
                {/* Детали (время/формат) */}
                {service.duration && (
                  <div className="inline-flex items-center gap-1.5 px-3 py-1 bg-gray-50 rounded-full text-xs font-medium text-gray-500 mb-2">
                    <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {service.duration}
                  </div>
                )}
              </div>

              {/* Футер карточки (Цена + Кнопка) */}
              <div className="pt-6 border-t border-gray-100 mt-auto">
                <div className="flex items-center justify-between">
                  <span className="text-2xl font-bold text-gray-900">
                    {service.price}
                  </span>
                  <button className="px-5 py-2.5 bg-rose-500 text-white text-sm font-medium rounded-2xl!
                   hover:bg-rose-600 duration-300 shadow-md hover:shadow-lg transition-all">
                    Записаться
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Кнопка действия внизу */}
        <div className="text-center mt-16">
          <Link
            to="/contact"
            className="inline-flex items-center gap-2 px-8 py-3 bg-white border-2 border-rose-300 text-rose-600! font-semibold rounded-2xl hover:border-rose-500! hover:bg-rose-50 hover:-translate-y-1 hover:shadow-lg hover:shadow-rose-100/50 transition-all duration-300"
            style={{ textDecoration: 'none' }}
          >
            Получить рекомендацию
            <svg className="w-5 h-5 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8l4 4m0 0l-4 4m4-4H3" />
            </svg>
          </Link>
        </div>
      </div>
    </section>
  );
};