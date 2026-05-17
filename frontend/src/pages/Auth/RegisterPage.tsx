import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!formData.agree) {
      setError('Необходимо согласие на обработку персональных данных');
      return;
    }

    setIsLoading(true);
    try {
      await register(formData.name, formData.email, formData.password);
      // После успешной регистрации можно сразу перенаправить на вход или в кабинет
      navigate('/login?registered=true'); 
    } catch (err: any) {
      // Здесь теперь отобразится конкретная ошибка, например "email уже зарегистрирован"
      setError(err.message || 'Ошибка регистрации');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full p-6 bg-white rounded-xl shadow-md">
        <h2 className="text-2xl font-bold text-center text-gray-800 mb-6">Регистрация</h2>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <div className="p-3 bg-red-50 text-red-600 text-sm rounded-lg border border-red-100">
              {error}
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700">Имя</label>
            <input 
              type="text" 
              required
              value={formData.name}
              onChange={(e) => setFormData({...formData, name: e.target.value})}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-rose-500 focus:ring-rose-500 p-2 border" 
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Email</label>
            <input 
              type="email" 
              required
              value={formData.email}
              onChange={(e) => setFormData({...formData, email: e.target.value})}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-rose-500 focus:ring-rose-500 p-2 border" 
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Пароль</label>
            <input 
              type="password" 
              required
              minLength={6}
              value={formData.password}
              onChange={(e) => setFormData({...formData, password: e.target.value})}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-rose-500 focus:ring-rose-500 p-2 border" 
            />
          </div>
          <div className="flex items-start">
            <input 
              type="checkbox" 
              id="terms" 
              checked={formData.agree}
              onChange={(e) => setFormData({...formData, agree: e.target.checked})}
              className="h-4 w-4 mt-1 text-rose-600 focus:ring-rose-500 border-gray-300 rounded" 
            />
            <label htmlFor="terms" className="ml-2 block text-sm text-gray-900">
              Согласен на обработку персональных данных
            </label>
          </div>
          <button 
            type="submit" 
            disabled={isLoading}
            className={`w-full py-2 px-4 bg-rose-500 text-white rounded-2xl! hover:bg-rose-600 transition-colors ${isLoading ? 'opacity-50 cursor-not-allowed' : ''}`}
          >
            {isLoading ? 'Регистрация...' : 'Зарегистрироваться'}
          </button>
          
          <p className="text-xs text-center text-gray-500 mt-4">
            Уже есть аккаунт?{' '}
            <Link to="/login" className="text-rose-500! hover:underline font-medium" style={{ textDecoration: 'none' }}>
              Войти
            </Link>
          </p>
        </form>
      </div>
    </div>
  );
};

export default RegisterPage;