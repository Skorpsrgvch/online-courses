import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { servicesService } from '../../api/services.service';
import type { Service } from '../../api/types';
import { renderParsedText } from '../../utils/textParser';

const ServicesPage: React.FC = () => {
    const [services, setServices] = useState<Service[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loadServices = async () => {
            try {
                const data = await servicesService.getAll();
                setServices(data);
            } catch (error) {
                console.error('Ошибка загрузки услуг:', error);
            } finally {
                setIsLoading(false);
            }
        };
        loadServices();
    }, []);

    const formatPrice = (price: number): string => {
        return new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'RUB',
            maximumFractionDigits: 0,
        }).format(price);
    };

    if (isLoading) {
        return (
            <div className="min-h-screen bg-gray-50 flex items-center justify-center">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-rose-500"></div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 py-12">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <Link to="/#courses" className="inline-flex items-start text-sm md:text-base lg:text-base font-medium text-gray-500! hover:text-rose-500! mb-4 transition-colors"
                          style={{ textDecoration: 'none' }}>
                            ← Назад на главную
                          </Link>

                {/* Заголовок страницы */}
                <div className="text-center mb-16">
                    <h1 className="text-4xl md:text-5xl font-serif font-bold text-gray-900 mb-6!">
                        Услуги и цены
                    </h1>
                    <p className="text-lg text-gray-600 max-w-2xl mx-auto">
                        Подробное описание всех доступных программ реабилитации и консультаций.
                        Выберите подходящий формат работы для восстановления здоровья.
                    </p>
                </div>

                {/* Сетка услуг */}
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                    {services.map((service) => {
                        return (
                            <div
                                key={service.id}
                                className="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-xl transition-all duration-300 flex flex-col"
                            >
                                {/* Шапка карточки */}
                                <div className="p-8 pb-6 border-b border-gray-50">
                                    <div className="flex justify-between items-center gap-4 mb-4">
                                        <h2 className="text-2xl font-serif font-bold text-gray-900 leading-tight">
                                            {service.title}
                                        </h2>
                                        {service.duration_minutes && (
                                            <span className="flex-shrink-0 inline-flex items-center gap-1.5 px-3 py-1 bg-rose-50 text-rose-700 rounded-full text-xs font-bold uppercase tracking-wide whitespace-nowrap">
                                                <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                                </svg>
                                                {service.duration_minutes}
                                            </span>
                                        )}
                                    </div>

                                    {/* Использование общей функции рендеринга */}
                                    <div className="text-gray-600 leading-relaxed text-base">
                                        {renderParsedText(service.description)}
                                    </div>
                                </div>

                                {/* Футер карточки */}
                                <div className="p-8 pt-6 mt-auto border-t border-gray-100 bg-white">
                                    <div className="flex items-center justify-between">
                                        <div>
                                            <span className="block text-xs text-gray-400 font-medium uppercase tracking-wider mb-1">Стоимость</span>
                                            <span className="text-3xl font-bold text-gray-900">
                                                {formatPrice(service.price)}
                                            </span>
                                        </div>
                                        <button 
                                        onClick={() => window.open('https://vk.com/im/convo/14374433?entrypoint=list_all&tab=all', '_blank')}
                                        className="px-8 py-3 bg-rose-500 text-white font-semibold rounded-2xl! hover:bg-rose-600 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300">
                                            Записаться
                                        </button>
                                    </div>
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>
        </div>
    );
};

export default ServicesPage;