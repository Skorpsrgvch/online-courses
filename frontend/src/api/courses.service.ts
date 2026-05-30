import apiClient from './axiosInstance';
import type {
  Course,
  CourseModule,
  Lesson,
  CreateCourseDto,
  CreateModuleDto,
  CourseFullResponse,
  FullCourseModule,
  FullCourseLesson,
  CreateLessonDto,
  UpdateFullCourseDto,
  PaymentResponse
} from './types';



const mapCourse = (data: any): Course => {
  return {
    id: data.id,
    title: data.title,
    description: data.description,
    is_public: data.is_public,
    price: data.price,
    author_id: data.author_id,
    is_active: data.is_active,
    cover_image_url: data.cover_image_url || '',
    contraindications: data.contraindications || '',
    recommendations: data.recommendations || '',
    target_audience: data.target_audience || '',
    course_basis: data.course_basis || '',
    class_basis: data.class_basis || '',
    bonuses: data.bonuses || [],

    is_purchased: data.is_purchased || false,
    progress_percent: data.progress_percent || 0,
    modules: data.modules || [],

    is_access_expired: data.is_access_expired,
    days_remaining: data.days_remaining,
    access_expires_at: data.access_expires_at,
  };
};

const mapLesson = (data: any): FullCourseLesson => ({
  id: data.id,
  module_id: data.module_id,
  title: data.title,
  description: data.description || '',
  video_embed_id: data.video_embed_id,
  private_key: data.private_key,
  order: data.order,
  is_completed: data.is_completed,
});

const mapModule = (data: any): FullCourseModule => ({
  id: data.id,
  course_id: data.course_id,
  title: data.title,
  order: data.order,
  is_locked: data.is_locked,
  unlock_date: data.unlock_date || null,
  lessons: data.lessons ? data.lessons.map(mapLesson) : [],
});


export const coursesService = {


  getAllCourses: async (): Promise<Course[]> => {
    const response = await apiClient.get<any[]>('/courses');
    return response.data.map(mapCourse);
  },

  getAllCoursesAdmin: async (): Promise<Course[]> => {
    const response = await apiClient.get<{ courses: any[] }>('/admin/courses/all');
    return response.data.courses.map(mapCourse);
  },


  getCourseFull: async (id: number): Promise<CourseFullResponse> => {
    const response = await apiClient.get<any>(`/courses/${id}/full`);
    const data = response.data;

    const courseData = data.course || data;
    const modulesData = data.modules || [];

    const mappedCourse = mapCourse(courseData);

    mappedCourse.is_access_expired = data.is_access_expired;
    mappedCourse.days_remaining = data.days_remaining;
    mappedCourse.access_expires_at = data.access_expires_at;

    return {
      course: mapCourse(courseData),
      modules: modulesData.map(mapModule),
      is_access_expired: data.is_access_expired,
      days_remaining: data.days_remaining,
      access_expires_at: data.access_expires_at,
    };
  },

  createCourse: async (data: CreateCourseDto): Promise<void> => {
    await apiClient.post('/admin/courses', data);
  },


  createCourseWithModules: async (data: any): Promise<void> => {
    await apiClient.post('/admin/courses/with-modules', data);
  },


  updateCourse: async (id: number, data: Partial<UpdateFullCourseDto>): Promise<void> => {
    await apiClient.put(`/admin/courses/${id}`, data);
  },


  updateFullCourse: async (id: number, data: UpdateFullCourseDto): Promise<void> => {
    await apiClient.put(`/admin/courses/${id}/full-update`, data);
  },


  toggleCourseStatus: async (id: number, isActive: boolean): Promise<void> => {
    await apiClient.patch(`/admin/courses/${id}/status`, { is_active: isActive });
  },


  getCourseModules: async (courseId: number): Promise<CourseModule[]> => {
    const response = await apiClient.get<any[]>(`/courses/${courseId}/modules`);
    return response.data.map((m: any) => ({
      id: m.ID,
      course_id: m.CourseID,
      title: m.Title,
      order: m.Order,
    }));
  },

  createModule: async (data: CreateModuleDto): Promise<void> => {
    await apiClient.post('/admin/modules', data);
  },

  updateModule: async (id: number, data: { title: string; order: number }): Promise<void> => {
    await apiClient.put(`/admin/modules/${id}`, data);
  },

  deleteModule: async (id: number): Promise<void> => {
    await apiClient.delete(`/admin/modules/${id}`);
  },

  getModuleLessons: async (moduleId: number): Promise<Lesson[]> => {
    const response = await apiClient.get<any>(`/modules/${moduleId}/lessons`);
    const lessonsData = response.data.lessons || response.data;

    if (!Array.isArray(lessonsData)) return [];

    return lessonsData.map((l: any) => ({
      id: l.ID,
      module_id: l.ModuleID,
      title: l.Title,
      description: l.Description || '',
      video_embed_id: l.VideoEmbedID || '',
      private_key: l.PrivateKey || null,
      order: l.Order,
    }));
  },

  createLesson: async (data: CreateLessonDto): Promise<void> => {
    await apiClient.post('/admin/lessons', data);
  },

  updateLesson: async (id: number, data: Partial<CreateLessonDto>): Promise<void> => {
    await apiClient.put(`/admin/lessons/${id}`, data);
  },

  deleteLesson: async (id: number): Promise<void> => {
    await apiClient.delete(`/admin/lessons/${id}`);
  },

  reorderLessons: async (moduleId: number, lessonIds: number[]): Promise<void> => {
    // Маршрут теперь ожидает :id, который мы передаем как moduleId
    await apiClient.put(`/admin/modules/${moduleId}/lessons/reorder`, { lesson_ids: lessonIds });
  },

  reorderModules: async (courseId: number, moduleIds: number[]): Promise<void> => {
    await apiClient.put(`/admin/courses/${courseId}/modules/reorder`, { module_ids: moduleIds });
  },

  /**  Прогресс и покупка  */
  markLessonComplete: async (lessonId: number): Promise<void> => {
    await apiClient.post(`/progress/lessons/${lessonId}/mark`);
  },

  createPayment: async (courseId: number, returnUrl: string): Promise<PaymentResponse> => {
    const response = await apiClient.post<PaymentResponse>('/payments', {
      course_id: courseId,
      return_url: returnUrl,
    });
    return response.data;
  },

  enrollFree: async (courseId: number, price: number): Promise<void> => {
  
  const response = await apiClient.post(`/courses/${courseId}/enroll-free`, {
    price: price 
  });
  
  return response.data;
},
};
