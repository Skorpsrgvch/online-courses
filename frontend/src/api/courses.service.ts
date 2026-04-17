import apiClient from './axiosInstance';
import type { Course, CourseModule, Lesson, CreateCourseDto, CreateModuleDto, CreateLessonDto, CourseFullResponse } from './types';

export const coursesService = {
  getAllCourses: async (): Promise<Course[]> => {
    const response = await apiClient.get<Course[]>('/courses');
    return response.data;
  },

  getCourseById: async (id: number): Promise<Course> => {
    const response = await apiClient.get<Course>(`/courses/${id}`);
    return response.data;
  },

  getCourseFull: async (id: number): Promise<CourseFullResponse> => {
    const response = await apiClient.get<CourseFullResponse>(`/courses/${id}/full`);
    return response.data;
  },

  // Создание курса (только admin)
  createCourse: async (data: CreateCourseDto): Promise<void> => {
    await apiClient.post('/courses', data);
  },

  // Создание курса с модулями и уроками одной транзакцией (только admin)
  createCourseWithModules: async (data: {
    title: string;
    description: string;
    is_public: boolean;
    price: number;
    cover_image_url?: string;
    modules: Array<{
      title: string;
      order: number;
      lessons: Array<{
        title: string;
        description: string;
        lesson_type: 'video' | 'article';
        video_embed_id?: string;
        article_content?: string;
        order: number;
      }>;
    }>;
  }): Promise<void> => {
    await apiClient.post('/courses/with-modules', data);
  },

  // Обновление курса (только admin)
  updateCourse: async (id: number, data: Partial<CreateCourseDto>): Promise<void> => {
    await apiClient.put(`/courses/${id}`, data);
  },

  // Удаление курса (только admin)
  deleteCourse: async (id: number): Promise<void> => {
    await apiClient.delete(`/courses/${id}`);
  },

  // === Модули ===

  getCourseModules: async (courseId: number): Promise<CourseModule[]> => {
    const response = await apiClient.get<CourseModule[]>(`/courses/${courseId}/modules`);
    return response.data;
  },

  createModule: async (data: CreateModuleDto): Promise<void> => {
    await apiClient.post('/modules', data);
  },

  updateModule: async (id: number, data: { title: string; order: number }): Promise<void> => {
    await apiClient.put(`/modules/${id}`, data);
  },

  deleteModule: async (id: number): Promise<void> => {
    await apiClient.delete(`/modules/${id}`);
  },

  // === Уроки ===

  getModuleLessons: async (moduleId: number): Promise<Lesson[]> => {
    const response = await apiClient.get<{ lessons: Lesson[] }>(`/modules/${moduleId}/lessons`);
    return response.data.lessons;
  },

  createLesson: async (data: CreateLessonDto): Promise<void> => {
    await apiClient.post('/lessons', data);
  },

  updateLesson: async (id: number, data: Partial<CreateLessonDto>): Promise<void> => {
    await apiClient.put(`/lessons/${id}`, data);
  },

  deleteLesson: async (id: number): Promise<void> => {
    await apiClient.delete(`/lessons/${id}`);
  },

  // === Прогресс ===

  markLessonComplete: async (lessonId: number): Promise<void> => {
    await apiClient.post(`/progress/lessons/${lessonId}/mark`);
  },
};
