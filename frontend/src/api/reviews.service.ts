import apiClient from './axiosInstance';
import type { Review, CreateReviewDto } from './types';

export interface ReviewsListResponse {
  reviews: Review[];
}

// Функция маппинга данных с бэкенда в формат фронтенда
const mapReview = (data: any): Review => {
  return {
    id: data.ID || data.id,
    user_id:  data.user_id,
    course_id: data.course_id,
    text:  data.text,
    rating:  data.rating,
    approved:  data.approved,
    created_at:  data.created_at,
    author_name:  data.author_name || 'Аноним', 
    course_title:  data.course_title, 
  };
};

export const reviewsService = {
  
  getApprovedReviews: async (): Promise<Review[]> => {
    return [];
  },

  submitReview: async (data: CreateReviewDto): Promise<void> => {
    await apiClient.post('/reviews', data);
  },

  approveReview: async (id: number): Promise<void> => {
    await apiClient.post(`/reviews/${id}/approve`);
  },

  rejectReview: async (id: number): Promise<void> => {
    await apiClient.delete(`/reviews/${id}`);
  },

  // Получить одобренные отзывы по курсу
  getCourseReviews: async (courseId: number): Promise<Review[]> => {
    const response = await apiClient.get<ReviewsListResponse>(`/courses/${courseId}/reviews`);
    
    // Применяем маппинг к каждому элементу массива
    if (Array.isArray(response.data.reviews)) {
      return response.data.reviews.map(mapReview);
    }
    // Если бэкенд вернул объект без ключа reviews (редко, но бывает)
    if (Array.isArray(response.data)) {
      return response.data.map(mapReview);
    }
    
    return [];
  },

  // Получить ожидающие модерации отзывы (только admin)
  getPendingReviews: async (): Promise<Review[]> => {
    const response = await apiClient.get<ReviewsListResponse>('/reviews/admin/pending');
    
    if (Array.isArray(response.data.reviews)) {
      return response.data.reviews.map(mapReview);
    }
    if (Array.isArray(response.data)) {
      return response.data.map(mapReview);
    }

    return [];
  },
};