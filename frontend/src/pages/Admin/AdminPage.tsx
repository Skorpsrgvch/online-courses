import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { coursesService, reviewsService, servicesService } from '../../api';
import apiClient from '../../api/axiosInstance';
import type { Course, Review, Service } from '../../api/types';
import { timeAgo } from '../../utils/dateUtils';

type AdminTab = 'courses' | 'hidden-courses' | 'reviews' | 'services' | 'grant-access';

interface UserSearchResult {
  id: number;
  email: string;
  name: string;
}

const AdminPage = () => {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<AdminTab>('courses');

  const [allCourses, setAllCourses] = useState<Course[]>([]);
  const [pendingReviews, setPendingReviews] = useState<Review[]>([]);
  const [servicesList, setServicesList] = useState<Service[]>([]);

  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Состояния для вкладки "Выдача доступа"
  const [emailQuery, setEmailQuery] = useState('');
  const [isSearching, setIsSearching] = useState(false);
  const [searchResults, setSearchResults] = useState<UserSearchResult[]>([]);
  const [selectedUser, setSelectedUser] = useState<UserSearchResult | null>(null);
  const [selectedUserCourseIds, setSelectedUserCourseIds] = useState<number[]>([]);
  const [selectedCourseId, setSelectedCourseId] = useState<number | ''>('');
  const [isGranting, setIsGranting] = useState(false);
  const [accessMessage, setAccessMessage] = useState<{ type: 'success' | 'error', text: string } | null>(null);

  // Состояния для модальных окон (Отклонение и Ошибка)
  const [isRejectModalOpen, setIsRejectModalOpen] = useState(false);
  const [rejectReason, setRejectReason] = useState('');
  const [selectedReviewId, setSelectedReviewId] = useState<number | null>(null);

  // Состояние для модалки ошибки
  const [errorModal, setErrorModal] = useState<{ isOpen: boolean; message: string }>({ isOpen: false, message: '' });

  // Функция загрузки данных при старте
  const loadData = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const [coursesData, reviewsData, servicesData] = await Promise.all([
        coursesService.getAllCoursesAdmin(),
        reviewsService.getPendingReviews(),
        servicesService.getAll(),
      ]);

      setAllCourses(Array.isArray(coursesData) ? coursesData : []);
      setPendingReviews(Array.isArray(reviewsData) ? reviewsData : []);
      const sortedServices = Array.isArray(servicesData) ? [...servicesData].sort((a, b) => a.id - b.id) : [];
      setServicesList(sortedServices);
    } catch (err: any) {
      console.error(err);
      setError(err.message || 'Не удалось загрузить данные');
    } finally {
      setIsLoading(false);
    }
  };

  // Отдельная функция только для перезагрузки отзывов
  const loadPendingReviews = async () => {
    try {
      const reviewsData = await reviewsService.getPendingReviews();
      setPendingReviews(Array.isArray(reviewsData) ? reviewsData : []);
    } catch (err: any) {
      console.error('Failed to reload reviews:', err);
      showErrorModal('Не удалось обновить список отзывов');
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const activeCourses = allCourses.filter(c => c.is_active !== false);
  const hiddenCourses = allCourses.filter(c => c.is_active === false);
  const availableCoursesForSelectedUser = selectedUser
    ? activeCourses.filter(course => !selectedUserCourseIds.includes(course.id))
    : activeCourses;

  const toggleCourseStatus = async (course: Course) => {
    const shouldHide = course.is_active !== false;
    const newStatus = !shouldHide;
    const actionText = shouldHide ? 'скрыть' : 'вернуть';

    if (!confirm(`Вы уверены, что хотите ${actionText} курс "${course.title}"?`)) return;

    try {
      await coursesService.toggleCourseStatus(course.id, newStatus);
      setAllCourses((prev) => prev.map(c =>
        c.id === course.id ? { ...c, is_active: newStatus } : c
      ));
      setError(null);
    } catch (err: any) {
      console.error(err);
      showErrorModal(err.response?.data?.message || err.message || `Ошибка операции с курсом`);
    }
  };

  const handleDeleteService = async (id: number, title: string) => {
    if (!confirm(`Вы уверены, что хотите удалить услугу "${title}"?`)) return;
    try {
      await servicesService.delete(id);
      setServicesList(prev => prev.filter(s => s.id !== id));
    } catch (err: any) {
      showErrorModal(err.message || 'Ошибка удаления услуги');
    }
  };

  // --- Логика Модальных окон ---

  const showErrorModal = (message: string) => {
    setErrorModal({ isOpen: true, message });
  };

  const closeErrorModal = () => {
    setErrorModal({ isOpen: false, message: '' });
  };

  // Открытие модалки отклонения
  const openRejectModal = (id: number) => {
    setSelectedReviewId(id);
    setRejectReason('');
    setIsRejectModalOpen(true);
  };

  // Закрытие модалки
  const closeRejectModal = () => {
    setIsRejectModalOpen(false);
    setSelectedReviewId(null);
    setRejectReason('');
  };

  const confirmReject = async () => {
    if (!selectedReviewId || !rejectReason.trim()) {
      showErrorModal('Пожалуйста, укажите причину отклонения (минимум 5 символов).');
      return;
    }

    try {
      // ВАЖНО: Используем PUT вместо DELETE
      await apiClient.put(`/admin/reviews/${selectedReviewId}`, {
        reason: rejectReason.trim(), // Передаем объект напрямую как тело запроса
      });

      loadPendingReviews();
      closeRejectModal();
    } catch (error: any) {
      const msg = error.response?.data?.error || error.response?.data?.message || 'Не удалось отклонить отзыв';
      showErrorModal(msg);
      console.error(error);
    }
  };


  const handleApproveReview = async (id: number) => {
    if (window.confirm('Одобрить этот отзыв?')) {
      try {
        // Роут: POST /reviews/:id/approve
        await apiClient.post(`/admin/reviews/${id}/approve`);
        loadPendingReviews();
      } catch (error: any) {
        const msg = error.response?.data?.message || error.response?.data?.error || 'Ошибка при одобрении';
        showErrorModal(msg);
      }
    }
  };

  const handleDeleteReview = async (id: number) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот отзыв?')) return;
    try {
      // Используем новый эндпоинт для админа
      await apiClient.delete(`/admin/reviews/${id}`);

      // Обновляем список
      loadPendingReviews();
    } catch (err: any) {
      showErrorModal(err.response?.data?.message || 'Ошибка при удалении отзыва');
    }
  };

  // Логика поиска пользователей
  useEffect(() => {
    const timer = setTimeout(() => {
      if (emailQuery.trim().length >= 2 && activeTab === 'grant-access' && emailQuery !== selectedUser?.email) {
        performSearch();
      } else if (emailQuery.trim().length === 0) {
        setSearchResults([]);
        setSelectedUser(null);
        setSelectedUserCourseIds([]);
      }
    }, 600);

    return () => clearTimeout(timer);
  }, [emailQuery, activeTab, selectedUser]);

  const performSearch = async () => {
    if (!emailQuery.trim()) return;
    setIsSearching(true);
    setSearchResults([]);
    setSelectedUser(null);
    setSelectedUserCourseIds([]);
    setAccessMessage(null);

    try {
      const response = await apiClient.post<UserSearchResult[]>('/admin/users/search', {
        email: emailQuery,
        limit: 10
      });
      setSearchResults(response.data);
    } catch (err: any) {
      console.error(err);
      showErrorModal(err.response?.data?.error || 'Ошибка при поиске пользователя.');
    } finally {
      setIsSearching(false);
    }
  };

  const handleSelectUser = async (user: UserSearchResult) => {
    setSelectedUser(user);
    setSearchResults([]);
    setEmailQuery(user.email);
    setSelectedCourseId('');
    setSelectedUserCourseIds([]);
    setAccessMessage(null);

    try {
      const response = await apiClient.get(`/admin/students/${user.id}`);
      const courseIds = Array.isArray(response.data?.courses)
        ? response.data.courses.map((course: any) => course.course_id)
        : [];
      setSelectedUserCourseIds(courseIds);
    } catch (err: any) {
      console.error(err);
      showErrorModal(err.response?.data?.error || 'Не удалось загрузить курсы пользователя.');
    }
  };

  // Логика выдачи доступа
  const handleGrantAccess = async () => {
    if (!selectedUser || !selectedCourseId) return;

    setIsGranting(true);
    setAccessMessage(null);

    try {
      await apiClient.post('/admin/access', {
        user_id: selectedUser.id,
        course_id: Number(selectedCourseId)
      });

      setAccessMessage({
        type: 'success',
        text: `Доступ к курсу успешно предоставлен пользователю ${selectedUser.email}`
      });

      setSelectedUser(null);
      setSelectedUserCourseIds([]);
      setSelectedCourseId('');
      setEmailQuery('');
      setSearchResults([]);

    } catch (err: any) {
      console.error(err);
      const errorMsg = err.response?.data?.error || err.response?.data?.message || 'Не удалось предоставить доступ';
      showErrorModal(errorMsg);
    } finally {
      setIsGranting(false);
    }
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

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Шапка */}
      <div className="bg-white border-b border-gray-200 py-6">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <a href="/" className="text-sm text-gray-500! hover:text-rose-500! transition-colors mb-2 inline-block" style={{ textDecoration: 'none' }}>
            ← На главную
          </a>
          <h1 className="text-2xl md:text-3xl font-serif font-bold text-gray-900">Админ-панель</h1>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        {error && (
          <div className="mb-4 p-3 bg-red-50 text-red-700 rounded-lg border border-red-100 flex items-center justify-between">
            <span>{error}</span>
            <button onClick={() => setError(null)} className="text-red-500 hover:text-red-700 ml-4">✕</button>
          </div>
        )}

        {/* Вкладки */}
        <div className="flex border-b border-gray-200 mb-6 overflow-x-auto">
          <button
            onClick={() => setActiveTab('courses')}
            className={`px-5 py-3 text-sm font-medium border-b-2 whitespace-nowrap transition-colors ${activeTab === 'courses' ? 'border-rose-500 text-rose-600' : 'border-transparent text-gray-500 hover:text-gray-700'}`}
          >
            Активные курсы ({activeCourses.length})
          </button>
          <button
            onClick={() => setActiveTab('hidden-courses')}
            className={`px-5 py-3 text-sm font-medium border-b-2 whitespace-nowrap transition-colors ${activeTab === 'hidden-courses' ? 'border-rose-500 text-rose-600' : 'border-transparent text-gray-500 hover:text-gray-700'}`}
          >
            Скрытые курсы ({hiddenCourses.length})
          </button>
          <button
            onClick={() => setActiveTab('services')}
            className={`px-5 py-3 text-sm font-medium border-b-2 whitespace-nowrap transition-colors ${activeTab === 'services' ? 'border-rose-500 text-rose-600' : 'border-transparent text-gray-500 hover:text-gray-700'}`}
          >
            Услуги ({servicesList.length})
          </button>
          <button
            onClick={() => setActiveTab('reviews')}
            className={`px-5 py-3 text-sm font-medium border-b-2 whitespace-nowrap transition-colors ${activeTab === 'reviews' ? 'border-rose-500 text-rose-600' : 'border-transparent text-gray-500 hover:text-gray-700'}`}
          >
            Отзывы ({pendingReviews.length})
          </button>
          <button
            onClick={() => setActiveTab('grant-access')}
            className={`px-5 py-3 text-sm font-medium border-b-2 whitespace-nowrap transition-colors ${activeTab === 'grant-access' ? 'border-rose-500 text-rose-600' : 'border-transparent text-gray-500 hover:text-gray-700'}`}
          >
            Выдача доступа
          </button>
        </div>

        {/* ===== Вкладка "Активные курсы" ===== */}
        {activeTab === 'courses' && (
          <div>
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-lg font-semibold text-gray-800">Управление курсами</h2>

              <button 
                  onClick={() => navigate('/admin/students')} 
                  className="px-4 py-2 bg-indigo-50 text-indigo-700 border border-indigo-200 rounded-lg hover:bg-indigo-100 hover:border-indigo-300 transition-all text-sm font-medium shadow-sm flex items-center gap-2"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                  Статистика учеников
                </button>
              <button onClick={() => navigate('/admin/courses/new')} className="px-4 py-2 bg-rose-500 text-white rounded-lg hover:bg-rose-600 transition-colors text-sm font-medium shadow-sm">
                + Создать курс
              </button>
            </div>
            {activeCourses.length === 0 ? (
              <EmptyState onCreate={() => navigate('/admin/courses/new')} icon="📚" title="Курсов пока нет" subtitle="Создайте первый учебный материал" btnText="Создать курс" />
            ) : (
              <div className="space-y-4">
                {activeCourses.map((course) => (
                  <CourseCard key={course.id} course={course} onEdit={() => navigate(`/admin/courses/${course.id}/edit`)} onToggle={() => toggleCourseStatus(course)} actionLabel="Скрыть" isHiddenView={false} />
                ))}
              </div>
            )}
          </div>
        )}

        {/* ===== Вкладка "Скрытые курсы" ===== */}
        {activeTab === 'hidden-courses' && (
          <div>
            <h2 className="text-lg font-semibold text-gray-800 mb-4">Архив скрытых курсов</h2>
            {hiddenCourses.length === 0 ? (
              <EmptyState onCreate={() => { }} icon="🗄️" title="Нет скрытых курсов" subtitle="Все курсы активны и видны пользователям" btnText={undefined} />
            ) : (
              <div className="space-y-4">
                {hiddenCourses.map((course) => (
                  <CourseCard key={course.id} course={course} onEdit={() => navigate(`/admin/courses/${course.id}/edit`)} onToggle={() => toggleCourseStatus(course)} actionLabel="Вернуть в каталог" isHiddenView={true} />
                ))}
              </div>
            )}
          </div>
        )}

        {/* ===== Вкладка "Услуги" ===== */}
        {activeTab === 'services' && (
          <div>
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-lg font-semibold text-gray-800">Управление услугами</h2>
              <button onClick={() => navigate('/admin/services/new')} className="px-4 py-2 bg-rose-500 text-white rounded-lg hover:bg-rose-600 transition-colors text-sm font-medium shadow-sm">
                + Создать услугу
              </button>
            </div>
            {servicesList.length === 0 ? (
              <EmptyState onCreate={() => navigate('/admin/services/new')} icon="💆‍♀️" title="Услуг пока нет" subtitle="Добавьте первую услугу для клиентов" btnText="Создать услугу" />
            ) : (
              <div className="space-y-4">
                {servicesList.map((service) => (
                  <div key={service.id} className="bg-white rounded-xl shadow-sm border border-gray-100 p-5 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-3 mb-2 flex-wrap">
                        <h3 className="text-lg font-bold text-gray-900 break-words">{service.title}</h3>
                        <span className="text-xs px-2 py-0.5 rounded-full font-medium bg-blue-50 text-blue-700 border border-blue-100">
                          {service.price} ₽
                        </span>
                        {service.duration_minutes && service.duration_minutes > 0 && (
                          <span className="text-xs px-2 py-0.5 rounded-full font-medium bg-gray-50 text-gray-600 border border-gray-100">
                            {service.duration_minutes} мин
                          </span>
                        )}
                      </div>
                      <p className="text-sm text-gray-600 break-words line-clamp-2">
                        {service.description?.split('|||')[0]}
                      </p>
                    </div>
                    <div className="flex flex-col sm:flex-row gap-2 w-full sm:w-auto shrink-0 pt-2 sm:pt-0 border-t sm:border-0 border-gray-100 mt-2 sm:mt-0">
                      <button onClick={() => navigate(`/admin/services/${service.id}/edit`)} className="w-full sm:w-auto px-4 py-2 text-sm font-medium text-blue-600 bg-blue-50 border border-blue-200 rounded-lg hover:bg-blue-100 transition-colors">
                        Редактировать
                      </button>
                      <button onClick={() => handleDeleteService(service.id, service.title)} className="w-full sm:w-auto px-4 py-2 text-sm font-medium text-red-600 bg-red-50 border border-red-200 rounded-lg hover:bg-red-100 transition-colors">
                        Удалить
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* ===== Вкладка "Отзывы" ===== */}
        {activeTab === 'reviews' && (
          <div>
            <h2 className="text-lg font-semibold text-gray-800 mb-4">Ожидают модерации</h2>
            {pendingReviews.length === 0 ? (
              <EmptyState onCreate={() => { }} icon="✅" title="Всё проверено" subtitle="Нет отзывов, ожидающих модерации" btnText={undefined} />
            ) : (
              <div className="space-y-4">
                {pendingReviews.map((review) => (
                  <div key={review.id} className="bg-white rounded-xl shadow-sm border border-gray-100 p-5">
                    <div className="flex flex-col sm:flex-row items-start justify-between gap-4">
                      <div className="flex-1 w-full">
                        <div className="mb-2">
                          <span className="inline-block px-2 py-1 bg-rose-50 text-rose-600 text-xs font-bold rounded-md border border-rose-100">
                            Курс: {review.course_title}
                          </span>
                        </div>

                        <div className="flex items-center gap-2 mb-2">
                          <div className="flex gap-0.5 text-yellow-400">
                            {[1, 2, 3, 4, 5].map((star) => (
                              <span key={star}>{star <= review.rating ? '★' : '☆'}</span>
                            ))}
                          </div>
                          <span className="text-xs text-gray-400">{timeAgo(review.created_at)}</span>
                        </div>
                        <p className="text-sm text-gray-700 leading-relaxed break-words">{review.text}</p>
                        <p className="text-xs text-gray-500 mt-2">Автор: {review.author_name || 'Аноним'}</p>
                      </div>

                      <div className="flex gap-2 w-full sm:w-auto">
                        <button
                          onClick={() => handleApproveReview(review.id)}
                          className="flex-1 sm:flex-none px-3 py-1.5 text-sm bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors shadow-sm"
                        >
                          ✓ Одобрить
                        </button>
                        <button
                          onClick={() => openRejectModal(review.id)}
                          className="flex-1 sm:flex-none px-3 py-1.5 text-sm text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors"
                        >
                          ✕ Отклонить
                        </button>
                        {/* Кнопка удаления */}
                        <button
                          onClick={() => handleDeleteReview(review.id)}
                          className="w-full sm:w-auto px-3 py-1.5 text-xs text-gray-500 border border-gray-200 rounded-lg hover:bg-gray-50 hover:text-gray-700 hover:border-gray-300 transition-colors flex items-center justify-center gap-1 mt-1 sm:mt-0"
                          title="Полностью удалить отзыв из базы"
                        >
                          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                          </svg>
                          Удалить
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}

            {/* ===== МОДАЛЬНОЕ ОКНО ОТКЛОНЕНИЯ ===== */}
            {isRejectModalOpen && (
              <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
                <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 animate-scaleIn">
                  <h3 className="text-xl font-bold text-gray-900 mb-2">Отклонение отзыва</h3>
                  <p className="text-sm text-gray-500 mb-4">
                    Укажите причину, по которой этот отзыв не будет опубликован.
                  </p>

                  <textarea
                    value={rejectReason}
                    onChange={(e) => setRejectReason(e.target.value)}
                    placeholder="Например: Отзыв содержит нецензурную лексику..."
                    className="w-full p-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-red-200 focus:border-red-400 min-h-[100px] text-sm"
                    autoFocus
                  />

                  <div className="flex gap-3 mt-6">
                    <button
                      onClick={closeRejectModal}
                      className="flex-1 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-xl hover:bg-gray-200 transition-colors"
                    >
                      Отмена
                    </button>
                    <button
                      onClick={confirmReject}
                      disabled={!rejectReason.trim()}
                      className="flex-1 px-4 py-2 text-sm font-medium text-white bg-red-500 rounded-xl hover:bg-red-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                      Отклонить
                    </button>
                  </div>
                </div>
              </div>
            )}

            {/* ===== МОДАЛЬНОЕ ОКНО ОШИБКИ ===== */}
            {errorModal.isOpen && (
              <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
                <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 animate-scaleIn border-l-4 border-red-500">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center text-red-600 flex-shrink-0">
                      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
                    </div>
                    <h3 className="text-xl font-bold text-gray-900">Ошибка операции</h3>
                  </div>

                  <p className="text-gray-600 mb-6 text-sm leading-relaxed bg-gray-50 p-3 rounded-lg border border-gray-100">
                    {errorModal.message}
                  </p>

                  <button
                    onClick={closeErrorModal}
                    className="w-full px-4 py-2.5 text-sm font-bold text-white bg-red-500 rounded-xl hover:bg-red-600 transition-colors shadow-md"
                  >
                    Понятно
                  </button>
                </div>
              </div>
            )}
          </div>
        )}

        {/* ===== НОВАЯ ВКЛАКА: Выдача доступа ===== */}
        {activeTab === 'grant-access' && (
          <div className="max-w-4xl mx-auto">
            <div className="mb-6">
              <h2 className="text-2xl font-serif font-bold text-gray-900">Ручная выдача доступа</h2>
              <p className="text-gray-600 mt-2">Найдите пользователя по email и предоставьте доступ к любому курсу вручную.</p>
            </div>

            {accessMessage && (
              <div className={`mb-6 p-4 rounded-lg border flex items-center gap-3 ${accessMessage.type === 'success'
                ? 'bg-green-50 border-green-200 text-green-800'
                : 'bg-red-50 border-red-200 text-red-800'
                }`}>
                <span>{accessMessage.type === 'success' ? '✅' : '⚠️'}</span>
                <span className="flex-1">{accessMessage.text}</span>
                <button onClick={() => setAccessMessage(null)} className="ml-auto opacity-70 hover:opacity-100">✕</button>
              </div>
            )}

            <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 md:p-8 space-y-6">

              {/* Шаг 1: Поиск */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  1. Поиск пользователя по email
                </label>
                <div className="relative">
                  <input
                    type="text"
                    value={emailQuery}
                    onChange={(e) => setEmailQuery(e.target.value)}
                    placeholder="Начните вводить email (минимум 2 символа)..."
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-500 focus:border-rose-500 outline-none transition-shadow pr-10"
                  />
                  {isSearching && (
                    <div className="absolute right-3 top-3.5">
                      <div className="w-5 h-5 border-2 border-rose-300 border-t-rose-500 rounded-full animate-spin"></div>
                    </div>
                  )}
                </div>

                {/* Исправлено: Добавлена проверка searchResults на null/undefined перед обращением к length */}
                {searchResults && searchResults.length > 0 && (
                  <div className="mt-2 border border-gray-200 rounded-lg overflow-hidden shadow-sm max-h-60 overflow-y-auto bg-white z-10 relative">
                    <ul className="divide-y divide-gray-100">
                      {searchResults.map((user) => (
                        <li
                          key={user.id}
                          onClick={() => handleSelectUser(user)}
                          className="p-3 cursor-pointer hover:bg-rose-50 transition-colors flex justify-between items-center"
                        >
                          <div>
                            <p className="font-medium text-gray-900">{user.name || 'Без имени'}</p>
                            <p className="text-sm text-gray-500">{user.email}</p>
                          </div>
                        </li>
                      ))}
                    </ul>
                  </div>
                )}

                {emailQuery.length >= 2 && !isSearching && (!searchResults || searchResults.length === 0) && !selectedUser && (
                  <p className="mt-2 text-sm text-gray-500 italic">Пользователи не найдены.</p>
                )}
              </div>

              {/* Шаг 2: Выбор курса */}
              <div className={`transition-opacity duration-300 ${selectedUser ? 'opacity-100' : 'opacity-50 pointer-events-none'}`}>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  2. Выберите курс для доступа
                </label>
                <select
                  value={selectedCourseId}
                  onChange={(e) => setSelectedCourseId(Number(e.target.value))}
                  disabled={!selectedUser}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-500 focus:border-rose-500 outline-none bg-white disabled:bg-gray-50"
                >
                  <option value="">-- Выберите курс --</option>
                  {availableCoursesForSelectedUser && availableCoursesForSelectedUser.length > 0 ? (
                    availableCoursesForSelectedUser.map((course) => (
                      <option key={course.id} value={course.id}>
                        {course.title} ({course.price === 0 ? 'Бесплатно' : `${course.price} ₽`})
                      </option>
                    ))
                  ) : (
                    <option disabled>{selectedUser ? 'У пользователя уже есть все активные курсы' : 'Нет активных курсов'}</option>
                  )}
                </select>
              </div>

              {/* Шаг 3: Подтверждение */}
              {selectedUser && selectedCourseId && (
                <div className="pt-6 border-t border-gray-100">
                  <div className="bg-blue-50 p-4 rounded-lg mb-6">
                    <p className="text-sm text-blue-800">
                      Вы собираетесь выдать доступ пользователю <strong>{selectedUser.email}</strong> к курсу:
                    </p>
                    <p className="text-lg font-bold text-blue-900 mt-1">
                      {availableCoursesForSelectedUser.find(c => c.id === selectedCourseId)?.title || 'Выбранный курс'}
                    </p>
                  </div>

                  <button
                    onClick={handleGrantAccess}
                    disabled={isGranting}
                    className="w-full sm:w-auto px-8 py-3 bg-rose-500 text-white font-semibold rounded-lg hover:bg-rose-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed shadow-md"
                  >
                    {isGranting ? 'Обработка...' : 'Подтвердить и выдать доступ'}
                  </button>
                </div>
              )}
            </div>
          </div>
        )}

      </div>
    </div>
  );
};

const CourseCard = ({ course, onEdit, onToggle, actionLabel, isHiddenView }: { course: Course, onEdit: () => void, onToggle: () => void, actionLabel: string, isHiddenView: boolean }) => {
  return (
    <div className={`bg-white rounded-xl shadow-sm border overflow-hidden transition-all ${isHiddenView ? 'border-gray-200 opacity-75 grayscale-[0.3]' : 'border-gray-100 hover:shadow-md'}`}>
      <div className="p-5 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div className="flex-1 min-w-0 w-full">
          <div className="flex items-center gap-3 mb-2 flex-wrap">
            <h3 className="text-lg font-bold text-gray-900 break-words w-full sm:w-auto">{course.title}</h3>
            <span className={`text-xs px-2 py-0.5 rounded-full font-medium shrink-0 ${course.is_public ? 'bg-green-100 text-green-700' : 'bg-blue-50 text-blue-700 border border-blue-100'}`}>
              {course.is_public ? 'Бесплатный' : `${course.price} ₽`}
            </span>
          </div>
          <p className="text-sm text-gray-600 break-words line-clamp-3 mb-3">{course.description}</p>
        </div>
        <div className="flex flex-col sm:flex-row gap-2 w-full sm:w-auto shrink-0 pt-2 sm:pt-0 border-t sm:border-0 border-gray-100 mt-2 sm:mt-0">
          <button onClick={onEdit} className="w-full sm:w-auto px-4 py-2 text-sm font-medium text-blue-600 bg-blue-50 border border-blue-200 rounded-lg hover:bg-blue-100 transition-colors">Редактировать</button>
          <button onClick={onToggle} className={`w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg transition-colors border ${isHiddenView ? 'text-green-600 bg-green-50 border-green-200 hover:bg-green-100' : 'text-red-600 bg-red-50 border-red-200 hover:bg-red-100'}`}>{actionLabel}</button>
        </div>
      </div>
    </div>
  );
};

const EmptyState = ({ icon, title, subtitle, btnText, onCreate }: any) => (
  <div className="bg-white p-12 rounded-xl shadow-sm border border-gray-100 text-center">
    <div className="text-5xl mb-4">{icon}</div>
    <h3 className="text-lg font-semibold text-gray-700 mb-2">{title}</h3>
    <p className="text-gray-500 mb-6">{subtitle}</p>
    {btnText && onCreate && (
      <button onClick={onCreate} className="px-6 py-2.5 bg-rose-500 text-white rounded-lg hover:bg-rose-600 transition-colors text-sm font-medium">{btnText}</button>
    )}
  </div>
);

export default AdminPage;
