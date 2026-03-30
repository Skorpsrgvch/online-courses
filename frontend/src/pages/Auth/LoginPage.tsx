const LoginPage = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full p-6 bg-white rounded-xl shadow-md">
        <h2 className="text-2xl font-bold text-center text-gray-800 mb-6">Вход</h2>
        <form className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Email</label>
            <input type="email" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-rose-500 focus:ring-rose-500 p-2 border" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Пароль</label>
            <input type="password" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-rose-500 focus:ring-rose-500 p-2 border" />
          </div>
          <button type="submit" className="w-full py-2 px-4 bg-rose-500 text-white rounded-md hover:bg-rose-600">
            Войти
          </button>
        </form>
      </div>
    </div>
  );
};

export default LoginPage;