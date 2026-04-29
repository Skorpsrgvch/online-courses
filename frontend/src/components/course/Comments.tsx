import React, { useState, useMemo } from 'react';
import { useAuth } from '../../context/AuthContext';
import { reviewsService } from '../../api/reviews.service';
import type { Review } from '../../api/types';
import { Button } from '../ui/Button';

interface CommentsProps {
  courseId: number;
  reviews: Review[] | null;
  onReviewSubmitted: () => void;
}

export const Comments: React.FC<CommentsProps> = ({ courseId, reviews, onReviewSubmitted }) => {
  const { isAuthenticated } = useAuth();
  const [newReviewText, setNewReviewText] = useState('');
  const [newReviewRating, setNewReviewRating] = useState(5);
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  const [currentSlide, setCurrentSlide] = useState(0);

  const reviewsArray = reviews || [];

  // Вычисляем средний рейтинг один раз при изменении отзывов
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
    if (!isAuthenticated) {
      alert('Пожалуйста, войдите в систему, чтобы оставить отзыв');
      return;
    }
    if (!newReviewText.trim()) return;

    setIsSubmitting(true);
    try {
      await reviewsService.submitReview({
        course_id: courseId,
        text: newReviewText,
        rating: newReviewRating,
      });
      setNewReviewText('');
      setNewReviewRating(5);
      onReviewSubmitted();
    } catch (err: any) {
      alert(err.message || 'Ошибка при отправке отзыва');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Компонент для отрисовки звезд (рейтинга)
  const renderStars = (rating: number, size = 'sm') => {
    const starSize = size === 'lg' ? 'text-xl' : 'text-sm';
    return (
      <div className={`flex ${starSize} text-yellow-400`}>
        {[1, 2, 3, 4, 5].map((star) => (
          <span key={star} className={star <= rating ? 'text-yellow-400' : 'text-gray-300'}>
            ★
          </span>
        ))}
      </div>
    );
  };

  return (
    <div className="space-y-8">
      {/* Форма отзыва */}
      {isAuthenticated && (
        <div className="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
          <h4 className="text-lg font-bold text-gray-800 mb-4">Оставить отзыв</h4>
          
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">Ваша оценка</label>
            <div className="flex gap-2">
              {[1, 2, 3, 4, 5].map((star) => (
                <button
                  key={star}
                  type="button"
                  onClick={() => setNewReviewRating(star)}
                  className={`text-3xl transition-transform hover:scale-110 focus:outline-none ${
                    star <= newReviewRating ? 'text-yellow-400 drop-shadow-sm' : 'text-gray-300'
                  }`}
                >
                  ★
                </button>
              ))}
            </div>
          </div>

          <textarea
            value={newReviewText}
            onChange={(e) => setNewReviewText(e.target.value)}
            placeholder="Поделитесь своими впечатлениями о курсе..."
            className="w-full p-4 border border-gray-200 rounded-xl focus:ring-2 focus:ring-rose-500 focus:border-transparent text-sm min-h-[100px] resize-y transition-shadow"
            required
          />
          
          <div className="mt-4 flex justify-end">
            <Button type="submit" isLoading={isSubmitting} variant="primary" className="rounded-xl px-8 py-2.5">
              Опубликовать
            </Button>
          </div>
        </div>
      )}

      {!isAuthenticated && (
        <div className="bg-blue-50/50 p-6 rounded-2xl text-center text-sm text-blue-800 border border-blue-100">
          <p className="mb-2">Хотите поделиться опытом?</p>
          <a href="/login" className="font-bold underline hover:text-blue-900 transition-colors">Войдите</a>, чтобы оставить отзыв
        </div>
      )}

      {/* Секция со слайдером отзывов */}
      <div className="bg-white rounded-3xl p-6 sm:p-10 shadow-lg border border-gray-100 relative overflow-hidden">
        
        {/* Заголовок и статистика */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between mb-8 gap-4">
          <div>
            <h3 className="text-2xl font-serif font-bold text-gray-900">Отзывы студентов</h3>
            {reviewsArray.length > 0 && (
              <div className="flex items-center gap-3 mt-2">
                <div className="flex items-center gap-1 bg-yellow-50 px-3 py-1 rounded-full border border-yellow-100">
                  <span className="text-yellow-500 text-lg">★</span>
                  <span className="font-bold text-gray-900 text-lg">{averageRating}</span>
                  <span className="text-gray-500 text-sm">из 5</span>
                </div>
                <span className="text-gray-400 text-sm">на основе {reviewsArray.length} отзывов</span>
              </div>
            )}
          </div>

          {/* Кнопки навигации (десктоп) */}
          {reviewsArray.length > 1 && (
            <div className="flex gap-2">
              <button
                onClick={prevSlide}
                className="p-3 rounded-full border border-gray-200 text-gray-600 hover:bg-rose-50 hover:text-rose-600 hover:border-rose-200 transition-all shadow-sm"
                aria-label="Предыдущий"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                </svg>
              </button>
              <button
                onClick={nextSlide}
                className="p-3 rounded-full border border-gray-200 text-gray-600 hover:bg-rose-50 hover:text-rose-600 hover:border-rose-200 transition-all shadow-sm"
                aria-label="Следующий"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </button>
            </div>
          )}
        </div>

        {reviewsArray.length === 0 ? (
          <div className="text-center py-12 text-gray-500 bg-gray-50 rounded-2xl border border-dashed border-gray-200">
            <div className="text-4xl mb-3">💬</div>
            <p className="font-medium">Отзывов пока нет</p>
            <p className="text-sm mt-1 text-gray-400">Будьте первым, кто оставит свой комментарий!</p>
          </div>
        ) : (
          <div className="relative min-h-[220px] flex flex-col justify-center">
            
            {/* Контент текущего отзыва */}
            <div className="transition-all duration-500 ease-in-out transform">
              <div className="flex flex-col items-center text-center max-w-3xl mx-auto">
                
                {/* Аватар и Имя */}
                <div className="mb-6">
                  <div className="w-16 h-16 bg-gradient-to-br from-rose-100 to-rose-200 rounded-full flex items-center justify-center text-rose-700 font-bold text-2xl shadow-md mx-auto mb-3 border-4 border-white">
                    {reviewsArray[currentSlide].author_name?.charAt(0).toUpperCase() || 'U'}
                  </div>
                  <h4 className="font-bold text-gray-900 text-lg">
                    {reviewsArray[currentSlide].author_name || 'Анонимный пользователь'}
                  </h4>
                  <div className="mt-2 flex justify-center">
                    {renderStars(reviewsArray[currentSlide].rating, 'sm')}
                  </div>
                </div>

                {/* Текст отзыва */}
                <blockquote className="relative">
                  <span className="absolute -top-6 -left-4 text-6xl text-rose-100 font-serif leading-none select-none">"</span>
                  <p className="text-gray-700 leading-relaxed text-lg sm:text-xl italic font-medium px-4">
                    {reviewsArray[currentSlide].text}
                  </p>
                  <span className="absolute -bottom-8 -right-4 text-6xl text-rose-100 font-serif leading-none select-none">"</span>
                </blockquote>

                {/* Дата */}
                <div className="mt-8 pt-6 border-t border-gray-100 w-full">
                  <span className="text-xs font-medium text-gray-400 uppercase tracking-wider">
                    {new Date(reviewsArray[currentSlide].created_at).toLocaleDateString('ru-RU', {
                      day: 'numeric', month: 'long', year: 'numeric'
                    })}
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
                    className={`h-2 rounded-full transition-all duration-300 ${
                      idx === currentSlide 
                        ? 'bg-rose-500 w-8' 
                        : 'bg-gray-200 w-2 hover:bg-gray-300'
                    }`}
                    aria-label={`Перейти к отзыву ${idx + 1}`}
                  />
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};