export interface User {
  id: string;
  email: string;
  name: string;
  role: 'user' | 'admin';
  avatar?: string;
  createdAt: string;
}

export interface Course {
  id: string;
  title: string;
  description: string;
  price: number;
  coverImage: string;
  modules: Module[];
  isPurchased: boolean;
  progress?: number; // 0-100
}

export interface Module {
  id: string;
  title: string;
  lessons: Lesson[];
}

export interface Lesson {
  id: string;
  title: string;
  type: 'video' | 'article';
  contentUrl?: string; // URL видео (RuTube) или HTML контент
  duration?: number; // в секундах
  isCompleted: boolean;
  moduleId: string;
}

export interface Review {
  id: string;
  authorName: string;
  text: string;
  rating: number;
  status: 'approved' | 'pending' | 'rejected';
  createdAt: string;
}

export interface AuthResponse {
  user: User;
  // Токены не возвращаем явно, так как они в httpOnly cookies
  message?: string;
}

export interface ApiError {
  message: string;
  code?: string;
  details?: Record<string, string[]>;
}