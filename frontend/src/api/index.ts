export { authService } from './auth.service';
export { coursesService } from './courses.service';
export { reviewsService } from './reviews.service';
export { servicesService } from './services.service';
export { userService } from './user.service';
export { apiClient } from './axiosInstance';

export type {
  User,
  Course,
  CourseModule,
  Lesson,
  LessonGroup,
  Review,
  AuthResponse,
  LoginDto,
  RegisterDto,
  CreateReviewDto,
  CreateCourseDto,
  CreateModuleDto,
  CreateLessonDto,
  CourseFullResponse,
  FullCourseModule,
  FullCourseLesson,
  UserProfile,
  UserCourseProgress,
  UserCoursesResponse,
} from './types';
