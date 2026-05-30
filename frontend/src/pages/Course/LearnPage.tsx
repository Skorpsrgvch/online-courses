import React, { useState, useEffect, useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { coursesService } from '../../api/courses.service';
import type { CourseFullResponse, FullCourseModule, FullCourseLesson } from '../../api/types';
import { VideoPlayer } from '../../components/course/VideoPlayer';
import { AccessExpirationAlert } from '../../components/course/AccessExpirationAlert';

const LearnPage = () => {
  const { id } = useParams<{ id: string }>();
  const courseId = Number(id);
  const navigate = useNavigate();

  const [courseData, setCourseData] = useState<CourseFullResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSavingProgress, setIsSavingProgress] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [selectedLesson, setSelectedLesson] = useState<{
    moduleIdx: number;
    lessonIdx: number;
  } | null>(null);

  const [expandedModules, setExpandedModules] = useState<Set<number>>(new Set([0]));
  const [completedLessonsIds, setCompletedLessonsIds] = useState<Set<number>>(new Set());

  // Состояние мобильного меню
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  // Состояния для модальных окон
  const [lockedModuleInfo, setLockedModuleInfo] = useState<{
    isOpen: boolean;
    unlockDate: Date | null;
    title: string;
  }>({ isOpen: false, unlockDate: null, title: '' });

  const [showCompletionModal, setShowCompletionModal] = useState(false);

  // Таймер для обновления обратного отсчета
  const [, forceUpdate] = useState(0);
  useEffect(() => {
    const timer = setInterval(() => forceUpdate(n => n + 1), 1000);
    return () => clearInterval(timer);
  }, []);

  // Хелпер для определения бонусного модуля
  const isModuleBonus = (module: FullCourseModule) => {
    return module.title.toLowerCase().includes('бонус');
  };

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

        if (data.modules && data.modules.length > 0) {
          let nextModuleIdx = 0;
          let nextLessonIdx = 0;
          let foundUncompleted = false;

          outerLoop:
          for (let m = 0; m < data.modules.length; m++) {
            const moduleItem = data.modules[m];

            // Пропускаем заблокированные (не бонусные) модули
            if (!isModuleBonus(moduleItem) && moduleItem.is_locked) continue;

            const lessons = moduleItem.lessons;
            if (!lessons || lessons.length === 0) continue;

            for (let l = 0; l < lessons.length; l++) {
              const lesson = lessons[l];
              if (!completedIdsFromBackend.has(lesson.id)) {
                nextModuleIdx = m;
                nextLessonIdx = l;
                foundUncompleted = true;
                break outerLoop;
              }
            }
          }

          if (!foundUncompleted && data.modules.length > 0) {
            const firstModuleIsBonus = isModuleBonus(data.modules[0]);
            if (!firstModuleIsBonus && data.modules[0].is_locked) {
              setSelectedLesson(null);
            } else {
              setSelectedLesson({ moduleIdx: nextModuleIdx, lessonIdx: nextLessonIdx });
            }
          } else {
            setSelectedLesson({ moduleIdx: nextModuleIdx, lessonIdx: nextLessonIdx });
          }

          setExpandedModules(prev => {
            const next = new Set(prev);
            const currentModule = data.modules[nextModuleIdx];
            if (currentModule && (!currentModule.is_locked || isModuleBonus(currentModule))) {
              next.add(nextModuleIdx);
            }
            return next;
          });
        }
      } catch (err) {
        console.error('Ошибка загрузки:', err);
        setError('Не удалось загрузить данные курса');
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
        if (completedLessonsIds.has(lesson.id)) completed++;
      });
    });
    const percent = total === 0 ? 0 : Math.round((completed / total) * 100);
    return { completed, total, percent };
  }, [courseData, completedLessonsIds]);

  const getRuTubeUrl = (lesson: FullCourseLesson) => {
    if (!lesson.video_embed_id) return '';
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
      setError('Не удалось сохранить прогресс.');
    } finally {
      setIsSavingProgress(false);
    }
  };

  const isFirstLessonOfCourse = selectedLesson ? (selectedLesson.moduleIdx === 0 && selectedLesson.lessonIdx === 0) : true;

  const isLastLessonOfCourse = selectedLesson && courseData?.modules
    ? (selectedLesson.moduleIdx === courseData.modules.length - 1 &&
      selectedLesson.lessonIdx === courseData.modules[courseData.modules.length - 1].lessons.length - 1)
    : false;


  const showLockInfo = (module: FullCourseModule) => {

    let dateObj: Date | null = null;

    if (module.unlock_date) {
      dateObj = new Date(module.unlock_date);
      if (isNaN(dateObj.getTime())) {

        dateObj = null;
      } else {
        dateObj.toISOString();
      }
    }

    setLockedModuleInfo({
      isOpen: true,
      unlockDate: dateObj,
      title: module.title
    });
  };
  const handleNextAction = () => {
    if (!selectedLesson || !courseData?.modules) {
      return;
    }

    const { moduleIdx, lessonIdx } = selectedLesson;
    const currentModule = courseData.modules[moduleIdx];

    // 1. Есть ли следующий урок в текущем модуле?
    if (lessonIdx < currentModule.lessons.length - 1) {
      selectLesson(moduleIdx, lessonIdx + 1);
      return;
    }

    // 2. Ищем следующий модуль
    let nextModuleIdx = moduleIdx + 1;
    let foundNextAvailableLesson = false;
    let firstLockedModuleAfterCurrent: FullCourseModule | null = null;

    while (nextModuleIdx < courseData.modules.length) {
      const nextMod = courseData.modules[nextModuleIdx];
      const isBonus = isModuleBonus(nextMod);


      if (!nextMod.is_locked || isBonus) {
        if (nextMod.lessons && nextMod.lessons.length > 0) {
          foundNextAvailableLesson = true;

          setExpandedModules(prev => {
            const next = new Set(prev);
            next.add(nextModuleIdx);
            return next;
          });

          selectLesson(nextModuleIdx, 0);
          return;
        }
      } else {
        firstLockedModuleAfterCurrent = nextMod;
        break;
      }

      nextModuleIdx++;
    }

    // 3. Результаты поиска
    if (!foundNextAvailableLesson) {
      if (firstLockedModuleAfterCurrent) {
        showLockInfo(firstLockedModuleAfterCurrent);
      } else {

        setShowCompletionModal(true);
      }
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
      let prevModuleIdx = moduleIdx - 1;
      const prevModule = courseData.modules[prevModuleIdx];
      if (prevModule.lessons && prevModule.lessons.length > 0) {
        const lastLessonIdx = prevModule.lessons.length - 1;
        setExpandedModules(prev => {
          const next = new Set(prev);
          next.add(prevModuleIdx);
          return next;
        });
        selectLesson(prevModuleIdx, lastLessonIdx);
      }
    }
  };

  const handleModuleClick = (modIdx: number, e?: React.MouseEvent) => {
    if (e) e.stopPropagation();

    const module = courseData?.modules[modIdx];
    if (!module) return;

    if (module.is_locked && !isModuleBonus(module)) {
      showLockInfo(module);
      return;
    }

    setExpandedModules((prev) => {
      const next = new Set(prev);
      if (next.has(modIdx)) next.delete(modIdx);
      else next.add(modIdx);
      return next;
    });
  };

  const selectLesson = (modIdx: number, lessonIdx: number) => {
    const module = courseData?.modules[modIdx];
    if (!module) return;

    if (module.is_locked && !isModuleBonus(module)) {
      showLockInfo(module);
      return;
    }

    setSelectedLesson({ moduleIdx: modIdx, lessonIdx: lessonIdx });
    setExpandedModules((prev) => {
      const next = new Set(prev);
      next.add(modIdx);
      return next;
    });
    setIsMobileMenuOpen(false);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const getTimeRemaining = (unlockDate: Date | null) => {
    if (!unlockDate) return null;
    const total = Date.parse(unlockDate.toString()) - Date.parse(new Date().toString());
    if (total <= 0) return null;

    const seconds = Math.floor((total / 1000) % 60);
    const minutes = Math.floor((total / 1000 / 60) % 60);
    const hours = Math.floor((total / (1000 * 60 * 60)) % 24);
    const days = Math.floor(total / (1000 * 60 * 60 * 24));

    return { total, days, hours, minutes, seconds };
  };

  const formatDateTime = (date: Date) => {
    return date.toLocaleDateString('ru-RU', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (isLoading || !courseData) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-rose-500 mx-auto mb-4"></div>
          <p className="text-gray-500">{error || 'Загрузка материалов курса...'}</p>
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

  const isCurrentLessonLocked = currentModule ? (currentModule.is_locked && !isModuleBonus(currentModule)) : false;

  const hasVideo = !!currentLesson?.video_embed_id && !isCurrentLessonLocked;
  const videoUrl = hasVideo ? getRuTubeUrl(currentLesson) : '';

  const progressColorClass = getProgressColor(progressStats.percent);

  const isExpired = courseData.is_access_expired ?? course?.is_access_expired ?? false;
  const daysRemaining = courseData.days_remaining ?? course?.days_remaining ?? 365;
  const expiresAt = courseData.access_expires_at ?? course?.access_expires_at ?? undefined;

  // Проверка истечения доступа
  if (isExpired) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
        <div className="bg-white p-8 rounded-2xl shadow-lg max-w-md w-full text-center border-t-4 border-red-500">
          <div className="text-5xl mb-4">🔒</div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Доступ истек</h2>
          <p className="text-gray-600 mb-6">
            Срок вашего доступа к этому курсу завершился{' '}
            {expiresAt ? new Date(expiresAt).toLocaleDateString('ru-RU') : 'недавно'}.
          </p>

          <AccessExpirationAlert
            isExpired={isExpired}
            daysRemaining={daysRemaining}
            expiresAt={expiresAt}
            onRenew={() => navigate(`/course/${courseId}`)}
          />

          <button
            onClick={() => navigate(`/course/${courseId}`)}
            className="mt-6 w-full py-3 bg-rose-500 text-white font-bold rounded-2xl! hover:bg-rose-600 transition-colors"
          >
            Вернуться к описанию курса
          </button>
        </div>
      </div>
    );
  }



  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      {/* ВЕРХНЯЯ ПАНЕЛЬ */}
      <div className="bg-white border-b border-gray-200 z-40 shadow-sm pt-4 pb-4">
        <div className="max-w-[1600px] mx-auto px-3 sm:px-6 lg:px-8 flex items-center justify-between gap-4">
          <div className="flex items-center gap-2 flex-shrink-0 min-w-[40px]">
            {/* КНОПКА ОТКРЫТИЯ МЕНЮ */}
            <button
              onClick={() => setIsMobileMenuOpen(true)}
              className="lg:hidden p-2 -ml-2 text-gray-500 hover:text-rose-600 hover:bg-rose-50 rounded-2xl transition-colors"
              aria-label="Открыть меню"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
            <button onClick={() => navigate(`/course/${courseId}`)} className="p-2 -ml-2 text-gray-500 hover:text-rose-600 font-medium flex items-center gap-1 transition-colors rounded-2xl hover:bg-gray-50">
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" /></svg>
              <span className="hidden sm:inline text-sm">Курс</span>
            </button>
          </div>
          <div className="flex-1 min-w-0 flex flex-col items-center justify-center max-w-xl px-2">
            <h3 className="text-xs sm:text-sm font-bold text-gray-900 text-center w-full mb-2 leading-tight break-words line-clamp-2">{course.title}</h3>
            <div className="w-full flex items-center gap-2 text-xs text-gray-500 mt-1">
              <span className={`whitespace-nowrap font-bold ${progressColorClass.split(' ')[0]}`}>{progressStats.percent}%</span>
              <div className="flex-1 h-2 bg-gray-100 rounded-2xl overflow-hidden">
                <div className={`h-full transition-all duration-500 ease-out rounded-2xl ${progressColorClass.split(' ')[1]}`} style={{ width: `${progressStats.percent}%` }} />
              </div>
              <span className="whitespace-nowrap hidden sm:inline">{progressStats.completed}/{progressStats.total} ур.</span>
            </div>
          </div>
          <div className="w-10 flex-shrink-0"></div>
        </div>
      </div>

      <div className="flex-grow max-w-[1600px] mx-auto w-full px-4 sm:px-6 lg:px-8 py-6 lg:py-8">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 lg:gap-8 h-full">

          {/* САЙДБАР (Desktop) */}
          <div className="hidden lg:block lg:col-span-4 xl:col-span-3">
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden sticky top-24 max-h-[calc(100vh-8rem)] flex flex-col">
              {/* Заголовок программы курса */}
              <div className="p-5 border-b border-gray-100 bg-gray-50/50">
                <h2 className="font-bold text-gray-800 text-lg leading-tight break-words">
                  Программа курса
                </h2>
              </div>

              <div className="overflow-y-auto flex-grow custom-scrollbar">
                {modules && modules.length > 0 ? (
                  modules.map((mod, modIdx) => {
                    const isBonus = isModuleBonus(mod);
                    const isLocked = mod.is_locked && !isBonus;
                    const lessonCount = mod.lessons ? mod.lessons.length : 0;

                    let modCompleted = 0;
                    mod.lessons?.forEach(l => { if (completedLessonsIds.has(l.id)) modCompleted++; });
                    const isModFullyCompleted = lessonCount > 0 && modCompleted === lessonCount;

                    return (
                      <div key={mod.id} className="border-b border-gray-100 last:border-0">
                        <button
                          onClick={(e) => handleModuleClick(modIdx, e)}
                          disabled={isLocked}
                          className={`w-full text-left px-5 py-4 transition-colors flex items-center justify-between group ${isLocked
                            ? 'bg-gray-100 cursor-not-allowed opacity-70 grayscale'
                            : isModFullyCompleted
                              ? 'bg-green-50 hover:bg-green-100'
                              : 'hover:bg-gray-50'
                            }`}
                        >
                          <div className="relative">
                            <div className="flex items-center gap-2 mb-1">
                              <span className={`text-xs font-bold uppercase tracking-wider ${isLocked ? 'text-gray-400' : isModFullyCompleted ? 'text-green-600' : 'text-rose-500'
                                }`}>
                                {isLocked ? 'Закрыто' : `Модуль ${modIdx + 1}`}
                              </span>
                              {isLocked && (
                                <svg className="w-3 h-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                                </svg>
                              )}
                              {!isLocked && isModFullyCompleted && (
                                <svg className="w-3 h-3 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                                </svg>
                              )}
                            </div>
                            <p className={`font-semibold text-sm line-clamp-2 ${isLocked ? 'text-gray-500' : isModFullyCompleted ? 'text-green-800' : 'text-gray-800 group-hover:text-rose-600'
                              }`}>
                              {mod.title}
                            </p>
                            {!isLocked && (
                              <p className="text-xs text-gray-500 mt-1">
                                {lessonCount} {declension(lessonCount, ['урок', 'урока', 'уроков'])}
                              </p>
                            )}
                          </div>
                          {!isLocked && (
                            <svg className={`w-4 h-4 text-gray-400 transition-transform ${expandedModules.has(modIdx) ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                            </svg>
                          )}
                        </button>

                        {!isLocked && expandedModules.has(modIdx) && mod.lessons && mod.lessons.length > 0 && (
                          <div className="bg-gray-50/50">
                            {mod.lessons.map((lesson, lessonIdx) => {
                              const isSelected = selectedLesson?.moduleIdx === modIdx && selectedLesson?.lessonIdx === lessonIdx;
                              const isCompleted = completedLessonsIds.has(lesson.id);
                              return (
                                <button
                                  key={lesson.id}
                                  onClick={() => selectLesson(modIdx, lessonIdx)}
                                  className={`w-full text-left px-5 py-3 pl-9 text-sm transition-all flex items-start gap-3 border-l-4 ${isSelected ? 'bg-white border-rose-500 text-rose-700 shadow-sm' :
                                    isCompleted ? 'bg-green-100 border-green-500 text-green-800' :
                                      'border-transparent text-gray-600 hover:bg-gray-100'
                                    }`}
                                >
                                  <span className={`mt-0.5 flex-shrink-0 ${isCompleted ? 'text-green-600' : isSelected ? 'text-rose-500' : 'text-gray-400'}`}>
                                    {isCompleted ? (
                                      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                                        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                                      </svg>
                                    ) : (
                                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                      </svg>
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
                ) : null}
              </div>
            </div>
          </div>

          {/* КОНТЕНТ */}
          <div className="lg:col-span-8 xl:col-span-9">
            {currentLesson && !isCurrentLessonLocked ? (
              <div className="space-y-6 animate-fade-in">
                <div className="bg-black rounded-2xl overflow-hidden shadow-xl relative w-full">
                  {hasVideo ? (
                    <VideoPlayer url={videoUrl} onProgress={() => { }} />
                  ) : (
                    <div className="aspect-video flex items-center justify-center text-white bg-gray-900">
                      <div className="text-center p-4">
                        <div className="text-6xl mb-4 opacity-50">🎥</div>
                        <p className="text-lg font-medium">Видео недоступно</p>
                      </div>
                    </div>
                  )}
                </div>

                <div className="bg-white rounded-2xl p-6 sm:p-8 shadow-sm border border-gray-100">
                  <div className="flex flex-col md:flex-row md:items-start justify-between gap-4 mb-6 pb-6 border-b border-gray-100">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-2 flex-wrap">
                        <span className="px-2 py-1 bg-rose-100 text-rose-700 text-xs font-bold rounded-2xl uppercase">Урок {selectedLesson ? selectedLesson.lessonIdx + 1 : 0}</span>
                        {completedLessonsIds.has(currentLesson.id) && (
                          <span className="px-2 py-1 bg-green-100 text-green-700 text-xs font-bold rounded-2xl! flex items-center gap-1">
                            <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" /></svg>
                            Пройден
                          </span>
                        )}
                      </div>
                      <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 font-serif mb-2 break-words">{currentLesson.title}</h2>
                      <p className="text-sm text-gray-500">Модуль: {currentModule?.title}</p>
                    </div>
                    <button onClick={markAsComplete} disabled={completedLessonsIds.has(currentLesson.id) || isSavingProgress} className={`px-6 py-3 rounded-2xl! font-medium transition-all shadow-md flex items-center gap-2 whitespace-nowrap self-start md:self-center ${completedLessonsIds.has(currentLesson.id) ? 'bg-green-100 text-green-700 cursor-default border border-green-200' : 'bg-green-600 hover:bg-green-700 text-white hover:shadow-lg hover:-translate-y-0.5'} ${isSavingProgress ? 'opacity-70 cursor-wait' : ''}`}>
                      {isSavingProgress ? 'Сохранение...' : completedLessonsIds.has(currentLesson.id) ? 'Пройдено' : 'Отметить как пройденное'}
                    </button>
                  </div>
                  {currentLesson.description && (
                    <div className="prose prose-rose max-w-none text-gray-600 leading-relaxed mb-8">
                      <h4 className="text-lg font-bold text-gray-800 mb-3 not-prose flex items-center gap-2">Описание урока</h4>
                      <div className="bg-gray-50 p-4 rounded-2xl border border-gray-100"><p className="whitespace-pre-line m-0">{currentLesson.description}</p></div>
                    </div>
                  )}
                  <div className="pt-6 border-t border-gray-100 flex justify-between items-center gap-4">
                    {!isFirstLessonOfCourse ? (
                      <button onClick={handlePrevAction} className="px-5 py-2.5 bg-white border border-gray-300 text-gray-700 hover:bg-gray-50 hover:border-gray-400 rounded-2xl! font-medium transition-all shadow-sm hover:shadow flex items-center gap-2 group">
                        <svg className="w-5 h-5 text-gray-400 group-hover:text-gray-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 17l-5-5m0 0l5-5m-5 5h12" /></svg>
                        <span className="hidden sm:inline">Назад</span>
                      </button>
                    ) : <div></div>}

                    {/* Кнопка Далее всегда видна, если это не последний урок всего курса */}
                    {!isLastLessonOfCourse ? (
                      <button onClick={handleNextAction} className="px-6 py-2.5 bg-rose-600 hover:bg-rose-700 text-white rounded-2xl! font-medium transition-colors shadow-md hover:shadow-lg hover:-translate-y-0.5 flex items-center gap-2">
                        <span className="hidden sm:inline">Следующий урок</span><span className="sm:hidden">Далее</span>
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" /></svg>
                      </button>
                    ) : (
                      <button onClick={() => setShowCompletionModal(true)} className="px-6 py-2.5 bg-green-600 hover:bg-green-700 text-white rounded-2xl font-medium transition-colors shadow-md hover:shadow-lg hover:-translate-y-0.5 flex items-center gap-2">
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                        Завершить курс
                      </button>
                    )}
                  </div>
                  {!isExpired && daysRemaining > 0 && (
                    <div className="mt-3 flex items-center gap-2 text-xs">
                      <span className={`flex-shrink-0 ${daysRemaining <= 7 ? 'text-red-500' :
                          daysRemaining <= 30 ? 'text-amber-500' :
                            'text-gray-400'
                        }`}>
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      </span>
                      <span className={`${daysRemaining <= 7 ? 'text-red-600 font-medium' :
                          daysRemaining <= 30 ? 'text-amber-600 font-medium' :
                            'text-gray-500'
                        }`}>
                        До конца доступа: <strong>{daysRemaining}</strong> {declension(daysRemaining, ['день', 'дня', 'дней'])}
                        {expiresAt && (
                          <span className="hidden sm:inline opacity-75 ml-1">
                            (до {new Date(expiresAt).toLocaleDateString('ru-RU', {
                              day: 'numeric',
                              month: 'long',
                              year: 'numeric'
                            })})
                          </span>
                        )}
                      </span>
                    </div>
                  )}
                </div>
              </div>
            ) : (
              <div className="h-full flex flex-col items-center justify-center text-center p-8 bg-white rounded-2xl border border-dashed border-gray-300 min-h-[400px]">
                {isCurrentLessonLocked ? (
                  <>
                    <div className="text-6xl mb-4">🔒</div>
                    <h3 className="text-xl font-bold text-gray-800 mb-2">Этот модуль еще закрыт</h3>
                    <p className="text-gray-500 max-w-md mb-6">Следующие уроки откроются автоматически.</p>
                    <button onClick={() => navigate(`/course/${courseId}`)} className="px-6 py-3 bg-rose-500 text-white rounded-2xl! font-medium hover:bg-rose-600 transition-colors">Вернуться к описанию курса</button>
                  </>
                ) : (
                  <>
                    <div className="text-6xl mb-4">👈</div>
                    <h3 className="text-xl font-bold text-gray-800 mb-2">Выберите урок</h3>
                    <p className="text-gray-500 max-w-md">Нажмите на любой доступный урок в списке слева.</p>
                  </>
                )}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* МОБИЛЬНОЕ МЕНЮ */}
      {isMobileMenuOpen && (
        <>
          <div className="fixed inset-0 bg-black/50 z-40 lg:hidden backdrop-blur-sm" onClick={() => setIsMobileMenuOpen(false)} />
          <div className="fixed inset-y-0 left-0 w-3/4 max-w-xs bg-white z-50 shadow-2xl lg:hidden flex flex-col pt-16 animate-slideInRight">
            <div className="p-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
              <h2 className="font-bold text-gray-800 text-sm truncate pr-2">Программа курса</h2>
              <button onClick={() => setIsMobileMenuOpen(false)} className="p-2 text-gray-500 hover:bg-gray-100 rounded-2xl shrink-0">
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <div className="overflow-y-auto flex-grow pb-20 custom-scrollbar">
              {modules && modules.map((mod, modIdx) => {
                const isBonus = isModuleBonus(mod);
                const isLocked = mod.is_locked && !isBonus;

                let modCompletedCount = 0;
                mod.lessons?.forEach(l => { if (completedLessonsIds.has(l.id)) modCompletedCount++; });
                const isModFullyCompletedMobile = mod.lessons && mod.lessons.length > 0 && modCompletedCount === mod.lessons.length;

                return (
                  <div key={mod.id} className="border-b border-gray-100">
                    <button
                      onClick={(e) => handleModuleClick(modIdx, e)}
                      disabled={isLocked}
                      className={`w-full text-left px-4 py-3 font-semibold text-sm flex justify-between items-center ${isLocked ? 'bg-gray-100 text-gray-500 cursor-not-allowed' :
                        isModFullyCompletedMobile ? 'bg-green-50 text-green-800' : 'bg-gray-50 text-gray-800'
                        }`}
                    >
                      <span className="truncate pr-2">
                        {isLocked ? `🔒 ${mod.title}` : `Модуль ${modIdx + 1}: ${mod.title || ''}`}
                      </span>
                      {!isLocked && isModFullyCompletedMobile && (
                        <svg className="w-4 h-4 text-green-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" /></svg>
                      )}
                    </button>
                    {!isLocked && expandedModules.has(modIdx) && mod.lessons && (
                      <div className="bg-white">
                        {mod.lessons.map((lesson, lIdx) => {
                          const isSelected = selectedLesson?.moduleIdx === modIdx && selectedLesson?.lessonIdx === lIdx;
                          const isCompleted = completedLessonsIds.has(lesson.id);
                          return (
                            <button
                              key={lesson.id}
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

      {/* МОДАЛКА ЗАБЛОКИРОВАННОГО МОДУЛЯ */}
      {lockedModuleInfo.isOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-fadeIn">
          <div className="bg-white rounded-3xl shadow-2xl max-w-md w-full p-8 text-center transform transition-all scale-100 animate-scaleIn relative overflow-hidden">
            <div className="absolute top-0 left-0 w-full h-2 bg-gradient-to-r from-gray-200 via-gray-400 to-gray-200"></div>

            <div className="w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-6">
              <svg className="w-10 h-10 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
            </div>

            <h3 className="text-2xl font-serif font-bold text-gray-900 mb-2">Модуль закрыт</h3>
            <p className="text-gray-600 mb-6 font-medium">{lockedModuleInfo.title}</p>

            <div className="bg-gray-50 rounded-2xl! p-6 mb-8 border border-gray-100">
              {lockedModuleInfo.unlockDate && !isNaN(lockedModuleInfo.unlockDate.getTime()) ? (
                <>
                  <p className="text-sm text-gray-500 mb-4 uppercase tracking-wide font-bold">Откроется:</p>
                  <p className="text-xl font-bold text-rose-600 mb-2">{formatDateTime(lockedModuleInfo.unlockDate)}</p>

                  {(() => {
                    const time = getTimeRemaining(lockedModuleInfo.unlockDate);
                    if (!time) return <p className="text-green-600 font-bold">Уже доступно!</p>;
                    return (
                      <div className="flex justify-center gap-3 mt-4">
                        {[
                          { val: time.days, label: 'дн' },
                          { val: time.hours, label: 'ч' },
                          { val: time.minutes, label: 'м' },
                          { val: time.seconds, label: 'с' }
                        ].map((item, idx) => (
                          <div key={idx} className="flex flex-col items-center">
                            <div className="w-12 h-12 bg-white rounded-xl shadow-sm border border-gray-200 flex items-center justify-center text-lg font-bold text-gray-800">
                              {String(item.val).padStart(2, '0')}
                            </div>
                            <span className="text-xs text-gray-400 mt-1 font-medium">{item.label}</span>
                          </div>
                        ))}
                      </div>
                    );
                  })()}
                </>
              ) : (
                <p className="text-gray-600">Дата открытия уточняется.</p>
              )}
            </div>

            <button
              onClick={() => setLockedModuleInfo({ ...lockedModuleInfo, isOpen: false })}
              className="w-full py-3.5 px-6 bg-gray-900 text-white font-bold rounded-2xl! hover:bg-gray-800 transition-colors shadow-lg"
            >
              Понятно
            </button>
          </div>
        </div>
      )}

      {/* МОДАЛКА ЗАВЕРШЕНИЯ КУРСА */}
      {showCompletionModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-fadeIn">
          <div className="bg-white rounded-3xl shadow-2xl max-w-md w-full p-8 text-center transform transition-all scale-100 animate-scaleIn">
            <div className="w-24 h-24 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6 animate-bounce">
              <svg className="w-12 h-12 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            </div>

            <h3 className="text-2xl font-serif font-bold text-gray-900 mb-3">Поздравляем!</h3>
            <p className="text-gray-600 mb-8 leading-relaxed">
              Вы успешно прошли все доступные уроки этого курса. Отличная работа!
            </p>

            <div className="flex flex-col gap-3">
              <button
                onClick={() => {
                  setShowCompletionModal(false);
                  navigate(`/course/${courseId}`);
                }}
                className="w-full py-3.5 px-6 bg-rose-500 text-white font-bold rounded-2xl! hover:bg-rose-600 transition-colors shadow-lg hover:shadow-xl transform hover:-translate-y-0.5"
              >
                Вернуться к курсу
              </button>
              <button
                onClick={() => setShowCompletionModal(false)}
                className="w-full py-3.5 px-6 bg-white text-gray-700 font-bold rounded-2xl! border-2 border-gray-200 hover:bg-gray-50 hover:border-gray-300 transition-all"
              >
                Продолжить просмотр
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default LearnPage;
