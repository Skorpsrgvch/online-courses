import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { userService, reviewsService } from '../../api';
import type { UserProfile, UserCourseProgress, Review } from '../../api/types';
import { formatDate } from '../../utils/dateUtils';
import apiClient from '../../api/axiosInstance'; // Для прямого вызова API удаления

const DashboardPage = () => {
  const navigate = useNavigate();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [courses, setCourses] = useState<UserCourseProgress[]>([]);
  const [myReviews, setMyReviews] = useState<Review[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'courses' | 'reviews' | 'profile'>('courses');

  // --- Состояния редактирования Профиля ---
  const [isEditing, setIsEditing] = useState(false);
  const [editName, setEditName] = useState('');
  const [editEmail, setEditEmail] = useState('');
  const [isSavingProfile, setIsSavingProfile] = useState(false);
  const [profileError, setProfileError] = useState<string | null>(null);

  // Состояние для процесса удаления отзыва
  const [deletingReviewId, setDeletingReviewId] = useState<number | null>(null);

  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const [profileData, coursesData, reviewsData] = await Promise.all([
          userService.getProfile(),
          userService.getCourses(),
          reviewsService.getMyReviews().catch(() => []),
        ]);
        
        setProfile(profileData);
        setCourses(coursesData.courses || []);
        setMyReviews(Array.isArray(reviewsData) ? reviewsData : []);
        
        setEditName(profileData.name || '');
        setEditEmail(profileData.email || '');
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

  const handleSaveProfile = async () => {
    if (!profile) return;
    setIsSavingProfile(true);
    setProfileError(null);

    try {
      const updates: any = {};
      if (editName !== profile.name) updates.name = editName;
      if (editEmail !== profile.email) updates.email = editEmail;

      if (Object.keys(updates).length === 0) {
        setIsEditing(false);
        return;
      }

      await userService.updateProfile(updates);
      setProfile({ ...profile, ...updates });
      setIsEditing(false);
    } catch (err: any) {
      setProfileError(err.message || 'Ошибка при сохранении');
    } finally {
      setIsSavingProfile(false);
    }
  };

  // Функция удаления отзыва
  const handleDeleteReview = async (reviewId: number, courseId: number) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот отзыв? Восстановить его будет невозможно.')) {
      return;
    }

    setDeletingReviewId(reviewId);
    try {
      // Используем прямой вызов API или метод сервиса, если он есть
      await apiClient.delete(`/reviews/${reviewId}`);
      
      // Удаляем отзыв из списка локально
      setMyReviews(prev => prev.filter(r => r.id !== reviewId));
    } catch (err: any) {
      setError(err.response?.data?.message || 'Не удалось удалить отзыв');
    } finally {
      setDeletingReviewId(null);
    }
  };

  const handleCardClick = (e: React.MouseEvent, courseId: number) => {
    if ((e.target as HTMLElement).closest('button')) return;
    navigate(`/course/${courseId}`);
  };

  const handleReviewClick = (e: React.MouseEvent, courseId: number) => {
    e.stopPropagation();
    navigate(`/course/${courseId}#reviews`);
  };

  const handleContinueClick = (e: React.MouseEvent, courseId: number) => {
    e.stopPropagation();
    navigate(`/course/${courseId}/learn`);
  };

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
          <button onClick={() => navigate('/')} className="text-rose-500 hover:underline font-medium">Вернуться на главную</button>
        </div>
      </div>
    );
  }

  const totalProgress = courses.length > 0
    ? Math.round(courses.reduce((sum, c) => sum + (c.progress_percent || 0), 0) / courses.length)
    : 0;
  const completedCourses = courses.filter(c => c.progress_percent === 100).length;

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 py-6">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <button onClick={() => navigate('/')} className="text-sm text-gray-500 hover:text-rose-500 transition-colors mb-2 inline-block"
            style={{ textDecoration: 'none', background: 'none', border: 'none', padding: 0, cursor: 'pointer' }}>
            ← На главную
          </button>
          <h1 className="text-2xl md:text-3xl font-serif font-bold text-gray-900">Личный кабинет</h1>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        {/* Профиль */}
        {profile && (
          <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-5 mb-6">
            <div className="flex items-center justify-between flex-wrap gap-4">
              <div className="flex items-center gap-4">
                <div className="w-14 h-14 bg-linear-to-br from-rose-400 to-rose-600 rounded-full flex items-center justify-center text-white font-bold text-xl shadow-lg">
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
                  <button
                    onClick={() => navigate('/admin')}
                    className="px-3 py-1.5 bg-gray-800 text-white text-xs rounded-full hover:bg-gray-700 transition-colors"
                  >
                    Панель администратора
                  </button>
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
                  {courses.reduce((s, c) => s + (c.completed_count || 0), 0)}
                </p>
                <p className="text-xs text-gray-500 mt-0.5">Уроков пройдено</p>
              </div>
            </div>
          </div>
        )}

        {/* Вкладки */}
        <div className="flex border-b border-gray-200 mb-6 overflow-x-auto">
          <button
            onClick={() => setActiveTab('courses')}
            className={`px-5 py-3 text-sm font-medium border-b-2 whitespace-nowrap transition-colors ${activeTab === 'courses'
                ? 'border-rose-500 text-rose-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
          >
            Мои курсы ({courses.length})
          </button>
          <button
            onClick={() => setActiveTab('reviews')}
            className={`px-5 py-3 text-sm font-medium border-b-2 whitespace-nowrap transition-colors ${activeTab === 'reviews'
                ? 'border-rose-500 text-rose-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
          >
            Мои отзывы ({myReviews.length})
          </button>
          <button
            onClick={() => setActiveTab('profile')}
            className={`px-5 py-3 text-sm font-medium border-b-2 whitespace-nowrap transition-colors ${activeTab === 'profile'
                ? 'border-rose-500 text-rose-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
          >
            Профиль
          </button>
        </div>

        {/* Вкладка "Мои курсы" */}
        {activeTab === 'courses' && (
          <div>
            {courses.length === 0 ? (
              <div className="bg-white p-12 rounded-xl shadow-sm border border-gray-100 text-center">
                <div className="text-5xl mb-4">📚</div>
                <h3 className="text-lg font-semibold text-gray-700 mb-2">У вас пока нет курсов</h3>
                <p className="text-gray-500 mb-4">Просмотрите каталог и выберите курс</p>
                <button
                  onClick={() => navigate('/courses')}
                  className="inline-block px-6 py-2.5 bg-rose-500 text-white rounded-lg! hover:bg-rose-600 transition-colors text-sm font-medium"
                >
                  Перейти в каталог
                </button>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {courses.map((course) => {
                  const isCompleted = course.progress_percent === 100;
                  const canReview = course.progress_percent > 80;
                  
                  // Проверяем, есть ли уже отзыв на этот курс
                  const hasReview = myReviews.some(r => r.course_id === course.id);

                  return (
                    <div
                      key={course.id}
                      onClick={(e) => handleCardClick(e, course.id)}
                      className="group bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-md hover:border-rose-200 transition-all cursor-pointer flex flex-col"
                    >
                      <div className="relative h-40 bg-linear-to-br from-rose-100 to-rose-200 overflow-hidden flex-shrink-0">
                        {course.cover_image_url ? (
                          <img src={course.cover_image_url} alt={course.title} className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
                        ) : (
                          <div className="flex items-center justify-center h-full"><span className="text-4xl">📖</span></div>
                        )}
                        {isCompleted && (
                          <div className="absolute top-3 right-3 px-2.5 py-1 bg-green-500 text-white text-xs rounded-full! font-bold shadow-sm">✓ Пройден</div>
                        )}
                      </div>

                      <div className="p-5 flex-grow flex flex-col">
                        <h3 className="font-bold text-gray-800 mb-2! line-clamp-2 leading-snug group-hover:text-rose-600 transition-colors">{course.title}</h3>
                        <p className="text-xs text-gray-500 mb-4! line-clamp-2 flex-grow">{course.description}</p>

                        <div className="mb-4">
                          <div className="flex justify-between text-xs mb-1">
                            <span className="text-gray-500">Прогресс</span>
                            <span className={`font-bold ${course.progress_percent === 100 ? 'text-green-600' : 'text-rose-600'}`}>{course.progress_percent}%</span>
                          </div>
                          <div className="w-full bg-gray-100 rounded-full! h-2 overflow-hidden">
                            <div className={`h-full rounded-full transition-all duration-700 ease-out ${course.progress_percent === 100 ? 'bg-green-500' : 'bg-rose-500'}`} style={{ width: `${course.progress_percent}%` }} />
                          </div>
                          <p className="text-[10px] text-gray-400 mt-1.5 text-right">{course.completed_count} из {course.total_lessons} уроков</p>
                        </div>

                        <div className="mt-auto pt-4 border-t border-gray-50 flex flex-col gap-2">
                          {isCompleted ? (
                            <>
                              <button onClick={(e) => handleContinueClick(e, course.id)} className="w-full py-2 bg-rose-50 text-rose-600 rounded-lg! text-sm font-semibold hover:bg-rose-100 transition-colors">Повторить курс</button>
                              {canReview && (
                                <button 
                                  onClick={(e) => handleReviewClick(e, course.id)} 
                                  className={`w-full py-2 rounded-lg! text-sm font-medium transition-colors border ${
                                    hasReview 
                                      ? 'bg-blue-50 border-blue-200 text-blue-600 hover:bg-blue-100' 
                                      : 'bg-white border-gray-200 text-gray-600 hover:border-rose-300 hover:text-rose-600'
                                  }`}
                                >
                                  {hasReview ? '✏️ Изменить отзыв' : '💬 Оставить отзыв'}
                                </button>
                              )}
                            </>
                          ) : (
                            <button onClick={(e) => handleContinueClick(e, course.id)} className="w-full py-2 bg-rose-500 text-white rounded-lg! text-sm font-semibold hover:bg-rose-600 transition-colors shadow-sm">{course.progress_percent > 0 ? 'Продолжить' : 'Начать обучение'}</button>
                          )}
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        )}

        {/* Вкладка "Мои отзывы" */}
        {activeTab === 'reviews' && (
          <div>
            {myReviews.length === 0 ? (
              <div className="bg-white p-12 rounded-xl shadow-sm border border-gray-100 text-center">
                <div className="text-5xl mb-4">💬</div>
                <h3 className="text-lg font-semibold text-gray-700 mb-2">Вы пока не оставляли отзывов</h3>
          
                <button
                  onClick={() => setActiveTab('courses')}
                  className="inline-block px-6 py-2.5 bg-rose-500 text-white rounded-lg! hover:bg-rose-600 transition-colors text-sm font-medium"
                >
                  Перейти к курсам
                </button>
              </div>
            ) : (
              <div className="space-y-4">
                {myReviews.map((review) => {
                  let statusBadge = null;
                  
                  if (review.approved) {
                    statusBadge = <span className="px-2.5 py-1 bg-green-100 text-green-700 text-xs font-bold rounded-full border border-green-200">Опубликован</span>;
                  } else if (review.rejection_reason) {
                    statusBadge = <span className="px-2.5 py-1 bg-red-100 text-red-700 text-xs font-bold rounded-full border border-red-200">Отклонен</span>;
                  } else {
                    statusBadge = <span className="px-2.5 py-1 bg-yellow-100 text-yellow-700 text-xs font-bold rounded-full border border-yellow-200">На модерации</span>;
                  }

                  const isDeleting = deletingReviewId === review.id;

                  return (
                    <div key={review.id} className="bg-white rounded-xl shadow-sm border border-gray-100 p-5 flex flex-col sm:flex-row gap-4 relative overflow-hidden">
                      {isDeleting && (
                        <div className="absolute inset-0 bg-white/80 z-10 flex items-center justify-center">
                          <div className="flex flex-col items-center">
                            <div className="w-6 h-6 border-2 border-rose-300 border-t-rose-500 rounded-full animate-spin mb-2"></div>
                            <span className="text-sm text-gray-500">Удаление...</span>
                          </div>
                        </div>
                      )}
                      
                      <div className="flex-1">
                        <div className="flex items-center justify-between mb-2 flex-wrap gap-2">
                          <div className="flex items-center gap-2">
                            <span className="text-xs font-bold text-rose-600 bg-rose-50 px-2 py-1 rounded-md border border-rose-100">
                              Курс: {review.course_title || 'Курс'}
                            </span>
                            {statusBadge}
                          </div>
                          <span className="text-xs text-gray-400">{formatDate(review.created_at)}</span>
                        </div>

                        <div className="flex items-center gap-1 mb-2">
                          {[1, 2, 3, 4, 5].map((star) => (
                            <span key={star} className={star <= review.rating ? 'text-yellow-400 text-lg' : 'text-gray-300 text-lg'}>★</span>
                          ))}
                        </div>

                        <p className="text-gray-700 leading-relaxed mb-3">{review.text}</p>

                        {review.rejection_reason && (
                          <div className="mt-3 p-3 bg-red-50 border border-red-100 rounded-lg">
                            <p className="text-xs font-bold text-red-800 mb-1">Причина отклонения:</p>
                            <p className="text-sm text-red-700">{review.rejection_reason}</p>
                          </div>
                        )}
                      </div>

                      <div className="flex sm:flex-col gap-2 border-t sm:border-t-0 sm:border-l border-gray-100 pt-3 sm:pt-0 sm:pl-4 mt-2 sm:mt-0 shrink-0">
                        <button 
                          onClick={() => navigate(`/course/${review.course_id}#reviews`)}
                          className="px-3 py-1.5 text-xs font-medium text-blue-600 bg-blue-50 border border-blue-200 rounded-2xl! hover:bg-blue-100 transition-colors whitespace-nowrap"
                        >
                          Изменить
                        </button>
                        <button 
                          onClick={() => handleDeleteReview(review.id, review.course_id)}
                          disabled={isDeleting}
                          className="px-3 py-1.5 text-xs font-medium text-red-600 bg-red-50 border border-red-200 rounded-2xl! hover:bg-red-100 transition-colors whitespace-nowrap disabled:opacity-50"
                        >
                          Удалить
                        </button>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        )}

        {/* Вкладка "Профиль" */}
        {activeTab === 'profile' && profile && (
          <div className="max-w-2xl">
            <div className="bg-white rounded-xl! shadow-sm border border-gray-100 p-6 space-y-6">
              <div className="flex justify-between items-center border-b border-gray-100 pb-4">
                <h3 className="text-lg font-semibold text-gray-800">Личные данные</h3>
                {!isEditing && (
                  <button onClick={() => { setIsEditing(true); setProfileError(null); }} className="text-sm text-rose-500 hover:text-rose-600 font-medium transition-colors">
                    Редактировать
                  </button>
                )}
              </div>

              {profileError && (
                <div className="p-3 bg-red-50 text-red-600 text-sm rounded-lg border border-red-100">
                  {profileError}
                </div>
              )}

              <div>
                <label className="block text-sm font-medium text-gray-500 mb-1">Имя</label>
                {isEditing ? (
                  <input type="text" value={editName} onChange={(e) => setEditName(e.target.value)} className="w-full px-3 py-2 border border-gray-300 rounded-lg! focus:ring-2 focus:ring-rose-500 focus:border-rose-500 outline-none transition-all" autoFocus />
                ) : (
                  <p className="text-gray-800 font-medium">{profile.name || 'Не указано'}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-500 mb-1">Email</label>
                {isEditing ? (
                  <input type="email" value={editEmail} onChange={(e) => setEditEmail(e.target.value)} className="w-full px-3 py-2 border border-gray-300 rounded-lg! focus:ring-2 focus:ring-rose-500 focus:border-rose-500 outline-none transition-all" />
                ) : (
                  <p className="text-gray-800">{profile.email}</p>
                )}
              </div>

              {isEditing && (
                <div className="flex gap-3 pt-4">
                  <button onClick={handleSaveProfile} disabled={isSavingProfile} className="px-6 py-2 bg-rose-500 text-white rounded-lg hover:bg-rose-600 disabled:opacity-50 transition-colors text-sm font-medium">
                    {isSavingProfile ? 'Сохранение...' : 'Сохранить изменения'}
                  </button>
                  <button onClick={() => { setIsEditing(false); setEditName(profile.name || ''); setEditEmail(profile.email); setProfileError(null); }} className="px-6 py-2 bg-gray-100 text-gray-600 rounded-lg! hover:bg-gray-200 transition-colors text-sm font-medium">
                    Отмена
                  </button>
                </div>
              )}

              <div className="pt-6 border-t border-gray-100">
                <div className="flex justify-between items-center mb-2">
                  <h4 className="text-md font-semibold text-gray-800">Безопасность</h4>
                  <button onClick={() => navigate('/password-recovery')} className="text-sm text-rose-500 hover:text-rose-600 font-medium transition-colors">
                    Изменить пароль
                  </button>
                </div>
                <p className="text-sm text-gray-500">Для смены пароля на почту будет отправлен проверочный код.</p>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default DashboardPage;
