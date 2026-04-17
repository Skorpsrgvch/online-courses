import { useState, useEffect } from 'react';
import { coursesService, reviewsService } from '../../api';
import type { Course, Review, CreateCourseDto, CreateLessonDto } from '../../api/types';
import { timeAgo } from '../../utils/dateUtils';

type AdminTab = 'courses' | 'reviews';

const AdminPage = () => {
  const [activeTab, setActiveTab] = useState<AdminTab>('courses');
  const [courses, setCourses] = useState<Course[]>([]);
  const [pendingReviews, setPendingReviews] = useState<Review[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Модалка создания курса
  const [showCreateCourse, setShowCreateCourse] = useState(false);
  const [newCourse, setNewCourse] = useState<CreateCourseDto>({
    title: '', description: '', is_public: false, price: 0, cover_image_url: '',
  });

  // Модалка создания модуля
  const [showCreateModule, setShowCreateModule] = useState(false);
  const [selectedCourseId, setSelectedCourseId] = useState<number | null>(null);
  const [newModule, setNewModule] = useState<{ title: string; order: number }>({ title: '', order: 1 });

  // Модалка создания урока
  const [showCreateLesson, setShowCreateLesson] = useState(false);
  const [newLesson, setNewLesson] = useState<CreateLessonDto>({
    module_id: 0, title: '', description: '', lesson_type: 'video',
    video_embed_id: '', article_content: '', order: 1,
  });

  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const [coursesData, reviewsData] = await Promise.all([
          coursesService.getAllCourses(),
          reviewsService.getPendingReviews(),
        ]);
        setCourses(Array.isArray(coursesData) ? coursesData : []);
        setPendingReviews(Array.isArray(reviewsData) ? reviewsData : []);
      } catch (err: any) {
        setError(err.message || 'Не удалось загрузить данные');
      } finally {
        setIsLoading(false);
      }
    };
    loadData();
  }, []);

  // === Курсы ===

  const handleCreateCourse = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await coursesService.createCourse(newCourse);
      setShowCreateCourse(false);
      setNewCourse({ title: '', description: '', is_public: false, price: 0, cover_image_url: '' });
      const updated = await coursesService.getAllCourses();
      setCourses(updated);
    } catch (err: any) {
      setError(err.message || 'Ошибка создания курса');
    }
  };

  const handleDeleteCourse = async (id: number) => {
    if (!confirm('Удалить курс?')) return;
    try {
      await coursesService.deleteCourse(id);
      setCourses((prev) => prev.filter((c) => c.id !== id));
    } catch (err: any) {
      setError(err.message || 'Ошибка удаления курса');
    }
  };

  const openCreateModule = (courseId: number) => {
    setSelectedCourseId(courseId);
    setNewModule({ title: '', order: 1 });
    setShowCreateModule(true);
  };

  const handleCreateModule = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedCourseId) return;
    try {
      await coursesService.createModule({
        course_id: selectedCourseId,
        title: newModule.title,
        order: newModule.order,
      });
      setShowCreateModule(false);
      // Перезагрузим курсы
      const updated = await coursesService.getAllCourses();
      setCourses(updated);
    } catch (err: any) {
      setError(err.message || 'Ошибка создания модуля');
    }
  };

  const openCreateLesson = (moduleId: number) => {
    setNewLesson({
      module_id: moduleId, title: '', description: '', lesson_type: 'video',
      video_embed_id: '', article_content: '', order: 1,
    });
    setShowCreateLesson(true);
  };

  const handleCreateLesson = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await coursesService.createLesson(newLesson);
      setShowCreateLesson(false);
      const updated = await coursesService.getAllCourses();
      setCourses(updated);
    } catch (err: any) {
      setError(err.message || 'Ошибка создания урока');
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
          <a href="/" className="text-sm text-gray-500 hover:text-rose-500 transition-colors mb-2 inline-block">
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
        <div className="flex border-b border-gray-200 mb-6">
          <button
            onClick={() => setActiveTab('courses')}
            className={`px-5 py-3 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'courses'
                ? 'border-rose-500 text-rose-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Курсы ({courses.length})
          </button>
          <button
            onClick={() => setActiveTab('reviews')}
            className={`px-5 py-3 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'reviews'
                ? 'border-rose-500 text-rose-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Отзывы на модерацию ({pendingReviews.length})
          </button>
        </div>

        {/* ===== Вкладка "Курсы" ===== */}
        {activeTab === 'courses' && (
          <div>
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-lg font-semibold text-gray-800">Управление курсами</h2>
              <button
                onClick={() => setShowCreateCourse(true)}
                className="px-4 py-2 bg-rose-500 text-white rounded-lg hover:bg-rose-600 transition-colors text-sm font-medium"
              >
                + Создать курс
              </button>
            </div>

            {courses.length === 0 ? (
              <div className="bg-white p-12 rounded-xl shadow-sm border border-gray-100 text-center">
                <div className="text-5xl mb-4">📚</div>
                <h3 className="text-lg font-semibold text-gray-700 mb-2">Курсов пока нет</h3>
                <p className="text-gray-500 mb-4">Создайте первый курс</p>
                <button
                  onClick={() => setShowCreateCourse(true)}
                  className="px-6 py-2.5 bg-rose-500 text-white rounded-lg hover:bg-rose-600 transition-colors text-sm font-medium"
                >
                  Создать курс
                </button>
              </div>
            ) : (
              <div className="space-y-4">
                {courses.map((course) => (
                  <div key={course.id} className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                    <div className="p-5 flex items-start justify-between gap-4">
                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-1">
                          <h3 className="text-lg font-semibold text-gray-800">{course.title}</h3>
                          <span className={`text-xs px-2 py-0.5 rounded-full ${
                            course.is_public ? 'bg-green-100 text-green-700' : 'bg-amber-100 text-amber-700'
                          }`}>
                            {course.is_public ? 'Бесплатный' : `${course.price} ₽`}
                          </span>
                        </div>
                        <p className="text-sm text-gray-500 line-clamp-2">{course.description}</p>
                      </div>
                      <button
                        onClick={() => handleDeleteCourse(course.id)}
                        className="flex-shrink-0 px-3 py-1.5 text-sm text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors"
                      >
                        Удалить
                      </button>
                    </div>

                    {/* Модули и уроки */}
                    <div className="border-t border-gray-100 px-5 py-3 bg-gray-50/50">
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-medium text-gray-600">Модули</span>
                        <button
                          onClick={() => openCreateModule(course.id)}
                          className="text-xs px-3 py-1 bg-rose-100 text-rose-600 rounded-lg hover:bg-rose-200 transition-colors"
                        >
                          + Модуль
                        </button>
                      </div>
                      {course.modules && course.modules.length > 0 ? (
                        <div className="space-y-2">
                          {course.modules.map((mod) => (
                            <div key={mod.id} className="bg-white rounded-lg border border-gray-100 p-3">
                              <div className="flex items-center justify-between mb-1">
                                <span className="text-sm font-medium text-gray-700">{mod.title}</span>
                                <button
                                  onClick={() => openCreateLesson(mod.id)}
                                  className="text-xs px-2 py-0.5 bg-blue-50 text-blue-600 rounded hover:bg-blue-100 transition-colors"
                                >
                                  + Урок
                                </button>
                              </div>
                              {mod.lessons && mod.lessons.length > 0 && (
                                <ul className="text-xs text-gray-500 space-y-0.5 pl-4 mt-1">
                                  {mod.lessons.map((lesson) => (
                                    <li key={lesson.id} className="flex items-center gap-1.5">
                                      <span className="text-gray-400">
                                        {lesson.lesson_type === 'video' ? '🎬' : '📄'}
                                      </span>
                                      {lesson.title}
                                    </li>
                                  ))}
                                </ul>
                              )}
                            </div>
                          ))}
                        </div>
                      ) : (
                        <p className="text-xs text-gray-400">Модулей пока нет</p>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* ===== Вкладка "Отзывы на модерацию" ===== */}
        {activeTab === 'reviews' && (
          <div>
            <h2 className="text-lg font-semibold text-gray-800 mb-4">Ожидают модерации</h2>
            {pendingReviews.length === 0 ? (
              <div className="bg-white p-12 rounded-xl shadow-sm border border-gray-100 text-center">
                <div className="text-5xl mb-4">✅</div>
                <h3 className="text-lg font-semibold text-gray-700 mb-2">Всё проверено</h3>
                <p className="text-gray-500">Нет отзывов, ожидающих модерации</p>
              </div>
            ) : (
              <div className="space-y-4">
                {pendingReviews.map((review) => (
                  <div key={review.id} className="bg-white rounded-xl shadow-sm border border-gray-100 p-5">
                    <div className="flex items-start justify-between gap-4">
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-2">
                          {/* Звёзды */}
                          <div className="flex gap-0.5">
                            {[1, 2, 3, 4, 5].map((star) => (
                              <span key={star} className={star <= review.rating ? 'text-yellow-400' : 'text-gray-300'}>
                                ★
                              </span>
                            ))}
                          </div>
                          <span className="text-xs text-gray-400">{timeAgo(review.created_at)}</span>
                        </div>
                        <p className="text-sm text-gray-700">{review.text}</p>
                      </div>
                      <div className="flex gap-2 flex-shrink-0">
                        <button
                          onClick={() => handleApproveReview(review.id)}
                          className="px-3 py-1.5 text-sm bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
                        >
                          ✓ Одобрить
                        </button>
                        <button
                          onClick={() => handleRejectReview(review.id)}
                          className="px-3 py-1.5 text-sm text-red-600 border border-red-200 rounded-lg hover:bg-red-50 transition-colors"
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

      {/* ===== Модалка: Создать курс ===== */}
      {showCreateCourse && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" onClick={() => setShowCreateCourse(false)}>
          <div className="bg-white rounded-xl shadow-xl max-w-lg w-full p-6" onClick={(e) => e.stopPropagation()}>
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Новый курс</h3>
            <form onSubmit={handleCreateCourse} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Название</label>
                <input
                  type="text"
                  value={newCourse.title}
                  onChange={(e) => setNewCourse({ ...newCourse, title: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Описание</label>
                <textarea
                  value={newCourse.description}
                  onChange={(e) => setNewCourse({ ...newCourse, description: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg resize-none focus:ring-2 focus:ring-rose-300"
                  rows={3}
                  required
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Цена (₽)</label>
                  <input
                    type="number"
                    value={newCourse.price}
                    onChange={(e) => setNewCourse({ ...newCourse, price: Number(e.target.value) })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300"
                    min={0}
                  />
                </div>
                <div className="flex items-center pt-6">
                  <label className="flex items-center gap-2 text-sm text-gray-700">
                    <input
                      type="checkbox"
                      checked={newCourse.is_public}
                      onChange={(e) => setNewCourse({ ...newCourse, is_public: e.target.checked })}
                      className="rounded border-gray-300 text-rose-500"
                    />
                    Бесплатный курс
                  </label>
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">URL обложки</label>
                <input
                  type="url"
                  value={newCourse.cover_image_url || ''}
                  onChange={(e) => setNewCourse({ ...newCourse, cover_image_url: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300"
                  placeholder="https://..."
                />
              </div>
              <div className="flex gap-3 pt-2">
                <button type="submit" className="flex-1 px-4 py-2 bg-rose-500 text-white rounded-lg hover:bg-rose-600 transition-colors font-medium">
                  Создать
                </button>
                <button type="button" onClick={() => setShowCreateCourse(false)} className="px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors">
                  Отмена
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* ===== Модалка: Создать модуль ===== */}
      {showCreateModule && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" onClick={() => setShowCreateModule(false)}>
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full p-6" onClick={(e) => e.stopPropagation()}>
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Новый модуль</h3>
            <form onSubmit={handleCreateModule} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Название</label>
                <input
                  type="text"
                  value={newModule.title}
                  onChange={(e) => setNewModule({ ...newModule, title: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300"
                  required
                />
              </div>
              <div className="flex gap-3 pt-2">
                <button type="submit" className="flex-1 px-4 py-2 bg-rose-500 text-white rounded-lg hover:bg-rose-600 transition-colors font-medium">
                  Создать
                </button>
                <button type="button" onClick={() => setShowCreateModule(false)} className="px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors">
                  Отмена
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* ===== Модалка: Создать урок ===== */}
      {showCreateLesson && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" onClick={() => setShowCreateLesson(false)}>
          <div className="bg-white rounded-xl shadow-xl max-w-lg w-full p-6" onClick={(e) => e.stopPropagation()}>
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Новый урок</h3>
            <form onSubmit={handleCreateLesson} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Название</label>
                <input
                  type="text"
                  value={newLesson.title}
                  onChange={(e) => setNewLesson({ ...newLesson, title: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Описание</label>
                <input
                  type="text"
                  value={newLesson.description}
                  onChange={(e) => setNewLesson({ ...newLesson, description: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Тип урока</label>
                <select
                  value={newLesson.lesson_type}
                  onChange={(e) => setNewLesson({ ...newLesson, lesson_type: e.target.value as 'video' | 'article' })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300"
                >
                  <option value="video">Видео</option>
                  <option value="article">Статья</option>
                </select>
              </div>
              {newLesson.lesson_type === 'video' ? (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Video Embed ID (RuTube)</label>
                  <input
                    type="text"
                    value={newLesson.video_embed_id || ''}
                    onChange={(e) => setNewLesson({ ...newLesson, video_embed_id: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300"
                    placeholder="abc123-def456"
                  />
                </div>
              ) : (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Содержание (HTML)</label>
                  <textarea
                    value={newLesson.article_content || ''}
                    onChange={(e) => setNewLesson({ ...newLesson, article_content: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg resize-none focus:ring-2 focus:ring-rose-300"
                    rows={5}
                  />
                </div>
              )}
              <div className="flex gap-3 pt-2">
                <button type="submit" className="flex-1 px-4 py-2 bg-rose-500 text-white rounded-lg hover:bg-rose-600 transition-colors font-medium">
                  Создать
                </button>
                <button type="button" onClick={() => setShowCreateLesson(false)} className="px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors">
                  Отмена
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminPage;
