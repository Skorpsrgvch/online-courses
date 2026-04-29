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
  BonusItem
} from './types';

// --- Мапперы (Преобразование данных с бэкенда) ---

const mapCourse = (data: any): Course => {

  return {
    id: data.ID,
    title: data.Title,
    description: data.Description,
    is_public: data.IsPublic,
    price: data.Price,
    author_id: data.AuthorID,
    is_active: data.IsActive,
    cover_image_url: data.CoverImageURL || '',
    contraindications: data.contraindications || '',
    recommendations: data.recommendations || '',
    target_audience: data.target_audience || '',
    course_basis: data.course_basis || '',
    class_basis: data.class_basis || '',
    bonuses: data.bonuses || [],
    is_purchased: data.is_purchased,
    progress: data.progress,
  };
};

const mapLesson = (data: any): FullCourseLesson => ({

  id: data.id,
  module_id: data.module_id,
  title: data.title ,
  description: data.description || '',
  video_embed_id: data.video_embed_id,
  private_key: data.private_key,
  order: data.order,
});

const mapModule = (data: any): FullCourseModule => ({
  id: data.id ,
  course_id: data.course_id,
  title: data.title ,
  order: data.order ,
  lessons: data.lessons ? data.lessons.map(mapLesson) : [],
});

// --- Сервис ---

export const coursesService = {
  getAllCourses: async (): Promise<Course[]> => {
    const response = await apiClient.get<any[]>('/courses');
    return response.data.map(mapCourse);
  },

  getCourseById: async (id: number): Promise<Course> => {
    const response = await apiClient.get<any>(`/courses/${id}`);
    return mapCourse(response.data);
  },

  getCourseFull: async (id: number): Promise<CourseFullResponse> => {
    const response = await apiClient.get<any>(`/courses/${id}/full`);
    const data = response.data;
    
    const courseData = data.course || data;
    const modulesData = data.modules || [];

    return {
      course: mapCourse(courseData),
      modules: modulesData.map(mapModule),
    };
  },

  createCourse: async (data: CreateCourseDto): Promise<void> => {
    await apiClient.post('/courses', data);
  },

  createCourseWithModules: async (data: any): Promise<void> => {
    await apiClient.post('/courses/with-modules', data);
  },

  updateCourse: async (id: number, data: Partial<CreateCourseDto>): Promise<void> => {
    await apiClient.put(`/courses/${id}`, data);
  },

  deleteCourse: async (id: number): Promise<void> => {
    await apiClient.delete(`/courses/${id}`);
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
    await apiClient.post('/modules', data);
  },

  updateModule: async (id: number, data: { title: string; order: number }): Promise<void> => {
    await apiClient.put(`/modules/${id}`, data);
  },

  deleteModule: async (id: number): Promise<void> => {
    await apiClient.delete(`/modules/${id}`);
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
    await apiClient.post('/lessons', data);
  },

  updateLesson: async (id: number, data: Partial<CreateLessonDto>): Promise<void> => {
    await apiClient.put(`/lessons/${id}`, data);
  },

  deleteLesson: async (id: number): Promise<void> => {
    await apiClient.delete(`/lessons/${id}`);
  },

  markLessonComplete: async (lessonId: number): Promise<void> => {
    await apiClient.post(`/progress/lessons/${lessonId}/mark`);
  },

  purchaseCourse: async (courseId: number) => {
    const response = await apiClient.post(`/courses/${courseId}/purchase`);
    return response.data;
  },
};