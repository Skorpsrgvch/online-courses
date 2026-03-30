import React from 'react';
import { Link } from 'react-router-dom';
import { Button } from '../ui/Button';

import { coursesService } from '../../api'; 
import { useApi } from '../../hooks/useAPI';
import { Course } from '../../api/types';
import { Skeleton } from '../ui/Skeleton'; // Нужно создать простой скелетон

const CourseCard: React.FC<{ course: Course }> = ({ course }) => (
  <div className="bg-white rounded-xl shadow-sm hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-100 flex flex-col">
    <div className="h-48 bg-gray-200 relative">
       {/* Здесь будет изображение course.coverImage */}
       <img src={course.coverImage || '/placeholder-course.jpg'} alt={course.title} className="w-full h-full object-cover" />
    </div>
    <div className="p-6 flex-1 flex flex-col">
      <h3 className="text-xl font-bold text-gray-800 mb-2">{course.title}</h3>
      <p className="text-gray-600 mb-4 line-clamp-3 flex-1">{course.description}</p>
      <div className="flex items-center justify-between mt-auto">
        <span className="text-2xl font-bold text-rose-500">{course.price} ₽</span>
        <Link to={`/course/${course.id}`}>
          <Button variant="outline" className="text-sm">Подробнее</Button>
        </Link>
      </div>
    </div>
  </div>
);

export const CourseCatalog: React.FC = () => {
  const { data: courses, loading, error } = useApi(coursesService.getAllCourses);

  if (error) return <div className="text-center text-red-500 p-10">Ошибка загрузки курсов</div>;

  return (
    <section id="courses" className="py-20 bg-gray-50">
      <div className="max-w-7xl mx-auto px-4">
        <h2 className="text-3xl font-serif font-bold text-center text-gray-800 mb-12">
          Обучающие программы
        </h2>
        
        {loading ? (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {[1, 2, 3].map(i => <Skeleton key={i} className="h-96" />)}
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {courses?.map(course => (
              <CourseCard key={course.id} course={course} />
            ))}
          </div>
        )}
      </div>
    </section>
  );
};