import React from 'react';
import { Link } from 'react-router-dom';

interface CourseCardProps {
  id: string;
  title: string;
  description: string;
  price: number;
  image: string;
  progress?: number;
}

const CourseCard: React.FC<CourseCardProps> = ({ 
  id, 
  title, 
  description, 
  price, 
  image,
  progress 
}) => {
  return (
    <div className="bg-white rounded-2xl overflow-hidden shadow-md hover:shadow-xl transition-shadow border border-gray-100">
      <div className="relative">
        <img 
          src={image} 
          alt={title} 
          className="w-full h-48 object-cover"
          loading="lazy"
        />
        {progress !== undefined && (
          <div className="absolute bottom-0 left-0 right-0 bg-black bg-opacity-60">
            <div className="h-1 bg-pink-400" style={{ width: `${progress}%` }}></div>
            <div className="text-white text-xs p-2">
              {progress}% завершено
            </div>
          </div>
        )}
        <div className="absolute top-4 right-4 bg-pink-100 text-pink-800 text-sm font-medium px-3 py-1 rounded-full">
          {price === 0 ? 'Бесплатно' : `${price.toLocaleString('ru-RU')} ₽`}
        </div>
      </div>
      <div className="p-6">
        <h3 className="text-xl font-bold text-gray-800 mb-2 line-clamp-2">{title}</h3>
        <p className="text-gray-600 mb-4 line-clamp-3">{description}</p>
        <Link 
          to={progress !== undefined ? `/course/${id}` : `/course/${id}/preview`}
          className="inline-block w-full bg-gradient-to-r from-pink-500 to-rose-500 text-white font-medium py-3 px-4 rounded-lg text-center hover:opacity-90 transition-opacity"
        >
          {progress !== undefined ? 'Продолжить обучение' : 'Подробнее'}
        </Link>
      </div>
    </div>
  );
};

export default CourseCard;