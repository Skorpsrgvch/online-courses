import PolicyPlaceholderPage from './PolicyPlaceholderPage';

const ConsentPage = () => (
  <PolicyPlaceholderPage
    title="Согласие на обработку персональных данных"
    text={
      <>
        <p className="mb-6">
          Заполняя форму подбора курса на сайте{' '}
          <a
            href="https://olga-pimchenko.ru/"
            className="text-rose-500 hover:underline"
            style={{ textDecoration: 'none' }}
          >
            https://olga-pimchenko.ru/
          </a>
          , и проставляя отметку в соответствующем поле, я действую свободно, своей волей и в своем интересе,
          подтверждаю дееспособность и в соответствии со статьей 9 Федерального закона от 27.07.2006 № 152-ФЗ
          «О персональных данных» даю Индивидуальному предпринимателю Пимченко Ольге Александровне
          (далее — Оператор) конкретное, информированное, сознательное и однозначное согласие на обработку
          персональных данных на следующих условиях. Настоящее согласие оформлено в виде самостоятельного
          документа и не является частью иных документов в соответствии с ч. 1 ст. 9 ФЗ-152 (в ред. ФЗ
          от 24.06.2025 № 156-ФЗ). Проставление отметки в поле формы является простой электронной подписью
          и равнозначно собственноручной подписи субъекта персональных данных. Факт выражения согласия
          фиксируется автоматически с сохранением даты, времени и IP-адреса устройства.
        </p>

        <h2 className="text-xl font-semibold mt-8 mb-4">1. Перечень персональных данных</h2>
        <ul className="list-none pl-6 space-y-2 mb-4 -ml-7">
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            ФИО;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            адрес электронной почты;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            номер телефона.
          </li>
        </ul>

        <h2 className="text-xl font-semibold mt-8 mb-4">2. Цели обработки</h2>
        <ul className="list-none pl-6 space-y-2 mb-4 -ml-7">
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            подбор подходящего курса на основании запроса субъекта персональных данных;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            направление информации о подобранном курсе на указанный адрес электронной почты или номер телефона;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            обработка обращений и ответы на запросы субъекта персональных данных.
          </li>
        </ul>

        <h2 className="text-xl font-semibold mt-8 mb-4">3. Перечень действий с персональными данными</h2>
        <p className="mb-4">
          Сбор, запись, систематизация, накопление, хранение, уточнение (обновление, изменение), извлечение,
          использование, передача (предоставление, доступ), блокирование, удаление, уничтожение. Обработка
          осуществляется с использованием средств автоматизации и без таковых.
        </p>

        <h2 className="text-xl font-semibold mt-8 mb-4">4. Передача третьим лицам</h2>
        <p className="mb-4">
          Персональные данные могут быть переданы исключительно в целях, указанных в разделе 2,
          следующим лицам:
        </p>
        <ul className="list-none pl-6 space-y-2 mb-4 -ml-7">
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            <strong>ООО «Яндекс»</strong> (ИНН 7736207543, сервис Яндекс.Почта) — для отправки
            информационных писем, уведомлений и рассылок на адрес электронной почты субъекта
            персональных данных;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            <strong>ООО НКО «ЮМани»</strong> (ИНН 7750005725, ЮKassa) — для обработки платежей,
            формирования чеков и проведения возвратов денежных средств;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            государственным органам — в случаях, предусмотренных законодательством Российской Федерации.
          </li>
        </ul>
        <p className="mb-4">
          Передача персональных данных за пределы Российской Федерации не осуществляется.
          Передача иным третьим лицам без отдельного согласия субъекта не производится.
        </p>

        <h2 className="text-xl font-semibold mt-8 mb-4">5. Срок хранения</h2>
        <p className="mb-4">
          Персональные данные хранятся не более 1 (одного) года с даты получения настоящего согласия.
          По истечении указанного срока или в случае отзыва согласия персональные данные уничтожаются
          в течение 30 (тридцати) календарных дней, если их дальнейшее хранение не является обязательным
          в силу закона.
        </p>

        <h2 className="text-xl font-semibold mt-8 mb-4">6. Порядок отзыва согласия</h2>
        <p className="mb-4">
          Настоящее согласие может быть отозвано в любое время путем направления заявления в произвольной форме
          на{' '}
          <a
            href="mailto:help-olga-pimchenko@yandex.ru"
            className="text-rose-500 hover:underline"
            style={{ textDecoration: 'none' }}
          >
            help-olga-pimchenko@yandex.ru
          </a>{' '}
          с темой «Отзыв согласия на обработку персональных данных». Оператор прекратит обработку и уничтожит
          данные в течение 30 (тридцати) календарных дней с даты получения заявления. Отзыв согласия не влечет
          правовых последствий для субъекта ПД, за исключением невозможности получения услуг, для оказания
          которых данная обработка необходима.
        </p>

        <h2 className="text-xl font-semibold mt-8 mb-4">7. Права субъекта персональных данных</h2>
        <ul className="list-none pl-6 space-y-2 mb-4 -ml-7">
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            получать информацию о составе, целях и сроках обработки персональных данных (ст. 14 ФЗ-152).
            Запрос рассматривается в течение 10 рабочих дней;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            получить копию обрабатываемых персональных данных по письменному запросу;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            требовать уточнения, блокирования или уничтожения данных, если они являются неполными, устаревшими
            или незаконно полученными (ст. 21 ФЗ-152);
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            отозвать согласие в порядке, указанном в разделе 6 (ст. 9 ФЗ-152);
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            обжаловать действия Оператора в Роскомнадзор: rkn.gov.ru, тел. 8-800-222-27-88 (бесплатно),
            или в суд по месту своего жительства либо по месту нахождения Оператора (ст. 17 ФЗ-152).
          </li>
        </ul>

        <h2 className="text-xl font-semibold mt-8 mb-4">8. Заключительные положения</h2>
        <p className="mb-4">
          Подтверждаю достоверность предоставленных персональных данных. С Политикой оператора в отношении
          обработки персональных данных ознакомлен(а).
        </p>
      </>
    }
  />
);

export default ConsentPage;