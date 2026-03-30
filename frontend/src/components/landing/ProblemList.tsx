import React from 'react';

const problems = [
  "Повторяющиеся инфекции мочевыводящих путей",
  "Последствия родов и грудного вскармливания",
  "Пролапс тазовых органов",
  "Гормональные изменения и менопауза",
  "Дисфункция тазового дна",
  "Хирургическая травма и восстановление"
];

export const ProblemList: React.FC = () => {
  return (
    <section className="py-20 bg-white">
      <div className="max-w-7xl mx-auto px-4">
        <h2 className="text-3xl font-serif font-bold text-center text-gray-800 mb-12">
          С какими проблемами я работаю
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {problems.map((problem, idx) => (
            <div key={idx} className="p-6 bg-rose-50 rounded-xl hover:shadow-md transition-shadow duration-300 border border-rose-100">
              <div className="w-10 h-10 bg-rose-200 rounded-full flex items-center justify-center mb-4 text-rose-600 font-bold">
                {idx + 1}
              </div>
              <h3 className="text-lg font-medium text-gray-800">{problem}</h3>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};