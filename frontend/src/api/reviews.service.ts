import apiClient from './axiosInstance';
import { Review } from './types';

export const reviewsService = {
  getApprovedReviews: async (): Promise<Review[]> => {
    const response = await apiClient.get<Review[]>('/reviews', { params: { status: 'approved' } });
    return response.data;
  },

  submitReview: async (text: string, rating: number): Promise<Review> => {
    const response = await apiClient.post<Review>('/reviews', { text, rating });
    return response.data;
  },
  
  // Только для админа
  getPendingReviews: async (): Promise<Review[]> => {
    const response = await apiClient.get<Review[]>('/reviews/admin/pending');
    return response.data;
  },

  approveReview: async (id: string): Promise<void> => {
    await apiClient.patch(`/reviews/admin/${id}/approve`);
  }
};