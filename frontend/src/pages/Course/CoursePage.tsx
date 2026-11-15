import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, useLocation, Link } from 'react-router-dom';
import { coursesService } from '../../api/courses.service';
import { reviewsService } from '../../api/reviews.service';
import { useAuth } from '../../context/AuthContext';
import type { CourseFullResponse, Review, BonusItem, PaymentResponse } from '../../api/types';
import { Comments } from '../../components/course/Comments';

const CoursePage = () => {
  const { id } = useParams<{ id: string }>();
  const courseId = Number(id);
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useAuth();

  const [courseData, setCourseData] = useState<CourseFullResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);

  const [isCheckoutOpen, setIsCheckoutOpen] = useState(false);
  const [agreeOffer, setAgreeOffer] = useState(false);
  const [agreeNewsletter, setAgreeNewsletter] = useState(false);
  const [checkoutEmail, setCheckoutEmail] = useState('');
  const [checkoutError, setCheckoutError] = useState<string | null>(null);

  const [approvedReviews, setApprovedReviews] = useState<Review[]>([]);
  const [myReview, setMyReview] = useState<Review | null>(null);

  const loadData = async () => {
    if (!courseId) return;
    // Не блокируем интерфейс полностью при фоновом обновлении, если данные уже есть
    if (!courseData) setIsLoading(true);
    setError(null);

    try {
      const [courseRes, publicReviewsData, allUserReviews] = await Promise.all([
        coursesService.getCourseFull(courseId),
        reviewsService.getCourseReviews(courseId).catch(() => []),
        user ? reviewsService.getMyReviews().catch(() => []) : Promise.resolve([])
      ]);

      setCourseData(courseRes);
      setApprovedReviews(Array.isArray(publicReviewsData) ? publicReviewsData.filter(r => r.approved) : []);
      const userReviewForThisCourse = allUserReviews.find(r => r.course_id === courseId) || null;
      setMyReview(userReviewForThisCourse);
    } catch (err: any) {
      console.error(err);
      setError(err.message || 'Не удалось загрузить курс');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, [courseId, user?.id]);

  // Обработка успешной оплаты и скролл
  useEffect(() => {
    if (!isLoading && courseData) {

      const hash = location.hash;
      const timer = setTimeout(() => {
        if (hash) {
          const element = document.querySelector(hash);
          if (element) element.scrollIntoView({ behavior: 'smooth', block: 'start' });
        } else {
          window.scrollTo({ top: 0, left: 0, behavior: 'auto' });
        }
      }, 100);
      return () => clearTimeout(timer);
    }
  }, [location, isLoading, courseData, navigate, courseId]);

  const requireAuth = () => {
    if (!user) {
      navigate('/login', { state: { from: location.pathname } });
      return false;
    }
    return true;
  };

  const handleActionClick = async () => {

    if (!requireAuth()) {
      return;
    }
    if (!courseData) {
      return;
    }

    const isFree = courseData.course.price === 0 || courseData.course.is_public;


    if (isFree && !courseData.course.is_purchased && user?.role !== 'admin') {
      try {
        setIsProcessing(true);
        const priceToSend = Number(courseData.course.price) || 0;

        await coursesService.enrollFree(courseId, priceToSend);

        await loadData();
        setIsProcessing(false);
      } catch (error: any) {
        console.error("Ошибка при зачислении:", error);
        setError("Не удалось оформить доступ.");
        setIsProcessing(false);
        return;
      }
    }

    navigate(`/course/${courseId}/learn`);
  };

  const handlePurchase = async () => {
    if (!requireAuth()) return;
    if (isProcessing) return;
    if (!courseData) return;

    setIsProcessing(true);
    try {
      const returnUrl = `${window.location.origin}/course/${courseId}/payment-success`;
      const paymentData: PaymentResponse = await coursesService.createPayment(courseId, returnUrl);

      if (paymentData.confirmation_url) {
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
      setError(errorMsg);
      setIsProcessing(false);
    }
  };

  const handleOpenCheckout = () => {
    if (!requireAuth()) return;
    if (!courseData) return;
    setCheckoutError(null);
    setAgreeOffer(false);
    setAgreeNewsletter(false);
    setCheckoutEmail(user?.email || '');
    setIsCheckoutOpen(true);
  };

  const handleConfirmCheckout = async () => {
    if (!agreeOffer) {
      setCheckoutError('Необходимо согласие с условиями оферты');
      return;
    }
    setCheckoutError(null);
    setIsCheckoutOpen(false);
    await handlePurchase(); // та же логика — редирект на ЮKassa
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
          <button onClick={() => navigate('/courses')} className="px-6 py-2 bg-rose-500 text-white rounded-2xl hover:bg-rose-600 transition-colors">
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
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>
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
    const items = normalizedText.split('|||').map(item => item.trim()).filter(item => item.length > 0);

    return (
      <section className="bg-white rounded-2xl p-6 sm:p-8 shadow-sm border border-gray-100 mb-8!">
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
                <span className={`text-sm md:text-base leading-relaxed font-medium ${textClass}`}>{cleanLine}</span>
              </div>
            );
          })}
        </div>
      </section>
    );
  };

  const { course, modules } = courseData;

  const progress = course.progress_percent || 0;
  const isFree = course.price === 0 || course.is_public;
  // Логика доступа: курс бесплатный ИЛИ куплен ИЛИ пользователь админ
  const userHasAccess = isFree || course.is_purchased || user?.role === 'admin';
  const isCompleted = progress === 100;

  const getActionButtonProps = () => {
    if (!userHasAccess) {
      return {
        text: `Купить за ${course.price} ₽`,
        color: 'bg-green-600 hover:bg-green-700 shadow-green-200',
        action: handleOpenCheckout,
        disabled: isProcessing
      };
    }

    if (isCompleted) {
      return {
        text: 'Повторить курс',
        color: 'bg-gray-500 hover:bg-gray-600 shadow-gray-200',
        action: handleActionClick,
        disabled: false
      };
    }

    return {
      text: progress > 0 ? 'Продолжить обучение' : 'Начать обучение',
      color: 'bg-rose-500 hover:bg-rose-600 shadow-rose-200',
      action: handleActionClick,
      disabled: false
    };
  };

  const buttonProps = getActionButtonProps();

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* HERO СЕКЦИЯ */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-350 mx-auto px-4 sm:px-6 lg:px-8 py-12 lg:py-16">
          <div className="mb-8">
            <button onClick={() => navigate('/courses')} className="inline-flex items-center gap-2 text-sm text-gray-500 hover:text-rose-500 transition-colors mb-8 font-medium">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" /></svg>
              Назад к каталогу
            </button>
          </div>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            <div className="order-2 lg:order-1">
              <h1 className="text-3xl lg:text-5xl font-serif font-bold text-gray-900 mb-6 leading-tight">{course.title}</h1>
              <p className="text-lg text-gray-600 mb-8 leading-relaxed">{course.description}</p>
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
              <button onClick={buttonProps.action} disabled={buttonProps.disabled} className={`w-full sm:w-auto px-10 py-4 text-lg font-bold text-white rounded-2xl! shadow-lg transition-all transform hover:-translate-y-1 hover:shadow-xl ${buttonProps.color} ${buttonProps.disabled ? 'opacity-50 cursor-not-allowed' : ''}`}>
                {isProcessing ? 'Обработка...' : buttonProps.text}
              </button>
            </div>
            <div className="order-1 lg:order-2 relative rounded-2xl! overflow-hidden shadow-2xl aspect-video max-h-125 bg-gray-100">
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
                <div key={idx} className="bg-white p-6 rounded-2xl! shadow-sm border border-gray-100 flex items-center gap-4 hover:shadow-md transition-shadow">
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

        <section id="reviews" className="pt-8 border-t border-gray-200 scroll-mt-24">
          <div className="flex items-center gap-3 mb-8">
            <h3 className="text-2xl font-serif font-bold text-gray-900">Отзывы клиентов</h3>
          </div>
          <Comments
            courseId={courseId}
            reviews={approvedReviews}
            myReview={myReview}
            hasAccess={userHasAccess}
            onReviewSubmitted={loadData}
          />
        </section>
      </div>


      {isCheckoutOpen && (
        <div
          className="fixed inset-0 z-50 overflow-y-auto bg-black/50 backdrop-blur-sm"
          onClick={() => setIsCheckoutOpen(false)}
        >
          <div className="flex w-full min-h-full items-center justify-center p-4 pb-8">
            <div
              className="bg-white rounded-2xl shadow-2xl w-full max-w-md p-6 sm:p-8 relative"
              onClick={(e) => e.stopPropagation()}
            >
              <button
                onClick={() => setIsCheckoutOpen(false)}
                className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 transition-colors"
                aria-label="Закрыть"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>

              <h2 className="text-2xl font-serif font-bold text-gray-900 mb-2 text-center">
                Оформление покупки
              </h2>
              <p className="text-center text-gray-500 text-sm mb-6">
                {course.title}
              </p>

              <div className="bg-rose-50 rounded-xl p-4 mb-6 text-center">
                <p className="text-sm text-rose-700 mb-1">Стоимость курса</p>
                <p className="text-3xl font-bold text-rose-600">{course.price} ₽</p>
              </div>

              {/* Email — пока не отправляется, но по ТЗ поле есть */}
              <div className="mb-5">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Email для чека и доступа
                </label>
                <input
                  type="email"
                  value={checkoutEmail}
                  onChange={(e) => setCheckoutEmail(e.target.value)}
                  placeholder="example@mail.ru"
                  className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400 outline-none transition-all"
                />
              </div>

              {/* Галочка оферты — обязательная */}
              <div className="flex items-start mb-3">
                <input
                  type="checkbox"
                  id="checkout-offer"
                  checked={agreeOffer}
                  onChange={(e) => setAgreeOffer(e.target.checked)}
                  className="h-4 w-4 mt-0.5 text-rose-600 focus:ring-rose-500 border-gray-300 rounded cursor-pointer"
                />
                <label htmlFor="checkout-offer" className="ml-2 block text-sm text-gray-600 leading-tight cursor-pointer">
                  Я принимаю условия{' '}
                  <Link
                    to="/offer"
                    target="_blank"
                    rel="noopener noreferrer"
                    onClick={(e) => e.stopPropagation()}
                    className="text-rose-500! hover:underline!"
                    style={{ textDecoration: 'none' }}
                  >
                    публичной оферты
                  </Link>{' '}
                  и подтверждаю ознакомление с{' '}
                  <Link
                    to="/privacy-policy"
                    target="_blank"
                    rel="noopener noreferrer"
                    onClick={(e) => e.stopPropagation()}
                    className="text-rose-500! hover:underline!"
                    style={{ textDecoration: 'none' }}
                  >
                    политикой обработки персональных данных
                  </Link>
                </label>
              </div>

              {/* Галочка рассылки — необязательная */}
              <div className="flex items-start mb-5">
                <input
                  type="checkbox"
                  id="checkout-newsletter"
                  checked={agreeNewsletter}
                  onChange={(e) => setAgreeNewsletter(e.target.checked)}
                  className="h-4 w-4 mt-0.5 text-rose-600 focus:ring-rose-500 border-gray-300 rounded cursor-pointer"
                />
                <label htmlFor="checkout-newsletter" className="ml-2 block text-sm text-gray-600 leading-tight cursor-pointer">
                  Я согласен(на) на получение{' '}
                  <Link
                    to="/privacy-policy"
                    target="_blank"
                    rel="noopener noreferrer"
                    onClick={(e) => e.stopPropagation()}
                    className="text-rose-500! hover:underline!"
                    style={{ textDecoration: 'none' }}
                  >{' '}
                  информационной и рекламной рассылки
                  </Link>
                   (необязательно)
                </label>
              </div>

              {checkoutError && (
                <div className="p-3 mb-4 bg-red-50 text-red-700 text-sm rounded-lg border border-red-100">
                  {checkoutError}
                </div>
              )}

              <button
                onClick={handleConfirmCheckout}
                disabled={isProcessing}
                className="w-full py-3 px-4 bg-green-600 text-white rounded-2xl! hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-bold shadow-md shadow-green-200"
              >
                {isProcessing ? 'Обработка...' : 'Купить курс'}
              </button>

              <button
                onClick={() => setIsCheckoutOpen(false)}
                className="w-full mt-2 py-2 text-sm text-gray-500 hover:text-gray-700 transition-colors"
              >
                Отмена
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default CoursePage;
