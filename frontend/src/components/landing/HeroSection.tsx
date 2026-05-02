import React from 'react';
import { Button } from '../ui/Button';




export const HeroSection: React.FC = () => {


  return (
    <section className="relative min-h-[90vh] flex items-center pb-8 md:pb-10 lg:pb-16  py-10 md:py-14 lg:py-22  overflow-hidden pt-20">

      {/* Декоративные пятна фона */}
      <div className="absolute top-10 left-10 w-72 h-72 bg-rose-200/30 rounded-full mix-blend-multiply filter blur-3xl animate-blob"></div>
      <div className="absolute bottom-10 right-10 w-72 h-72 bg-lavender-200/30 rounded-full mix-blend-multiply filter blur-3xl animate-blob animation-delay-2000"></div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">

          <div className="text-center lg:text-left space-y-8">
            <div className="space-y-4">
              <h1 className="text-4xl md:text-5xl lg:text-6xl font-serif font-bold text-gray-900 leading-tight">
                Восстановление <br />
                <span className="text-transparent bg-clip-text bg-linear-to-r from-rose-500 to-pink-600">
                  женского здоровья
                </span>
              </h1>

              <p className="text-lg md:text-xl text-gray-600 leading-relaxed max-w-xl mx-auto lg:mx-0 font-light">
                Профессиональная помощь и обучающие материалы по здоровью тазового дна.
                Верните уверенность в себе без лекарств и операций.
              </p>
            </div>


            <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-center pt-4">
              <Button
                onClick={() => {
                  const section = document.querySelector('#courses');
                  if (section) {
                    section.scrollIntoView({ behavior: 'smooth' }); // Плавная прокрутка
                  }
                }}
                className="px-8 py-4 text-lg rounded-2xl! shadow-lg shadow-rose-200 hover:shadow-rose-300 hover:-translate-y-1 transition-all duration-300"
              >
                Начать обучение
              </Button>

              {/* Кнопка записи */}
              <button
                onClick={() => {
                  const section = document.querySelector('#courses');
                  if (section) {
                    section.scrollIntoView({ behavior: 'smooth' }); // Плавная прокрутка
                  }
                }}
                className="px-8 py-4 text-lg font-medium text-gray-700 bg-white border-2 border-gray-200 rounded-2xl! 
                           hover:border-rose-300 hover:text-rose-500 hover:shadow-md hover:shadow-rose-100 
                           transition-all duration-300 transform"
              >
                Записаться на консультацию
              </button>
            </div>

            
          </div>

          <div className="order-first lg:order-last relative group">


            <div className="relative rounded-4xl overflow-hidden  aspect-4/3 lg:aspect-4/5 max-h-150 shadow-2xl shadow-rose-200">
              <img
                src="/images/HeroMain.jpg"
                alt="Hero"
                className="w-full h-full object-cover transform group-hover:scale-105 transition-transform duration-700"
              />
            </div>
          </div>

        </div>
      </div>
    </section>
  );
};