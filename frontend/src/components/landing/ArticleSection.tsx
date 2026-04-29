import React from 'react';
import { Link } from 'react-router-dom';
import { Button } from '../ui/Button';




const articlesData = [
  {
    id: 1,
    title: "10 советов как уменьшить боль в тазу",
    image: "/images/Course1.jpg",
    questions: [
      "Как быстро снять острую боль?",
      "Что важно знать о первой помощи?",
      "Почему нельзя терпеть дискомфорт?"
    ],
    link: "/articles/1"
  },
  {
    id: 2,
    title: "Мифы о здоровье тазового дна",
    image: "/images/Course2.jpg",
    questions: [
      "Правда ли, что упражнения Кегеля всем полезны?",
      "Влияет ли осанка на боли в тазу?",
      "Когда нужно идти к специалисту?"
    ],
    link: "/articles/2"
  },
  {
    id: 3,
    title: "Восстановление после родов: с чего начать",
    image: "/images/Course1.jpg",
    questions: [
      "Можно ли качать пресс сразу?",
      "Как понять, есть ли диастаз?",
      "Почему важно обратиться к тазовому терапевту?"
    ],
    link: "/articles/3"
  }
];

interface ArticleCardProps {
  title: string;
  image: string;
  questions: string[];
  link: string;
}

const ArticleCard: React.FC<ArticleCardProps> = ({ title, image, questions, link }) => {
  return (
    <div className=" rounded-2xl shadow-lg overflow-hidden border border-gray-100 hover:shadow-xl hover:shadow-rose-100/40 transition-all duration-300 transform hover:-translate-y-1 flex flex-col h-full group">
      {/* Изображение */}
      <div className="relative h-48 overflow-hidden">
        <img 
          src={image} 
          alt={title} 
          className="w-full h-full object-cover transform group-hover:scale-105 transition-transform duration-700"
        />
        <div className="absolute top-3 right-3 bg-white/90 backdrop-blur-sm px-3 py-1 rounded-full text-xs font-bold text-rose-500 uppercase tracking-wide shadow-sm">
          Статья
        </div>
      </div>

      {/* Контент */}
      <div className="p-5 flex-grow flex flex-col">
        <h3 className="text-lg font-serif font-bold text-gray-900 mb-4 leading-snug line-clamp-2 min-h-14">
          {title}
        </h3>

        <div className="mb-6 flex-grow">
          {/* Заголовок "В этой статье" удален */}
          <ul className="space-y-3 pl-0 -ml-7">
            {questions.map((q, index) => (
              <li key={index} className="flex items-start gap-3 text-gray-700 text-sm md:text-base leading-relaxed font-medium">
                <span className="flex-shrink-0 w-1.5 h-1.5 mt-1.5 rounded-full bg-rose-400 shadow-sm"></span>
                <span>{q}</span>
              </li>
            ))}
          </ul>
        </div>

        <div className="mt-auto pt-4 border-t border-gray-50">
          <Link to={link}>
            <Button variant="outline" className="w-full py-2.5 rounded-xl! text-sm font-semibold border-rose-200 text-rose-600 hover:bg-rose-50 hover:border-rose-300 hover:shadow-md transition-all">
              Читать статью
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
};

export const ArticleSection: React.FC = () => {
  return (
    <section id="articles" className="py-12 md:py-16 lg:py-24 relative z-10 ">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Заголовок секции */}
        <div className="text-center mb-12 md:mb-16 max-w-3xl mx-auto">
          <h2 className="text-3xl md:text-4xl font-serif font-bold text-gray-900 mb-4">
            Полезные материалы
          </h2>
          <p className="text-base text-gray-600 leading-relaxed">
            Ответы на важные вопросы о женском здоровье, разбор мифов и практические советы от эксперта.
          </p>
        </div>

        {/* Сетка статей */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
          {articlesData.map((article) => (
            <ArticleCard 
              key={article.id}
              {...article}
            />
          ))}
        </div>

        {/* Кнопка "Все статьи" */}
        <div className="text-center mt-12">
          <button className="inline-flex items-center gap-2 px-8 py-3 bg-white border-2 border-rose-300 text-rose-600! font-semibold rounded-2xl! hover:border-rose-500! hover:bg-rose-50 hover:-translate-y-1 hover:shadow-lg hover:shadow-rose-100/50 transition-all duration-300 group">
            Смотреть все статьи
            <svg className="w-4 h-4 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8l4 4m0 0l-4 4m4-4H3" />
            </svg>
          </button>
        </div>

      </div>
    </section>
  );
};