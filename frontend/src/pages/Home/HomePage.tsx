const HomePage = () => {
  return (
    <div className="min-h-screen bg-pink-50 flex items-center justify-center">
      <div className="text-center p-8">
        <h1 className="text-4xl font-bold text-gray-800 mb-4">
          Восстановление женского здоровья
        </h1>
        <p className="text-lg text-gray-600">
          Профессиональная помощь и обучающие материалы по здоровью тазового дна.
        </p>
        <button className="mt-6 px-6 py-3 bg-rose-400 text-white rounded-lg hover:bg-rose-500 transition">
          Начать обучение
        </button>
      </div>
    </div>
  );
};

export default HomePage;