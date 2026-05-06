import React, { useEffect, useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { coursesService } from '../../api/courses.service';
import type { Course } from '../../api/types';

const CourseCatalog: React.FC = () => {
  const navigate = useNavigate();
  const [courses, setCourses] = useState<Course[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchCourses = async () => {
      try {
        setLoading(true);
        const data = await coursesService.getAllCourses();
        // Фильтруем только активные курсы
        const activeCourses = data.filter(c => c.is_active);
        setCourses(activeCourses);
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

  // Берем только первые 3 курса для отображения на главной
  const displayCourses = courses.slice(0, 3);

  if (loading) {
    return (
      <section id="courses" className="py-20 relative z-10">
        <div className="max-w-7xl mx-auto px-4 text-center">
          <div className="inline-block w-10 h-10 border-4 border-rose-200 border-t-rose-500 rounded-full animate-spin"></div>
          <p className="mt-4 text-gray-500">Загрузка программ...</p>
        </div>
      </section>
    );
  }

  if (error) {
    return (
      <section id="courses" className="py-20 relative z-10 ">
        <div className="max-w-7xl mx-auto px-4 text-center text-red-500 bg-red-50 p-6 rounded-xl">
          <p>{error}</p>
          <button onClick={() => window.location.reload()} className="mt-4 text-sm underline hover:text-red-700">
            Попробовать снова
          </button>
        </div>
      </section>
    );
  }

  return (
    <section id="courses" className="py-10 md:py-14 lg:py-16 relative z-10">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Заголовок секции */}
        <div className="text-center mb-14 max-w-3xl mx-auto">
          <h2 className="text-3xl md:text-4xl font-serif font-bold text-gray-900 mb-4">
            Обучающие программы
          </h2>
          <p className="text-base text-gray-600 leading-relaxed">
            Выберите направление, которое поможет вам вернуть уверенность и здоровье. 
          </p>
        </div>

        {/* Сетка курсов */}
        {displayCourses.length === 0 ? (
          <div className="text-center text-gray-500 py-10">Курсы пока не добавлены</div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
            {displayCourses.map((course) => (
              <div
                key={course.id}
                className="bg-white rounded-2xl shadow-lg overflow-hidden border border-gray-100 hover:shadow-2xl hover:shadow-rose-100/50 transition-all duration-300 transform hover:-translate-y-1 flex flex-col h-full group"
              >
                {/* Изображение */}
                <div className="relative h-56 overflow-hidden bg-gray-100">
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
                      <svg className="w-12 h-12 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                      </svg>
                    </div>
                  )}
                  <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent"></div>
                  <div className="absolute bottom-3 left-4">
                    <h3 className="text-xl font-serif font-bold text-white leading-tight drop-shadow-md line-clamp-2">
                      {course.title}
                    </h3>
                  </div>
                </div>

                {/* Контент карточки */}
                <div className="p-5 flex-grow flex flex-col">
                  <div className="mb-4 flex-grow">
                    <p className="text-gray-600 text-sm leading-relaxed line-clamp-3 min-h-[3.75rem]">
                      {course.description || 'Описание курса скоро появится...'}
                    </p>
                  </div>

                  {/* Блок с ценой и кнопкой */}
                  <div className="mt-auto pt-4 border-t border-gray-100">
                    <div className="flex items-baseline justify-between mb-4">
                      <span className="text-sm text-gray-500 font-medium">Стоимость:</span>
                      <span className="text-xl font-bold text-gray-900">
                        {course.price === 0 ? (
                          <span className="text-gray-900">Бесплатно</span>
                        ) : (
                          `${course.price} ₽`
                        )}
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

        {/* Кнопка "Все курсы" - показываем, если всего курсов больше 3 */}
        {courses.length > 3 && (
          <div className="text-center mt-12">
            <Link 
              to="/courses" 
              className="inline-flex items-center gap-2 px-8 py-3 bg-white border-2 border-rose-300 text-rose-600! font-semibold rounded-2xl hover:border-rose-500! hover:bg-rose-50 hover:-translate-y-1 hover:shadow-lg hover:shadow-rose-100/50 transition-all duration-300 group" 
              style={{ textDecoration: 'none' }}
            >
              Смотреть все программы
              <svg className="w-5 h-5 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8l4 4m0 0l-4 4m4-4H3" />
              </svg>
            </Link>
          </div>
        )}
      </div>
    </section>
  );
};

export default CourseCatalog;