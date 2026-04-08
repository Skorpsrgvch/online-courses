import React from 'react';
import { Link } from 'react-router-dom';
import { Button } from '../ui/Button';

import course1Photo from '../../assets/images/Course1.jpg';
import course2Photo from '../../assets/images/Course2.jpg';

// Заглушка изображения
const defaultImg = "https://images.unsplash.com/photo-1544367563-12123d8965cd?q=80&w=2070&auto=format&fit=crop";

const coursesData = [
  {
    id: 1,
    title: "Восстановление после родов",
    image: course1Photo,
    price: 4900,
    problems: [
      "Остался животик, диастаз и сутулость",
      "Опущение внутренних органов",
      "Боли в пояснице и тазу",
      "Недержание при нагрузках"
    ],
    link: "/course/1"
  },
  {
    id: 2,
    title: "Здоровье тазового дна",
    image: course2Photo,
    price: 5900,
    problems: [
      "Хронические боли в области таза",
      "Дискомфорт при близости",
      "Подготовка к родам",
      "Пролапс (начальные стадии)"
    ],
    link: "/course/2"
  },
  {
    id: 3,
    title: "Женское здоровье и гормоны",
    image: defaultImg, // Замените на course3Photo при наличии
    price: 3900,
    problems: [
      "Нерегулярный цикл и ПМС",
      "Гормональный дисбаланс",
      "Лишний вес и отеки",
      "Снижение либидо и энергии"
    ],
    link: "/course/3"
  }
];

interface CourseCardProps {
  title: string;
  image: string;
  price: number;
  problems: string[];
  link: string;
}

const CourseCard: React.FC<CourseCardProps> = ({ title, image, price, problems, link }) => {
  return (
    <div className="bg-white rounded-2xl shadow-lg overflow-hidden border border-gray-100 hover:shadow-2xl hover:shadow-rose-100/50 transition-all duration-300 transform hover:-translate-y-1 flex flex-col h-full">
      {/* Изображение (уменьшена высота) */}
      <div className="relative h-56 overflow-hidden group">
        <img 
          src={image} 
          alt={title} 
          className="w-full h-full object-cover transform group-hover:scale-110 transition-transform duration-700"
          onError={(e) => {
            (e.target as HTMLImageElement).src = defaultImg;
          }}
        />
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent"></div>
        <div className="absolute bottom-3 left-4">
          <h3 className="text-xl font-serif font-bold text-white leading-tight drop-shadow-md">{title}</h3>
        </div>
      </div>

      {/* Контент карточки */}
      <div className="p-5 flex-grow flex flex-col">
        <div className="mb-4 flex-grow">
          <h4 className="text-xs font-bold text-rose-500 uppercase tracking-wider mb-3">
            Этот курс для вас, если:
          </h4>
          <ul className="space-y-2 pl-0 -ml-7 w-full">
            {problems.map((problem, index) => (
              <li key={index} className="flex items-start gap-2.5 text-gray-600 text-sm leading-relaxed">
                <span className="flex-shrink-0 w-1.5 h-1.5 mt-1.5 rounded-full bg-rose-400"></span>
                <span>{problem}</span>
              </li>
            ))}
          </ul>
        </div>

        {/* Блок с ценой и кнопкой */}
        <div className="mt-auto pt-4 border-t border-gray-100">
          <div className="flex items-baseline justify-between mb-4">
            <span className="text-sm text-gray-500 font-medium">Стоимость курса:</span>
            <span className="text-2xl font-bold text-gray-900">{price} ₽</span>
          </div>
          
          <Link to={link}>
            <Button className="w-full py-2.5 !rounded-xl text-sm font-semibold shadow-md hover:shadow-lg transition-all">
              Подробнее о программе
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
};

export const CourseCatalog: React.FC = () => {
  return (
    <section id="courses" className="py-20 relative z-10">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Заголовок секции */}
        <div className="text-center mb-14 max-w-3xl mx-auto">
          <h2 className="text-3xl md:text-4xl font-serif font-bold text-gray-900 mb-4">
            Обучающие программы
          </h2>
          <p className="text-base text-gray-600 leading-relaxed">
            Выберите направление, которое поможет вам вернуть уверенность и здоровье. 
            Каждая программа основана на доказательной медицине.
          </p>
        </div>

        {/* Сетка курсов */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
          {coursesData.map((course) => (
            <CourseCard 
              key={course.id}
              {...course}
            />
          ))}
        </div>

      </div>
    </section>
  );
};