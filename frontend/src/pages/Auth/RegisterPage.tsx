import React, { useState, useEffect } from 'react';
import {  useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

const RegisterPage = () => {
  const { register } = useAuth();
  const navigate = useNavigate();
  
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    agree: false
  });

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!formData.agree) {
      setError('Необходимо согласие на обработку персональных данных');
      return;
    }

    setIsLoading(true);
    try {
      // Регистрируем и сразу логиним (если логика useAuth это делает)
      await register(formData.name, formData.email, formData.password);
      
      // Редирект в личный кабинет сразу после регистрации
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.message || 'Ошибка регистрации');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4 py-12">
      <div className="max-w-md w-full">
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 sm:p-8">
          <h2 className="text-2xl font-serif font-bold text-center text-gray-800 mb-2">
            Регистрация
          </h2>
          <p className="text-center text-gray-500 text-sm mb-6">
            Создайте аккаунт для доступа к курсам
          </p>
          
          <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
              <div className="p-3 bg-red-50 text-red-700 text-sm rounded-lg border border-red-100">
                {error}
              </div>
            )}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Имя</label>
              <input 
                type="text" 
                required
                value={formData.name}
                onChange={(e) => setFormData({...formData, name: e.target.value})}
                className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400 outline-none transition-all"
                placeholder="Елена"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <input 
                type="email" 
                required
                value={formData.email}
                onChange={(e) => setFormData({...formData, email: e.target.value})}
                className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400 outline-none transition-all"
                placeholder="example@mail.ru"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Пароль</label>
              <input 
                type="password" 
                required
                minLength={6}
                value={formData.password}
                onChange={(e) => setFormData({...formData, password: e.target.value})}
                className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400 outline-none transition-all"
                placeholder="Минимум 6 символов"
              />
            </div>

            <div className="flex items-start">
              <input 
                type="checkbox" 
                id="terms" 
                checked={formData.agree}
                onChange={(e) => setFormData({...formData, agree: e.target.checked})}
                className="h-4 w-4 text-rose-600 focus:ring-rose-500 border-gray-300 rounded cursor-pointer" 
              />
              <label htmlFor="terms" className="ml-2 block text-sm text-gray-600 leading-tight cursor-pointer">
                Я согласна на обработку <span className="text-rose-500 hover:underline">персональных данных</span>
              </label>
            </div>

            <button 
              type="submit" 
              disabled={isLoading}
              className="w-full py-3 px-4 bg-rose-500 text-white rounded-2xl! hover:bg-rose-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium shadow-md shadow-rose-200"
            >
              {isLoading ? 'Регистрация...' : 'Зарегистрироваться'}
            </button>
            
            <p className="text-center text-xs text-gray-500 mt-6 pt-6 border-t border-gray-100">
              Уже есть аккаунт?{' '}
              <button
                type="button"
                onClick={() => navigate('/login')}
                className="text-rose-500 hover:underline font-medium"
              >
                Войти
              </button>
            </p>
          </form>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage;