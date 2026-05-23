import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

const LoginPage = () => {
  const { login } = useAuth();
  const navigate = useNavigate();

  const [formData, setFormData] = useState({ email: '', password: '' });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    try {
      await login(formData.email, formData.password);
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.message || 'Ошибка входа');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4 py-12">
      <div className="max-w-md w-full">

        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 sm:p-8">
          <h2 className="text-2xl font-serif font-bold text-center text-gray-800 mb-2">
            Вход в аккаунт
          </h2>
          <p className="text-center text-gray-500 text-sm mb-6">
            Введите данные для продолжения
          </p>

          <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
              <div className="p-3 bg-red-50 text-red-600 text-sm rounded-lg border border-red-100 animate-pulse">
                {error}
              </div>
            )}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <input
                type="email"
                required
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400 outline-none transition-all"
                placeholder="example@mail.ru"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Пароль</label>
              <input
                type="password"
                required
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400 outline-none transition-all"
                placeholder="••••••••"
              />
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className={`w-full py-3 px-4 mb-4! bg-rose-500 text-white rounded-2xl! hover:bg-rose-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium shadow-md shadow-rose-200`}
            >
              {isLoading ? 'Загрузка...' : 'Войти'}
            </button>

            <div className="pt-2 space-y-2">
              <p className="text-xs text-center text-gray-500">
                Нет аккаунта?{' '}
                <button
                  type="button"
                  onClick={() => navigate('/register')}
                  className="text-rose-500 hover:underline font-medium"
                >
                  Зарегистрироваться
                </button>
              </p>

              <p className="text-xs text-center text-gray-400">
                <button
                  type="button"
                  onClick={() => navigate('/password-recovery')}
                  className="text-rose-500 hover:underline font-medium"
                >
                  Забыли пароль?
                </button>
              </p>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;