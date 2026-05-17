import React, { useEffect, useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { coursesService } from '../../api/courses.service';
import type { Course } from '../../api/types';

const AllCoursesPage: React.FC = () => {
  const navigate = useNavigate();
  const [courses, setCourses] = useState<Course[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchCourses = async () => {
      try {
        setLoading(true);
        const data = await coursesService.getAllCourses();
        setCourses(data.filter(c => c.is_active));
        setError(null);
      } catch (err) {
        console.error('Ошибка загрузки курсов:', err);
        setError('Не удалось загрузить каталог курсов');
      } finally {
        setLoading(false);
      }
    };

    fetchCourses();
  }, []);

  const handleDetailsClick = (id: number) => {
    navigate(`/course/${id}`);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="w-12 h-12 border-4 border-rose-300 border-t-rose-500 rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-500">Загрузка курсов...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center bg-white p-8 rounded-xl shadow-sm max-w-md">
          <h2 className="text-xl font-semibold text-red-600 mb-2">Ошибка</h2>
          <p className="text-gray-600 mb-4">{error}</p>
          <Link to="/" className="text-rose-500! hover:underline" style={{ textDecoration: 'none' }}>Вернуться на главную</Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Заголовок страницы */}
        <Link to="/#courses" className="inline-flex items-start text-sm md:text-base lg:text-base font-medium text-gray-500! hover:text-rose-500! mb-4 transition-colors"
          style={{ textDecoration: 'none' }}>
            ← Назад на главную
          </Link>
        <div className="text-center mb-12">
          
          <h1 className="text-3xl md:text-4xl font-serif font-bold text-gray-900 mb-4">
            Все обучающие программы
          </h1>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            Полный каталог курсов по восстановлению женского здоровья
          </p>
        </div>

        {/* Сетка курсов */}
        {courses.length === 0 ? (
          <div className="text-center text-gray-500 py-12">Курсы пока не добавлены</div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {courses.map((course) => (
              <div
                key={course.id}
                className="bg-white rounded-2xl shadow-lg overflow-hidden border border-gray-100 hover:shadow-2xl hover:shadow-rose-100/50 transition-all duration-300 transform hover:-translate-y-1 flex flex-col h-full"
              >
                {/* Изображение */}
                <div className="relative h-56 overflow-hidden group bg-gray-100">
                  {course.cover_image_url ? (
                    <img 
                      src={course.cover_image_url} 
                      alt={course.title} 
                      className="w-full h-full object-cover transform group-hover:scale-110 transition-transform duration-700"
                      onError={(e) => {
                        (e.target as HTMLImageElement).src = "/images/Course2.jpg";
                      }}
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-gray-400 bg-gray-50">
                      Нет фото
                    </div>
                  )}
                  <div className="absolute inset-0 bg-linear-to-t from-black/60 to-transparent"></div>
                  <div className="absolute bottom-3 left-4">
                    <h3 className="text-xl font-serif font-bold text-white leading-tight drop-shadow-md">
                      {course.title}
                    </h3>
                  </div>
                </div>

                {/* Контент карточки */}
                <div className="p-5 flex-grow flex flex-col">
                  <div className="mb-4 flex-grow">
                    <p className="text-gray-600 text-sm leading-relaxed line-clamp-3">
                      {course.description}
                    </p>
                  </div>

                  {/* Блок с ценой и кнопкой */}
                  <div className="mt-auto pt-4 border-t border-gray-100">
                    <div className="flex items-baseline justify-between mb-4">
                      <span className="text-sm text-gray-500 font-medium">Стоимость:</span>
                      <span className="text-xl font-bold text-gray-900">
                        {course.price === 0 ? 'Бесплатно' : `${course.price} ₽`}
                      </span>
                    </div>
                    
                    <button 
                      onClick={() => handleDetailsClick(course.id)}
                      className="w-full py-2.5 px-4 bg-rose-500 text-white text-sm font-semibold rounded-2xl! shadow-md hover:bg-rose-600 hover:shadow-lg transition-all duration-300 cursor-pointer"
                    >
                      Подробнее о программе
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default AllCoursesPage;