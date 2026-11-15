import PolicyPlaceholderPage from './PolicyPlaceholderPage';

const NewsletterConsentPage = () => (
  <PolicyPlaceholderPage
    title="Согласие на получение информационной и рекламной рассылки"
    text={
      <>
        <p className="mb-6">
          Проставляя галочку или иной знак в чек-боксе «согласен на получение рекламных рассылок»
          Вы даете свое согласие на получение информационной и рекламной рассылки от
          Индивидуального предпринимателя Пимченко Ольги Александровны (ИНН: 320801525917, ОГРНИП:
          326320000038641), касающейся информации о продуктах, услугах и всех дополнительных
          сведений, имеющих отношение к получению доступа, консультаций, а также всех сведений о
          других продуктах Администрации. Информационная и рекламная рассылка будет направляться с
          электронной почты Администрации на вашу электронную почту, которую вы оставили при
          заполнении заявки на сайте:{' '}
          <a
            href="https://olga-pimchenko.ru/"
            className="text-rose-500 hover:underline"
            style={{ textDecoration: 'none' }}
          >
            https://olga-pimchenko.ru/
          </a>{' '}
          или в виде смс-сообщения на номер телефона, указанный вами при заполнении данных в форме
          на сайте, или в другие мессенджеры.
        </p>

        <p className="mb-4">
          <strong>Цель обработки персональных данных:</strong> информирование пользователя о
          новостях, об услугах.
        </p>

        <p className="mb-4">
          В ходе обработки с персональными данными будут совершены следующие действия:
        </p>
        <ul className="list-none pl-6 space-y-2 mb-4 -ml-7">
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            сбор;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            запись;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            систематизация;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            накопление;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            хранение;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            уточнение (обновление, изменение);
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            использование;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            передача (предоставление, доступ);
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            блокирование;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            удаление;
          </li>
          <li className="relative pl-6">
            <span className="absolute left-0 text-rose-500 font-medium">—</span>
            уничтожение.
          </li>
        </ul>

        <p className="mb-4">
          Вы вправе в любой момент отписаться от рассылки, указав на это в ответном письме
          Администрации, либо направив отказ от рассылки на электронную почту Администрации:{' '}
          <a
            href="mailto:help-olga-pimchenko@yandex.ru"
            className="text-rose-500 hover:underline"
            style={{ textDecoration: 'none' }}
          >
            help-olga-pimchenko@yandex.ru
          </a>
          .
        </p>

        <p className="mb-4">
          В случае отзыва Согласия Администрация вправе продолжить обработку персональных данных без
          него при наличии оснований, указанных в пунктах 2 – 11 части 1 статьи 6, части 2 статьи 10
          и части 2 статьи 11 Федерального закона №152-ФЗ «О персональных данных» от 27.07.2006 г.
          Дата начала действия согласия – дата проставления галочки или иного знака.
        </p>
      </>
    }
  />
);

export default NewsletterConsentPage;