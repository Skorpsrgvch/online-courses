const DashboardPage = () => {
  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <h1 className="text-3xl font-bold text-gray-800 mb-6">Личный кабинет</h1>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-xl shadow-sm">
          <h3 className="text-xl font-semibold mb-2">Мои курсы</h3>
          <p className="text-gray-600">Список ваших активных курсов появится здесь.</p>
        </div>
        <div className="bg-white p-6 rounded-xl shadow-sm">
          <h3 className="text-xl font-semibold mb-2">Профиль</h3>
          <p className="text-gray-600">Управление личными данными.</p>
        </div>
        <div className="bg-white p-6 rounded-xl shadow-sm">
          <h3 className="text-xl font-semibold mb-2">Настройки</h3>
          <p className="text-gray-600">Настройки безопасности и приватности.</p>
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;