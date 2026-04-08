import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { Modal } from '../ui/Modal';

export const Footer: React.FC = () => {
  const [activeModal, setActiveModal] = useState<string | null>(null);

  const openModal = (type: string) => setActiveModal(type);
  const closeModal = () => setActiveModal(null);

  const currentYear = new Date().getFullYear();

  // Ссылки на соцсети (замените # на реальные ссылки)
  const socialLinks = [
    {
      name: 'Telegram', url: 'https://t.me/olga_pimchenko', icon: (
        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M11.944 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 0 1 .171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.48.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z" /></svg>
      )
    },
    {
      name: 'Max', url: '#', icon: (
        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M24 4.557c-.883.392-1.832.656-2.828.775 1.017-.609 1.798-1.574 2.165-2.724-.951.564-2.005.974-3.127 1.195-.897-.957-2.178-1.555-3.594-1.555-3.179 0-5.515 2.966-4.797 6.045-4.091-.205-7.719-2.165-10.148-5.144-1.29 2.213-.669 5.108 1.523 6.574-.806-.026-1.566-.247-2.229-.616-.054 2.281 1.581 4.415 3.949 4.89-.693.188-1.452.232-2.224.084.626 1.956 2.444 3.379 4.6 3.419-2.07 1.623-4.678 2.348-7.29 2.04 2.179 1.397 4.768 2.212 7.548 2.212 9.142 0 14.307-7.721 13.995-14.646.962-.695 1.797-1.562 2.457-2.549z" /></svg>
      )
    },
    {
      name: 'RuTube', url: '#', icon: (
        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z" /></svg>
      )
    }
  ];

  return (
    <footer id="footer" className="bg-gray-900 text-gray-300 pt-16 pb-8 border-t border-gray-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">

          {/* Бренд и Описание */}
          <div className="col-span-1">
            <Link
              to="/"
              className="text-2xl font-serif font-bold text-white mb-4 block hover:text-rose-400 transition-colors duration-300 no-underline"
              style={{ textDecoration: 'none' }}
              onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
            >
              Woman<span className="text-rose-400">Formula</span>
            </Link>
            <p className="text-sm leading-relaxed text-gray-400 mb-6">
              Профессиональная платформа для восстановления женского здоровья.
              Научный подход, забота и поддержка на каждом этапе.
            </p>

            {/* Соцсети (Кружки) */}
            <div className="flex space-x-3">
              {socialLinks.map((social) => (
                <a
                  key={social.name}
                  href={social.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="w-9 h-9 bg-gray-800 rounded-full flex items-center justify-center text-gray-400 hover:bg-rose-500 hover:text-white transition-all duration-300 transform hover:-translate-y-1"
                  aria-label={social.name}
                >
                  {social.icon}
                </a>
              ))}
            </div>
          </div>
          {/* Навигация */}
          <div>
            <h3 className="text-white font-semibold mb-6 uppercase tracking-wider text-sm">Навигация</h3>
            <ul className="space-y-3 text-sm pl-0 -ml-7">
              {[
                { name: 'Главная', href: '/' },
                { name: 'О специалисте', href: '/#about' },
                { name: 'Курсы', href: '/#courses' },
                { name: 'Статьи', href: '/#articles' },
                { name: 'Видео', href: '/#videos' },
                { name: 'Отзывы', href: '/#reviews' }
              ].map((link) => (
                <li key={link.name}>
                  <a
                    href={link.href}
                    className="relative !text-white hover:!text-rose-400 font-medium text-base transition-colors duration-300 decoration-transparent group block"
                    style={{ textDecoration: 'none' }}
                  >
                    {link.name}

                  </a>
                </li>
              ))}
            </ul>
          </div>


          {/* Контакты */}
          <div>
            <h3 className="text-white font-semibold mb-6 pb-3 uppercase tracking-wider text-sm">Контакты</h3>
            <ul className="space-y-4 text-sm pl-0 -ml-7">
              <li className="flex items-start gap-3">
                <svg className="w-5 h-5 !text-blue-500 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 24 24"><path d="M.057 24l1.687-6.163c-1.041-1.804-1.588-3.849-1.587-5.946.003-6.556 5.338-11.891 11.893-11.891 3.181.001 6.167 1.24 8.413 3.488 2.245 2.248 3.481 5.236 3.48 8.414-.003 6.557-5.338 11.892-11.893 11.892-1.99-.001-3.951-.5-5.688-1.448l-6.305 1.654zm6.597-3.807c1.676.995 3.276 1.591 5.392 1.592 5.448 0 9.886-4.434 9.889-9.885.002-5.462-4.415-9.89-9.881-9.892-5.452 0-9.887 4.434-9.889 9.884-.001 2.225.651 3.891 1.746 5.634l-.999 3.648 3.742-.981zm11.387-5.464c-.074-.124-.272-.198-.57-.347-.297-.149-1.758-.868-2.031-.967-.272-.099-.47-.149-.669.149-.198.297-.768.967-.941 1.165-.173.198-.347.223-.644.074-.297-.149-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.297-.347.446-.521.151-.172.2-.296.3-.495.099-.198.05-.372-.025-.521-.075-.148-.669-1.611-.916-2.206-.242-.579-.487-.501-.669-.51l-.57-.01c-.198 0-.52.074-.792.372s-1.04 1.016-1.04 2.479 1.065 2.876 1.213 3.074c.149.198 2.095 3.2 5.076 4.487.709.306 1.263.489 1.694.626.712.226 1.36.194 1.872.118.571-.085 1.758-.719 2.006-1.413.248-.695.248-1.29.173-1.414z" /></svg>
                <a href="https://wa.me/79266177171"
                  className="hover:text-rose-400! !text-white  transition-colors"
                  style={{ textDecoration: 'none' }}>+7 999 99 99 99</a>
              </li>
              <li className="flex items-start gap-3">
                <svg className="w-5 h-5 !text-blue-500 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 24 24"><path d="M11.944 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 0 1 .171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.48.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z" /></svg>
                <a href="https://t.me/olga_pimchenko"
                  className="hover:!text-rose-400 !text-white  transition-colors"
                  style={{ textDecoration: 'none' }}>@olga_pimchenko</a>
              </li>
              <li className="flex items-start gap-3">
                <svg className="w-5 h-5 !text-blue-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
                <a href="mailto:womanformula@gmail.com"
                  className="hover:!text-rose-400 !text-white transition-colors"
                  style={{ textDecoration: 'none' }}>womanformula@gmail.com</a>
              </li>
            </ul>
          </div>

          {/* Документы */}
          <div>
            <h3 className="text-white font-semibold mb-6 pb-3 uppercase  tracking-wider text-sm">Документы</h3>
            <ul className="space-y-3 text-sm pl-0 -ml-7">
              <li>
                <button onClick={() => openModal('privacy')} className="hover:text-rose-400 transition-colors text-left w-full">
                  Политика обработки персональных данных
                </button>
              </li>
              <li>
                <button onClick={() => openModal('consent')} className="hover:text-rose-400 transition-colors text-left w-full">
                  Согласие на обработку ПДн
                </button>
              </li>
              <li>
                <button onClick={() => openModal('offer')} className="hover:text-rose-400 transition-colors text-left w-full">
                  Договор оферты
                </button>
              </li>
              <li>
                <button onClick={() => openModal('terms')} className="hover:text-rose-400 transition-colors text-left w-full">
                  Пользовательское соглашение
                </button>
              </li>
              <li>
                <button onClick={() => openModal('ads')} className="hover:text-rose-400 transition-colors text-left w-full">
                  Согласие на рассылку
                </button>
              </li>
            </ul>
          </div>
        </div>

        {/* Копирайт */}
        <div className="border-t border-gray-800 pt-8 flex flex-col md:flex-row justify-between items-center text-xs text-gray-500">
          <p>&copy; {currentYear} WomanFormula. Все права защищены.</p>
          <p className="mt-2 md:mt-0 text-center md:text-right">
            Сайт носит информационный характер и не является публичной офертой.
          </p>
        </div>
      </div>

      {/* Модальные окна для документов (Заглушки текста) */}
      <Modal isOpen={activeModal === 'privacy'} onClose={closeModal} title="Политика обработки персональных данных">
        <div className="prose prose-sm max-w-none text-gray-600">
          <p>Здесь должен быть полный текст политики конфиденциальности в соответствии с ФЗ-152...</p>
        </div>
      </Modal>
      <Modal isOpen={activeModal === 'consent'} onClose={closeModal} title="Согласие на обработку ПДн">
        <p className="text-gray-600">Текст согласия на обработку персональных данных...</p>
      </Modal>
      <Modal isOpen={activeModal === 'offer'} onClose={closeModal} title="Договор оферты">
        <p className="text-gray-600">Текст договора публичной оферты...</p>
      </Modal>
      <Modal isOpen={activeModal === 'terms'} onClose={closeModal} title="Пользовательское соглашение">
        <p className="text-gray-600">Текст пользовательского соглашения...</p>
      </Modal>
      <Modal isOpen={activeModal === 'ads'} onClose={closeModal} title="Согласие на рассылку">
        <p className="text-gray-600">Условия получения рекламных и информационных рассылок...</p>
      </Modal>
    </footer>
  );
};