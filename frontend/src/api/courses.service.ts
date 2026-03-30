import apiClient from './axiosInstance';
import { Course, Lesson } from './types';

export const coursesService = {
  getAllCourses: async (): Promise<Course[]> => {
    const response = await apiClient.get<Course[]>('/courses');
    return response.data;
  },

  getCourseById: async (id: string): Promise<Course> => {
    const response = await apiClient.get<Course>(`/courses/${id}`);
    return response.data;
  },

  getLessonById: async (courseId: string, lessonId: string): Promise<Lesson> => {
    const response = await apiClient.get<Lesson>(`/courses/${courseId}/lessons/${lessonId}`);
    return response.data;
  },

  markLessonComplete: async (courseId: string, lessonId: string): Promise<void> => {
    await apiClient.post(`/courses/${courseId}/lessons/${lessonId}/complete`);
  },

  purchaseCourse: async (courseId: string): Promise<void> => {
    // Здесь может быть интеграция с платежной системой
    await apiClient.post(`/courses/${courseId}/purchase`);
  }
};