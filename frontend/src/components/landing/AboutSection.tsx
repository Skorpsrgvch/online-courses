import React from 'react';




export const AboutSection: React.FC = () => {
  return (
    <section id="about" className="py-10 md:py-14 lg:py-16 relative overflow-hidden">
      {/* Декоративные фоновые элементы */}
      
      <div className="absolute bottom-0 left-0 w-96 h-96 bg-lavender-100/50 rounded-full mix-blend-multiply filter blur-3xl opacity-60 translate-y-1/2 -translate-x-1/2"></div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        
        {/* items-start выравнивает колонки по верхнему краю, чтобы текст и фото начинались на одном уровне */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-20 items-start">
          
          {/* Блок с текстом */}
          <div className="order-1 lg:order-2 flex flex-col pt-8 lg:pt-0">
            
            {/* Эстетичный заголовок */}
            <div className="mb-6 text-center lg:text-left">
              <h2 className="text-4xl md:text-5xl font-serif font-bold text-gray-900 leading-tight mb-4">
                Ольга Пимченко
              </h2>
              
              {/* Декоративная линия */}
              <div className="flex items-center gap-4 mb-6 justify-center lg:justify-center">
                <span className="block h-px w-12 bg-rose-300"></span>
                <span className="inline-block px-5 py-1.5 bg-rose-50 text-rose-600 border border-rose-100 rounded-full text-sm font-semibold tracking-wide uppercase">
                  Автор курсов
                </span>
                <span className="block h-px w-12 bg-rose-300"></span>
              </div>
              
              
            </div>

            {/* Список регалий */}
            <div className="space-y-3 max-w-lg mx-auto lg:mx-0 mb-10">
              {[
                "Физический терапевт",
                "Специалист по послеродовой и тазовой реабилитации",
                "Специалист по реабилитации при РМЖ",
                "Инструктор-методист АФК",
                "Мама 3-х детей",
                "Студент-медик"
              ].map((item, index) => (
                <div key={index} className="flex items-start gap-4 p-3 rounded-xl hover:bg-white/60 transition-colors duration-300 group">
                  <div className="shrink-0 w-8 h-8 rounded-full bg-rose-100 flex items-center justify-center text-rose-500 group-hover:bg-rose-500 group-hover:text-white transition-colors duration-300 shadow-sm">
                    <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                  </div>
                  <span className="text-gray-700 font-medium text-lg leading-snug">{item}</span>
                </div>
              ))}
            </div>

            {/* Кнопка */}
            <div className="flex justify-center ">
              <a 
                href="/courses" 
                className="inline-flex items-center gap-2 px-8 py-3 bg-white border-2 border-rose-300 text-rose-600! font-semibold rounded-2xl hover:border-rose-500! hover:bg-rose-50 hover:-translate-y-1 hover:shadow-lg hover:shadow-rose-100/50 transition-all duration-300 group"
                style={{ textDecoration: 'none' }}
              >
                Посмотреть программы обучения
                <svg className="w-5 h-5 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M17 8l4 4m0 0l-4 4m4-4H3" />
                </svg>
              </a>
            </div>
          </div>
          
          {/* Блок с фото */}
          {/* Добавлен mt-0 для десктопа, чтобы фото начиналось строго с верха сетки */}
          <div className="relative order-2 lg:order-1 group mt-8 lg:mt-0">
            <div className="relative aspect-4/5 max-w-md mx-auto lg:max-w-full overflow-hidden rounded-4xl! shadow-2xl shadow-rose-200/50">
              <img 
                src="/images/AboutSection.jpg" 
                alt="Ольга Пимченко" 
                className="w-full h-full object-cover transform group-hover:scale-105 transition-transform duration-700 ease-out"
              />
              {/* Декоративная рамка */}
              <div className="absolute inset-0 border-2 border-white/30 rounded-4xl! pointer-events-none"></div>
              
              {/* Легкий градиент снизу для объема */}
              <div className="absolute inset-0 bg-linear-to-t from-black/10 to-transparent pointer-events-none"></div>
            </div>
          </div>

        </div>
      </div>
    </section>
  );
};