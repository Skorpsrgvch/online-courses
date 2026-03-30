import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '../ui/Button';

export const HeroSection: React.FC = () => {
  const navigate = useNavigate();

  return (
    <section className="relative bg-gradient-to-br from-rose-50 to-lavender-100 pt-20 pb-32 overflow-hidden">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        <div className="text-center max-w-3xl mx-auto">
          <h1 className="text-4xl md:text-6xl font-serif font-bold text-gray-900 mb-6 leading-tight">
            Восстановление женского здоровья <span className="text-rose-500">без лекарств</span>
          </h1>
          <p className="text-lg md:text-xl text-gray-600 mb-10 leading-relaxed">
            Научный подход к здоровью тазового дна. Избавьтесь от дискомфорта, восстановите уверенность в себе и верните радость движения под руководством эксперта.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button onClick={() => navigate('/#courses')} className="text-lg px-8 py-4">
              Начать обучение
            </Button>
            <Button variant="outline" onClick={() => navigate('/#about')} className="text-lg px-8 py-4">
              Узнать больше
            </Button>
          </div>
        </div>
      </div>
      {/* Декоративные элементы фона */}
      <div className="absolute top-0 left-0 w-full h-full overflow-hidden pointer-events-none opacity-30">
        <div className="absolute -top-24 -right-24 w-96 h-96 bg-rose-200 rounded-full mix-blend-multiply filter blur-3xl"></div>
        <div className="absolute -bottom-24 -left-24 w-96 h-96 bg-lavender-200 rounded-full mix-blend-multiply filter blur-3xl"></div>
      </div>
    </section>
  );
};