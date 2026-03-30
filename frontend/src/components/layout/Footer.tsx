import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { Modal } from '../ui/Modal';

export const Footer: React.FC = () => {
  const [activeModal, setActiveModal] = useState<string | null>(null);

  const openModal = (type: string) => setActiveModal(type);
  const closeModal = () => setActiveModal(null);

  const currentYear = new Date().getFullYear();

  return (
    <footer className="bg-gray-900 text-gray-300 pt-16 pb-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-12 mb-12">
          {/* Бренд */}
          <div className="col-span-1 md:col-span-1">
            <Link to="/" className="text-2xl font-serif font-bold text-white mb-4 block">
              Woman<span className="text-rose-400">Formula</span>
            </Link>
            <p className="text-sm leading-relaxed text-gray-400 mb-6">
              Профессиональная платформа для восстановления женского здоровья. 
              Научный подход, забота и поддержка на каждом этапе.
            </p>
            <div className="flex space-x-4">
              {/* Соцсети (заглушки) */}
              {['tg', 'vk', 'wa'].map((social) => (
                <a key={social} href="#" className="w-8 h-8 bg-gray-800 rounded-full flex items-center justify-center hover:bg-rose-500 transition-colors">
                  <span className="sr-only">{social}</span>
                  <div className="w-4 h-4 bg-current rounded-sm opacity-50"></div>
                </a>
              ))}
            </div>
          </div>

          {/* Навигация */}
          <div>
            <h3 className="text-white font-semibold mb-4 uppercase tracking-wider text-sm">Навигация</h3>
            <ul className="space-y-2 text-sm">
              <li><Link to="/#about" className="hover:text-rose-400 transition-colors">О специалисте</Link></li>
              <li><Link to="/#courses" className="hover:text-rose-400 transition-colors">Курсы</Link></li>
              <li><Link to="/#reviews" className="hover:text-rose-400 transition-colors">Отзывы</Link></li>
              <li><Link to="/#about" className="hover:text-rose-400 transition-colors">Статьи</Link></li>
              <li><Link to="/#about" className="hover:text-rose-400 transition-colors">Видео</Link></li>
              <li><Link to="/login" className="hover:text-rose-400 transition-colors">Вход</Link></li>
            </ul>
          </div>

          {/* Помощь */}
          <div>
            <h3 className="text-white font-semibold mb-4 uppercase tracking-wider text-sm">Поддержка</h3>
            <ul className="space-y-2 text-sm">
              <li><a href="mailto:support@bloomcare.ru" className="hover:text-rose-400 transition-colors">support@womanformula.ru</a></li>
              <li><a href="tel:+79990000000" className="hover:text-rose-400 transition-colors">+7 (999) 000-00-00</a></li>
              <li className="pt-2 text-gray-500 text-xs">Пн-Пт: 10:00 - 19:00 (МСК)</li>
            </ul>
          </div>

          {/* Документы */}
          <div>
            <h3 className="text-white font-semibold mb-4 uppercase tracking-wider text-sm">Документы</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <button onClick={() => openModal('privacy')} className="hover:text-rose-400 transition-colors text-left">
                  Политика конфиденциальности
                </button>
              </li>
              <li>
                <button onClick={() => openModal('terms')} className="hover:text-rose-400 transition-colors text-left">
                  Пользовательское соглашение
                </button>
              </li>
              <li>
                <button onClick={() => openModal('cookies')} className="hover:text-rose-400 transition-colors text-left">
                  Cookie-политика
                </button>
              </li>
            </ul>
          </div>
        </div>

        <div className="border-t border-gray-800 pt-8 flex flex-col md:flex-row justify-between items-center text-xs text-gray-500">
          <p>&copy; {currentYear} WomanFormula. Все права защищены.</p>
          <p className="mt-2 md:mt-0">
            Сайт носит информационный характер и не является публичной офертой.
          </p>
        </div>
      </div>

      {/* Модальные окна для документов */}
      <Modal isOpen={activeModal === 'privacy'} onClose={closeModal} title="Политика конфиденциальности">
        <div className="prose prose-sm max-w-none text-gray-600">
          <p>Здесь должен быть полный текст политики конфиденциальности в соответствии с ФЗ-152...</p>
          <p>Мы собираем только необходимые данные для оказания услуг...</p>
        </div>
      </Modal>
      
      <Modal isOpen={activeModal === 'terms'} onClose={closeModal} title="Пользовательское соглашение">
        <p className="text-gray-600">Текст пользовательского соглашения...</p>
      </Modal>

      <Modal isOpen={activeModal === 'cookies'} onClose={closeModal} title="Использование Cookie">
        <p className="text-gray-600">Информация о том, какие cookie мы используем и зачем...</p>
      </Modal>
    </footer>
  );
};