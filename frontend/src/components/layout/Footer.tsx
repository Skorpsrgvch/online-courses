import React, { useState, useEffect, type FormEvent } from 'react';
import { Link } from 'react-router-dom';
import { Modal } from '../ui/Modal';
import apiClient from '../../api/axiosInstance';

export const Footer: React.FC = () => {
  const [activeModal, setActiveModal] = useState<string | null>(null);
  const [supportForm, setSupportForm] = useState({
    name: '',
    email: '',
    question: '',
  });
  const [supportStatus, setSupportStatus] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
  const [isSupportSubmitting, setIsSupportSubmitting] = useState(false);

  const openModal = (type: string) => {
    setActiveModal(type);
    // Блокируем скролл фона
    document.body.style.overflow = 'hidden';
  };

  const closeModal = () => {
    setActiveModal(null);
    setSupportStatus(null);
    // Возвращаем скролл
    document.body.style.overflow = '';
  };

  const handleSupportSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setSupportStatus(null);
    setIsSupportSubmitting(true);

    try {
      await apiClient.post('/support/contact', supportForm);
      setSupportStatus({ type: 'success', text: 'Обращение отправлено. Мы ответим вам на почту.' });
      setSupportForm({ name: '', email: '', question: '' });
    } catch (err: any) {
      setSupportStatus({ type: 'error', text: err.message || 'Не удалось отправить обращение.' });
    } finally {
      setIsSupportSubmitting(false);
    }
  };

  // Закрываем модалку по Escape
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && activeModal) {
        closeModal();
      }
    };

    if (activeModal) {
      document.addEventListener('keydown', handleEscape);
    }

    return () => {
      document.removeEventListener('keydown', handleEscape);
      document.body.style.overflow = '';
    };
  }, [activeModal]);

  const currentYear = new Date().getFullYear();

  const socialLinks = [
    {
      name: 'Telegram', url: 'https://t.me/olga_pimchenko', icon: (
        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M11.944 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 0 1 .171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.48.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z" /></svg>
      )
    },
    {
      name: 'VK', url: '#', icon: (
        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M15.684 0H8.316C1.592 0 0 1.592 0 8.316v7.368C0 22.408 1.592 24 8.316 24h7.368C22.408 24 24 22.408 24 15.684V8.316C24 1.592 22.408 0 15.684 0zm3.692 17.123h-1.744c-.66 0-.864-.525-2.05-1.719-1.033-1.033-1.49-1.17-1.744-1.17-.356 0-.458.102-.458.593v1.575c0 .424-.135.678-1.253.678-1.846 0-3.896-1.122-5.335-3.202C4.624 10.857 4 8.673 4 8.232c0-.254.102-.491.593-.491h1.744c.44 0 .61.203.779.678.863 2.494 2.303 4.678 2.896 4.678.22 0 .322-.102.322-.66V9.735c-.068-1.19-.695-1.29-.695-1.71 0-.203.17-.407.44-.407h2.743c.373 0 .508.203.508.643v3.473c0 .372.17.508.271.508.22 0 .407-.136.813-.542 1.254-1.406 2.15-3.573 2.15-3.573.119-.254.322-.491.762-.491h1.744c.525 0 .644.27.525.643-.22 1.017-2.354 4.031-2.354 4.031-.186.305-.254.44 0 .78.186.254.796.779 1.203 1.253.745.847 1.32 1.558 1.473 2.05.17.475-.085.72-.576.72z" /></svg>
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
              className="block mb-3 hover:opacity-80 transition-opacity duration-300 no-underline"
              onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
              style={{ textDecoration: 'none' }}
            >
              <div className="flex flex-col leading-none gap-0">
                <span className="text-2xl font-serif font-light tracking-tight text-white">
                  Ольга Пимченко
                </span>
                <span className="text-2xl font-serif font-light text-pink-600 tracking-tight -mt-0.5">
                  Тазовое здоровье
                </span>
              </div>
            </Link>
            <p className="text-sm leading-relaxed text-gray-400 mb-6">
              Профессиональная платформа для восстановления женского здоровья.
              Научный подход, забота и поддержка на каждом этапе.
            </p>

            {/* Соцсети */}
            <div className="flex space-x-3">
              {socialLinks.map((social) => (
                <a
                  key={social.name}
                  href={social.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="w-9 h-9 bg-gray-800 rounded-full flex items-center justify-center text-gray-400 hover:bg-rose-500 hover:text-white transition-all duration-300 transform hover:-translate-y-1 focus:outline-none focus:ring-2 focus:ring-rose-500"
                  aria-label={social.name}
                  tabIndex={activeModal ? -1 : 0}
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
                { name: 'Услуги', href: '/services' },
                { name: 'Курсы', href: '/courses' }
              ].map((link) => (
                <li key={link.name}>
                  <a
                    href={link.href}
                    className="text-white hover:text-rose-400! font-medium text-base transition-colors duration-300 block focus:outline-none focus:text-rose-400"
                    style={{ textDecoration: 'none' }}
                    tabIndex={activeModal ? -1 : 0}
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
                <svg className="w-5 h-5 text-blue-500 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 24 24"><path d="M.057 24l1.687-6.163c-1.041-1.804-1.588-3.849-1.587-5.946.003-6.556 5.338-11.891 11.893-11.891 3.181.001 6.167 1.24 8.413 3.488 2.245 2.248 3.481 5.236 3.48 8.414-.003 6.557-5.338 11.892-11.893 11.892-1.99-.001-3.951-.5-5.688-1.448l-6.305 1.654zm6.597-3.807c1.676.995 3.276 1.591 5.392 1.592 5.448 0 9.886-4.434 9.889-9.885.002-5.462-4.415-9.89-9.881-9.892-5.452 0-9.887 4.434-9.889 9.884-.001 2.225.651 3.891 1.746 5.634l-.999 3.648 3.742-.981zm11.387-5.464c-.074-.124-.272-.198-.57-.347-.297-.149-1.758-.868-2.031-.967-.272-.099-.47-.149-.669.149-.198.297-.768.967-.941 1.165-.173.198-.347.223-.644.074-.297-.149-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.297-.347.446-.521.151-.172.2-.296.3-.495.099-.198.05-.372-.025-.521-.075-.148-.669-1.611-.916-2.206-.242-.579-.487-.501-.669-.51l-.57-.01c-.198 0-.52.074-.792.372s-1.04 1.016-1.04 2.479 1.065 2.876 1.213 3.074c.149.198 2.095 3.2 5.076 4.487.709.306 1.263.489 1.694.626.712.226 1.36.194 1.872.118.571-.085 1.758-.719 2.006-1.413.248-.695.248-1.29.173-1.414z" /></svg>
                <a href="https://wa.me/79266177171"
                  className="text-white hover:text-rose-400! transition-colors focus:outline-none focus:text-rose-400" style={{ textDecoration: 'none' }}
                  tabIndex={activeModal ? -1 : 0}>+7 999 99 99 99</a>
              </li>
              <li className="flex items-start gap-3">
                <svg className="w-5 h-5 text-blue-500 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 24 24"><path d="M11.944 0A12 12 0 0 0 0 12a12 12 0 0 0 12 12 12 12 0 0 0 12-12A12 12 0 0 0 12 0a12 12 0 0 0-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 0 1 .171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.48.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z" /></svg>
                <a href="https://t.me/olga_pimchenko"
                  className="text-white hover:text-rose-400! transition-colors focus:outline-none focus:text-rose-400" style={{ textDecoration: 'none' }}
                  tabIndex={activeModal ? -1 : 0}>@olga_pimchenko</a>
              </li>
              <li className="flex items-start gap-3">
                <svg className="w-5 h-5 text-blue-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
                <a href="mailto:womanformula@gmail.com"
                  className="text-white hover:text-rose-400! transition-colors focus:outline-none focus:text-rose-400" style={{ textDecoration: 'none' }}
                  tabIndex={activeModal ? -1 : 0}>womanformula@gmail.com</a>
              </li>
              <li className="flex items-start gap-3">
                <svg className="w-5 h-5 text-blue-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 4v-4z" /></svg>
                <button
                  type="button"
                  onClick={() => openModal('support')}
                  className="text-white hover:text-rose-400! transition-colors focus:outline-none focus:text-rose-400 text-left"
                  tabIndex={activeModal ? -1 : 0}
                >
                  Обратиться в поддержку
                </button>
              </li>
            </ul>
          </div>

          {/* Документы */}
          <div>
            <h3 className="text-white font-semibold mb-6 pb-3 uppercase tracking-wider text-sm">Документы</h3>
            <ul className="space-y-3 text-sm pl-0 -ml-7">
              {[
                { name: 'privacy', label: 'Политика обработки персональных данных' },
                { name: 'consent', label: 'Согласие на обработку ПДн' },
                { name: 'offer', label: 'Договор оферты' },
                { name: 'terms', label: 'Пользовательское соглашение' },
                { name: 'ads', label: 'Согласие на рассылку' }
              ].map((doc) => (
                <li key={doc.name}>
                  <button
                    onClick={() => openModal(doc.name)}
                    className="text-white hover:text-rose-400! transition-colors text-left w-full focus:outline-none focus:text-rose-400"
                    tabIndex={activeModal ? -1 : 0}
                  >
                    {doc.label}
                  </button>
                </li>
              ))}
            </ul>
          </div>
        </div>

        {/* Копирайт */}
        <div className="border-t border-gray-800 pt-8 flex flex-col md:flex-row justify-between items-center text-xs text-gray-500">
          <p>&copy; {currentYear} Ольга Пимченко Тазовое здоровье. Все права защищены.</p>
          <p className="mt-2 md:mt-0 text-center md:text-right">
            Сайт носит информационный характер и не является публичной офертой.
          </p>
        </div>
      </div>

      {/* Модальные окна */}
      <Modal isOpen={activeModal === 'support'} onClose={closeModal} title="Обращение в поддержку">
        <form onSubmit={handleSupportSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="support-name">
              Имя
            </label>
            <input
              id="support-name"
              type="text"
              value={supportForm.name}
              onChange={(e) => setSupportForm((prev) => ({ ...prev, name: e.target.value }))}
              required
              maxLength={120}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-500 focus:border-rose-500 outline-none"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="support-email">
              Email
            </label>
            <input
              id="support-email"
              type="email"
              value={supportForm.email}
              onChange={(e) => setSupportForm((prev) => ({ ...prev, email: e.target.value }))}
              required
              maxLength={255}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-500 focus:border-rose-500 outline-none"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="support-question">
              Ваш вопрос
            </label>
            <textarea
              id="support-question"
              value={supportForm.question}
              onChange={(e) => setSupportForm((prev) => ({ ...prev, question: e.target.value }))}
              required
              maxLength={5000}
              rows={5}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-500 focus:border-rose-500 outline-none resize-y"
            />
          </div>

          {supportStatus && (
            <div className={`p-3 rounded-lg text-sm border ${supportStatus.type === 'success'
              ? 'bg-green-50 text-green-700 border-green-200'
              : 'bg-red-50 text-red-700 border-red-200'
              }`}>
              {supportStatus.text}
            </div>
          )}

          <button
            type="submit"
            disabled={isSupportSubmitting}
            className="w-full px-6 py-3 bg-rose-500 text-white font-semibold rounded-lg hover:bg-rose-600 transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
          >
            {isSupportSubmitting ? 'Отправка...' : 'Отправить'}
          </button>
        </form>
      </Modal>

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
