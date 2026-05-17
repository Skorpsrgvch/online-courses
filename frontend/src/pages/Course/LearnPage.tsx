import React, { useState, useEffect, useMemo } from 'react';
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
  const [isSavingProgress, setIsSavingProgress] = useState(false);

  const [selectedLesson, setSelectedLesson] = useState<{
    moduleIdx: number;
    lessonIdx: number;
  } | null>(null);

  const [expandedModules, setExpandedModules] = useState<Set<number>>(new Set([0]));
  const [completedLessonsIds, setCompletedLessonsIds] = useState<Set<number>>(new Set());
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  useEffect(() => {
    const loadData = async () => {
      if (!courseId) return;
      setIsLoading(true);
      try {
        const data = await coursesService.getCourseFull(courseId);
        setCourseData(data);

        const completedIdsFromBackend = new Set<number>();
        data.modules.forEach(mod => {
          mod.lessons.forEach(lesson => {
            if (lesson.is_completed) {
              completedIdsFromBackend.add(lesson.id);
            }
          });
        });
        setCompletedLessonsIds(completedIdsFromBackend);

        // ЛОГИКА АВТОМАТИЧЕСКОГО ВЫБОРА УРОКА
        if (data.modules && data.modules.length > 0) {
          let nextModuleIdx = 0;
          let nextLessonIdx = 0;
          let foundUncompleted = false;

          // Ищем первый непройденный урок
          outerLoop:
          for (let m = 0; m < data.modules.length; m++) {
            const lessons = data.modules[m].lessons;
            if (!lessons || lessons.length === 0) continue;

            for (let l = 0; l < lessons.length; l++) {
              const lesson = lessons[l];
              // Если урок не помечен как пройденный бэкендом
              if (!completedIdsFromBackend.has(lesson.id)) {
                nextModuleIdx = m;
                nextLessonIdx = l;
                foundUncompleted = true;
                break outerLoop; // Прерываем оба цикла
              }
            }
          }

          // Если все уроки пройдены, остаемся на первом уроке (для повторения)
          // Если нашли непройденный - выбираем его
          setSelectedLesson({ moduleIdx: nextModuleIdx, lessonIdx: nextLessonIdx });

          // Раскрываем модуль, в котором находится выбранный урок
          setExpandedModules(prev => {
            const next = new Set(prev);
            next.add(nextModuleIdx);
            return next;
          });
        }
      } catch (err) {
        console.error('Ошибка загрузки:', err);
        alert('Не удалось загрузить данные курса');
      } finally {
        setIsLoading(false);
      }
    };
    loadData();
  }, [courseId]);

  const progressStats = useMemo(() => {
    if (!courseData || !courseData.modules) return { completed: 0, total: 0, percent: 0 };

    let total = 0;
    let completed = 0;

    courseData.modules.forEach(mod => {
      mod.lessons.forEach(lesson => {
        total++;
        if (completedLessonsIds.has(lesson.id)) {
          completed++;
        }
      });
    });

    const percent = total === 0 ? 0 : Math.round((completed / total) * 100);
    return { completed, total, percent };
  }, [courseData, completedLessonsIds]);

  const getRuTubeUrl = (lesson: FullCourseLesson) => {
    if (!lesson.video_embed_id) return '';

    // ВНИМАНИЕ: Убран пробел после embed/
    const base = `https://rutube.ru/video/embed/${lesson.video_embed_id}/`;

    if (lesson.private_key) {
      return `${base}?p=${lesson.private_key}&quality=auto`;
    }

    return base;
  };

  const getProgressColor = (percent: number) => {
    if (percent >= 80) return 'text-green-600 bg-green-500';
    if (percent >= 40) return 'text-yellow-600 bg-yellow-500';
    return 'text-red-600 bg-red-500';
  };

  const declension = (number: number, titles: [string, string, string]) => {
    const cases = [2, 0, 1, 1, 1, 2];
    return titles[(number % 100 > 4 && number % 100 < 20) ? 2 : cases[(number % 10 < 5) ? number % 10 : 5]];
  };

  const markAsComplete = async () => {
    if (!currentLesson || !currentLesson.id) return;
    if (completedLessonsIds.has(currentLesson.id)) return;

    setIsSavingProgress(true);
    try {
      await coursesService.markLessonComplete(currentLesson.id);
      const newCompleted = new Set(completedLessonsIds);
      newCompleted.add(currentLesson.id);
      setCompletedLessonsIds(newCompleted);
    } catch (error) {
      console.error('Failed to save progress:', error);
      alert('Не удалось сохранить прогресс.');
    } finally {
      setIsSavingProgress(false);
    }
  };

  const isFirstLessonOfCourse = selectedLesson
    ? (selectedLesson.moduleIdx === 0 && selectedLesson.lessonIdx === 0)
    : true;

  const isLastLessonOfCourse = selectedLesson && courseData?.modules
    ? (selectedLesson.moduleIdx === courseData.modules.length - 1 &&
      selectedLesson.lessonIdx === courseData.modules[courseData.modules.length - 1].lessons.length - 1)
    : false;

  const handleNextAction = () => {
    if (!selectedLesson || !courseData?.modules) return;
    const { moduleIdx, lessonIdx } = selectedLesson;
    const currentModule = courseData.modules[moduleIdx];

    if (lessonIdx < currentModule.lessons.length - 1) {
      selectLesson(moduleIdx, lessonIdx + 1);
      return;
    }

    if (moduleIdx < courseData.modules.length - 1) {
      const nextModule = courseData.modules[moduleIdx + 1];
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

    if (window.confirm('Поздравляем! Вы прошли весь курс! Вернуться на страницу курса?')) {
      navigate(`/course/${courseId}`);
    }
  };



  const handlePrevAction = () => {
    if (!selectedLesson || !courseData?.modules) return;
    const { moduleIdx, lessonIdx } = selectedLesson;

    if (lessonIdx > 0) {
      selectLesson(moduleIdx, lessonIdx - 1);
      return;
    }

    if (moduleIdx > 0) {
      const prevModule = courseData.modules[moduleIdx - 1];
      if (prevModule.lessons && prevModule.lessons.length > 0) {
        const lastLessonIdx = prevModule.lessons.length - 1;
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
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  if (isLoading || !courseData) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-rose-500 mx-auto mb-4"></div>
          <p className="text-gray-500">Загрузка материалов курса...</p>
        </div>
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

  const hasVideo = !!currentLesson?.video_embed_id;
  const videoUrl = hasVideo ? getRuTubeUrl(currentLesson) : '';

  const progressColorClass = getProgressColor(progressStats.percent);

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">

      {/* === ВЕРХНЯЯ ПАНЕЛЬ === */}
      <div className="bg-white border-b border-gray-200 sticky top-0 z-40 shadow-sm pt-4 pb-4">
        <div className="max-w-[1600px] mx-auto px-3 sm:px-6 lg:px-8 flex items-center justify-between gap-4">

          {/* Левая часть */}
          <div className="flex items-center gap-2 flex-shrink-0 min-w-[40px]">
            <button
              onClick={() => setIsMobileMenuOpen(true)}
              className="lg:hidden p-2 -ml-2 text-gray-500 hover:text-rose-600 hover:bg-rose-50 rounded-2xl transition-colors"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
            <button
              onClick={() => navigate(`/course/${courseId}`)}
              className="p-2 -ml-2 text-gray-500 hover:text-rose-600 font-medium flex items-center gap-1 transition-colors rounded-2xl hover:bg-gray-50"
              title="Вернуться к описанию курса"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
              <span className="hidden sm:inline text-sm">Курс</span>
            </button>
          </div>

          {/* Центральная часть */}
          <div className="flex-1 min-w-0 flex flex-col items-center justify-center max-w-xl px-2">
            <h3 className="text-xs sm:text-sm font-bold text-gray-900 text-center w-full mb-2 leading-tight break-words line-clamp-2">
              {course.title}
            </h3>

            {/* Индикатор прогресса с отступом */}
            <div className="w-full flex items-center gap-2 text-xs text-gray-500 mt-1">
              <span className={`whitespace-nowrap font-bold ${progressColorClass.split(' ')[0]}`}>
                {progressStats.percent}%
              </span>
              <div className="flex-1 h-2 bg-gray-100 rounded-2xl overflow-hidden">
                <div
                  className={`h-full transition-all duration-500 ease-out rounded-2xl ${progressColorClass.split(' ')[1]}`}
                  style={{ width: `${progressStats.percent}%` }}
                />
              </div>
              <span className="whitespace-nowrap hidden sm:inline">
                {progressStats.completed}/{progressStats.total} ур.
              </span>
            </div>
          </div>

          <div className="w-10 flex-shrink-0"></div>
        </div>
      </div>

      <div className="flex-grow max-w-[1600px] mx-auto w-full px-4 sm:px-6 lg:px-8 py-6 lg:py-8">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 lg:gap-8 h-full">

          {/* === САЙДБАР (Desktop) === */}
          <div className="hidden lg:block lg:col-span-4 xl:col-span-3">
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden sticky top-24 max-h-[calc(100vh-8rem)] flex flex-col">
              <div className="p-5 border-b border-gray-100 bg-gray-50">
                <h2 className="font-bold text-gray-800 text-lg">Программа курса</h2>
                <div className="flex items-center justify-between mt-1">
                  <p className="text-xs text-gray-500">
                    {modules?.length || 0} {declension(modules?.length || 0, ['модуль', 'модуля', 'модулей'])}
                  </p>
                  {progressStats.percent === 100 && (
                    <span className="text-xs font-bold text-green-600 bg-green-100 px-2 py-1 rounded-2xl flex items-center gap-1">
                      <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" /></svg>
                      Пройден
                    </span>
                  )}
                </div>
              </div>

              <div className="overflow-y-auto flex-grow custom-scrollbar">
                {modules && modules.length > 0 ? (
                  modules.map((mod, modIdx) => {
                    const modKey = mod.id ?? `mod-${modIdx}`;
                    const lessonCount = mod.lessons ? mod.lessons.length : 0;

                    let modCompleted = 0;
                    mod.lessons?.forEach(l => {
                      if (completedLessonsIds.has(l.id)) modCompleted++;
                    });
                    const isModFullyCompleted = lessonCount > 0 && modCompleted === lessonCount;

                    return (
                      <div key={modKey} className="border-b border-gray-100 last:border-0">
                        <button
                          onClick={() => toggleModule(modIdx)}
                          className={`w-full text-left px-5 py-4 transition-colors flex items-center justify-between group ${isModFullyCompleted ? 'bg-green-50 hover:bg-green-100' : 'hover:bg-gray-50'}`}
                        >
                          <div>
                            <div className="flex items-center gap-2 mb-1">
                              <span className={`text-xs font-bold uppercase tracking-wider ${isModFullyCompleted ? 'text-green-600' : 'text-rose-500'}`}>
                                Модуль {modIdx + 1}
                              </span>
                              {isModFullyCompleted && (
                                <svg className="w-3 h-3 text-green-600" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" /></svg>
                              )}
                            </div>
                            <p className={`font-semibold text-sm line-clamp-2 ${isModFullyCompleted ? 'text-green-800' : 'text-gray-800 group-hover:text-rose-600'}`}>
                              {mod.title || 'Без названия'}
                            </p>
                            <p className="text-xs text-gray-500 mt-1">
                              {lessonCount} {declension(lessonCount, ['урок', 'урока', 'уроков'])}
                            </p>
                          </div>
                          <svg className={`w-4 h-4 text-gray-400 transition-transform ${expandedModules.has(modIdx) ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
                        </button>

                        {expandedModules.has(modIdx) && mod.lessons && mod.lessons.length > 0 && (
                          <div className={`${isModFullyCompleted ? 'bg-green-50/30' : 'bg-gray-50/50'}`}>
                            {mod.lessons.map((lesson, lessonIdx) => {
                              const lesKey = lesson.id ?? `les-${modIdx}-${lessonIdx}`;
                              const isSelected = selectedLesson?.moduleIdx === modIdx && selectedLesson?.lessonIdx === lessonIdx;
                              const isCompleted = completedLessonsIds.has(lesson.id);

                              return (
                                <button
                                  key={lesKey}
                                  onClick={() => selectLesson(modIdx, lessonIdx)}
                                  className={`w-full text-left px-5 py-3 pl-9 text-sm transition-all flex items-start gap-3 border-l-4 ${isSelected
                                    ? 'bg-white border-rose-500 text-rose-700 shadow-sm'
                                    : isCompleted
                                      ? 'bg-green-100 border-green-500 text-green-800'
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
                  <div className="p-8 text-center text-gray-500 text-sm">Нет доступных уроков</div>
                )}
              </div>
            </div>
          </div>

          {/* === МОБИЛЬНАЯ ШТОРКА === */}
          {isMobileMenuOpen && (
            <>
              <div className="fixed inset-0 bg-black/50 z-40 lg:hidden backdrop-blur-sm" onClick={() => setIsMobileMenuOpen(false)} />
              <div className="fixed inset-y-0 left-0 w-3/4 max-w-xs bg-white z-50 shadow-2xl lg:hidden flex flex-col pt-16">
                <div className="p-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
                  <h2 className="font-bold text-gray-800 text-sm truncate pr-2">Программа курса</h2>
                  <button onClick={() => setIsMobileMenuOpen(false)} className="p-2 text-gray-500 hover:bg-gray-100 rounded-2xl flex-shrink-0">
                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
                  </button>
                </div>
                <div className="overflow-y-auto flex-grow pb-20">
                  {modules && modules.map((mod, modIdx) => {
                    const modKey = mod.id ?? `mod-m-${modIdx}`;

                    let modCompletedCount = 0;
                    mod.lessons?.forEach(l => {
                      if (completedLessonsIds.has(l.id)) modCompletedCount++;
                    });
                    const isModFullyCompletedMobile = mod.lessons && mod.lessons.length > 0 && modCompletedCount === mod.lessons.length;

                    return (
                      <div key={modKey} className="border-b border-gray-100">
                        <button
                          onClick={() => toggleModule(modIdx)}
                          className={`w-full text-left px-4 py-3 font-semibold text-sm flex justify-between items-center ${isModFullyCompletedMobile ? 'bg-green-50 text-green-800' : 'bg-gray-50 text-gray-800'}`}
                        >
                          <span className="truncate pr-2">Модуль {modIdx + 1}: {mod.title || ''}</span>
                          {isModFullyCompletedMobile && (
                            <svg className="w-4 h-4 text-green-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" /></svg>
                          )}
                        </button>
                        {expandedModules.has(modIdx) && mod.lessons && (
                          <div className="bg-white">
                            {mod.lessons.map((lesson, lIdx) => {
                              const lesKey = lesson.id ?? `les-m-${modIdx}-${lIdx}`;
                              const isSelected = selectedLesson?.moduleIdx === modIdx && selectedLesson?.lessonIdx === lIdx;
                              const isCompleted = completedLessonsIds.has(lesson.id);
                              return (
                                <button
                                  key={lesKey}
                                  onClick={() => selectLesson(modIdx, lIdx)}
                                  className={`block w-full text-left px-6 py-3 text-sm border-l-4 ${isSelected ? 'border-rose-500 bg-rose-50 text-rose-700 font-bold' :
                                    isCompleted ? 'border-green-400 bg-green-50 text-green-800' : 'border-transparent text-gray-600'
                                    }`}
                                >
                                  <div className="flex items-center gap-2">
                                    {isCompleted && <span className="text-green-600">✓</span>}
                                    <span className="truncate">{lesson.title}</span>
                                  </div>
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

          {/* === КОНТЕНТ === */}
          <div className="lg:col-span-8 xl:col-span-9">
            {currentLesson ? (
              <div className="space-y-6 animate-fade-in">

                <div className="bg-black rounded-2xl overflow-hidden shadow-xl relative w-full">
                  {hasVideo ? (
                    <VideoPlayer
                      url={videoUrl}
                      onProgress={(percent) => {
                        // Логика сохранения прогресса при просмотре
                        console.log(`Просмотрено: ${percent}%`);
                      }}
                    />
                  ) : (
                    <div className="aspect-video flex items-center justify-center text-white bg-gray-900">
                      <div className="text-center p-4">
                        <div className="text-6xl mb-4 opacity-50">🎥</div>
                        <p className="text-lg font-medium">Видео недоступно</p>
                        <p className="text-sm opacity-70 mt-2">Или урок содержит только текст</p>
                      </div>
                    </div>
                  )}
                </div>

                {/* Блок управления */}
                <div className="bg-white rounded-2xl p-6 sm:p-8 shadow-sm border border-gray-100">
                  <div className="flex flex-col md:flex-row md:items-start justify-between gap-4 mb-6 pb-6 border-b border-gray-100">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-2 flex-wrap">
                        <span className="px-2 py-1 bg-rose-100 text-rose-700 text-xs font-bold rounded-2xl uppercase">
                          Урок {selectedLesson ? selectedLesson.lessonIdx + 1 : 0}
                        </span>
                        {completedLessonsIds.has(currentLesson.id) && (
                          <span className="px-2 py-1 bg-green-100 text-green-700 text-xs font-bold rounded-2xl! flex items-center gap-1">
                            <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" /></svg>
                            Пройден
                          </span>
                        )}
                      </div>
                      <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 font-serif mb-2 break-words">
                        {currentLesson.title}
                      </h2>
                      <p className="text-sm text-gray-500">
                        Модуль: {currentModule?.title}
                      </p>
                    </div>

                    <button
                      onClick={markAsComplete}
                      disabled={completedLessonsIds.has(currentLesson.id) || isSavingProgress}
                      className={`px-6 py-3 rounded-2xl! font-medium transition-all shadow-md flex items-center gap-2 whitespace-nowrap self-start md:self-center ${completedLessonsIds.has(currentLesson.id)
                        ? 'bg-green-100 text-green-700 cursor-default border border-green-200'
                        : 'bg-green-600 hover:bg-green-700 text-white hover:shadow-lg hover:-translate-y-0.5'
                        } ${isSavingProgress ? 'opacity-70 cursor-wait' : ''}`}
                    >
                      {isSavingProgress ? (
                        <>
                          <svg className="animate-spin h-5 w-5 text-current" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                          </svg>
                          Сохранение...
                        </>
                      ) : completedLessonsIds.has(currentLesson.id) ? (
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
                      <h4 className="text-lg font-bold text-gray-800 mb-3 not-prose flex items-center gap-2">
                        <svg className="w-5 h-5 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
                        Описание урока
                      </h4>
                      <div className="bg-gray-50 p-4 rounded-2xl border border-gray-100">
                        <p className="whitespace-pre-line m-0">{currentLesson.description}</p>
                      </div>
                    </div>
                  )}

                  {/* Навигация */}
                  <div className="pt-6 border-t border-gray-100 flex justify-between items-center gap-4">
                    {!isFirstLessonOfCourse ? (
                      <button
                        onClick={handlePrevAction}
                        className="px-5 py-2.5 bg-white border border-gray-300 text-gray-700 hover:bg-gray-50 hover:border-gray-400 rounded-2xl! font-medium transition-all shadow-sm hover:shadow flex items-center gap-2 group"
                      >
                        <svg className="w-5 h-5 text-gray-400 group-hover:text-gray-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 17l-5-5m0 0l5-5m-5 5h12" />
                        </svg>
                        <span className="hidden sm:inline">Назад</span>
                      </button>
                    ) : (
                      <div></div>
                    )}

                    {!isLastLessonOfCourse ? (
                      <button
                        onClick={handleNextAction}
                        className="px-6 py-2.5 bg-rose-600 hover:bg-rose-700 text-white rounded-2xl! font-medium transition-colors shadow-md hover:shadow-lg hover:-translate-y-0.5 flex items-center gap-2"
                      >
                        <span className="hidden sm:inline">Следующий урок</span>
                        <span className="sm:hidden">Далее</span>
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                        </svg>
                      </button>
                    ) : (
                      <button
                        onClick={() => navigate(`/course/${courseId}`)}
                        className="px-6 py-2.5 bg-green-600 hover:bg-green-700 text-white rounded-2xl font-medium transition-colors shadow-md hover:shadow-lg hover:-translate-y-0.5 flex items-center gap-2"
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
                  Нажмите на любой урок в списке слева (или в меню на мобильных), чтобы начать просмотр.
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