import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { Button } from '../ui/Button';

export const Header: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();
  const [isScrolled, setIsScrolled] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  // Отслеживаем скролл
  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 10); // Чуть раньше начинаем скролл эффект
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const handleLogout = async () => {
    await logout();
    navigate('/');
    setIsMobileMenuOpen(false);
  };

  const navLinks = [
    { name: 'О специалисте', href: '/#about' },
    { name: 'Курсы', href: '/#courses' },
    { name: 'Отзывы', href: '/#reviews' },
    { name: 'Статьи', href: '/#articles' },
    { name: 'Видео', href: '/#videos' },
  ];

  return (
    // header теперь всегда на всю ширину, фон меняется через классы
    <header
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 bg-white ${
        isScrolled ? 'shadow-md py-3' : 'shadow-none py-3'
      }`}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-12"> 
          
          {/* Логотип */}
          <Link to="/" className="group flex items-center gap-3 z-50 relative ">
            <div className="w-10 h-10 bg-gradient-to-br from-rose-400 to-rose-600 rounded-full flex items-center justify-center text-white font-serif font-bold text-xl shadow-lg group-hover:shadow-rose-300 group-hover:scale-105 transition-all duration-300">
              W
            </div>
            <span className={`text-2xl font-serif font-bold tracking-tight transition-colors duration-300 ${
              isScrolled ? 'text-gray-800' : 'text-gray-800'
            }`}>
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
          <div className="hidden md:flex items-center space-x-4">
            {isAuthenticated ? (
              <div className="flex items-center gap-4">
                <Link
                  to="/dashboard"
                  className="flex items-center gap-2 text-gray-700 hover:text-rose-600 font-medium transition-colors"
                >
                  <div className="w-9 h-9 bg-rose-100 rounded-full flex items-center justify-center text-rose-600 font-bold border-2 border-white shadow-sm">
                    {user?.name?.charAt(0).toUpperCase() || 'U'}
                  </div>
                </Link>
                {user?.role === 'admin' && (
                  <Link to="/admin">
                    <span className="px-2 py-1 bg-gray-800 text-white text-xs rounded-full">Admin</span>
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
                <Link to="/login">
                  <button className="px-6 py-2.5 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-full hover:border-rose-300 hover:text-rose-500 hover:shadow-md transition-all duration-300 transform hover:-translate-y-0.5">
                    Вход
                  </button>
                </Link>
                <Link to="/register">
                  <button className="px-6 py-2.5 text-sm font-medium text-white bg-gradient-to-r from-rose-500 to-rose-600 rounded-full hover:from-rose-600 hover:to-rose-700 hover:shadow-lg hover:shadow-rose-300/50 transition-all duration-300 transform hover:-translate-y-0.5">
                    Регистрация
                  </button>
                </Link>
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
                className="block px-4 py-3 text-lg font-medium text-gray-700 hover:text-rose-600 hover:bg-rose-50 rounded-xl transition-colors"
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
                  <button
                    onClick={handleLogout}
                    className="w-full text-left px-4 py-3 text-base font-medium text-red-500 hover:bg-red-50 rounded-xl"
                  >
                    Выйти
                  </button>
                </>
              ) : (
                <div className="grid grid-cols-2 gap-3">
                  <Link
                    to="/login"
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="text-center px-4 py-3 border border-gray-200 text-gray-700 rounded-full font-medium hover:border-rose-300 hover:text-rose-500 transition-all"
                  >
                    Вход
                  </Link>
                  <Link
                    to="/register"
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="text-center px-4 py-3 bg-rose-500 text-white rounded-full font-medium shadow-lg shadow-rose-200 hover:bg-rose-600 transition-all"
                  >
                    Регистрация
                  </Link>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </header>
  );
};