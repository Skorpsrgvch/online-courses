import React, { useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import { reviewsService } from '../../api';
import type { Review } from '../../api/types';
import { timeAgo } from '../../utils/dateUtils';

interface CommentsProps {
  courseId: number;
  reviews?: Review[];
  onReviewSubmitted?: () => void;
}

/**
 * Компонент комментариев/отзывов к курсу.
 */
export const Comments: React.FC<CommentsProps> = ({ courseId, reviews = [], onReviewSubmitted }) => {
  const { isAuthenticated } = useAuth();
  const [text, setText] = useState('');
  const [rating, setRating] = useState(5);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!text.trim()) return;

    setIsSubmitting(true);
    setError(null);
    try {
      await reviewsService.submitReview({
        course_id: courseId,
        text: text.trim(),
        rating,
      });
      setText('');
      setRating(5);
      onReviewSubmitted?.();
    } catch (err: any) {
      setError(err.message || 'Не удалось отправить отзыв');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Форма отзыва */}
      {isAuthenticated && (
        <form onSubmit={handleSubmit} className="bg-white p-5 rounded-xl shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold text-gray-800 mb-3">Оставить отзыв</h3>

          {/* Рейтинг */}
          <div className="flex items-center gap-1 mb-3">
            {[1, 2, 3, 4, 5].map((star) => (
              <button
                key={star}
                type="button"
                onClick={() => setRating(star)}
                className={`text-2xl ${star <= rating ? 'text-yellow-400' : 'text-gray-300'}`}
              >
                ★
              </button>
            ))}
          </div>

          {/* Текст отзыва */}
          <textarea
            value={text}
            onChange={(e) => setText(e.target.value)}
            placeholder="Расскажите о вашем опыте..."
            className="w-full p-3 border border-gray-200 rounded-lg resize-none focus:ring-2 focus:ring-rose-300 focus:border-rose-400"
            rows={4}
            minLength={10}
            required
          />

          {error && (
            <div className="mt-2 text-sm text-red-600">{error}</div>
          )}

          <button
            type="submit"
            disabled={isSubmitting || text.trim().length < 10}
            className="mt-3 px-6 py-2 bg-rose-500 text-white rounded-lg hover:bg-rose-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {isSubmitting ? 'Отправка...' : 'Отправить отзыв'}
          </button>
        </form>
      )}

      {/* Список отзывов */}
      <div className="space-y-4">
        {reviews.length === 0 ? (
          <p className="text-gray-500 text-center py-6">Отзывов пока нет. Будьте первым!</p>
        ) : (
          reviews.map((review) => (
            <div key={review.id} className="bg-white p-4 rounded-xl border border-gray-100">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-1">
                  {[1, 2, 3, 4, 5].map((star) => (
                    <span key={star} className={star <= review.rating ? 'text-yellow-400' : 'text-gray-300'}>
                      ★
                    </span>
                  ))}
                </div>
                <span className="text-xs text-gray-400">{timeAgo(review.created_at)}</span>
              </div>
              <p className="text-gray-700 text-sm">{review.text}</p>
            </div>
          ))
        )}
      </div>
    </div>
  );
};
