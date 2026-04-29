import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { coursesService } from '../../api/courses.service';
import type { CourseFullResponse, FullCourseModule, FullCourseLesson } from '../../api/types';
import { VideoPlayer } from '../../components/course/VideoPlayer';

const LearnPage = () => {
  const { id } = useParams<{ id: string }>();
  const courseId = Number(id);
  const navigate = useNavigate();

  const [courseData, setCourseData] = useState<CourseFullResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Состояние для выбранного урока
  const [selectedLesson, setSelectedLesson] = useState<{
    moduleIdx: number;
    lessonIdx: number;
  } | null>(null);

  // Состояние для разворачивания модулей (по умолчанию первый открыт)
  const [expandedModules, setExpandedModules] = useState<Set<number>>(new Set([0]));

  // Состояние для пройденных уроков (хранит ID уроков)
  const [completedLessons, setCompletedLessons] = useState<Set<number>>(new Set());

  // Состояние для мобильного меню (шторка)
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  // Загрузка данных и прогресса
  useEffect(() => {
    const loadData = async () => {
      if (!courseId) return;
      try {
        const data = await coursesService.getCourseFull(courseId);
        setCourseData(data);

        // Загрузка прогресса из localStorage
        const savedProgress = localStorage.getItem(`course_${courseId}_progress`);
        if (savedProgress) {
          const parsedIds = JSON.parse(savedProgress) as number[];
          setCompletedLessons(new Set(parsedIds));
        }

        // Автовыбор первого урока
        if (data.modules && data.modules.length > 0) {
          const firstModuleLessons = data.modules[0].lessons;
          if (firstModuleLessons && firstModuleLessons.length > 0) {
            setSelectedLesson({ moduleIdx: 0, lessonIdx: 0 });
          }
        }
      } catch (err) {
        console.error('Ошибка загрузки:', err);
      } finally {
        setIsLoading(false);
      }
    };
    loadData();
  }, [courseId]);

  // Отметка урока как пройденного
  const markAsComplete = () => {
    if (!currentLesson) return;

    const newCompleted = new Set(completedLessons);
    newCompleted.add(currentLesson.id);
    setCompletedLessons(newCompleted);

    // Сохраняем в localStorage
    if (courseData) {
      localStorage.setItem(`course_${courseData.course.id}_progress`, JSON.stringify(Array.from(newCompleted)));
    }
  };

  const isFirstLessonOfCourse = selectedLesson
    ? (selectedLesson.moduleIdx === 0 && selectedLesson.lessonIdx === 0)
    : true;

  // Логика перехода к следующему элементу
  const handleNextAction = () => {
    if (!selectedLesson || !modules) return;

    const { moduleIdx, lessonIdx } = selectedLesson;
    const currentModule = modules[moduleIdx];

    // 1. Если есть следующий урок в ЭТОМ модуле
    if (lessonIdx < currentModule.lessons.length - 1) {
      selectLesson(moduleIdx, lessonIdx + 1);
      return;
    }

    // 2. Если урок последний в модуле, ищем следующий модуль
    if (moduleIdx < modules.length - 1) {
      const nextModule = modules[moduleIdx + 1];
      if (nextModule.lessons && nextModule.lessons.length > 0) {
        setExpandedModules(prev => {
          const next = new Set(prev);
          next.add(moduleIdx + 1);
          return next;
        });
        selectLesson(moduleIdx + 1, 0);
        return;
      }
    }

    // 3. Конец курса
    alert('Поздравляем! Вы прошли весь курс!');
    navigate(`/course/${courseId}`);
  };

  // Логика перехода к предыдущему элементу
  const handlePrevAction = () => {
    if (!selectedLesson || !modules) return;

    const { moduleIdx, lessonIdx } = selectedLesson;

    // 1. Если есть предыдущий урок в ЭТОМ модуле
    if (lessonIdx > 0) {
      selectLesson(moduleIdx, lessonIdx - 1);
      return;
    }

    // 2. Если это первый урок модуля (но не самого первого), идем в конец предыдущего модуля
    if (moduleIdx > 0) {
      const prevModule = modules[moduleIdx - 1];
      if (prevModule.lessons && prevModule.lessons.length > 0) {
        const lastLessonIdx = prevModule.lessons.length - 1;

        // Раскрываем предыдущий модуль
        setExpandedModules(prev => {
          const next = new Set(prev);
          next.add(moduleIdx - 1);
          return next;
        });

        selectLesson(moduleIdx - 1, lastLessonIdx);
        return;
      }
    }
  };

  const toggleModule = (idx: number) => {
    setExpandedModules((prev) => {
      const next = new Set(prev);
      if (next.has(idx)) next.delete(idx);
      else next.add(idx);
      return next;
    });
  };

  const selectLesson = (modIdx: number, lessonIdx: number) => {
    setSelectedLesson({ moduleIdx: modIdx, lessonIdx: lessonIdx });
    setExpandedModules((prev) => {
      const next = new Set(prev);
      next.add(modIdx);
      return next;
    });
    setIsMobileMenuOpen(false);
  };

  if (isLoading || !courseData) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-rose-500"></div>
      </div>
    );
  }

  const { course, modules } = courseData;

  const currentModule = selectedLesson && modules[selectedLesson.moduleIdx]
    ? modules[selectedLesson.moduleIdx]
    : null;

  const currentLesson = selectedLesson && currentModule && currentModule.lessons[selectedLesson.lessonIdx]
    ? currentModule.lessons[selectedLesson.lessonIdx]
    : null;

  // Проверка: является ли текущий урок последним во всем курсе
  const isLastLessonOfCourse = selectedLesson && modules
    ? (selectedLesson.moduleIdx === modules.length - 1 &&
      selectedLesson.lessonIdx === modules[modules.length - 1].lessons.length - 1)
    : false;

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">

      {/* === ВЕРХНЯЯ ПАНЕЛЬ === */}
      <div className="bg-white border-b border-gray-200 sticky top-0 z-40 shadow-sm">
        <div className="max-w-[1600px] mx-auto px-3 sm:px-6 lg:px-8 h-16 flex items-center justify-between gap-3">
          <div className="flex items-center gap-2 flex-shrink-0 min-w-[80px]">
            <button
              onClick={() => setIsMobileMenuOpen(true)}
              className="lg:hidden p-2 -ml-2 text-gray-500 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-colors"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
            <button
              onClick={() => navigate(`/course/${courseId}`)}
              className="p-2 -ml-2 text-gray-500 hover:text-rose-600 font-medium flex items-center gap-1 transition-colors rounded-lg hover:bg-gray-50"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
              <span className="hidden sm:inline text-sm">Курс</span>
            </button>
          </div>

          <h3 className="text-xs sm:text-sm md:text-base font-bold text-gray-900 text-center break-words leading-snug min-w-0 flex-1 px-2">
            {course.title}
          </h3>

          <div className="w-10 flex-shrink-0"></div>
        </div>
      </div>

      <div className="flex-grow max-w-[1600px] mx-auto w-full px-4 sm:px-6 lg:px-8 py-6 lg:py-8">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 lg:gap-8 h-full">

          {/* === САЙДБАР (Программа курса) === */}
          <div className="hidden lg:block lg:col-span-4 xl:col-span-3">
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden sticky top-24 max-h-[calc(100vh-8rem)] flex flex-col">
              <div className="p-5 border-b border-gray-100 bg-gray-50">
                <h2 className="font-bold text-gray-800 text-lg">Программа курса</h2>
                <p className="text-xs text-gray-500 mt-1">{modules?.length || 0} модулей</p>
              </div>

              <div className="overflow-y-auto flex-grow custom-scrollbar">
                {modules && modules.length > 0 ? (
                  modules.map((mod, modIdx) => {
                    const modKey = mod.id ?? `mod-${modIdx}`;
                    const lessonCount = mod.lessons ? mod.lessons.length : 0;

                    return (
                      <div key={modKey} className="border-b border-gray-100 last:border-0">
                        <button
                          onClick={() => toggleModule(modIdx)}
                          className="w-full text-left px-5 py-4 hover:bg-gray-50 transition-colors flex items-center justify-between group"
                        >
                          <div>
                            <span className="text-xs font-bold text-rose-500 uppercase tracking-wider block mb-1">
                              Модуль {modIdx + 1}
                            </span>
                            <p className="font-semibold text-gray-800 text-sm group-hover:text-rose-600 line-clamp-2">
                              {mod.title || 'Без названия'}
                            </p>
                            <p className="text-xs text-gray-500 mt-1">{lessonCount} уроков</p>
                          </div>
                          <svg className={`w-4 h-4 text-gray-400 transition-transform ${expandedModules.has(modIdx) ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
                        </button>

                        {expandedModules.has(modIdx) && mod.lessons && mod.lessons.length > 0 && (
                          <div className="bg-gray-50/50">
                            {mod.lessons.map((lesson, lessonIdx) => {
                              const lesKey = lesson.id ?? `les-${modIdx}-${lessonIdx}`;
                              const isSelected = selectedLesson?.moduleIdx === modIdx && selectedLesson?.lessonIdx === lessonIdx;
                              const isCompleted = completedLessons.has(lesson.id);

                              return (
                                <button
                                  key={lesKey}
                                  onClick={() => selectLesson(modIdx, lessonIdx)}
                                  className={`w-full text-left px-5 py-3 pl-9 text-sm transition-all flex items-start gap-3 border-l-4 ${isSelected
                                      ? 'bg-white border-rose-500 text-rose-700 shadow-sm'
                                      : isCompleted
                                        ? 'bg-green-50 border-green-400 text-green-800' // Зеленый стиль для пройденных
                                        : 'border-transparent text-gray-600 hover:bg-gray-100'
                                    }`}
                                >
                                  <span className={`mt-0.5 flex-shrink-0 ${isCompleted ? 'text-green-600' : (isSelected ? 'text-rose-500' : 'text-gray-400')
                                    }`}>
                                    {isCompleted ? (
                                      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" /></svg>
                                    ) : (
                                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                                    )}
                                  </span>
                                  <span className="line-clamp-2 font-medium">{lesson.title}</span>
                                </button>
                              );
                            })}
                          </div>
                        )}
                      </div>
                    );
                  })
                ) : (
                  <div className="p-8 text-center text-gray-500 text-sm">Программа загружается...</div>
                )}
              </div>
            </div>
          </div>

          {/* === МОБИЛЬНАЯ ШТОРКА === */}
          {isMobileMenuOpen && (
            <>
              <div className="fixed inset-0 bg-black/50 z-40 lg:hidden backdrop-blur-sm" onClick={() => setIsMobileMenuOpen(false)} />
              <div className="fixed inset-y-0 left-0 w-3/4 max-w-xs bg-white z-50 shadow-2xl lg:hidden flex flex-col pt-20">
                <div className="p-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
                  <h2 className="font-bold text-gray-800">Программа курса</h2>
                  <button onClick={() => setIsMobileMenuOpen(false)} className="p-2 text-gray-500 hover:bg-gray-100 rounded-full">
                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
                  </button>
                </div>
                <div className="overflow-y-auto flex-grow">
                  {modules && modules.map((mod, modIdx) => {
                    const modKey = mod.id ?? `mod-m-${modIdx}`;
                    return (
                      <div key={modKey} className="border-b border-gray-100">
                        <button onClick={() => toggleModule(modIdx)} className="w-full text-left px-4 py-3 bg-gray-50 font-semibold text-sm text-gray-800">
                          Модуль {modIdx + 1}: {mod.title || ''}
                        </button>
                        {expandedModules.has(modIdx) && mod.lessons && (
                          <div className="bg-white">
                            {mod.lessons.map((lesson, lIdx) => {
                              const lesKey = lesson.id ?? `les-m-${modIdx}-${lIdx}`;
                              const isSelected = selectedLesson?.moduleIdx === modIdx && selectedLesson?.lessonIdx === lIdx;
                              const isCompleted = completedLessons.has(lesson.id);
                              return (
                                <button
                                  key={lesKey}
                                  onClick={() => selectLesson(modIdx, lIdx)}
                                  className={`block w-full text-left px-6 py-2 text-sm ${isSelected ? 'text-rose-600 font-bold bg-rose-50' :
                                      isCompleted ? 'text-green-700 bg-green-50' : 'text-gray-600'
                                    }`}
                                >
                                  {lesson.title} {isCompleted && '✓'}
                                </button>
                              );
                            })}
                          </div>
                        )}
                      </div>
                    );
                  })}
                </div>
              </div>
            </>
          )}

          {/* === КОНТЕНТ (Видео + Описание) === */}
          <div className="lg:col-span-8 xl:col-span-9">
            {currentLesson ? (
              <div className="space-y-6 animate-fade-in">

                {/* Видео плеер */}
                <div className="bg-black rounded-2xl overflow-hidden shadow-xl aspect-video relative group">
                  {currentLesson.video_embed_id ? (
                    <VideoPlayer
                      url={
                        currentLesson.private_key
                          ? `https://rutube.ru/video/embed/${currentLesson.video_embed_id}/?p=${currentLesson.private_key}`
                          : `https://rutube.ru/video/embed/${currentLesson.video_embed_id}/`
                      }
                      onProgress={() => { }}
                    />
                  ) : (
                    <div className="absolute inset-0 flex items-center justify-center text-white">
                      <div className="text-center">
                        <div className="text-6xl mb-4 opacity-50">🎥</div>
                        <p>Видео недоступно</p>
                      </div>
                    </div>
                  )}
                </div>

                {/* Блок управления и Описание */}
                <div className="bg-white rounded-2xl p-6 sm:p-8 shadow-sm border border-gray-100">
                  <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6 pb-6 border-b border-gray-100">
                    <div>
                      <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 font-serif mb-2">
                        {currentLesson.title}
                      </h2>
                      <p className="text-sm text-gray-500">
                        Урок {selectedLesson ? selectedLesson.lessonIdx + 1 : 0} из {currentModule?.lessons.length || 0} в модуле {selectedLesson ? selectedLesson.moduleIdx + 1 : 0}
                      </p>
                    </div>

                    <button
                      onClick={markAsComplete}
                      disabled={completedLessons.has(currentLesson.id)}
                      className={`px-6 py-3 rounded-xl font-medium transition-all shadow-md flex items-center gap-2 whitespace-nowrap ${completedLessons.has(currentLesson.id)
                          ? 'bg-green-100 text-green-700 cursor-default'
                          : 'bg-green-600 hover:bg-green-700 text-white hover:shadow-lg'
                        }`}
                    >
                      {completedLessons.has(currentLesson.id) ? (
                        <>
                          <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" /></svg>
                          Пройдено
                        </>
                      ) : (
                        <>
                          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>
                          Отметить как пройденное
                        </>
                      )}
                    </button>
                  </div>

                  {currentLesson.description && (
                    <div className="prose prose-rose max-w-none text-gray-600 leading-relaxed mb-8">
                      <h4 className="text-lg font-bold text-gray-800 mb-2 not-prose">Описание урока:</h4>
                      <p className="whitespace-pre-line">{currentLesson.description}</p>
                    </div>
                  )}

                  {/* Навигация */}
                  <div className="pt-6 border-t border-gray-100 flex justify-between items-center">
                    
                    {/* Кнопка НАЗАД (Слева) */}
                    {!isFirstLessonOfCourse && (
                      <button
                        onClick={handlePrevAction}
                        className="px-6 py-3 bg-white border border-gray-300 text-gray-700 hover:bg-gray-50 hover:border-gray-400 rounded-xl font-medium transition-all shadow-sm hover:shadow flex items-center gap-2"
                      >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 17l-5-5m0 0l5-5m-5 5h12" />
                        </svg>
                        Назад
                      </button>
                    )}
                    
                    {/* Пустой блок для балансировки, если кнопки "Назад" нет */}
                    {isFirstLessonOfCourse && <div></div>}

                    {/* Кнопка ДАЛЕЕ / ЗАВЕРШИТЬ (Справа) */}
                    {!isLastLessonOfCourse ? (
                      <button
                        onClick={handleNextAction}
                        className="px-8 py-3 bg-rose-600 hover:bg-rose-700 text-white rounded-xl font-medium transition-colors shadow-md hover:shadow-lg flex items-center gap-2"
                      >
                        Далее
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                        </svg>
                      </button>
                    ) : (
                      <button
                        onClick={() => navigate(`/course/${courseId}`)}
                        className="px-8 py-3 bg-green-600 hover:bg-green-700 text-white rounded-xl font-medium transition-colors shadow-md hover:shadow-lg flex items-center gap-2"
                      >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        Завершить курс
                      </button>
                    )}
                  </div>
                </div>

              </div>
            ) : (
              <div className="h-full flex flex-col items-center justify-center text-center p-8 bg-white rounded-2xl border border-dashed border-gray-300 min-h-[400px]">
                <div className="text-6xl mb-4">👈</div>
                <h3 className="text-xl font-bold text-gray-800 mb-2">Выберите урок</h3>
                <p className="text-gray-500 max-w-md">
                  Нажмите на любой урок в списке слева, чтобы начать просмотр.
                </p>
              </div>
            )}
          </div>

        </div>
      </div>
    </div>
  );
};

export default LearnPage;