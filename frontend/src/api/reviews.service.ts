import apiClient from './axiosInstance';
import type { Review, CreateReviewDto } from './types';

export interface ReviewsListResponse {
  reviews: Review[];
}

export const reviewsService = {
  // ВНИМАНИЕ: Эндпоинт GET /reviews отсутствует на бэкенде.
  // Заглушка — вернуть пустой массив до реализации бэкенда.
  getApprovedReviews: async (): Promise<Review[]> => {
    // TODO: реализовать GET /reviews?approved=true на бэкенде
    return [];
  },

  submitReview: async (data: CreateReviewDto): Promise<void> => {
    await apiClient.post('/reviews', data);
  },

  // Одобрение отзыва (только admin)
  approveReview: async (id: number): Promise<void> => {
    await apiClient.post(`/reviews/${id}/approve`);
  },

  // Отклонение отзыва (только admin)
  rejectReview: async (id: number): Promise<void> => {
    await apiClient.delete(`/reviews/${id}`);
  },

  // Получить одобренные отзывы по курсу
  getCourseReviews: async (courseId: number): Promise<Review[]> => {
    const response = await apiClient.get<ReviewsListResponse>(`/courses/${courseId}/reviews`);
    return response.data.reviews;
  },

  // Получить ожидающие модерации отзывы (только admin)
  getPendingReviews: async (): Promise<Review[]> => {
    const response = await apiClient.get<ReviewsListResponse>('/reviews/admin/pending');
    return response.data.reviews;
  },
};
