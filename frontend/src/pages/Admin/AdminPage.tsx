import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { coursesService, reviewsService } from '../../api';
import type { Course, Review } from '../../api/types';
import { timeAgo } from '../../utils/dateUtils';

type AdminTab = 'courses' | 'hidden-courses' | 'reviews';

// --- Вспомогательные функции (вынесены наверх) ---

const getTotalLessons = (course: Course) => {
  if (!course.modules || !Array.isArray(course.modules)) return 0;
  return course.modules.reduce((acc, m) => acc + (m.lessons?.length || 0), 0);

};

const AdminPage = () => {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<AdminTab>('courses');

  const [allCourses, setAllCourses] = useState<Course[]>([]);
  const [pendingReviews, setPendingReviews] = useState<Review[]>([]);

  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const [coursesData, reviewsData] = await Promise.all([
          coursesService.getAllCoursesAdmin(), 
          reviewsService.getPendingReviews(),
        ]);

        setAllCourses(Array.isArray(coursesData) ? coursesData : []);
        setPendingReviews(Array.isArray(reviewsData) ? reviewsData : []);
      } catch (err: any) {
        console.error(err);
        setError(err.message || 'Не удалось загрузить данные');
      } finally {
        setIsLoading(false);
      }
    };
    loadData();
  }, []);

  // Фильтрация на клиенте
  const activeCourses = allCourses.filter(c => c.is_active !== false);
  const hiddenCourses = allCourses.filter(c => c.is_active === false);

  // === Действия с курсами ===

  const toggleCourseStatus = async (course: Course) => {
  const shouldHide = course.is_active !== false;
  const newStatus = !shouldHide; // true если возвращаем, false если скрываем
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
    setError(err.response?.data?.message || err.message || `Ошибка операции с курсом`);
  }
};

  // === Отзывы ===
  const handleApproveReview = async (id: number) => {
    try {
      await reviewsService.approveReview(id);
      setPendingReviews((prev) => prev.filter((r) => r.id !== id));
    } catch (err: any) {
      setError(err.message || 'Ошибка одобрения');
    }
  };

  const handleRejectReview = async (id: number) => {
    try {
      await reviewsService.rejectReview(id);
      setPendingReviews((prev) => prev.filter((r) => r.id !== id));
    } catch (err: any) {
      setError(err.message || 'Ошибка отклонения');
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
        {/* Ошибка */}
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
            className={`px-5 py-3 text-sm font-medium border-b-2 whitespace-nowrap transition-colors ${activeTab === 'courses' ? 'border-rose-500 text-rose-600' : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
          >
            Активные курсы ({activeCourses.length})
          </button>
          <button
            onClick={() => setActiveTab('hidden-courses')}
            className={`px-5 py-3 text-sm font-medium border-b-2 whitespace-nowrap transition-colors ${activeTab === 'hidden-courses' ? 'border-rose-500 text-rose-600' : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
          >
            Скрытые курсы ({hiddenCourses.length})
          </button>
          <button
            onClick={() => setActiveTab('reviews')}
            className={`px-5 py-3 text-sm font-medium border-b-2 whitespace-nowrap transition-colors ${activeTab === 'reviews' ? 'border-rose-500 text-rose-600' : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
          >
            Отзывы на модерацию ({pendingReviews.length})
          </button>
        </div>

        {/* ===== Вкладка "Активные курсы" ===== */}
        {activeTab === 'courses' && (
          <div>
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-lg font-semibold text-gray-800">Управление курсами</h2>
              <button
                onClick={() => navigate('/admin/courses/new')}
                className="px-4 py-2 bg-rose-500 text-white rounded-lg hover:bg-rose-600 transition-colors text-sm font-medium shadow-sm"
              >
                + Создать курс
              </button>
            </div>

            {activeCourses.length === 0 ? (
              <EmptyState onCreate={() => navigate('/admin/courses/new')} icon="📚" title="Курсов пока нет" subtitle="Создайте первый учебный материал" btnText="Создать курс" />
            ) : (
              <div className="space-y-4">
                {activeCourses.map((course) => (
                  <CourseCard
                    key={course.id}
                    course={course}
                    onEdit={() => navigate(`/admin/courses/${course.id}/edit`)}
                    onToggle={() => toggleCourseStatus(course)}
                    actionLabel="Скрыть"
                    isHiddenView={false}
                  />
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
                  <CourseCard
                    key={course.id}
                    course={course}
                    onEdit={() => navigate(`/admin/courses/${course.id}/edit`)}
                    onToggle={() => toggleCourseStatus(course)}
                    actionLabel="Вернуть в каталог"
                    isHiddenView={true}
                  />
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
                          onClick={() => handleRejectReview(review.id)}
                          className="flex-1 sm:flex-none px-3 py-1.5 text-sm text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors"
                        >
                          ✕ Отклонить
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

// Компонент карточки курса
const CourseCard = ({
  course,
  onEdit,
  onToggle,
  actionLabel,
  isHiddenView
}: {
  course: Course,
  onEdit: () => void,
  onToggle: () => void,
  actionLabel: string,
  isHiddenView: boolean
}) => {
 

  return (
    <div className={`bg-white rounded-xl shadow-sm border overflow-hidden transition-all ${isHiddenView ? 'border-gray-200 opacity-75 grayscale-[0.3]' : 'border-gray-100 hover:shadow-md'}`}>
      {/* Основная часть: Адаптивная верстка */}
      <div className="p-5 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">

        {/* Левая часть: Информация */}
        <div className="flex-1 min-w-0 w-full">
          <div className="flex items-center gap-3 mb-2 flex-wrap">
            <h3 className="text-lg font-bold text-gray-900 break-words w-full sm:w-auto">{course.title}</h3>
            <span className={`text-xs px-2 py-0.5 rounded-full font-medium shrink-0 ${course.is_public
                ? 'bg-green-100 text-green-700'
                : 'bg-blue-50 text-blue-700 border border-blue-100'
              }`}>
              {course.is_public ? 'Бесплатный' : `${course.price} ₽`}
            </span>
          </div>

          <p className="text-sm text-gray-600 break-words line-clamp-3 mb-3">
            {course.description}
          </p>

        </div>

        {/* Правая часть: Кнопки */}
        <div className="flex flex-col sm:flex-row gap-2 w-full sm:w-auto shrink-0 pt-2 sm:pt-0 border-t sm:border-0 border-gray-100 mt-2 sm:mt-0">
          <button
            onClick={onEdit}
            className="w-full sm:w-auto px-4 py-2 text-sm font-medium text-blue-600 bg-blue-50 border border-blue-200 rounded-lg hover:bg-blue-100 transition-colors"
          >
            Редактировать
          </button>
          <button
            onClick={onToggle}
            className={`w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg transition-colors border ${isHiddenView
                ? 'text-green-600 bg-green-50 border-green-200 hover:bg-green-100'
                : 'text-red-600 bg-red-50 border-red-200 hover:bg-red-100'
              }`}
          >
            {actionLabel}
          </button>
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
      <button onClick={onCreate} className="px-6 py-2.5 bg-rose-500 text-white rounded-lg hover:bg-rose-600 transition-colors text-sm font-medium">
        {btnText}
      </button>
    )}
  </div>
);

export default AdminPage;