import React from 'react';
import { Link } from 'react-router-dom';

const ConsentPage = () => {
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col font-sans text-gray-800">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 sticky top-0 z-40">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <Link to="/" className="text-xl font-serif font-bold text-rose-600 hover:text-rose-700 transition-colors">
            RAD COP
          </Link>
          <Link 
            to="/privacy-policy" 
            className="text-sm text-gray-500 hover:text-rose-600 transition-colors font-medium"
          >
            Политика конфиденциальности
          </Link>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-grow py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto bg-white rounded-2xl shadow-sm border border-gray-100 p-8 md:p-12">
          
          <h1 className="text-3xl md:text-4xl font-serif font-bold text-gray-900 mb-8 text-center leading-tight">
            Согласие на обработку персональных данных
          </h1>

          <div className="prose prose-rose max-w-none text-gray-600 leading-relaxed space-y-6">
            <p>
              Настоящим, оставляя заявку на услуги Оператора либо совершая иные действия, связанные с передачей своих персональных данных Оператору, на сайте и его поддоменах <strong>https://radcop.online</strong>, действуя свободно, своей волей и в своем интересе, а также подтверждая свою дееспособность, я даю свое полное, конкретное, информированное, сознательное и однозначное согласие <strong>ПРОИЗВОДСТВЕННОМУ КООПЕРАТИВУ «РАПИД ЭССЕССМЕНТ ДЕЛИВЕРИ КООПЕРЕИТИВ»</strong>, ИНН 5024223968, адрес: 143404, Московская обл., г.о. Красногорск, г. Красногорск, ул. Дачная, дом 11А, офис 14/37, комната 15 (далее — Оператор) на обработку персональных данных в соответствии с Политикой в отношении обработки персональных данных на следующих условиях:
            </p>

            <ol className="list-decimal pl-6 space-y-4 marker:font-bold marker:text-rose-600">
              <li>
                Данное Согласие дается на обработку персональных данных с использованием средств автоматизации.
              </li>
              
              <li>
                Согласие дается на обработку следующих моих персональных данных и персональных данных лиц, не достигших полной дееспособности, законными представителями которых является Субъект ПД:
                <ul className="list-[2.1] pl-6 mt-2 space-y-1">
                  <li><span className="font-medium text-gray-800">2.1. Общие персональные данные:</span></li>
                  <ul className="list-[2.1.1] pl-6 space-y-1">
                    <li>Фамилия, имя, отчество;</li>
                    <li>Адрес электронной почты;</li>
                    <li>Номер телефона;</li>
                    <li>Банковские реквизиты;</li>
                    <li>Номер расчетного счета;</li>
                    <li>Никнейм в мессенджере Telegram.</li>
                  </ul>
                </ul>
              </li>

              <li>
                <span className="font-medium text-gray-800">Цели обработки персональных данных:</span>
                <ul className="list-[3.1] pl-6 mt-2 space-y-1">
                  <li>Заключение и исполнение гражданско-правового договора (доступ к курсам, оплата услуг);</li>
                  <li>Продвижение услуг Оператора на рынке;</li>
                  <li>Осуществление рассылки сообщений информационного и рекламного характера.</li>
                </ul>

                <div className="mt-4 overflow-x-auto">
                  <table className="min-w-full border border-gray-200 text-sm">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="border border-gray-200 px-4 py-2 text-left font-semibold text-gray-700">Категория Субъекта ПД</th>
                        <th className="border border-gray-200 px-4 py-2 text-left font-semibold text-gray-700">Категория персональных данных</th>
                        <th className="border border-gray-200 px-4 py-2 text-left font-semibold text-gray-700">Цель</th>
                        <th className="border border-gray-200 px-4 py-2 text-left font-semibold text-gray-700">Способы обработки ПД</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr>
                        <td className="border border-gray-200 px-4 py-2 align-top">Пользователи Сайта</td>
                        <td className="border border-gray-200 px-4 py-2 align-top">Cookie-файлы, IP-адрес, а также иные данные, собираемые посредством метрических программ</td>
                        <td className="border border-gray-200 px-4 py-2 align-top">Продвижение услуг на рынке</td>
                        <td className="border border-gray-200 px-4 py-2 align-top text-xs">сбор, запись, систематизация, накопление, хранение, уточнение, извлечение, передача, использование, обезличивание, блокирование, удаление, уничтожение</td>
                      </tr>
                      <tr>
                        <td className="border border-gray-200 px-4 py-2 align-top">Контрагенты, представители контрагентов, клиенты</td>
                        <td className="border border-gray-200 px-4 py-2 align-top">Фамилия, имя, отчество, Адрес электронной почты, Номер телефона, Никнейм в Telegram, Банковские реквизиты, номер расчетного счета</td>
                        <td className="border border-gray-200 px-4 py-2 align-top">Подготовка, заключение и исполнение гражданско-правового договора</td>
                        <td className="border border-gray-200 px-4 py-2 align-top text-xs">сбор, запись, систематизация, накопление, хранение, уточнение, извлечение, передача, распространение, использование, обезличивание, блокирование, удаление, уничтожение</td>
                      </tr>
                      <tr>
                        <td className="border border-gray-200 px-4 py-2 align-top">Пользователи Сайта, контрагенты, клиенты</td>
                        <td className="border border-gray-200 px-4 py-2 align-top">Фамилия, имя, отчество, Адрес электронной почты, Номер телефона, Никнейм в Telegram</td>
                        <td className="border border-gray-200 px-4 py-2 align-top">Осуществление рассылки сообщений информационного и рекламного характера</td>
                        <td className="border border-gray-200 px-4 py-2 align-top text-xs">сбор, запись, систематизация, накопление, хранение, уточнение, извлечение, передача, использование, обезличивание, блокирование, удаление, уничтожение</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </li>

              <li>
                Я знаю, что Оператор обрабатывает персональные данные Субъекта только в случае их заполнения и (или) отправки Пользователем самостоятельно на Сайте, в переписке в мессенджере.
              </li>

              <li>
                В ходе обработки над персональными данными будут совершены следующие действия: сбор; запись; систематизация; накопление; хранение; уточнение (обновление, изменение); извлечение; использование; передача (предоставление, доступ, распространение); блокирование; удаление; уничтожение.
              </li>

              <li>
                В отношении несовершеннолетних лиц дача согласия на обработку персональных данных реализуется их законными представителями.
              </li>

              <li>
                Обработка персональных данных начинается с момента получения персональных данных. После заполнения всех данных, лицо, желающее получить Услугу (доступ к курсу), подтверждает правильность и достоверность указанных им данных и выражает желание подать заявку путем активации поля такого типа как «Далее», «Оплатить», «Зарегистрироваться» или иного, аналогичного ему по функциональному назначению. Любая информация, которая передана Оператору посредством информационно-телекоммуникационной сети «Интернет» либо иным способом, считается конфиденциальной.
              </li>

              <li>
                Оператором обеспечивается сохранность полученных персональных данных Субъектов персональных данных, за исключением добровольного предоставления Субъектом персональных данных информации о себе для общего доступа неограниченному кругу лиц.
              </li>

              <li>
                Срок, в течение которого действует согласие субъекта персональных данных: до достижения целей обработки персональных данных.
              </li>

              <li>
                Я ознакомлен(-а) с Политикой и Согласием посредством выражения согласия на Сайте. Я подтверждаю, что действую по своей воле и в своём интересе. Данное Согласие является конкретным, информированным, сознательным и однозначным.
              </li>

              <li>
                Согласие может быть отозвано субъектом персональных данных или его представителем путем направления письменного заявления по адресу электронной почты: <a href="mailto:inbox@radcop.online" className="text-rose-600 hover:underline font-medium">inbox@radcop.online</a>.
              </li>
            </ol>

            <div className="mt-12 pt-8 border-t border-gray-200 text-sm text-gray-500">
              <p>Дата последнего обновления: 14.10.2025</p>
            </div>
          </div>
        </div>
      </main>

      {/* Footer Simple */}
      <footer className="bg-white border-t border-gray-200 py-8">
        <div className="max-w-4xl mx-auto px-4 text-center text-sm text-gray-500">
          <p>&copy; {new Date().getFullYear()} РАД КОП. Все права защищены.</p>
        </div>
      </footer>
    </div>
  );
};

export default ConsentPage;