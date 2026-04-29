export interface User {
  id: number;
  email: string;
  name: string;
  role: 'user' | 'admin';
  created_at?: string;
}


export interface BonusItem {
  title: string;
  description: string;
  icon: string; 
}

export interface CourseModule {
  id: number;
  course_id: number;
  title: string;
  order: number;
}

export interface Lesson {
  id: number;
  module_id: number;
  title: string;
  description: string;
  video_embed_id: string;
  private_key?: string | null;
  order: number;
}

export interface Course {
  id: number;
  title: string;
  description: string;
  is_public: boolean;
  price: number;
  author_id: number;
  is_active: boolean;
  cover_image_url?: string;
  
  // Старые поля
  contraindications?: string; 
  recommendations?: string;     
  bonuses?: BonusItem[];     
  
  // НОВЫЕ ПОЛЯ
  target_audience?: string;   // "Курс для вас, если"
  course_basis?: string;      // "Курс включает в себя"
  class_basis?: string;       // "Основа занятий"

  modules?: LessonGroup[];
  is_purchased?: boolean;
  progress?: number;
}

export interface LessonGroup {
  id: number;
  title: string;
  order: number;
  lessons: Lesson[];
}

export interface Review {
  id: number;
  user_id: number;
  course_id: number;
  text: string;
  rating: number;
  approved: boolean;
  created_at: string;
  author_name?: string;
}

export interface AuthResponse {
  user_id: number;
  email: string;
  name: string;
  role: string;
  access_token?: string;
  refresh_token?: string;
  expires_in?: number;
  message?: string;
}

export interface ApiError {
  message: string;
  code?: string;
  details?: Record<string, string[]>;
}

export interface LoginDto {
  email: string;
  password: string;
}

export interface RegisterDto {
  email: string;
  password: string;
  full_name: string;
}

export interface CreateReviewDto {
  course_id: number;
  text: string;
  rating: number;
}

export interface CreateCourseDto {
  title: string;
  description: string;
  is_public: boolean;
  price: number;
  cover_image_url?: string;
  
  // Старые поля
  contraindications?: string;
  recommendations?: string;
  bonuses?: BonusItem[];

  // НОВЫЕ ПОЛЯ
  target_audience?: string;
  course_basis?: string;
  class_basis?: string;
}

export interface CreateModuleDto {
  course_id: number;
  title: string;
  order: number;
}

export interface CreateLessonDto {
  module_id: number;
  title: string;
  description: string;
  video_embed_id?: string;
  private_key?: string | null;
  order: number;
}

export interface FullCourseLesson {
  id: number;
  module_id: number;
  title: string;
  description: string;
  video_embed_id: string;
  private_key?: string | null;
  order: number;
}

export interface FullCourseModule {
  id: number;
  course_id: number;
  title: string;
  order: number;
  lessons: FullCourseLesson[];
}

export interface CourseFullResponse {
  course: Course;
  modules: FullCourseModule[];
}

export interface UserProfile {
  id: number;
  email: string;
  name: string;
  role: 'user' | 'admin';
  created_at: string;
}

export interface UserCourseProgress {
  id: number;
  title: string;
  description: string;
  price: number;
  is_public: boolean;
  cover_image_url: string;
  completed_count: number;
  total_lessons: number;
  progress_percent: number;
}

export interface UserCoursesResponse {
  courses: UserCourseProgress[];
}