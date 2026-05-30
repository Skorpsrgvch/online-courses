import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { servicesService } from '../../api/services.service';
import type { Service } from '../../api/types';

type FormMode = 'create' | 'edit';

interface ServiceFormPageProps {
    mode: FormMode;
}

const ServiceFormPage: React.FC<ServiceFormPageProps> = ({ mode }) => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();

    // Используем переданный mode, а не вычисляем его из URL
    const isEdit = mode === 'edit';

    const [formData, setFormData] = useState({
        title: '',
        description: '',
        price: 0,
        duration_minutes: 0,
    });

    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [successMsg, setSuccessMsg] = useState<string | null>(null);

    useEffect(() => {
        // Загружаем данные только если режим редактирования и есть ID
        if (isEdit && id) {
            const loadService = async () => {
                setIsLoading(true);
                try {
                    const data = await servicesService.getById(Number(id));
                    setFormData({
                        title: data.title,
                        description: data.description,
                        price: data.price,
                        duration_minutes: data.duration_minutes || 0,
                    });
                } catch (err: any) {
                    setError('Не удалось загрузить услугу');
                    console.error(err);
                } finally {
                    setIsLoading(false);
                }
            };
            loadService();
        }
    }, [id, isEdit]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: name === 'price' || name === 'duration_minutes' ? Number(value) : value,
        }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError(null);
        setSuccessMsg(null);

        try {
            if (isEdit) {
                if (!id) throw new Error('ID услуги не указан');
                await servicesService.update(Number(id), formData);
                setSuccessMsg('Услуга успешно обновлена!');
            } else {
                await servicesService.create(formData);
                setSuccessMsg('Услуга успешно создана!');
                // Очистить форму при создании
                setFormData({ title: '', description: '', price: 0, duration_minutes: 0 });
            }

            // Через 1.5 сек возвращаемся в админку
            setTimeout(() => navigate('/admin'), 1500);
        } catch (err: any) {
            setError(err.message || 'Ошибка сохранения');
        } finally {
            setIsLoading(false);
        }
    };

    if (isLoading && isEdit && !formData.title) {
        return (
            <div className="min-h-screen bg-gray-50 flex items-center justify-center">
                <div className="text-center">
                    <div className="w-12 h-12 border-4 border-rose-300 border-t-rose-500 rounded-full animate-spin mx-auto mb-4"></div>
                    <p className="text-gray-500">Загрузка...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 py-10">
            <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="mb-6">
                    <a href="/admin" className="text-sm text-gray-500 hover:text-rose-500 transition-colors"
                    style={{ textDecoration: 'none' }}
                    >← Назад в админ-панель</a>
                    <h1 className="text-2xl font-serif font-bold text-gray-900 mt-2">
                        {isEdit ? 'Редактирование услуги' : 'Новая услуга'}
                    </h1>
                </div>

                <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 md:p-8">
                    {error && (
                        <div className="mb-6 p-3 bg-red-50 text-red-700 rounded-lg border border-red-100 text-sm">
                            {error}
                        </div>
                    )}
                    {successMsg && (
                        <div className="mb-6 p-3 bg-green-50 text-green-700 rounded-lg border border-green-100 text-sm">
                            {successMsg}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-6">
                        {/* Название */}
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Название услуги</label>
                            <input
                                type="text"
                                name="title"
                                required
                                value={formData.title}
                                onChange={handleChange}
                                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-500 focus:border-rose-500 outline-none transition-shadow"
                                placeholder="Например: Индивидуальная консультация"
                            />
                        </div>

                        {/* Описание */}
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Описание
                                <span className="text-xs text-gray-500 ml-2">(Используйте ||| для разделения абзацев)</span>
                            </label>
                            <textarea
                                name="description"
                                required
                                rows={6}
                                value={formData.description}
                                onChange={handleChange}
                                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-500 focus:border-rose-500 outline-none transition-shadow"
                                placeholder="'Комплексная программа физической реабилитации для женщин после лечения рака молочной железы. 
                                              Поможет восстановить подвижность плеча, справиться с отёками и рубцами, вернуть уверенность в своём теле. 
                                              Индивидуальный подход с учётом диагноза и этапа лечения.|||
                                              Составление комплексов упражнений с учетом диагноза и проведенного лечения, а также вашего запроса.|||
                                              Восстановление функциональной активности с учетом риска лимфостаза.|||
                                              Работа с самыми частыми жалобами и состояниями на фоне лечения:|||
                                              - Постмастэктомический синдром (отёки, ограничения движения),|||
                                              - Web-syndrome (тяж), отеки, грубые рубцы,|||
                                              - Контрактура плечевого сустава (снижение объёма движения).|||
                                              Кинезиотейпирование.'"
                            />
                        </div>

                        {/* Цена и Длительность */}
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Стоимость (₽)</label>
                                <input
                                    type="number"
                                    name="price"
                                    required
                                    value={formData.price}
                                    onChange={handleChange}
                                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-500 focus:border-rose-500 outline-none transition-shadow"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Длительность (мин)</label>
                                <input
                                    type="number"
                                    name="duration_minutes"
                                    value={formData.duration_minutes}
                                    onChange={handleChange}
                                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-500 focus:border-rose-500 outline-none transition-shadow"
                                />
                            </div>
                        </div>

                        {/* Кнопки */}
                        <div className="flex items-center gap-4 pt-4 border-t border-gray-100">
                            <button
                                type="submit"
                                disabled={isLoading}
                                className="px-6 py-2.5 bg-rose-500 text-white font-medium rounded-lg hover:bg-rose-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                {isLoading ? 'Сохранение...' : (isEdit ? 'Сохранить изменения' : 'Создать услугу')}
                            </button>
                            <button
                                type="button"
                                onClick={() => navigate('/admin')}
                                className="px-6 py-2.5 text-gray-600 font-medium hover:text-gray-800 transition-colors"
                            >
                                Отмена
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default ServiceFormPage;