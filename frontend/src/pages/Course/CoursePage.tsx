import { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import { coursesService, reviewsService } from '../../api';
import type { CourseFullResponse, Review } from '../../api/types';
import { VideoPlayer } from '../../components/course/VideoPlayer';
import { ArticleViewer } from '../../components/course/ArticleViewer';
import { Comments } from '../../components/course/Comments';

const CoursePage = () => {
  const { id } = useParams<{ id: string }>();
  const courseId = Number(id);

  const [courseData, setCourseData] = useState<CourseFullResponse | null>(null);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Выбранный урок: [moduleIndex, lessonIndex]
  const [selectedLesson, setSelectedLesson] = useState<{ moduleIdx: number; lessonIdx: number } | null>(null);
  // Развёрнутые модули (по умолчанию первый развёрнут)
  const [expandedModules, setExpandedModules] = useState<Set<number>>(new Set([0]));

  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const [courseRes, reviewsData] = await Promise.all([
          coursesService.getCourseFull(courseId),
          reviewsService.getCourseReviews(courseId),
        ]);
        setCourseData(courseRes);
        setReviews(reviewsData);
        // Автовыбор первого урока
        if (courseRes.modules.length > 0 && courseRes.modules[0].lessons.length > 0) {
          setSelectedLesson({ moduleIdx: 0, lessonIdx: 0 });
        }
      } catch (err: any) {
        setError(err.message || 'Не удалось загрузить курс');
      } finally {
        setIsLoading(false);
      }
    };

    if (courseId) loadData();
  }, [courseId]);

  const toggleModule = useCallback((idx: number) => {
    setExpandedModules((prev) => {
      const next = new Set(prev);
      if (next.has(idx)) {
        next.delete(idx);
      } else {
        next.add(idx);
      }
      return next;
    });
  }, []);

  const selectLesson = useCallback((moduleIdx: number, lessonIdx: number) => {
    setSelectedLesson({ moduleIdx, lessonIdx });
    // Разворачиваем модуль при выборе урока
    setExpandedModules((prev) => {
      const next = new Set(prev);
      next.add(moduleIdx);
      return next;
    });
  }, []);

  const handleReviewSubmitted = useCallback(async () => {
    const updated = await reviewsService.getCourseReviews(courseId);
    setReviews(updated);
  }, [courseId]);

  const handleLessonComplete = useCallback(async (lessonId: number) => {
    try {
      await coursesService.markLessonComplete(lessonId);
    } catch (err) {
      console.error('Ошибка при отметке урока:', err);
    }
  }, []);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="w-12 h-12 border-4 border-rose-300 border-t-rose-500 rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-500">Загрузка курса...</p>
        </div>
      </div>
    );
  }

  if (error || !courseData) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center bg-white p-8 rounded-xl shadow-sm max-w-md">
          <div className="text-5xl mb-4">😔</div>
          <h2 className="text-xl font-semibold text-gray-800 mb-2">Ошибка загрузки</h2>
          <p className="text-gray-600 mb-4">{error || 'Курс не найден'}</p>
          <a href="/" className="text-rose-500 hover:underline font-medium">Вернуться на главную</a>
        </div>
      </div>
    );
  }

  const { course, modules } = courseData;
  const currentModule = selectedLesson ? modules[selectedLesson.moduleIdx] : null;
  const currentLesson = selectedLesson && currentModule ? currentModule.lessons[selectedLesson.lessonIdx] : null;

  // Всего уроков и выполненных
  const totalLessons = modules.reduce((sum, m) => sum + m.lessons.length, 0);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Шапка курса */}
      <div className="bg-white border-b border-gray-200 py-6">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <a href="/" className="text-sm text-gray-500 hover:text-rose-500 transition-colors mb-2 inline-block">
            ← Назад к каталогу
          </a>
          <h1 className="text-2xl md:text-3xl font-serif font-bold text-gray-900">{course.title}</h1>
          <p className="text-gray-600 mt-2 text-sm md:text-base max-w-3xl">{course.description}</p>
          <div className="flex items-center gap-4 mt-3 text-sm text-gray-500">
            <span>{modules.length} модулей</span>
            <span>•</span>
            <span>{totalLessons} уроков</span>
            {course.price > 0 && !course.is_public && (
              <>
                <span>•</span>
                <span className="font-semibold text-gray-900">{course.price} ₽</span>
              </>
            )}
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">

          {/* Левая колонка: Контент урока */}
          <div className="lg:col-span-2 space-y-6">
            {currentLesson ? (
              <>
                {/* Заголовок урока */}
                <div className="bg-white p-5 rounded-xl shadow-sm border border-gray-100">
                  <h2 className="text-xl font-semibold text-gray-800">{currentLesson.title}</h2>
                  {currentLesson.description && (
                    <p className="text-gray-600 text-sm mt-1">{currentLesson.description}</p>
                  )}
                </div>

                {/* Видеоплеер или статья */}
                {currentLesson.lesson_type === 'video' ? (
                  <VideoPlayer
                    url={`https://rutube.ru/play/embed/${currentLesson.video_embed_id}`}
                    onProgress={() => {}}
                  />
                ) : currentLesson.lesson_type === 'article' ? (
                  <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
                    <ArticleViewer content={currentLesson.article_content} />
                  </div>
                ) : (
                  <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 text-center text-gray-500 py-12">
                    <div className="text-4xl mb-3">📝</div>
                    <p>Тип урока не определён</p>
                  </div>
                )}

                {/* Кнопка "Отметить как пройденный" */}
                <div className="flex justify-end">
                  <button
                    onClick={() => handleLessonComplete(currentLesson.id)}
                    className="px-6 py-2 bg-rose-500 text-white rounded-lg hover:bg-rose-600 transition-colors text-sm font-medium"
                  >
                    ✓ Отметить как пройденный
                  </button>
                </div>
              </>
            ) : (
              <div className="bg-white p-12 rounded-xl shadow-sm border border-gray-100 text-center">
                <div className="text-5xl mb-4">📚</div>
                <h3 className="text-lg font-semibold text-gray-700 mb-2">Выберите урок</h3>
                <p className="text-gray-500 text-sm">Нажмите на урок в списке справа, чтобы начать обучение</p>
              </div>
            )}

            {/* Отзывы */}
            <div className="bg-white p-5 rounded-xl shadow-sm border border-gray-100">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">Отзывы</h3>
              <Comments courseId={courseId} reviews={reviews} onReviewSubmitted={handleReviewSubmitted} />
            </div>
          </div>

          {/* Правая колонка: Список модулей и уроков */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-xl shadow-sm border border-gray-100 sticky top-20 overflow-hidden">
              <div className="p-4 border-b border-gray-100 bg-gray-50">
                <h3 className="font-semibold text-gray-800">Содержание курса</h3>
              </div>
              <div className="max-h-[calc(100vh-12rem)] overflow-y-auto">
                {modules.map((mod, modIdx) => (
                  <div key={mod.id} className="border-b border-gray-100 last:border-b-0">
                    {/* Заголовок модуля */}
                    <button
                      onClick={() => toggleModule(modIdx)}
                      className="w-full text-left px-4 py-3 hover:bg-gray-50 transition-colors flex items-center justify-between"
                    >
                      <div>
                        <span className="text-xs text-rose-500 font-medium uppercase tracking-wider">
                          Модуль {modIdx + 1}
                        </span>
                        <p className="text-sm font-medium text-gray-800 mt-0.5">{mod.title}</p>
                        <p className="text-xs text-gray-500">{mod.lessons.length} уроков</p>
                      </div>
                      <svg
                        className={`w-5 h-5 text-gray-400 transition-transform ${
                          expandedModules.has(modIdx) ? 'rotate-180' : ''
                        }`}
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>

                    {/* Уроки модуля */}
                    {expandedModules.has(modIdx) && (
                      <div className="bg-gray-50/50">
                        {mod.lessons.map((lesson, lessonIdx) => {
                          const isSelected =
                            selectedLesson?.moduleIdx === modIdx && selectedLesson?.lessonIdx === lessonIdx;
                          return (
                            <button
                              key={lesson.id}
                              onClick={() => selectLesson(modIdx, lessonIdx)}
                              className={`w-full text-left px-4 py-2.5 pl-8 text-sm transition-colors flex items-start gap-2 ${
                                isSelected
                                  ? 'bg-rose-50 text-rose-700 border-l-2 border-rose-500'
                                  : 'hover:bg-gray-100 text-gray-700'
                              }`}
                            >
                              <span className="flex-shrink-0 mt-0.5">
                                {lesson.lesson_type === 'video' ? (
                                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                  </svg>
                                ) : (
                                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                  </svg>
                                )}
                              </span>
                              <span className="line-clamp-2">{lesson.title}</span>
                            </button>
                          );
                        })}
                      </div>
                    )}
                  </div>
                ))}

                {modules.length === 0 && (
                  <div className="p-6 text-center text-gray-500 text-sm">
                    Уроки пока не добавлены
                  </div>
                )}
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  );
};

export default CoursePage;
