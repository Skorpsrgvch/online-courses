import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { authService } from '../../api/auth.service';


export const Header: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();
  

  const [isScrolled, setIsScrolled] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);


  useEffect(() => {
    const handleScroll = () => setIsScrolled(window.scrollY > 20);
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const handleLogout = async () => {
    try {
      await authService.logout();
      logout();
      navigate('/');
    } catch (err) {
      console.error("Критическая ошибка при выходе:", err);
      localStorage.clear();
      logout();
      navigate('/');
    }
  };

  const openLogin = () => {
    navigate('/login');
    setIsMobileMenuOpen(false);
  };

  const openRegister = () => {
    navigate('/register');
    setIsMobileMenuOpen(false);
  };

  const navLinks = [
    { name: 'О специалисте', href: '/#about' },
    { name: 'Услуги', href: '/services' },
    { name: 'Курсы', href: '/courses' },
    { name: 'Контакты', href: '/#footer' },
  ];

  return (
    <header
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 bg-white ${isScrolled ? 'shadow-md py-3' : 'shadow-none py-3'
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
              className="w-10 h-10 bg-linear-to-br from-rose-400 to-rose-600 rounded-full flex items-center justify-center text-white font-serif font-bold text-xl shadow-lg transition-all duration-300"
            >
              W
            </div>
            <span className="text-2xl font-serif font-bold tracking-tight text-gray-800 transition-colors duration-300">
              Woman<span className="text-rose-500">Formula</span>
            </span>
          </Link>

          
          <nav className="hidden lg:flex items-center space-x-8">
            {navLinks.map((link) => (
              <a
                key={link.name}
                href={link.href}
                className="relative text-gray-600! hover:text-rose-500! font-medium text-lg transition-colors duration-300 decoration-transparent group"
                style={{ textDecoration: 'none' }}
              >
                {link.name}
                <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-rose-500 transition-all duration-300 group-hover:w-full"></span>
              </a>
            ))}
          </nav>

          
          <div className="hidden lg:flex items-center gap-2 lg:gap-4">
            {isAuthenticated ? (
              <div className="flex items-center gap-2 lg:gap-3">
                <Link
                  to="/dashboard"
                  className="flex items-center gap-2 text-gray-700! hover:text-rose-600! font-medium transition-colors"
                  style={{ textDecoration: 'none' }}
                >
                  <div className="w-8 h-8 lg:w-9 lg:h-9 bg-rose-100 rounded-full flex items-center justify-center text-rose-600! font-bold border-2 border-white shadow-sm">
                    {user?.name?.charAt(0).toUpperCase() || 'U'}
                  </div>
                  <span className="text-sm" style={{ textDecoration: 'none' }}>{user?.name || 'Кабинет'}</span>
                </Link>
                {user?.role === 'admin' && (
                  <Link to="/admin" className="px-2 py-1 lg:px-3 lg:py-1.5 bg-gray-800 text-white text-[10px] lg:text-xs rounded-full hover:bg-gray-700 transition-colors whitespace-nowrap"
                    style={{ textDecoration: 'none' }} >
                    Админ-панель
                  </Link>
                )}
                <button
                  onClick={handleLogout}
                  className="p-2 text-gray-600 border border-gray-300 rounded-full hover:border-rose-500 hover:text-rose-500 transition-all duration-300 bg-white/50 hover:bg-white"
                  title="Выйти"
                >
                  <svg className="w-4 h-4 lg:w-5 lg:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                  </svg>
                </button>
              </div>
            ) : (
              <>
                <button
                  onClick={openLogin}
                  className="px-6 py-2.5 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-2xl! hover:border-rose-300 hover:text-rose-500 hover:shadow-md transition-all duration-300 transform">
                  Вход
                </button>

                <button
                  onClick={openRegister}
                  className="px-6 py-2.5 text-sm font-medium text-white bg-rose-500 rounded-2xl! hover:bg-rose-600 hover:shadow-lg transition-all duration-300 transform">
                  Регистрация
                </button>
              </>
            )}
          </div>

          <button
            className="lg:hidden text-gray-600 focus:outline-none z-50 relative"
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

      {/* Мобильное меню (теперь работает и для планшетов < 1024px) */}
      {isMobileMenuOpen && (
        <div className="lg:hidden absolute top-full left-0 w-full bg-white/95 backdrop-blur-xl border-b border-gray-100 shadow-xl animate-fade-in-down">
          <div className="px-4 pt-4 pb-8 space-y-2">
            {navLinks.map((link) => (
              <a
                key={link.name}
                href={link.href}
                onClick={() => setIsMobileMenuOpen(false)}
                className="block px-4 py-3 text-lg font-medium text-gray-700! hover:text-rose-600! hover:bg-rose-50 rounded-xl transition-colors"
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
                    onClick={openLogin}
                    className="text-center px-4 py-3 border border-gray-200 text-gray-700! rounded-2xl! font-medium hover:border-rose-300 hover:text-rose-500 transition-all bg-white"
                  >
                    Вход
                  </button>
                  <button
                    onClick={openRegister}
                    className="text-center px-4 py-3 bg-rose-500 text-white rounded-2xl! font-medium shadow-lg shadow-rose-200 hover:bg-rose-600 transition-all"
                  >
                    Регистрация
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </header>
  );
};