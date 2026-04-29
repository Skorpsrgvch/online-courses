import React, { useState, useEffect, useCallback } from 'react';

// Данные отзывов
const reviewsData = [
  {
    id: 1,
    name: "Елена, 34 года",
    role: "Мама в декрете",
    text: "После вторых родов я чувствовала себя разбитой. Боли в пояснице не давали спокойно играть с детьми. Курс Ольги стал спасением. Уже через 2 недели я забыла о боли, а через месяц увидела, как подтянулся живот.",
    image: "/images/HeroMain.jpg"
  },
  {
    id: 2,
    name: "Марина, 42 года",
    role: "Предприниматель",
    text: "Долго стеснялась проблемы с недержанием. Ольга объяснила, что это решаемо без операций. Я прошла программу и теперь могу смело смеяться и заниматься спортом. Спасибо за профессионализм!",
    image: "/images/HeroMain.jpg"
  },
  {
    id: 3,
    name: "Анна, 29 лет",
    role: "Спортсменка",
    text: "Занимаюсь кроссфитом, и однажды почувствовала дискомфорт. Оказалось, опущение. С программой восстановления я вернулась в зал уже через месяц, но теперь делаю все правильно и безопасно.",
    image: "/images/HeroMain.jpg"
  },
  {
    id: 4,
    name: "Светлана, 38 лет",
    role: "Многодетная мама",
    text: "Три беременности подряд сильно повлияли на здоровье. Курсы Ольги дали понимание, как работать со своим телом. Без запугивания, четко по делу. Рекомендую каждой женщине!",
    image: "/images/HeroMain.jpg"
  }
];

export const ReviewSlider: React.FC = () => {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isAnimating, setIsAnimating] = useState(false);

  useEffect(() => {
    const timer = setInterval(() => handleNext(), 6000);
    return () => clearInterval(timer);
  }, [currentIndex]);

  const handleNext = useCallback(() => {
    if (isAnimating) return;
    setIsAnimating(true);
    setCurrentIndex((prev) => (prev + 1) % reviewsData.length);
    setTimeout(() => setIsAnimating(false), 500);
  }, [isAnimating]);

  const handlePrev = useCallback(() => {
    if (isAnimating) return;
    setIsAnimating(true);
    setCurrentIndex((prev) => (prev - 1 + reviewsData.length) % reviewsData.length);
    setTimeout(() => setIsAnimating(false), 500);
  }, [isAnimating]);

  const goToSlide = (index: number) => {
    if (isAnimating) return;
    setIsAnimating(true);
    setCurrentIndex(index);
    setTimeout(() => setIsAnimating(false), 500);
  };

  return (
    <section id="reviews" className="py-12 md:py-16 lg:py-20 relative z-10 bg-white overflow-hidden">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Заголовок */}
        <div className="text-center mb-10 md:mb-14">
          <h2 className="text-3xl md:text-4xl font-serif font-bold text-gray-900 mb-3">
            Истории восстановления
          </h2>
          <p className="text-sm md:text-base text-gray-600 max-w-2xl mx-auto">
            Реальные отзывы женщин, которые вернули уверенность в себе.
          </p>
        </div>

        {/* Контейнер слайдера */}
        <div className="relative max-w-4xl mx-auto">
          
          {/* Карточка отзыва */}
          <div className="overflow-hidden rounded-2xl bg-rose-50/40 border border-rose-100 shadow-lg min-h-[280px] md:min-h-[320px]">
            <div 
              className="flex transition-transform duration-500 ease-in-out h-full"
              style={{ transform: `translateX(-${currentIndex * 100}%)` }}
            >
              {reviewsData.map((review) => (
                <div key={review.id} className="w-full flex-shrink-0 p-6 md:p-8 flex flex-col md:flex-row items-center md:items-start gap-6 md:gap-8">
                  
                  {/* Левая колонка: Фото + Инфо */}
                  <div className="flex flex-col items-center text-center md:w-48 md:flex-shrink-0">
                    <div className="w-24 h-24 md:w-28 md:h-28 relative mb-4">
                      <div className="absolute inset-0 bg-rose-200 rounded-full transform rotate-6 opacity-60"></div>
                      <img 
                        src={review.image} 
                        alt={review.name} 
                        className="w-full h-full object-cover rounded-full relative z-10 border-4 border-white shadow-sm"
                      />
                    </div>
                    
                    {/* Имя и роль под фото */}
                    <div>
                      <h4 className="font-bold text-gray-900 text-base md:text-lg leading-tight">{review.name}</h4>
                      <p className="text-rose-500 text-xs md:text-sm font-medium mt-1">{review.role}</p>
                    </div>
                  </div>

                  {/* Правая колонка: Звезды + Текст */}
                  <div className="flex-grow flex flex-col justify-center">
                    {/* Звезды над текстом */}
                    <div className="flex gap-1 mb-3 md:mb-4 justify-center md:justify-start">
                      {[1, 2, 3, 4, 5].map((star) => (
                        <svg key={star} className="w-4 h-4 md:w-5 md:h-5 text-yellow-400 fill-current" viewBox="0 0 20 20">
                          <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                        </svg>
                      ))}
                    </div>
                    
                    <p className="text-sm md:text-base text-gray-700 italic leading-relaxed font-light text-center md:text-left">
                      "{review.text}"
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Кнопки навигации (уменьшены отступы) */}
          <button 
            onClick={handlePrev}
            className="absolute left-0 top-1/2 -translate-y-1/2 -translate-x-2 md:-translate-x-8 w-10 h-10 bg-white rounded-full shadow-md flex items-center justify-center text-gray-600 hover:text-rose-600 hover:scale-110 transition-all duration-300 border border-gray-100 z-20 focus:outline-none"
            aria-label="Предыдущий отзыв"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
          </button>

          <button 
            onClick={handleNext}
            className="absolute right-0 top-1/2 -translate-y-1/2 translate-x-2 md:translate-x-8 w-10 h-10 bg-white rounded-full shadow-md flex items-center justify-center text-gray-600 hover:text-rose-600 hover:scale-110 transition-all duration-300 border border-gray-100 z-20 focus:outline-none"
            aria-label="Следующий отзыв"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
          </button>

          {/* Индикаторы */}
          <div className="flex justify-center gap-2 mt-6">
            {reviewsData.map((_, index) => (
              <button
                key={index}
                onClick={() => goToSlide(index)}
                className={`h-2 rounded-full transition-all duration-300 ${
                  index === currentIndex ? 'bg-rose-500 w-6' : 'bg-gray-300 w-2 hover:bg-gray-400'
                }`}
                aria-label={`Перейти к отзыву ${index + 1}`}
              />
            ))}
          </div>

        </div>
      </div>
    </section>
  );
};