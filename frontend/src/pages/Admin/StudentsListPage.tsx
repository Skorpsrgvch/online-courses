import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { adminService } from '../../api/admin.service';
import type { StudentStat } from '../../api/types';
import { useAuth } from '../../context/AuthContext';

const StudentsListPage = () => {
    const { user } = useAuth();
    const navigate = useNavigate();
    const [students, setStudents] = useState<StudentStat[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Проверка прав администратора
        if (!user || user.role !== 'admin') {
            navigate('/');
            return;
        }

        const loadData = async () => {
            try {
                const data = await adminService.getTopStudents(20);
                setStudents(data);
            } catch (err) {
                console.error('Ошибка загрузки статистики:', err);
            } finally {
                setLoading(false);
            }
        };

        loadData();
    }, [navigate, user]);

    const formatDate = (dateStr?: string | null) => {
        if (!dateStr) return '—';
        return new Date(dateStr).toLocaleDateString('ru-RU', {
            day: 'numeric', month: 'short', year: 'numeric'
        });
    };

    if (loading) {
        return (
            <div className="min-h-screen bg-gray-50 flex items-center justify-center">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-rose-500"></div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 p-6">
            <div className="max-w-7xl mx-auto">
                <header className="mb-8 flex justify-between items-center">
                    <div>
                        <h1 className="text-3xl font-serif font-bold text-gray-900">Клиенты</h1>
                        <p className="text-gray-500 mt-1">Рейтинг активности и прогресса</p>
                    </div>
                    <button
                        onClick={() => navigate('/admin')}
                        className="px-4 py-2 bg-white border border-gray-300 rounded-lg text-sm font-medium hover:bg-gray-50"
                    >
                        Назад в админку
                    </button>
                </header>

                {/* KPI Карточки (опционально можно рассчитать агрегаты) */}
                <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
                    <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
                        <p className="text-sm text-gray-500 font-medium">Всего клиентов</p>
                        <p className="text-2xl font-bold text-gray-900 mt-2">{students.length}</p>
                    </div>
                    {/* Можно добавить другие метрики, если бэкенд их отдаст отдельно */}
                </div>

                {/* Таблица */}
                <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full text-left">
                            <thead className="bg-gray-50 border-b border-gray-100">
                                <tr>
                                    <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Клиент</th>
                                    <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider text-center">Прогресс</th>
                                    <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider text-center">Курсы</th>
                                    <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider text-center">Уроки</th>
                                    <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Посл. активность</th>
                                    <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider text-right">Действия</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100">
                                {students && Array.isArray(students) && students.length > 0 ? (
                                    students.map((student) => (
                                        <tr
                                        key={student.id}
                                        onClick={() => navigate(`/admin/students/${student.id}`)}
                                        className="hover:bg-gray-50 cursor-pointer transition-colors"
                                    >
                                        <td className="px-6 py-4">
                                            <div className="flex flex-col">
                                                <span className="font-semibold text-gray-900">{student.name}</span>
                                                <span className="text-sm text-gray-500">{student.email}</span>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 text-center">
                                            <div className="flex items-center justify-center gap-2">
                                                <div className="w-24 h-2 bg-gray-100 rounded-full overflow-hidden">
                                                    <div
                                                        className={`h-full rounded-full ${student.progress_percent === 100 ? 'bg-green-500' :
                                                            student.progress_percent > 50 ? 'bg-blue-500' : 'bg-yellow-500'
                                                            }`}
                                                        style={{ width: `${student.progress_percent}%` }}
                                                    />
                                                </div>
                                                <span className="text-sm font-bold text-gray-700">{student.progress_percent}%</span>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 text-center text-sm text-gray-600">
                                            <span className="font-bold text-gray-900">{student.completed_courses}</span> / {student.total_courses}
                                        </td>
                                        <td className="px-6 py-4 text-center text-sm text-gray-600">
                                            <span className="font-bold text-gray-900">{student.completed_lessons}</span> / {student.total_lessons}
                                        </td>
                                        <td className="px-6 py-4 text-sm text-gray-500">
                                            {formatDate(student.last_activity_at)}
                                        </td>
                                        <td className="px-6 py-4 text-right">
                                            <button className="text-rose-600 hover:text-rose-700 font-medium text-sm">
                                                Подробнее →
                                            </button>
                                        </td>
                                    </tr>
                                    ))
                                ) : (
                                    <div className="text-center py-10 text-gray-500">
                                        {loading ? 'Загрузка...' : 'Клиенты не найдены или произошла ошибка загрузки.'}
                                    </div>
                                )}
                            
                                {students.length === 0 && (
                                    <tr>
                                        <td colSpan={6} className="px-6 py-12 text-center text-gray-500">
                                            Нет данных об учениках
                                        </td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default StudentsListPage;