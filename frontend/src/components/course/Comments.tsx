import React, { useState, useMemo, useEffect } from 'react';
import { useAuth } from '../../context/AuthContext';
import { reviewsService } from '../../api/reviews.service';
import type { Review } from '../../api/types';
import { Button } from '../ui/Button';

interface CommentsProps {
  courseId: number;
  reviews: Review[] | null;
  myReview?: Review | null; // Отзыв текущего пользователя
  hasAccess?: boolean;     // Куплен ли курс
  onReviewSubmitted: () => void; // Колбэк для обновления данных у родителя
}

export const Comments: React.FC<CommentsProps> = ({
  courseId,
  reviews,
  myReview,
  hasAccess = false,
  onReviewSubmitted
}) => {
  const { isAuthenticated } = useAuth();

  // Состояния формы
  const [newReviewText, setNewReviewText] = useState('');
  const [newReviewRating, setNewReviewRating] = useState(5);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Режим редактирования (активен только если есть myReview)
  const [isEditing, setIsEditing] = useState(false);

  // Состояние для модального окна успеха
  const [showSuccessModal, setShowSuccessModal] = useState(false);

  // Состояние для ошибки валидации
  const [validationError, setValidationError] = useState<string | null>(null);

  const [currentSlide, setCurrentSlide] = useState(0);
  const reviewsArray = reviews || [];

  // Подстановка данных при входе в режим редактирования
  useEffect(() => {
    if (isEditing && myReview) {
      setNewReviewText(myReview.text);
      setNewReviewRating(myReview.rating);
      setValidationError(null);
    }
  }, [isEditing, myReview]);

  // Сброс формы, если отзыв появился/исчез (например, после успешной отправки)
  useEffect(() => {
    if (!isEditing && myReview) {
      // Если мы не редактируем, но отзыв есть, сбрасываем форму в чистое состояние
      // (на случай, если пользователь был в процессе ввода нового)
      setNewReviewText('');
      setNewReviewRating(5);
    }
  }, [myReview, isEditing]);

  const averageRating = useMemo(() => {
    if (reviewsArray.length === 0) return 0;
    const sum = reviewsArray.reduce((acc, review) => acc + review.rating, 0);
    return (sum / reviewsArray.length).toFixed(1);
  }, [reviewsArray]);

  const nextSlide = () => {
    if (reviewsArray.length <= 1) return;
    setCurrentSlide((prev) => (prev === reviewsArray.length - 1 ? 0 : prev + 1));
  };

  const prevSlide = () => {
    if (reviewsArray.length <= 1) return;
    setCurrentSlide((prev) => (prev === 0 ? reviewsArray.length - 1 : prev - 1));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    setValidationError(null);

    if (!isAuthenticated) {
      setValidationError('Пожалуйста, войдите в систему');
      return;
    }

    const trimmedText = newReviewText.trim();
    if (trimmedText.length < 10) {
      setValidationError('Текст отзыва должен содержать минимум 10 символов.');
      return;
    }

    setIsSubmitting(true);
    try {
      await reviewsService.submitReview({
        course_id: courseId,
        text: trimmedText,
        rating: newReviewRating,
      });

      // Сбрасываем режим редактирования
      setIsEditing(false);

      // Показываем модалку
      setShowSuccessModal(true);

      // Уведомляем родителя об обновлении (он перезагрузит myReview)
      onReviewSubmitted();

    } catch (err: any) {
      setValidationError(err.message || 'Ошибка при отправке отзыва.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setNewReviewText(e.target.value);
    if (validationError) setValidationError(null);
  };

  const cancelEdit = () => {
    setIsEditing(false);
    setValidationError(null);
  };

  const renderStars = (rating: number, size = 'sm') => {
    const starSize = size === 'lg' ? 'text-xl' : 'text-sm';
    return (
      <div className={`flex ${starSize} text-yellow-400`}>
        {[1, 2, 3, 4, 5].map((star) => (
          <span key={star} className={star <= rating ? 'text-yellow-400' : 'text-gray-300'}>★</span>
        ))}
      </div>
    );
  };

  return (
    <div className="space-y-8">

      {/* БЛОК ВЗАИМОДЕЙСТВИЯ (Форма или Статус) */}
      {isAuthenticated && hasAccess && (
        <div className="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">


          {/* СЦЕНАРИЙ 2: Курс куплен, отзыва НЕТ -> Форма создания */}
          {hasAccess && !myReview && (
            <>
              <h4 className="text-lg font-bold text-gray-800 mb-4">Оставить отзыв</h4>
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">Ваша оценка</label>
                <div className="flex gap-2">
                  {[1, 2, 3, 4, 5].map((star) => (
                    <button
                      key={star}
                      type="button"
                      onClick={() => setNewReviewRating(star)}
                      className={`text-3xl transition-transform hover:scale-110 focus:outline-none ${star <= newReviewRating ? 'text-yellow-400 drop-shadow-sm' : 'text-gray-300'
                        }`}
                    >
                      ★
                    </button>
                  ))}
                </div>
              </div>

              <div className="mb-2">
                <div className="flex justify-between items-center mb-2">
                  <label className="block text-sm font-medium text-gray-700">Текст отзыва</label>
                  <span className="text-xs text-gray-500 font-medium">Минимум 10 символов</span>
                </div>

                <textarea
                  value={newReviewText}
                  onChange={handleTextChange}
                  placeholder="Поделитесь своими впечатлениями о курсе..."
                  className={`w-full p-4 border rounded-xl! focus:ring-2 focus:border-transparent text-sm min-h-[100px] resize-y transition-all duration-200 ${validationError
                      ? 'border-red-300 bg-red-50 focus:ring-red-200 text-red-900 placeholder-red-300'
                      : 'border-gray-200 focus:ring-rose-200'
                    }`}
                />

                {validationError && (
                  <p className="mt-2 text-sm text-red-600 flex items-center gap-1">
                    <span>⚠️</span> {validationError}
                  </p>
                )}
              </div>

              <div className="mt-4 flex justify-end">
                <Button
                  type="submit"
                  isLoading={isSubmitting}
                  variant="primary"
                  className="rounded-xl! px-8 py-2.5"
                  onClick={handleSubmit}
                >
                  Опубликовать
                </Button>
              </div>
            </>
          )}

          {/* СЦЕНАРИЙ 3: Курс куплен, отзыв ЕСТЬ -> Просмотр или Редактирование */}
          {hasAccess && myReview && (
            <>
              {!isEditing ? (
                /* --- Вид просмотра (Заглушка) --- */
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <h4 className="text-lg font-bold text-gray-800">Ваш отзыв</h4>
                    <span className="text-xs text-gray-400">
                      {new Date(myReview.created_at).toLocaleDateString()}
                    </span>
                  </div>

                  <div className="p-5 bg-gray-50 rounded-xl border border-gray-100">
                    <div className="flex items-center gap-2 mb-3">
                      <span className="font-bold text-gray-900">Вы</span>
                      {renderStars(myReview.rating)}
                    </div>

                    {/* Статус модерации */}
                    <div className="mb-3">
                      {myReview.approved ? (
                        <span className="inline-flex items-center gap-1 px-3 py-1 bg-green-100 text-green-700 text-xs font-bold rounded-full">
                          ✓ Опубликован
                        </span>
                      ) : myReview.rejection_reason ? (
                        <span className="inline-flex items-center gap-1 px-3 py-1 bg-red-100 text-red-700 text-xs font-bold rounded-full">
                          ✕ Отклонен
                        </span>
                      ) : (
                        <span className="inline-flex items-center gap-1 px-3 py-1 bg-yellow-100 text-yellow-700 text-xs font-bold rounded-full">
                          ⏳ На модерации
                        </span>
                      )}
                    </div>

                    <p className="text-gray-700 whitespace-pre-line leading-relaxed mb-4">{myReview.text}</p>

                    {myReview.rejection_reason && (
                      <div className="p-3 bg-red-50 border border-red-100 rounded-lg">
                        <p className="text-xs font-bold text-red-800 mb-1">Причина отклонения:</p>
                        <p className="text-sm text-red-700">{myReview.rejection_reason}</p>
                      </div>
                    )}
                  </div>

                  <div className="flex gap-3">
                    <button
                      onClick={() => setIsEditing(true)}
                      className="px-5 py-2.5 text-sm font-medium text-blue-600 bg-blue-50 border border-blue-200 rounded-xl hover:bg-blue-100 transition-colors flex items-center gap-2 shadow-sm"
                    >
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                      Исправить отзыв
                    </button>
                  </div>
                </div>
              ) : (
                /* --- Вид редактирования (Форма с данными) --- */
                <>
                  <h4 className="text-lg font-bold text-gray-800 mb-4">Редактирование отзыва</h4>
                  <div className="mb-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">Ваша оценка</label>
                    <div className="flex gap-2">
                      {[1, 2, 3, 4, 5].map((star) => (
                        <button
                          key={star}
                          type="button"
                          onClick={() => setNewReviewRating(star)}
                          className={`text-3xl transition-transform hover:scale-110 focus:outline-none ${star <= newReviewRating ? 'text-yellow-400 drop-shadow-sm' : 'text-gray-300'
                            }`}
                        >
                          ★
                        </button>
                      ))}
                    </div>
                  </div>

                  <div className="mb-2">
                    <label className="block text-sm font-medium text-gray-700 mb-2">Текст отзыва</label>
                    <textarea
                      value={newReviewText}
                      onChange={handleTextChange}
                      className={`w-full p-4 border rounded-xl! focus:ring-2 focus:border-transparent text-sm min-h-[100px] resize-y transition-all duration-200 ${validationError
                          ? 'border-red-300 bg-red-50 focus:ring-red-200'
                          : 'border-gray-200 focus:ring-rose-200'
                        }`}
                    />
                    {validationError && (
                      <p className="mt-2 text-sm text-red-600">{validationError}</p>
                    )}
                  </div>

                  <div className="mt-6 flex justify-end gap-3">
                    <button
                      type="button"
                      onClick={cancelEdit}
                      disabled={isSubmitting}
                      className="px-6 py-2.5 text-sm font-medium text-gray-600 bg-gray-100 rounded-xl hover:bg-gray-200 transition-colors disabled:opacity-50"
                    >
                      Отмена
                    </button>
                    <Button
                      type="submit"
                      isLoading={isSubmitting}
                      variant="primary"
                      className="rounded-xl! px-8 py-2.5"
                      onClick={handleSubmit}
                    >
                      Сохранить изменения
                    </Button>
                  </div>
                </>
              )}
            </>
          )}
        </div>
      )}

      {/* Слайдер чужих отзывов */}
      <div className="bg-white rounded-3xl p-6 sm:p-10 shadow-lg border border-gray-100 relative overflow-hidden">
        {/* Заголовок и кнопки навигации в одной сетке */}
        <div className="grid grid-cols-[auto_1fr_auto] items-center gap-4 mb-8">

          {/* Кнопка "Назад" (скрыта на мобильных, если отзывов мало) */}
          <div className={`${reviewsArray.length <= 1 ? 'invisible' : ''}`}>
            <button
              onClick={prevSlide}
              className="p-3 rounded-full border border-gray-200 text-gray-600 hover:bg-rose-50 hover:text-rose-600 hover:border-rose-200 transition-all shadow-sm"
              aria-label="Предыдущий отзыв"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
            </button>
          </div>

          {/* Заголовок по центру */}
          <div className="text-center">
            <h3 className="text-2xl font-serif font-bold text-gray-900">Отзывы клиентов</h3>
            
          </div>

          {/* Кнопка "Вперед" (скрыта на мобильных, если отзывов мало) */}
          <div className={`${reviewsArray.length <= 1 ? 'invisible' : ''}`}>
            <button
              onClick={nextSlide}
              className="p-3 rounded-full border border-gray-200 text-gray-600 hover:bg-rose-50 hover:text-rose-600 hover:border-rose-200 transition-all shadow-sm"
              aria-label="Следующий отзыв"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        </div>

        {/* Контент слайдера */}
        {reviewsArray.length === 0 ? (
          <div className="text-center py-12 text-gray-500 bg-gray-50 rounded-2xl border border-dashed border-gray-200">
            <div className="text-4xl mb-3">💬</div>
            <p className="font-medium">Отзывов пока нет</p>
            <p className="text-sm mt-1 text-gray-400">Будьте первым!</p>
          </div>
        ) : (
          <div className="relative min-h-[220px] flex flex-col justify-center">
            <div className="transition-all duration-500 ease-in-out transform">
              <div className="flex flex-col items-center text-center max-w-3xl mx-auto px-4">
                <div className="mb-6">
                  <div className="w-16 h-16 bg-gradient-to-br from-rose-100 to-rose-200 rounded-full flex items-center justify-center text-rose-700 font-bold text-2xl shadow-md mx-auto mb-3 border-4 border-white">
                    {reviewsArray[currentSlide].author_name?.charAt(0).toUpperCase() || 'U'}
                  </div>
                  <h4 className="font-bold text-gray-900 text-lg">{reviewsArray[currentSlide].author_name || 'Анонимный пользователь'}</h4>
                  <div className="mt-2 flex justify-center">{renderStars(reviewsArray[currentSlide].rating, 'sm')}</div>
                </div>
                <blockquote className="relative">
                  <span className="absolute -top-6 -left-4 text-6xl text-rose-100 font-serif leading-none select-none">"</span>
                  <p className="text-gray-700 leading-relaxed text-lg sm:text-xl italic font-medium px-4">{reviewsArray[currentSlide].text}</p>
                  <span className="absolute -bottom-8 -right-4 text-6xl text-rose-100 font-serif leading-none select-none">"</span>
                </blockquote>
                <div className="mt-8 pt-6 border-t border-gray-100 w-full">
                  <span className="text-xs font-medium text-gray-400 uppercase tracking-wider">
                    {new Date(reviewsArray[currentSlide].created_at).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })}
                  </span>
                </div>
              </div>
            </div>

            {/* Индикаторы (точки) */}
            {reviewsArray.length > 1 && (
              <div className="flex justify-center gap-2 mt-8">
                {reviewsArray.map((_, idx) => (
                  <button
                    key={idx}
                    onClick={() => setCurrentSlide(idx)}
                    className={`h-2 rounded-2xl! transition-all duration-300 ${idx === currentSlide ? 'bg-rose-500 w-8' : 'bg-gray-200 w-2 hover:bg-gray-300'}`}
                    aria-label={`Перейти к отзыву ${idx + 1}`}
                  />
                ))}
              </div>
            )}
          </div>
        )}
      </div>

      {/* Модальное окно успеха */}
      {showSuccessModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm animate-fadeIn">
          <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full p-8 text-center transform transition-all scale-100 animate-scaleIn">
            <div className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
              <svg className="w-10 h-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            </div>

            <h3 className="text-2xl font-serif font-bold text-gray-900 mb-3">
              {myReview && isEditing ? 'Изменения сохранены!' : 'Спасибо за отзыв!'}
            </h3>
            <p className="text-gray-600 mb-8 leading-relaxed">
              {myReview && isEditing
                ? 'Ваш отзыв обновлен и снова отправлен на модерацию.'
                : 'Ваш отзыв успешно отправлен. Он появится после проверки модератором.'}
            </p>

            <button
              onClick={() => setShowSuccessModal(false)}
              className="w-full py-3 px-6 bg-rose-500 text-white font-bold rounded-2xl! hover:bg-rose-600 transition-colors shadow-lg"
            >
              Закрыть
            </button>
          </div>
        </div>
      )}
    </div>
  );
};
