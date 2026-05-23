import React from 'react';
import { useNavigate } from 'react-router-dom';

const PrivacyPolicyPage = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Шапка с кнопкой назад */}
      <div className="bg-white border-b border-gray-200 py-6 sticky top-0 z-10 shadow-sm">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <button 
            onClick={() => navigate('/')} 
            className="inline-flex items-center gap-2 text-sm text-gray-500 hover:text-rose-500 transition-colors font-medium"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            На главную
          </button>
        </div>
      </div>

      {/* Основной контент */}
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-10 lg:py-16">
        
        {/* Заголовок */}
        <div className="mb-10 text-center">
          <h1 className="text-3xl md:text-4xl font-serif font-bold text-gray-900 mb-4">
            Политика конфиденциальности
          </h1>
          <p className="text-lg text-gray-600 font-medium">
            РАД КОП (РАПИД ЭССЕССМЕНТ ДЕЛИВЕРИ КООПЕРЕИТИВ)
          </p>
          <p className="text-sm text-gray-400 mt-2">
            Последнее обновление: 14.10.2025
          </p>
        </div>

        {/* Вступление */}
        <div className="bg-white rounded-2xl p-6 sm:p-8 shadow-sm border border-gray-100 mb-8">
          <p className="text-gray-700 leading-relaxed mb-4">
            ПРОИЗВОДСТВЕННЫЙ КООПЕРАТИВ «РАПИД ЭССЕССМЕНТ ДЕЛИВЕРИ КООПЕРЕИТИВ» (далее - администрация сайта <strong>radcop.online</strong>, Оператор, мы), ИНН 5024223968, адрес: 143404, Московская обл., г.о. Красногорск, г. Красногорск, ул. Дачная, дом 11А, офис 14/37, комната 15, придаёт большое значение защите вашей частной жизни и безопасности ваших персональных данных.
          </p>
          <p className="text-gray-700 leading-relaxed mb-4">
            Политика предназначена для информирования вас о наших действиях по сбору, обработке и защите ваших персональных данных для достижения нами заявленных целей обработки персональных данных на сайте radcop.online.
          </p>
          <p className="text-gray-700 leading-relaxed">
            Мы соблюдаем требования российского законодательства в области персональных данных. При обработке персональных данных мы придерживаемся принципов, изложенных в ст. 5 Федерального закона от 27.07.2006 г №152-ФЗ «О персональных данных» (далее - 152-ФЗ).
          </p>
        </div>

        {/* Термины */}
        <section className="mb-10">
          <h2 className="text-2xl font-serif font-bold text-gray-900 mb-6 border-l-4 border-rose-500 pl-4">
            Термины
          </h2>
          <div className="space-y-4 text-gray-700 leading-relaxed">
            <TermItem term="Персональные данные" def="любая информация, относящаяся к прямо или косвенно определенному или определяемому физическому лицу (субъекту персональных данных);" />
            <TermItem term="Оператор персональных данных (оператор)" def="государственный орган, муниципальный орган, юридическое или физическое лицо, самостоятельно или совместно с другими лицами организующие и (или) осуществляющие обработку персональных данных, а также определяющие цели обработки персональных данных, состав персональных данных, подлежащих обработке, действия (операции), совершаемые с персональными данными;" />
            <TermItem term="Обработка персональных данных" def="любое действие (операция) или совокупность действий (операций) с персональными данными, совершаемых с использованием средств автоматизации или без их использования. Включает: сбор, запись, систематизацию, накопление, хранение, уточнение (обновление, изменение), извлечение, использование, передачу (предоставление, доступ), блокирование, удаление, уничтожение." />
            <TermItem term="Автоматизированная обработка" def="обработка персональных данных с помощью средств вычислительной техники;" />
            <TermItem term="Предоставление персональных данных" def="действия, направленные на раскрытие персональных данных определенному лицу или определенному кругу лиц;" />
            <TermItem term="Блокирование персональных данных" def="временное прекращение обработки персональных данных (за исключением случаев, если обработка необходима для уточнения персональных данных);" />
            <TermItem term="Уничтожение персональных данных" def="действия, в результате которых становится невозможным восстановить содержание персональных данных в информационной системе персональных данных и (или) в результате которых уничтожаются материальные носители персональных данных;" />
            <TermItem term="Информационная система персональных данных" def="совокупность содержащихся в базах данных персональных данных и обеспечивающих их обработку информационных технологий и технических средств." />
          </div>
        </section>

        {/* 1. Сфера применения */}
        <section className="mb-10">
          <h2 className="text-2xl font-serif font-bold text-gray-900 mb-4">1. Сфера применения</h2>
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
            <p className="text-gray-700 leading-relaxed mb-4">
              Политика предназначена для информирования вас о наших действиях по обработке и защите ваших персональных данных для достижения нами заявленных целей обработки персональных данных на сайте radcop.online.
            </p>
            <p className="text-gray-700 leading-relaxed italic text-sm bg-blue-50 p-3 rounded-lg border border-blue-100">
              Обратите внимание: наш сайт может содержать ссылки на ресурсы других поставщиков услуг, которые мы не контролируем и на которые не распространяется действие Политики.
            </p>
          </div>
        </section>

        {/* 2. Несовершеннолетние */}
        <section className="mb-10">
          <h2 className="text-2xl font-serif font-bold text-gray-900 mb-4">2. Сбор персональных данных несовершеннолетних</h2>
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
            <p className="text-gray-700 leading-relaxed">
              Наш сайт не предназначен для обработки персональных данных несовершеннолетних. Если у вас есть основания полагать, что несовершеннолетний предоставил нам свои персональные данные через сайт, то просим вас сообщить нам об этом, написав на почту <a href="mailto:inbox@radcop.online.ru" className="text-rose-500 hover:underline font-medium">inbox@radcop.online.ru</a>.
            </p>
          </div>
        </section>

        {/* 3. Цели обработки */}
        <section className="mb-10">
          <h2 className="text-2xl font-serif font-bold text-gray-900 mb-6">3. Для чего мы обрабатываем ваши персональные данные</h2>
          
          <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
            <div className="p-6 border-b border-gray-100 bg-gray-50">
              <h3 className="text-lg font-bold text-gray-800 flex items-center gap-2">
                <span className="w-8 h-8 rounded-full bg-rose-100 text-rose-600 flex items-center justify-center text-sm font-bold">1</span>
                Форма обратной связи
              </h3>
            </div>
            <div className="p-6 space-y-4">
              <InfoRow label="Цель обработки" value="Заказ обратного звонка." />
              <InfoRow label="Обрабатываемые данные" value="Ваше имя; номер телефона." />
              <InfoRow label="Действия с данными" value="Сбор, запись, систематизация, накопление, хранение, уточнение (обновление, изменение), извлечение, использование, передача (предоставление, доступ), блокирование, удаление, уничтожение." />
              <InfoRow label="Основание" value="Согласие на обработку персональных данных." />
              <InfoRow label="Срок хранения" value="По достижении цели обработки (до окончания звонка и переписки), либо до отзыва согласия." />
              
              <div className="mt-4 pt-4 border-t border-gray-100">
                <p className="text-sm text-gray-600">
                  <strong className="text-gray-800">Что происходит после:</strong> Мы уничтожим ваши персональные данные путем удаления из информационной системы в течение 30 дней после обсуждения проекта или получения отзыва согласия.
                </p>
              </div>
            </div>
          </div>
        </section>

        {/* 4. Права субъекта */}
        <section className="mb-10">
          <h2 className="text-2xl font-serif font-bold text-gray-900 mb-6">4. Каковы ваши права?</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            {[
              "Право на доступ к персональным данным",
              "Право на уточнение персональных данных",
              "Право на обжалование действий или бездействия",
              "Право на обжалование решений на основе авто-обработки",
              "Право на отзыв согласия"
            ].map((item, idx) => (
              <div key={idx} className="bg-white p-4 rounded-xl border border-gray-100 shadow-sm flex items-start gap-3">
                <span className="text-green-500 mt-1">✓</span>
                <span className="text-sm font-medium text-gray-700">{item}</span>
              </div>
            ))}
          </div>

          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 space-y-6">
            <div>
              <h3 className="font-bold text-gray-900 mb-2">Как реализовать права?</h3>
              <p className="text-gray-700 text-sm leading-relaxed mb-3">
                Вы можете направить запрос следующими способами:
              </p>
              <ul className="list-disc list-inside text-sm text-gray-700 space-y-2 ml-2">
                <li>Почтовым отправлением по адресу: <br/> <span className="font-medium">143404, Московская обл., г.о. Красногорск, г. Красногорск, ул. Дачная, дом 11А, офис 14/37, комната 15</span></li>
                <li>На электронную почту: <a href="mailto:inbox@radcop.online.ru" className="text-rose-500 hover:underline">inbox@radcop.online.ru</a></li>
              </ul>
              <p className="text-xs text-gray-500 mt-3">
                В запросе необходимо указать сведения о документе, удостоверяющем личность (тип, серия, номер, кем и когда выдан), ваше ФИО, информацию о взаимоотношениях с нами и подпись.
              </p>
            </div>
            
            <div className="pt-4 border-t border-gray-100">
               <p className="text-sm text-gray-600">
                 Информация предоставляется бесплатно. В случае явно необоснованных или чрезмерных запросов мы можем отказать или взять плату за предоставление информации.
               </p>
            </div>
          </div>
        </section>

        {/* 5. Cookie */}
        <section className="mb-10">
          <h2 className="text-2xl font-serif font-bold text-gray-900 mb-6">5. Cookie и автоматическое логирование</h2>
          <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
            <div className="p-6 bg-gray-50 border-b border-gray-100">
              <p className="text-sm text-gray-600">
                Файлы cookie – это небольшие текстовые файлы, хранящиеся на вашем устройстве и содержащие информацию о вашей активности. Мы используем их для улучшения качества контента и работы сайта.
              </p>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-sm text-left">
                <thead className="text-xs text-gray-500 uppercase bg-gray-50 border-b border-gray-100">
                  <tr>
                    <th className="px-6 py-3 font-medium">Наименование</th>
                    <th className="px-6 py-3 font-medium">Провайдер</th>
                    <th className="px-6 py-3 font-medium">Срок</th>
                    <th className="px-6 py-3 font-medium">Описание</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  <CookieRow name="_ym_isad" provider="radcop.online" time="19 ч." desc="Определяют наличие блокировки рекламы" />
                  <CookieRow name="_ym_visorc" provider="radcop.online" time="30 мин." desc="Корректная работа вебвизора" />
                  <CookieRow name="_ym_uid" provider="radcop.online" time="1 год" desc="Различение посетителей" />
                  <CookieRow name="BX_USER_ID" provider="radcop.online" time="10 лет" desc="Работа сервисов сайта" />
                  <CookieRow name="_ym_d" provider="radcop.online" time="1 год" desc="Дата входа на сайт" />
                </tbody>
              </table>
            </div>
          </div>
        </section>

        {/* 6. Аналитика */}
        <section className="mb-10">
          <h2 className="text-2xl font-serif font-bold text-gray-900 mb-4">6. Аналитика данных</h2>
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
            <p className="text-gray-700 leading-relaxed">
              Собираемые данные позволяют выполнять статистический анализ для улучшения сайта. Такие файлы хранят обезличенные данные, собираемые анонимно с помощью систем аналитики, в частности <strong>Яндекс.Метрика</strong>. Это помогает нам понимать популярность разделов и удобство использования сайта.
            </p>
          </div>
        </section>

        {/* 7. Безопасность */}
        <section className="mb-10">
          <h2 className="text-2xl font-serif font-bold text-gray-900 mb-4">7. Безопасность данных</h2>
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
            <p className="text-gray-700 leading-relaxed mb-4">
              Персональные данные считаются конфиденциальной информацией и защищены от потери, изменения или несанкционированного доступа согласно законодательству РФ.
            </p>
            <p className="text-gray-700 leading-relaxed">
              Все наши сотрудники соблюдают внутренние правила обработки данных. Мы внедрили технические и организационные меры защиты с учетом современного уровня техники и рисков, связанных с обработкой.
            </p>
          </div>
        </section>

        {/* 8. Изменения */}
        <section className="mb-10">
          <h2 className="text-2xl font-serif font-bold text-gray-900 mb-4">8. Изменение политики</h2>
          <div className="bg-blue-50 rounded-xl p-6 border border-blue-100">
            <p className="text-blue-800 leading-relaxed">
              Мы оставляем за собой право вносить изменения в Политику. Просим регулярно проверять обновления. О существенных изменениях мы будем уведомлять всеми доступными способами.
            </p>
          </div>
        </section>

        {/* 9. Контакты */}
        <section className="mb-16">
          <h2 className="text-2xl font-serif font-bold text-gray-900 mb-4">9. Контакты</h2>
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 flex flex-col sm:flex-row items-start sm:items-center gap-4">
            <div className="w-12 h-12 bg-rose-100 rounded-full flex items-center justify-center text-rose-600 flex-shrink-0">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
            </div>
            <div>
              <p className="text-gray-600 text-sm mb-1">Для вопросов, связанных с обработкой и защитой персональных данных:</p>
              <a href="mailto:inbox@radcop.online" className="text-xl font-bold text-rose-600 hover:text-rose-700 transition-colors">inbox@radcop.online</a>
            </div>
          </div>
        </section>

      </div>
    </div>
  );
};

// Вспомогательные компоненты для чистоты кода

const TermItem = ({ term, def }: { term: string; def: string }) => (
  <div className="flex flex-col sm:flex-row gap-2 sm:gap-4">
    <dt className="font-bold text-gray-900 min-w-[200px]">{term}:</dt>
    <dd className="text-gray-700">{def}</dd>
  </div>
);

const InfoRow = ({ label, value }: { label: string; value: string }) => (
  <div className="border-b border-gray-50 last:border-0 pb-2 last:pb-0">
    <p className="text-xs font-bold text-gray-500 uppercase tracking-wide mb-1">{label}</p>
    <p className="text-sm text-gray-800 leading-snug">{value}</p>
  </div>
);

const CookieRow = ({ name, provider, time, desc }: { name: string; provider: string; time: string; desc: string }) => (
  <tr className="hover:bg-gray-50 transition-colors">
    <td className="px-6 py-4 font-medium text-gray-900">{name}</td>
    <td className="px-6 py-4 text-gray-600">{provider}</td>
    <td className="px-6 py-4 text-gray-600 whitespace-nowrap">{time}</td>
    <td className="px-6 py-4 text-gray-600">{desc}</td>
  </tr>
);

export default PrivacyPolicyPage;