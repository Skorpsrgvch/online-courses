import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import { coursesService } from '../../api/courses.service';
import { reviewsService } from '../../api/reviews.service';
import { useAuth } from '../../context/AuthContext';
import type { CourseFullResponse, Review, BonusItem, PaymentResponse } from '../../api/types';
import { Comments } from '../../components/course/Comments';

const CoursePage = () => {
  const { id } = useParams<{ id: string }>();
  const courseId = Number(id);
  const navigate = useNavigate();
  const location = useLocation(); // Для определения, вернулись ли мы с оплаты
  const { user, isAuthenticated } = useAuth();

  const [courseData, setCourseData] = useState<CourseFullResponse | null>(null);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);

  // Функция загрузки данных
  const loadData = async (forceCheckPurchase = false) => {
    if (!courseId) return;

    setIsLoading(true);
    setError(null);
    try {
      
      const courseRes = await coursesService.getCourseFull(courseId);
      
      const reviewsData = await reviewsService.getCourseReviews(courseId).catch(() => []);

      setCourseData(courseRes);
      setReviews(reviewsData);
      
    } catch (err: any) {
      console.error(err);
      setError(err.message || 'Не удалось загрузить курс');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const loadData = async () => {
      if (!courseId) return;

      setIsLoading(true);
      setError(null);
      try {
        const [courseRes, reviewsData] = await Promise.all([
          coursesService.getCourseFull(courseId),
          reviewsService.getCourseReviews(courseId).catch(() => []),
        ]);

        setCourseData(courseRes);
        setReviews(reviewsData);
      } catch (err: any) {
        console.error(err);
        setError(err.message || 'Не удалось загрузить курс');
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, [courseId]);


  useEffect(() => {

    if (location.pathname.includes('/payment-success')) {

      const timer = setTimeout(() => {

        const refreshData = async () => {
          try {
            const courseRes = await coursesService.getCourseFull(courseId);
            setCourseData(courseRes);
            

          } catch (err) {
            console.error("Failed to refresh course data after payment", err);
          }
        };
        
        refreshData();
      }, 1000); 

      return () => clearTimeout(timer);
    }
  }, [location.pathname, courseId, navigate]);

  const handleActionClick = async () => {
    if (isProcessing) return;

    if (!courseData || !user) {
      alert('Пожалуйста, войдите в систему для покупки.');
      navigate('/login');
      return;
    }

    const isFree = courseData.course.price === 0 || courseData.course.is_public;
    const isPurchased = courseData.course.is_purchased || user?.role === 'admin';

    if (isFree || isPurchased) {
      // Переход к обучению
      navigate(`/course/${courseId}/learn`);
    } else {
      // Если вдруг статус не обновился, но пользователь утверждает, что купил - можно предложить обновить
      const confirmRefresh = window.confirm('Статус оплаты еще не обновился. Обновить страницу?');
      if (confirmRefresh) {
        window.location.reload();
      }
    }
  };

  // Реализация оплаты
  const handlePurchase = async () => {
    if (isProcessing) return;
    
    if (!courseData || !user) {
      alert('Пожалуйста, войдите в систему для покупки.');
      navigate('/login');
      return;
    }

    setIsProcessing(true);

    try {
      const returnUrl = `${window.location.origin}/course/${courseId}/payment-success`;
      
      const paymentData: PaymentResponse = await coursesService.createPayment(courseId, returnUrl);

      if (paymentData.confirmation_url) {
        // Перенаправление на шлюз ЮKassa
        window.location.href = paymentData.confirmation_url;
      } else {
        throw new Error('Не удалось получить ссылку на оплату');
      }
    } catch (error: any) {
      console.error('Ошибка при создании платежа:', error);
      
      let errorMsg = 'Не удалось создать платеж. Попробуйте позже.';
      if (error.response?.status === 409) {
        errorMsg = 'Вы уже приобрели этот курс. Обновляем страницу...';
        setTimeout(() => window.location.reload(), 2000);
      } else if (error.message) {
        errorMsg = error.message;
      }
      
      alert(errorMsg);
      setIsProcessing(false);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="w-12 h-12 border-4 border-rose-300 border-t-rose-500 rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-500">Загрузка информации о курсе...</p>
        </div>
      </div>
    );
  }

  if (error || !courseData) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
        <div className="text-center bg-white p-8 rounded-xl shadow-sm max-w-md">
          <div className="text-5xl mb-4">😔</div>
          <h2 className="text-xl font-semibold text-gray-800 mb-2">Ошибка доступа</h2>
          <p className="text-gray-600 mb-4">{error || 'Курс не найден'}</p>
          <button
            onClick={() => navigate('/courses')}
            className="px-6 py-2 bg-rose-500 text-white rounded-2xl! hover:bg-rose-600 transition-colors"
          >
            Вернуться в каталог
          </button>
        </div>
      </div>
    );
  }

  const renderSplitList = (text: string, variant: 'check' | 'dot' = 'check') => {
    if (!text) return null;
    const items = text.split('|||').map(i => i.trim()).filter(i => i.length > 0);

    if (items.length === 0) return null;

    return (
      <ul className="space-y-4">
        {items.map((item, idx) => (
          <li key={idx} className="flex items-start gap-3">
            {variant === 'check' ? (
              <span className="flex-shrink-0 w-6 h-6 rounded-full bg-green-100 text-green-600 flex items-center justify-center mt-0.5">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              </span>
            ) : (
              <span className="flex-shrink-0 w-2 h-2 rounded-full bg-rose-400 mt-2.5"></span>
            )}
            <span className="text-gray-700 leading-relaxed text-base">{item}</span>
          </li>
        ))}
      </ul>
    );
  };

  const renderList = (text: string, type: 'bad' | 'good') => {
    if (!text) return null;

    const icon = type === 'bad' ? '❌' : '✅';
    const bgClass = type === 'bad' ? 'bg-red-50' : 'bg-emerald-50';
    const textClass = type === 'bad' ? 'text-red-900' : 'text-emerald-900';
    const borderClass = type === 'bad' ? 'border-red-100' : 'border-emerald-100';
    const title = type === 'bad' ? 'Противопоказания' : 'Рекомендации';
    const IconSvg = type === 'bad'
      ? () => <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>
      : () => <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>;

    const normalizedText = text.replace(/\. /g, '|||').replace(/\n/g, '|||');
    const items = normalizedText.split('|||')
      .map(item => item.trim())
      .filter(item => item.length > 0);

    return (
      <section className="bg-white rounded-2xl p-6 sm:p-8 shadow-sm border border-gray-100 mb-8">
        <div className="flex items-center gap-3 mb-6">
          <div className={`p-2 ${type === 'bad' ? 'bg-red-100 text-red-600' : 'bg-emerald-100 text-emerald-600'} rounded-full`}>
            <IconSvg />
          </div>
          <h3 className="text-xl font-bold text-gray-900">{title}</h3>
        </div>

        <div className="space-y-3">
          {items.map((line, idx) => {
            const cleanLine = line.replace(/^[❌✅]\s*/, '').replace(/^\.\s*/, '');
            if (!cleanLine) return null;
            return (
              <div key={idx} className={`flex items-start gap-3 p-4 rounded-xl ${bgClass} border ${borderClass} transition-transform hover:scale-[1.01]`}>
                <span className="text-2xl flex-shrink-0 mt-0.5 select-none">{icon}</span>
                <span className={`text-sm md:text-base leading-relaxed font-medium ${textClass}`}>
                  {cleanLine}
                </span>
              </div>
            );
          })}
        </div>
      </section>
    );
  };

  const { course, modules } = courseData;


  const isFree = course.price === 0 || course.is_public;

  const isPurchased = course.is_purchased || user?.role === 'admin';

  const getActionButtonProps = () => {
    if (isFree || isPurchased) {
      return {
        text: isPurchased && !isFree ? 'Продолжить обучение' : 'Начать обучение',
        color: 'bg-rose-500 hover:bg-rose-600 shadow-rose-200',
        action: handleActionClick
      };
    }
    return {
      text: `Купить за ${course.price} ₽`,
      color: 'bg-green-600 hover:bg-green-700 shadow-green-200',
      action: handlePurchase,
      disabled: isProcessing
    };
  };

  const buttonProps = getActionButtonProps();

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* HERO СЕКЦИЯ */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-350 mx-auto px-4 sm:px-6 lg:px-8 py-12 lg:py-16">
          <div className="mb-8">
            <button
              onClick={() => navigate('/courses')}
              className="inline-flex items-center gap-2 text-sm text-gray-500 hover:text-rose-500 transition-colors mb-8 font-medium"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
              Назад к каталогу
            </button>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            <div className="order-2 lg:order-1">
              <h1 className="text-3xl lg:text-5xl font-serif font-bold text-gray-900 mb-6 leading-tight">
                {course.title}
              </h1>
              <p className="text-lg text-gray-600 mb-8 leading-relaxed">
                {course.description}
              </p>

              <div className="flex flex-wrap gap-3 mb-8">
                <span className="px-4 py-2 bg-rose-50 text-rose-700 rounded-full text-sm font-medium flex items-center gap-2">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
                  Онлайн формат
                </span>
                <span className="px-4 py-2 bg-blue-50 text-blue-700 rounded-full text-sm font-medium flex items-center gap-2">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" /></svg>
                  {modules ? modules.reduce((acc, m) => acc + (m.lessons?.length || 0), 0) : 0} уроков
                </span>
                
              </div>

              <button
                onClick={buttonProps.action}
                disabled={buttonProps.disabled}
                className={`w-full sm:w-auto px-10 py-4 text-lg font-bold text-white rounded-2xl! shadow-lg transition-all transform hover:-translate-y-1 hover:shadow-xl ${buttonProps.color} ${buttonProps.disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
              >
                {isProcessing ? 'Обработка...' : buttonProps.text}
              </button>
              
              
            </div>

            <div className="order-1 lg:order-2 relative rounded-2xl overflow-hidden shadow-2xl aspect-video lg:aspect-square max-h-125 bg-gray-100">
              {course.cover_image_url ? (
                <img src={course.cover_image_url} alt={course.title} className="w-full h-full object-cover" />
              ) : (
                <div className="w-full h-full flex items-center justify-center text-gray-400 bg-gray-200">
                  <span className="text-lg">Изображение курса</span>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-250 mx-auto px-4 sm:px-6 lg:px-8 py-12 space-y-12">
        {/* Остальные секции (target_audience, course_basis и т.д.) без изменений */}
        {course.target_audience && (
          <section>
            <h2 className="text-3xl font-serif font-bold text-gray-900 mb-8! text-center">Курс для вас, если</h2>
            <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
              {renderSplitList(course.target_audience, 'check')}
              <div className="mt-10 flex justify-center">
                <button onClick={buttonProps.action} disabled={buttonProps.disabled} className={`px-8 py-3 text-base font-bold text-white rounded-2xl! shadow-md transition-transform hover:-translate-y-0.5 ${buttonProps.color}`}>
                  {buttonProps.text}
                </button>
              </div>
            </div>
          </section>
        )}

        {course.course_basis && (
          <section>
            <h2 className="text-3xl font-serif font-bold text-gray-900 mb-8! text-center">Курс включает в себя</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {course.course_basis.split('|||').filter(i => i.trim()).map((item, idx) => (
                <div key={idx} className="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 flex items-center gap-4 hover:shadow-md transition-shadow">
                  <div className="w-12 h-12 rounded-full bg-rose-50 text-rose-500 flex items-center justify-center flex-shrink-0">
                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                  </div>
                  <span className="font-medium text-gray-800 leading-snug">{item}</span>
                </div>
              ))}
            </div>
            <div className="mt-10 flex justify-center">
              <button onClick={buttonProps.action} disabled={buttonProps.disabled} className={`px-8 py-3 text-base font-bold text-white rounded-2xl! shadow-md transition-transform hover:-translate-y-0.5 ${buttonProps.color}`}>
                {buttonProps.text}
              </button>
            </div>
          </section>
        )}

        {course.class_basis && (
          <section>
            <h2 className="text-3xl font-serif font-bold text-gray-900 mb-8! text-center">Основа занятий</h2>
            <div className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
              {renderSplitList(course.class_basis, 'dot')}
            </div>
          </section>
        )}

        {course.recommendations && renderList(course.recommendations, 'good')}
        {course.contraindications && renderList(course.contraindications, 'bad')}

        {course.bonuses && course.bonuses.length > 0 && (
          <section className="pt-4">
            <h3 className="text-2xl font-serif font-bold text-gray-900 mb-10! text-center">🎁 Бонусы при покупке</h3>
            <div className="flex flex-wrap justify-center gap-6">
              {course.bonuses.map((bonus: BonusItem, idx: number) => {
                let icon = '🎁';
                if (bonus.icon === 'file') icon = '📄';
                if (bonus.icon === 'video') icon = '🎥';
                if (bonus.icon === 'chat') icon = '💬';
                if (bonus.icon === 'audio') icon = '🎧';
                if (bonus.icon === 'check') icon = '✅';
                return (
                  <div key={idx} className="bg-white p-6 rounded-xl shadow-md border border-gray-100 text-center hover:shadow-lg transition-shadow group">
                    <div className="text-5xl mb-4 transform group-hover:scale-110 transition-transform duration-300">{icon}</div>
                    <h4 className="font-bold text-gray-800 mb-2 text-lg">{bonus.title}</h4>
                    <p className="text-sm text-gray-600 leading-relaxed">{bonus.description}</p>
                  </div>
                );
              })}
            </div>
          </section>
        )}

        <section className="pt-8 border-t border-gray-200">
          <div className="flex items-center gap-3 mb-8">
            <h3 className="text-2xl font-serif font-bold text-gray-900">Отзывы студентов</h3>
          </div>
          <Comments courseId={courseId} reviews={reviews} onReviewSubmitted={() => loadData()} />
        </section>
      </div>
    </div>
  );
};

export default CoursePage;