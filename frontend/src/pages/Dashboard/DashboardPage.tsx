import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { userService } from '../../api';
import type { UserProfile, UserCourseProgress } from '../../api/types';
import { formatDate } from '../../utils/dateUtils';

const DashboardPage = () => {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [courses, setCourses] = useState<UserCourseProgress[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'courses' | 'profile'>('courses');

  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const [profileData, coursesData] = await Promise.all([
          userService.getProfile(),
          userService.getCourses(),
        ]);
        setProfile(profileData);
        setCourses(coursesData.courses);
      } catch (err: any) {
        setError(err.message || 'Не удалось загрузить данные');
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, []);

  const handleLogout = useCallback(() => {
    localStorage.removeItem('currentUser');
    window.location.href = '/';
  }, []);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="w-12 h-12 border-4 border-rose-300 border-t-rose-500 rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-500">Загрузка...</p>
        </div>
      </div>
    );
  }

  if (error && !profile) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center bg-white p-8 rounded-xl shadow-sm max-w-md">
          <div className="text-5xl mb-4">😔</div>
          <h2 className="text-xl font-semibold text-gray-800 mb-2">Ошибка загрузки</h2>
          <p className="text-gray-600 mb-4">{error}</p>
          <a href="/" className="text-rose-500 hover:underline font-medium">Вернуться на главную</a>
        </div>
      </div>
    );
  }

  // Статистика
  const totalProgress = courses.length > 0
    ? Math.round(courses.reduce((sum, c) => sum + c.progress_percent, 0) / courses.length)
    : 0;
  const completedCourses = courses.filter(c => c.progress_percent === 100).length;

  return (
    <div className="min-h-screen bg-gray-50">
      
      <div className="bg-white border-b border-gray-200 py-6">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <a href="/" className="text-sm text-gray-500! hover:text-rose-500! transition-colors mb-2 inline-block"
          style={{ textDecoration: 'none' }}>
            ← На главную
          </a>
          <h1 className="text-2xl md:text-3xl font-serif font-bold text-gray-900">Личный кабинет</h1>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        {/* Профиль — карточка сверху */}
        {profile && (
          <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-5 mb-6">
            <div className="flex items-center justify-between flex-wrap gap-4">
              <div className="flex items-center gap-4">
                <div className="w-14 h-14 bg-linear-to-br from-rose-400 to-rose-600 rounded-full flex items-center justify-center text-white font-bold text-xl shadow-lg"
                style={{ textDecoration: 'none' }}>
                  {profile.name?.charAt(0).toUpperCase() || profile.email.charAt(0).toUpperCase()}
                </div>
                <div>
                  <h2 className="text-lg font-semibold text-gray-800">
                    {profile.name || 'Пользователь'}
                  </h2>
                  <p className="text-sm text-gray-500">{profile.email}</p>
                  <p className="text-xs text-gray-400 mt-0.5">
                    На платформе с {formatDate(profile.created_at)}
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                {profile.role === 'admin' && (
                  <Link
                    to="/admin"
                    className="px-3 py-1.5 bg-gray-800 text-white text-xs rounded-full hover:bg-gray-700 transition-colors"
                  >
                    Панель администратора
                  </Link>
                )}
                <button
                  onClick={handleLogout}
                  className="px-4 py-2 text-sm text-gray-600 border border-gray-300 rounded-lg hover:border-red-400 hover:text-red-500 transition-colors"
                >
                  Выйти
                </button>
              </div>
            </div>

            {/* Статистика */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-4 pt-4 border-t border-gray-100">
              <div className="text-center">
                <p className="text-2xl font-bold text-rose-500">{courses.length}</p>
                <p className="text-xs text-gray-500 mt-0.5">Куплено курсов</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-bold text-green-500">{completedCourses}</p>
                <p className="text-xs text-gray-500 mt-0.5">Завершено</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-bold text-rose-500">{totalProgress}%</p>
                <p className="text-xs text-gray-500 mt-0.5">Средний прогресс</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-bold text-gray-700">
                  {courses.reduce((s, c) => s + c.completed_count, 0)}
                </p>
                <p className="text-xs text-gray-500 mt-0.5">Уроков пройдено</p>
              </div>
            </div>
          </div>
        )}

        {/* Вкладки */}
        <div className="flex border-b border-gray-200 mb-6">
          <button
            onClick={() => setActiveTab('courses')}
            className={`px-5 py-3 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'courses'
                ? 'border-rose-500 text-rose-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Мои курсы ({courses.length})
          </button>
          <button
            onClick={() => setActiveTab('profile')}
            className={`px-5 py-3 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'profile'
                ? 'border-rose-500 text-rose-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Профиль
          </button>
        </div>

        {/* Содержимое вкладки "Мои курсы" */}
        {activeTab === 'courses' && (
          <div>
            {courses.length === 0 ? (
              <div className="bg-white p-12 rounded-xl shadow-sm border border-gray-100 text-center">
                <div className="text-5xl mb-4">📚</div>
                <h3 className="text-lg font-semibold text-gray-700 mb-2">У вас пока нет курсов</h3>
                <p className="text-gray-500 mb-4">
                  Просмотрите каталог и выберите курс, который вам подходит
                </p>
                <Link
                  to="/#courses"
                  className="inline-block px-6 py-2.5 bg-rose-500 text-white rounded-lg hover:bg-rose-600 transition-colors text-sm font-medium"
                >
                  Перейти в каталог
                </Link>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {courses.map((course) => (
                  <Link
                    key={course.id}
                    to={`/course/${course.id}`}
                    className="block bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-md hover:border-rose-200 transition-all"
                  >
                    {/* Обложка */}
                    <div className="relative h-36 bg-gradient-to-br from-rose-100 to-rose-200 overflow-hidden">
                      {course.cover_image_url ? (
                        <img
                          src={course.cover_image_url}
                          alt={course.title}
                          className="w-full h-full object-cover"
                        />
                      ) : (
                        <div className="flex items-center justify-center h-full">
                          <span className="text-4xl">📖</span>
                        </div>
                      )}
                      {course.progress_percent === 100 && (
                        <div className="absolute top-2 right-2 px-2 py-1 bg-green-500 text-white text-xs rounded-full font-medium">
                          ✓ Завершён
                        </div>
                      )}
                    </div>

                    {/* Информация */}
                    <div className="p-4">
                      <h3 className="font-semibold text-gray-800 mb-1 truncate">{course.title}</h3>
                      <p className="text-xs text-gray-500 mb-3 line-clamp-2">{course.description}</p>

                      {/* Прогресс-бар */}
                      <div className="flex items-center gap-2">
                        <div className="flex-1 bg-gray-200 rounded-full h-2 overflow-hidden">
                          <div
                            className={`h-full rounded-full transition-all duration-500 ${
                              course.progress_percent === 100 ? 'bg-green-500' : 'bg-rose-500'
                            }`}
                            style={{ width: `${course.progress_percent}%` }}
                          />
                        </div>
                        <span className="text-xs font-medium text-gray-600 min-w-[36px] text-right">
                          {course.progress_percent}%
                        </span>
                      </div>
                      <p className="text-xs text-gray-400 mt-1.5">
                        {course.completed_count} из {course.total_lessons} уроков
                      </p>
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Содержимое вкладки "Профиль" */}
        {activeTab === 'profile' && profile && (
          <div className="max-w-2xl">
            <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 space-y-5">
              <h3 className="text-lg font-semibold text-gray-800 border-b border-gray-100 pb-3">
                Личные данные
              </h3>

              <div>
                <label className="block text-sm font-medium text-gray-500 mb-1">Имя</label>
                <p className="text-gray-800">{profile.name || 'Не указано'}</p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-500 mb-1">Email</label>
                <p className="text-gray-800">{profile.email}</p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-500 mb-1">Роль</label>
                <p className="text-gray-800">
                  {profile.role === 'admin' ? (
                    <span className="px-2 py-0.5 bg-gray-800 text-white text-xs rounded-full">Администратор</span>
                  ) : (
                    'Пользователь'
                  )}
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-500 mb-1">Дата регистрации</label>
                <p className="text-gray-800">{formatDate(profile.created_at)}</p>
              </div>

              <div className="pt-4 border-t border-gray-100">
                <p className="text-sm text-gray-400 italic">
                  Редактирование профиля и смена пароля будут доступны в ближайшее время.
                </p>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default DashboardPage;
