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
  
  // Состояние для модального окна успеха
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  
  // Состояние для ошибки валидации
  const [validationError, setValidationError] = useState<string | null>(null);

  const [currentSlide, setCurrentSlide] = useState(0);

  const reviewsArray = reviews || [];

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
    
    // Сброс ошибок перед новой проверкой
    setValidationError(null);

    if (!isAuthenticated) {
      alert('Пожалуйста, войдите в систему, чтобы оставить отзыв');
      return;
    }

    // Проверка на минимальную длину (10 символов)
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
      
      // Очистка формы
      setNewReviewText('');
      setNewReviewRating(5);
      
      // Показываем модальное окно
      setShowSuccessModal(true);
      
    } catch (err: any) {
      setValidationError(err.message || 'Ошибка при отправке отзыва. Попробуйте позже.');
      alert(err.message || 'Ошибка при отправке отзыва');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setNewReviewText(e.target.value);
    if (validationError) {
      setValidationError(null);
    }
  };

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

          <div className="mb-2">
            <div className="flex justify-between items-center mb-2">
              <label className="block text-sm font-medium text-gray-700">
                Текст отзыва
              </label>
              <span className="text-xs text-gray-500 font-medium">
                Минимум 10 символов
              </span>
            </div>
            
            <textarea
              value={newReviewText}
              onChange={handleTextChange}
              placeholder="Поделитесь своими впечатлениями о курсе..."
              className={`w-full p-4 border rounded-xl! focus:ring-2 focus:border-transparent text-sm min-h-[100px] resize-y transition-all duration-200 ${
                validationError 
                  ? 'border-red-300 bg-red-50 focus:ring-red-200 text-red-900 placeholder-red-300' 
                  : 'border-gray-200 focus:ring-rose-200'
              }`}
              required
            />
            
            {validationError && (
              <p className="mt-2 text-sm text-red-600 flex items-center gap-1 animate-pulse">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {validationError}
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
        </div>
      )}

      {!isAuthenticated && (
        <div className="bg-blue-50/50 p-6 rounded-2xl! text-center text-sm text-blue-800 border border-blue-100">
          <p className="mb-2">Хотите поделиться опытом?</p>
          <a href="/login" className="font-bold underline hover:text-blue-900 transition-colors">Войдите</a>, чтобы оставить отзыв
        </div>
      )}

      {/* Секция со слайдером отзывов */}
      <div className="bg-white rounded-3xl p-6 sm:p-10 shadow-lg border border-gray-100 relative overflow-hidden">
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

          {reviewsArray.length > 1 && (
            <div className="flex gap-2">
              <button onClick={prevSlide} className="p-3 rounded-full border border-gray-200 text-gray-600 hover:bg-rose-50 hover:text-rose-600 hover:border-rose-200 transition-all shadow-sm" aria-label="Предыдущий">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" /></svg>
              </button>
              <button onClick={nextSlide} className="p-3 rounded-full border border-gray-200 text-gray-600 hover:bg-rose-50 hover:text-rose-600 hover:border-rose-200 transition-all shadow-sm" aria-label="Следующий">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" /></svg>
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
            <div className="transition-all duration-500 ease-in-out transform">
              <div className="flex flex-col items-center text-center max-w-3xl mx-auto">
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
            {reviewsArray.length > 1 && (
              <div className="flex justify-center gap-2 mt-8">
                {reviewsArray.map((_, idx) => (
                  <button key={idx} onClick={() => setCurrentSlide(idx)} className={`h-2 rounded-2xl! transition-all duration-300 ${idx === currentSlide ? 'bg-rose-500 w-8' : 'bg-gray-200 w-2 hover:bg-gray-300'}`} aria-label={`Перейти к отзыву ${idx + 1}`} />
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
            
            <h3 className="text-2xl font-serif font-bold text-gray-900 mb-3">Спасибо за отзыв!</h3>
            <p className="text-gray-600 mb-8 leading-relaxed">
              Ваш отзыв успешно отправлен. Он появится на странице курса после прохождения модерации.
            </p>
            
            <button
              onClick={() => {
                setShowSuccessModal(false);
                onReviewSubmitted();
              }}
              className="w-full py-3 px-6 bg-rose-500 text-white font-bold rounded-xl hover:bg-rose-600 transition-colors shadow-lg hover:shadow-xl transform hover:-translate-y-0.5"
            >
              Закрыть
            </button>
          </div>
        </div>
      )}
    </div>
  );
};