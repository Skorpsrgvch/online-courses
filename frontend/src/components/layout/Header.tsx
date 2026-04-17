import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { Button } from '../ui/Button';
import { Modal } from '../ui/Modal';
import { Input } from '../ui/Input';

export const Header: React.FC = () => {
  const { user, isAuthenticated, logout, login, register } = useAuth();
  const navigate = useNavigate();
  const [isScrolled, setIsScrolled] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  
  // Состояния для модальных окон
  const [isLoginOpen, setIsLoginOpen] = useState(false);
  const [isRegisterOpen, setIsRegisterOpen] = useState(false);

  // Состояния для форм
  const [loginForm, setLoginForm] = useState({ email: '', password: '' });
  const [registerForm, setRegisterForm] = useState({ name: '', email: '', password: '', agree: false });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleScroll = () => setIsScrolled(window.scrollY > 20);
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const handleLogout = async () => {
    await logout();
    navigate('/');
    setIsMobileMenuOpen(false);
  };

  const handleLoginSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);
    try {
      await login(loginForm.email, loginForm.password);
      setIsLoginOpen(false);
      setLoginForm({ email: '', password: '' });
    } catch (err: any) {
      setError(err.message || 'Ошибка входа');
    } finally {
      setIsLoading(false);
    }
  };

  const handleRegisterSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!registerForm.agree) {
      setError('Необходимо согласие на обработку данных');
      return;
    }
    setIsLoading(true);
    setError(null);
    try {
      await register(registerForm.name, registerForm.email, registerForm.password);
      setIsRegisterOpen(false);
      setRegisterForm({ name: '', email: '', password: '', agree: false });
    } catch (err: any) {
      setError(err.message || 'Ошибка регистрации');
    } finally {
      setIsLoading(false);
    }
  };

  const navLinks = [
    { name: 'О специалисте', href: '/#about' },
    { name: 'Курсы', href: '/#courses' },
    { name: 'Статьи', href: '/#articles' },
    { name: 'Видео', href: '/#videos' },
    { name: 'Контакты', href: '/#footer' },
  ];

  return (
    <header
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 bg-white ${
        isScrolled ? 'shadow-md py-3' : 'shadow-none py-3'
      }`}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-12"> 
          
          {/* Логотип */}
          <Link 
            to="/" 
            className="group flex items-center gap-3 z-50 relative no-underline"
            style={{ textDecoration: 'none' }}
          >
            <div 
              className="w-10 h-10 bg-gradient-to-br from-rose-400 to-rose-600 rounded-full flex items-center justify-center text-white font-serif font-bold text-xl shadow-lg transition-all duration-300"
            >
              W
            </div>
            <span className="text-2xl font-serif font-bold tracking-tight text-gray-800 transition-colors duration-300">
              Woman<span className="text-rose-500">Formula</span>
            </span>
          </Link>

          {/* Десктопная навигация */}
          <nav className="hidden md:flex items-center space-x-8">
            {navLinks.map((link) => (
              <a
                key={link.name}
                href={link.href}
                className="relative !text-gray-600 hover:!text-rose-500 font-medium text-lg transition-colors duration-300 decoration-transparent group"
                style={{ textDecoration: 'none' }} 
              >
                {link.name}
                {/* Красная линия при наведении */}
                <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-rose-500 transition-all duration-300 group-hover:w-full"></span>
              </a>
            ))}
          </nav>

          {/* Кнопки авторизации (Desktop) */}
          <div className="hidden md:flex items-center gap-3">
            {isAuthenticated ? (
              <div className="flex items-center gap-4">
                <Link
                  to="/dashboard"
                  className="flex items-center gap-2 text-gray-700 hover:text-rose-600 font-medium transition-colors"
                >
                  <div className="w-9 h-9 bg-rose-100 rounded-full flex items-center justify-center text-rose-600 font-bold border-2 border-white shadow-sm">
                    {user?.name?.charAt(0).toUpperCase() || 'U'}
                  </div>
                  <span className="text-sm">{user?.name || 'Кабинет'}</span>
                </Link>
                {user?.role === 'admin' && (
                  <Link to="/admin" className="px-3 py-1.5 bg-gray-800 text-white text-xs rounded-full hover:bg-gray-700 transition-colors">
                    Админ-панель
                  </Link>
                )}
                <button
                  onClick={handleLogout}
                  className="px-5 py-2 text-sm font-medium text-gray-600 border border-gray-300 rounded-full hover:border-rose-500 hover:text-rose-500 transition-all duration-300 bg-white/50 hover:bg-white"
                >
                  Выйти
                </button>
              </div>
            ) : (
              <>
                <button 
                  onClick={() => setIsLoginOpen(true)}
                  className="px-6 py-2.5 text-sm font-medium text-gray-700 bg-white border border-gray-200 !rounded-2xl hover:border-rose-300 hover:text-rose-500 hover:shadow-md transition-all duration-300 transform">
                    Вход
                  </button>
              
                <button 
                  onClick={() => setIsRegisterOpen(true)}
                  className="px-6 py-2.5 text-sm font-medium text-white bg-rose-500 !rounded-2xl hover:bg-rose-600 hover:shadow-lg transition-all duration-300 transform">
                    Регистрация
                  </button>
              </>
            )}
          </div>

          {/* Мобильная кнопка */}
          <button
            className="md:hidden text-gray-600 focus:outline-none z-50 relative"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          >
            <svg className="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              {isMobileMenuOpen ? (
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              ) : (
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              )}
            </svg>
          </button>
        </div>
      </div>

      {/* Мобильное меню (выпадает поверх всего) */}
      {isMobileMenuOpen && (
        <div className="md:hidden absolute top-full left-0 w-full bg-white/95 backdrop-blur-xl border-b border-gray-100 shadow-xl animate-fade-in-down">
          <div className="px-4 pt-4 pb-8 space-y-2">
            {navLinks.map((link) => (
              <a
                key={link.name}
                href={link.href}
                onClick={() => setIsMobileMenuOpen(false)}
                className="block px-4 py-3 text-lg font-medium !text-gray-700 hover:!text-rose-600 hover:bg-rose-50 rounded-xl transition-colors"
                style={{ textDecoration: 'none' }}
              >
                {link.name}
              </a>
            ))}
            <div className="border-t border-gray-100 my-4 pt-4 space-y-3">
              {isAuthenticated ? (
                <>
                  <Link
                    to="/dashboard"
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="block px-4 py-3 text-base font-medium text-gray-700 hover:text-rose-600 bg-gray-50 rounded-xl"
                  >
                    Личный кабинет
                  </Link>
                  {user?.role === 'admin' && (
                    <Link
                      to="/admin"
                      onClick={() => setIsMobileMenuOpen(false)}
                      className="block px-4 py-3 text-base font-medium text-gray-700 hover:text-rose-600 bg-gray-50 rounded-xl"
                    >
                      Админ-панель
                    </Link>
                  )}
                  <button
                    onClick={handleLogout}
                    className="w-full text-left px-4 py-3 text-base font-medium text-red-500 hover:bg-red-50 rounded-xl"
                  >
                    Выйти
                  </button>
                </>
              ) : (
                <div className="grid grid-cols-2 gap-3">
                  <button
                    onClick={() => { setIsLoginOpen(true); setIsMobileMenuOpen(false); }}
                    className="text-center px-4 py-3 border border-gray-200 !text-gray-700 rounded-2xl font-medium hover:border-rose-300 hover:text-rose-500 transition-all bg-white"
                  >
                    Вход
                  </button>
                  <button
                    onClick={() => { setIsRegisterOpen(true); setIsMobileMenuOpen(false); }}
                    className="text-center px-4 py-3 bg-rose-500 text-white rounded-2xl font-medium shadow-lg shadow-rose-200 hover:bg-rose-600 transition-all"
                  >
                    Регистрация
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Модальное окно Входа */}
      <Modal isOpen={isLoginOpen} onClose={() => { setIsLoginOpen(false); setError(null); }} title="Вход в аккаунт">
        <form onSubmit={handleLoginSubmit} className="space-y-4 mt-2">
          {error && (
            <div className="p-3 bg-red-50 text-red-600 text-sm rounded-lg border border-red-100">
              {error}
            </div>
          )}
          <Input
            label="Email"
            type="email"
            value={loginForm.email}
            onChange={(e) => setLoginForm({ ...loginForm, email: e.target.value })}
            required
            placeholder="example@mail.ru"
          />
          <Input
            label="Пароль"
            type="password"
            value={loginForm.password}
            onChange={(e) => setLoginForm({ ...loginForm, password: e.target.value })}
            required
            placeholder="••••••••"
          />
          <div className="flex justify-center pt-2">
            <Button type="submit" isLoading={isLoading} className="w-full !rounded-2xl sm:w-auto">
              Войти
            </Button>
          </div>
          <p className="text-xs text-center text-gray-500 mt-4">
            Нет аккаунта?{' '}
            <button type="button" onClick={() => { setIsLoginOpen(false); setIsRegisterOpen(true); }} className="text-rose-500 hover:underline font-medium">
              Зарегистрироваться
            </button>
          </p>
          <p className="text-xs text-center text-gray-400 mt-2">
            <Link to="/password-recovery" onClick={() => setIsLoginOpen(false)} className="text-rose-400 hover:text-rose-500 hover:underline">
              Забыли пароль?
            </Link>
          </p>
        </form>
      </Modal>

      {/* Модальное окно Регистрации */}
      <Modal isOpen={isRegisterOpen} onClose={() => { setIsRegisterOpen(false); setError(null); }} title="Регистрация">
        <form onSubmit={handleRegisterSubmit} className="space-y-4 mt-2">
          {error && (
            <div className="p-3 bg-red-50 text-red-600 text-sm rounded-lg border border-red-100">
              {error}
            </div>
          )}
          <Input
            label="Ваше имя"
            type="text"
            value={registerForm.name}
            onChange={(e) => setRegisterForm({ ...registerForm, name: e.target.value })}
            required
            placeholder="Елена"
          />
          <Input
            label="Email"
            type="email"
            value={registerForm.email}
            onChange={(e) => setRegisterForm({ ...registerForm, email: e.target.value })}
            required
            placeholder="example@mail.ru"
          />
          <Input
            label="Пароль"
            type="password"
            value={registerForm.password}
            onChange={(e) => setRegisterForm({ ...registerForm, password: e.target.value })}
            required
            placeholder="••••••••"
            minLength={6}
          />
          
          <div className="flex items-start gap-2 pt-2">
            <input
              type="checkbox"
              id="agree"
              checked={registerForm.agree}
              onChange={(e) => setRegisterForm({ ...registerForm, agree: e.target.checked })}
              className="w-4.5 h-4.5 mt-0.5 text-rose-500 border-gray-300 rounded focus:ring-rose-500 cursor-pointer"
            />
            <label htmlFor="agree" className="text-sm  text-gray-600 leading-tight">
              Я согласна на обработку <button type="button" className="text-rose-500 hover:underline">персональных данных</button> и принимаю <button type="button" className="text-rose-500 hover:underline">условия соглашения</button>.
            </label>
          </div>

          <div className="flex justify-center pt-2">
            <Button type="submit" isLoading={isLoading}  className="w-full !rounded-2xl sm:w-auto">
              Создать аккаунт
            </Button>
          </div>
          <p className="text-xs text-center text-gray-500 mt-4">
            Уже есть аккаунт?{' '}
            <button type="button" onClick={() => { setIsRegisterOpen(false); setIsLoginOpen(true); }} className="text-rose-500 hover:underline font-medium">
              Войти
            </button>
          </p>
        </form>
      </Modal>
    </header>
  );
};