-- =============================================================================
-- 1. ПОЛЬЗОВАТЕЛИ
-- =============================================================================
INSERT INTO users (email, full_name, role, password_hash, created_at)
VALUES 
    ('admin@test.ru', 'Администратор', 'admin', '$2a$10$JH/Kx8uNH1qus2TfXupBRuVGeEZ90AW7vnYjt0nbDzOalRvuSnHIO', NOW()),
    ('user1@test.ru', 'Тестовый Пользователь1', 'user', '$2a$10$/60ZTbYEJNC5kyDGD0LzNuzAv0zq39I.xYTs5xNuH./n2Wzi6v95a', NOW()),
    ('user2@test.ru', 'Тестовый Пользователь2', 'user', '$2a$10$/60ZTbYEJNC5kyDGD0LzNuzAv0zq39I.xYTs5xNuH./n2Wzi6v95a', NOW()),
    ('user3@test.ru', 'Тестовый Пользователь3', 'user', '$2a$10$/60ZTbYEJNC5kyDGD0LzNuzAv0zq39I.xYTs5xNuH./n2Wzi6v95a', NOW()),
    ('user4@test.ru', 'Тестовый Пользователь4', 'user', '$2a$10$/60ZTbYEJNC5kyDGD0LzNuzAv0zq39I.xYTs5xNuH./n2Wzi6v95a', NOW())
ON CONFLICT (email) DO NOTHING;

-- =============================================================================
-- 2. КУРСЫ (Обновлено с новыми полями)
-- =============================================================================

-- Курс 1: Женский курс (Бесплатный)
INSERT INTO courses (title, description, is_public, price, author_id, is_active, cover_image_url, contraindications, recommendations, target_audience, course_basis, class_basis, bonuses)
SELECT 
    'Женский курс', 
    'Онлайн-курс для женщин, которые хотят улучшить своё здоровье и справиться с деликатными проблемами', 
    true, 0,
    (SELECT id FROM users WHERE email = 'admin@test.ru'), 
    true, 
    '/images/Course1.jpg',
    '❌ Беременность.|||❌ Первые 8 недель после родов.|||❌ Острые воспалительные процессы и острые боли.',
    '✅ Исключить тяжелые нагрузки.|||✅ Вставать и ложиться через бок.|||✅ Исключить подъем тяжестей свыше 3 кг.',
    'Вы чувствуете необходимость восстанавливаться после родов, даже если прошло много месяцев.|||У вас есть проблемы деликатного характера: подтекание мочи, хлюпающие звуки.|||Вам диагностировали опущение внутренних органов.|||Чувствуете тяжесть или давление внизу живота.|||Вас беспокоят боли в тазобедренных суставах, крестце, пояснице.|||У вас было кесарево сечение или другие операции.|||Вы хотите подготовиться к беременности.|||Вы хотите улучшить осанку и повысить энергию.|||У вас нет жалоб, но хочется заботиться о здоровье.',
    'Более 10 комплексов упражнений.|||Комплекс самомассажа живота.|||Занятие «Работа со шрамом после КС».|||Приятные бонусы.',
    'Комплексы ЛФК при дисфункциях МТД и ОДА.|||Реабилитационные техники.|||Техники Мио-Фасциального Релиза (МФР).|||Техники Пост-Изометрической Релаксации (ПИР).|||Пилатес.|||Йогатерапия.',
    '[{"title": "МЕНОПАУЗА", "description": "Особенности периода", "icon": "video"}, {"title": "ПЕРЕД СНОМ", "description": "Практика 2-3 мин", "icon": "video"}, {"title": "САМОМАССАЖ", "description": "Улучшает лимфоток", "icon": "video"}]'::jsonb
WHERE NOT EXISTS (SELECT 1 FROM courses WHERE title = 'Женский курс');

-- Курс 2: Восстановление после родов (Платный)
INSERT INTO courses (title, description, is_public, price, author_id, is_active, cover_image_url, contraindications, recommendations, target_audience, course_basis, class_basis, bonuses)
SELECT 
    'Восстановление после родов', 
    'Комплексная программа восстановления после родов.', 
    false, 4900,
    (SELECT id FROM users WHERE email = 'admin@test.ru'), 
    true, 
    '/images/Course2.jpg',
    '❌ Кесарево сечение (менее 2 месяцев).|||❌ Активное кровотечение.',
    '✅ Начинать занятия после осмотра гинеколога.|||✅ Слушать свое тело.',
    'Прошло мало времени после родов.|||Есть диастаз или боли в спине.|||Хочется вернуть форму безопасно.',
    'Поэтапная программа восстановления.|||Работа с диастазом.|||Массаж шрамов.',
    'Дыхательные практики.|||Активация глубоких мышц живота.|||Безопасные упражнения для спины.',
    '[{"title": "Вебинар", "description": "Разбор ошибок", "icon": "video"}]'::jsonb
WHERE NOT EXISTS (SELECT 1 FROM courses WHERE title = 'Восстановление после родов');

-- Курс 3: Здоровая спина мамы (Платный)
INSERT INTO courses (title, description, is_public, price, author_id, is_active, cover_image_url, contraindications, recommendations, target_audience, course_basis, class_basis, bonuses)
SELECT
    'Здоровая спина мамы',
    'Избавьтесь от боли в спине и пояснице после родов.',
    false, 2900,
    (SELECT id FROM users WHERE email = 'admin@test.ru'),
    true,
    '/images/Course3.jpg',
    '❌ Острая боль в позвоночнике.|||❌ Травмы спины.',
    '✅ Делать упражнения плавно.|||✅ Использовать коврик.',
    'Болит спина от ношения ребенка на руках.|||Появилась сутулость.|||Чувствуете скованность в пояснице.',
    'Упражнения у стены.|||Работа с мышцами кора.|||Расслабление грудного отдела.',
    'МФР с мячом.|||Укрепление ягодичных мышц.|||Правильная биомеханика движений.',
    '[{"title": "Чек-лист", "description": "Эргономика быта", "icon": "file"}]'::jsonb
WHERE NOT EXISTS (SELECT 1 FROM courses WHERE title = 'Здоровая спина мамы');

-- Курс 4: Сила интимных мышц (Платный)
INSERT INTO courses (title, description, is_public, price, author_id, is_active, cover_image_url, contraindications, recommendations, target_audience, course_basis, class_basis, bonuses)
SELECT
    'Сила интимных мышц',
    'Эффективные техники для укрепления мышц тазового дна.',
    false, 3500,
    (SELECT id FROM users WHERE email = 'admin@test.ru'),
    true,
    '/images/Course4.jpg',
    '❌ Воспалительные процессы.|||❌ Первые 2 месяца после операций.',
    '✅ Регулярность занятий.|||✅ Полное расслабление между подходами.',
    'Есть подтекание мочи при кашле или смехе.|||Хочется улучшить качество интимной жизни.|||Профилактика опущений.',
    'Анатомия мышц тазового дна.|||Техника изолированного сокращения.|||Координация с дыханием.',
    'Быстрые сокращения.|||Длительное удержание.|||Техника "Лифт".',
    '[{"title": "Аудиогид", "description": "Дыхательные практики", "icon": "audio"}]'::jsonb
WHERE NOT EXISTS (SELECT 1 FROM courses WHERE title = 'Сила интимных мышц');

-- Курс 5: Легкая беременность (Бесплатный)
INSERT INTO courses (title, description, is_public, price, author_id, is_active, cover_image_url, contraindications, recommendations, target_audience, course_basis, class_basis, bonuses)
SELECT
    'Легкая беременность',
    'Безопасная гимнастика для подготовки тела к родам.',
    true, 0,
    (SELECT id FROM users WHERE email = 'admin@test.ru'),
    true,
    '/images/Course1.jpg',
    '❌ Угроза прерывания.|||❌ Предлежание плаценты.',
    '✅ Заниматься в хорошем самочувствии.|||✅ Следить за пульсом.',
    'Хочется подготовить тело к родам.|||Беспокоят отеки или боли в спине.|||Нужно научиться правильно дышать.',
    'Упражнения для каждого триместра.|||Подготовка промежности.|||Дыхание в родах.',
    'Щадящая гимнастика.|||Работа с тазом.|||Расслабление.',
    '[{"title": "Гайд", "description": "Список вещей в роддом", "icon": "file"}]'::jsonb
WHERE NOT EXISTS (SELECT 1 FROM courses WHERE title = 'Легкая беременность');

-- =============================================================================
-- 3. МОДУЛИ И УРОКИ: ЖЕНСКИЙ КУРС
-- =============================================================================
-- Неделя 1
INSERT INTO modules (course_id, title, "order") SELECT id, 'Неделя 1', 1 FROM courses WHERE title = 'Женский курс';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order")
SELECT m.id, 'Дыхание. Грудная клетка.', 'Первое занятие посвящено дыханию.', '7e8c5681487cd383d46a513aee7d0601', 'nJeUFf03iW39UEyULnRh2w', 1 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 1';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order")
SELECT m.id, 'Работа с перикардом.', 'Мячик не нужен.', '70ef5b8c45818ee28152174035fcf490', '9h4lav1VzpqYCpRsZggeww', 2 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 1';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order")
SELECT m.id, 'Работа с языком.', 'Нужна салфетка.', 'a34d5132da9445515f6948846a6b8f17', 'PvmACTvui-OzC2Y1_sVgvg', 3 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 1';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order")
SELECT m.id, 'Практика МТД (часть 2).', 'Нужна опора под таз.', '5bdb9c6e72a90f207714b86088055930', 'g7uHW3rEFXySf-iKwfMhWw', 4 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 1';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order")
SELECT m.id, 'Женская практика на каждый день.', 'Нужен болстер.', 'ced130ddbd12a767f04e8e0f5b73eb33', '7d63pyz-HbvUAdQHJc7nGw', 5 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 1';

-- Неделя 2
INSERT INTO modules (course_id, title, "order") SELECT id, 'Неделя 2', 2 FROM courses WHERE title = 'Женский курс';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'Антидиастаз.', 'Нужен ремень.', 'e44cd2c3ff35ef42295da830f480030e', 'BwqaCaxopg0y6mvrTzNe3w', 1 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 2';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'Стабильный таз.', 'Нужны мячики.', '850e83ac57e788316f71366755a64dd5', 'KApx6z5iK0AcLE55kAZO_g', 2 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 2';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'Свободный позвоночник.', 'Нужен ремень.', '5c3b90c21f87403581ccd51d8f6f41e7', '2felxxdRIdIdybZJCfSSYA', 3 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 2';

-- Неделя 3
INSERT INTO modules (course_id, title, "order") SELECT id, 'Неделя 3', 3 FROM courses WHERE title = 'Женский курс';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'Интеграция.', 'Сокращения МТД.', '209de268fb6ab2e4d2298bf018751bdb', 'GSw-B5n8Fy4gP1YxtI1ndA', 1 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 3';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'Терапевтические проходки.', 'Нужен блок.', '4bc97662acb81aa004c24674760255b8', 'DS8cprMTxk5wc1usE-znHg', 2 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 3';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'Работа у стены.', 'Нужен болстер.', 'c5e463257d3abc31daba0a582c6d9c29', '33r-JQgOlLt2nqp2trrscQ', 3 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 3';

-- Неделя 4
INSERT INTO modules (course_id, title, "order") SELECT id, 'Неделя 4', 4 FROM courses WHERE title = 'Женский курс';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'Правка копчика.', 'Нужен болстер.', '951325b973a117f33a534cf6f53a649a', 'sNl5tNFIL3C-BEsc9i334g', 1 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 4';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'На четвереньках.', '', '82454a9bcd472f5b2a9292ce62fdea0c', 'PH9jPx1hwMBSmZXA9IyiCA', 2 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 4';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'На стуле.', 'Нужен шарф.', 'ab617f1ce081dfd08bb672f49eafddb9', 'Rvo08xFghYk5nl4swRnn5A', 3 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 4';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'Стопы.', 'Нужен мячик.', 'b1cd3963f4e27239c0344b811162b2ba', 'jQtWg6U4frOd-vjoUUyS3w', 4 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Неделя 4';

-- Работа со шрамами
INSERT INTO modules (course_id, title, "order") SELECT id, 'Работа со шрамами', 5 FROM courses WHERE title = 'Женский курс';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'Самомассаж.', '', '54ee29ae0f51d5997562c5495ae35bcc', 'BDz729GpOPAL7ply5ESVKA', 1 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Работа со шрамами';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'Второй этап.', '', 'cce7262fc30fb1f0c4bb9d65999b3593', 'sUGh73pzwFLav450ugIADg', 2 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Работа со шрамами';

-- Бонусы
INSERT INTO modules (course_id, title, "order") SELECT id, 'Бонусы', 6 FROM courses WHERE title = 'Женский курс';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'МЕНОПАУЗА.', '', 'da3452f6f267137e8098402f08b92cec', 'IRbf1CgojA-7kgE7kCrPEA', 1 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Бонусы';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'ПЕРЕД СНОМ.', 'Опора под таз.', 'b50f41b7fcb2cf03123c2d088d576f9a', 'V68r3_qWr7f0_9OaWd_ntg', 2 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Бонусы';
INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") SELECT m.id, 'Самомассаж живота.', 'От запоров.', '832fa6bf43064cbdd1302a201fd885db', 'AkmrgD5VemgGkJbXA894lA', 3 FROM modules m JOIN courses c ON m.course_id = c.id WHERE c.title = 'Женский курс' AND m.title = 'Бонусы';

-- =============================================================================
-- 4. МОДУЛИ И УРОКИ: ВОССТАНОВЛЕНИЕ ПОСЛЕ РОДОВ
-- =============================================================================
DO $$
DECLARE c_id int; m_id int;
BEGIN
    SELECT id INTO c_id FROM courses WHERE title = 'Восстановление после родов';
    
    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Первые дни', 1) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Дыхание и таз.', 'База.', '7e8c5681487cd383d46a513aee7d0601', 'nJeUFf03iW39UEyULnRh2w', 1),
    (m_id, 'Активация низа живота.', 'Безопасно.', '70ef5b8c45818ee28152174035fcf490', '9h4lav1VzpqYCpRsZggeww', 2),
    (m_id, 'Расслабление спины.', 'Снятие зажимов.', 'a34d5132da9445515f6948846a6b8f17', 'PvmACTvui-OzC2Y1_sVgvg', 3);

    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Диастаз', 2) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Что такое диастаз.', 'Теория.', 'e44cd2c3ff35ef42295da830f480030e', 'BwqaCaxopg0y6mvrTzNe3w', 1),
    (m_id, 'Запрещенные упражнения.', 'Чего избегать.', '850e83ac57e788316f71366755a64dd5', 'KApx6z5iK0AcLE55kAZO_g', 2),
    (m_id, 'Безопасный пресс.', 'Техники.', '5c3b90c21f87403581ccd51d8f6f41e7', '2felxxdRIdIdybZJCfSSYA', 3);

    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Шрам от КС', 3) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Массаж шрама.', 'Этап 1.', '54ee29ae0f51d5997562c5495ae35bcc', 'BDz729GpOPAL7ply5ESVKA', 1),
    (m_id, 'Мобилизация тканей.', 'Этап 2.', 'cce7262fc30fb1f0c4bb9d65999b3593', 'sUGh73pzwFLav450ugIADg', 2),
    (m_id, 'Интеграция в движение.', 'Финал.', '209de268fb6ab2e4d2298bf018751bdb', 'GSw-B5n8Fy4gP1YxtI1ndA', 3);

    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Функционал', 4) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Подъем ребенка.', 'Биомеханика.', '4bc97662acb81aa004c24674760255b8', 'DS8cprMTxk5wc1usE-znHg', 1),
    (m_id, 'Прогулка с коляской.', 'Осанка.', '951325b973a117f33a534cf6f53a649a', 'sNl5tNFIL3C-BEsc9i334g', 2),
    (m_id, 'Бытовые дела.', 'Эргономика.', 'b1cd3963f4e27239c0344b811162b2ba', 'jQtWg6U4frOd-vjoUUyS3w', 3);
END $$;

-- =============================================================================
-- 5. МОДУЛИ И УРОКИ: ЗДОРОВАЯ СПИНА МАМЫ
-- =============================================================================
DO $$
DECLARE c_id int; m_id int;
BEGIN
    SELECT id INTO c_id FROM courses WHERE title = 'Здоровая спина мамы';
    
    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Мобилизация', 1) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Кошка-Корова.', 'База.', '7e8c5681487cd383d46a513aee7d0601', 'nJeUFf03iW39UEyULnRh2w', 1),
    (m_id, 'Скручивания лежа.', 'Грудной отдел.', '70ef5b8c45818ee28152174035fcf490', '9h4lav1VzpqYCpRsZggeww', 2),
    (m_id, 'Раскрытие груди.', 'Против сутулости.', 'a34d5132da9445515f6948846a6b8f17', 'PvmACTvui-OzC2Y1_sVgvg', 3);

    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Стабилизация', 2) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Планка у стены.', 'Безопасно.', 'e44cd2c3ff35ef42295da830f480030e', 'BwqaCaxopg0y6mvrTzNe3w', 1),
    (m_id, 'Ягодичный мост.', 'Укрепление.', '850e83ac57e788316f71366755a64dd5', 'KApx6z5iK0AcLE55kAZO_g', 2),
    (m_id, 'Птица-Собака.', 'Баланс.', '5c3b90c21f87403581ccd51d8f6f41e7', '2felxxdRIdIdybZJCfSSYA', 3);

    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Расслабление', 3) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'МФР с мячом.', 'Спина.', '54ee29ae0f51d5997562c5495ae35bcc', 'BDz729GpOPAL7ply5ESVKA', 1),
    (m_id, 'Дыхание животом.', 'Снятие стресса.', 'cce7262fc30fb1f0c4bb9d65999b3593', 'sUGh73pzwFLav450ugIADg', 2),
    (m_id, 'Шавасана.', 'Полный релакс.', '209de268fb6ab2e4d2298bf018751bdb', 'GSw-B5n8Fy4gP1YxtI1ndA', 3);

    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Быт', 4) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Как носить малыша.', 'Правила.', '4bc97662acb81aa004c24674760255b8', 'DS8cprMTxk5wc1usE-znHg', 1),
    (m_id, 'Кормление без боли.', 'Позы.', '951325b973a117f33a534cf6f53a649a', 'sNl5tNFIL3C-BEsc9i334g', 2),
    (m_id, 'Сон и отдых.', 'Организация.', 'b1cd3963f4e27239c0344b811162b2ba', 'jQtWg6U4frOd-vjoUUyS3w', 3);
END $$;

-- =============================================================================
-- 6. МОДУЛИ И УРОКИ: СИЛА ИНТИМНЫХ МЫШЦ
-- =============================================================================
DO $$
DECLARE c_id int; m_id int;
BEGIN
    SELECT id INTO c_id FROM courses WHERE title = 'Сила интимных мышц';
    
    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Анатомия', 1) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Где они?', 'Поиск мышц.', '7e8c5681487cd383d46a513aee7d0601', 'nJeUFf03iW39UEyULnRh2w', 1),
    (m_id, 'Типы волокон.', 'Теория.', '70ef5b8c45818ee28152174035fcf490', '9h4lav1VzpqYCpRsZggeww', 2),
    (m_id, 'Диафрагма и МТД.', 'Связь.', 'a34d5132da9445515f6948846a6b8f17', 'PvmACTvui-OzC2Y1_sVgvg', 3);

    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Навык', 2) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Изолированное сокращение.', 'Техника.', 'e44cd2c3ff35ef42295da830f480030e', 'BwqaCaxopg0y6mvrTzNe3w', 1),
    (m_id, 'Полное расслабление.', 'Важно.', '850e83ac57e788316f71366755a64dd5', 'KApx6z5iK0AcLE55kAZO_g', 2),
    (m_id, 'Координация.', 'Вдох-выдох.', '5c3b90c21f87403581ccd51d8f6f41e7', '2felxxdRIdIdybZJCfSSYA', 3);

    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Сила', 3) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Быстрые сокращения.', 'Тренировка.', '54ee29ae0f51d5997562c5495ae35bcc', 'BDz729GpOPAL7ply5ESVKA', 1),
    (m_id, 'Длительное удержание.', 'Выносливость.', 'cce7262fc30fb1f0c4bb9d65999b3593', 'sUGh73pzwFLav450ugIADg', 2),
    (m_id, 'Лифт.', 'Продвинутый уровень.', '209de268fb6ab2e4d2298bf018751bdb', 'GSw-B5n8Fy4gP1YxtI1ndA', 3);

    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Интеграция', 4) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'The Knack.', 'Защита при нагрузке.', '4bc97662acb81aa004c24674760255b8', 'DS8cprMTxk5wc1usE-znHg', 1),
    (m_id, 'В быту.', 'Применение.', '951325b973a117f33a534cf6f53a649a', 'sNl5tNFIL3C-BEsc9i334g', 2),
    (m_id, 'В спорте.', 'Защита.', 'b1cd3963f4e27239c0344b811162b2ba', 'jQtWg6U4frOd-vjoUUyS3w', 3);
END $$;

-- =============================================================================
-- 7. МОДУЛИ И УРОКИ: ЛЕГКАЯ БЕРЕМЕННОСТЬ
-- =============================================================================
DO $$
DECLARE c_id int; m_id int;
BEGIN
    SELECT id INTO c_id FROM courses WHERE title = 'Легкая беременность';
    
    INSERT INTO modules (course_id, title, "order") VALUES (c_id, '1 Триместр', 1) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Адаптация.', 'Щадящий режим.', '7e8c5681487cd383d46a513aee7d0601', 'nJeUFf03iW39UEyULnRh2w', 1),
    (m_id, 'Дыхание.', 'От токсикоза.', '70ef5b8c45818ee28152174035fcf490', '9h4lav1VzpqYCpRsZggeww', 2),
    (m_id, 'Расслабление.', 'Сон.', 'a34d5132da9445515f6948846a6b8f17', 'PvmACTvui-OzC2Y1_sVgvg', 3);

    INSERT INTO modules (course_id, title, "order") VALUES (c_id, '2 Триместр', 2) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Укрепление ног.', 'Против отеков.', 'e44cd2c3ff35ef42295da830f480030e', 'BwqaCaxopg0y6mvrTzNe3w', 1),
    (m_id, 'Грудной отдел.', 'Подготовка.', '850e83ac57e788316f71366755a64dd5', 'KApx6z5iK0AcLE55kAZO_g', 2),
    (m_id, 'Таз.', 'Мобильность.', '5c3b90c21f87403581ccd51d8f6f41e7', '2felxxdRIdIdybZJCfSSYA', 3);

    INSERT INTO modules (course_id, title, "order") VALUES (c_id, '3 Триместр', 3) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Облегчение спины.', 'С опорой.', '54ee29ae0f51d5997562c5495ae35bcc', 'BDz729GpOPAL7ply5ESVKA', 1),
    (m_id, 'Промежность.', 'Подготовка к родам.', 'cce7262fc30fb1f0c4bb9d65999b3593', 'sUGh73pzwFLav450ugIADg', 2),
    (m_id, 'Позиции для родов.', 'Практика.', '209de268fb6ab2e4d2298bf018751bdb', 'GSw-B5n8Fy4gP1YxtI1ndA', 3);

    INSERT INTO modules (course_id, title, "order") VALUES (c_id, 'Бонусы', 4) RETURNING id INTO m_id;
    INSERT INTO lessons (module_id, title, description, video_embed_id, private_key, "order") VALUES
    (m_id, 'Схватки.', 'Дыхание.', '4bc97662acb81aa004c24674760255b8', 'DS8cprMTxk5wc1usE-znHg', 1),
    (m_id, 'Партнерские роды.', 'Советы.', '951325b973a117f33a534cf6f53a649a', 'sNl5tNFIL3C-BEsc9i334g', 2),
    (m_id, 'Первый час.', 'Контакт.', 'b1cd3963f4e27239c0344b811162b2ba', 'jQtWg6U4frOd-vjoUUyS3w', 3);
END $$;

-- =============================================================================
-- 8. ПОКУПКИ И ОТЗЫВЫ
-- =============================================================================

-- Покупка
INSERT INTO user_purchases (user_id, course_id)
SELECT u.id, c.id
FROM users u, courses c
WHERE u.email = 'user1@test.ru'
  AND c.title = 'Восстановление после родов'
  AND NOT EXISTS (SELECT 1 FROM user_purchases p WHERE p.user_id = u.id AND p.course_id = c.id);

-- Отзывы (упрощенная вставка по одному, чтобы избежать ошибок)
INSERT INTO reviews (user_id, course_id, text, rating, approved, created_at)
SELECT u.id, c.id, 'Отличный курс! Очень помогло восстановить силы.', 5, true, NOW()
FROM users u, courses c WHERE u.email = 'user@test.ru' AND c.title = 'Женский курс'
ON CONFLICT (user_id, course_id) DO UPDATE SET text = EXCLUDED.text;

INSERT INTO reviews (user_id, course_id, text, rating, approved, created_at)
SELECT u.id, c.id, 'Наконец-то прошла боль в пояснице!', 5, true, NOW()
FROM users u, courses c WHERE u.email = 'user@test.ru' AND c.title = 'Здоровая спина мамы'
ON CONFLICT (user_id, course_id) DO UPDATE SET text = EXCLUDED.text;

INSERT INTO reviews (user_id, course_id, text, rating, approved, created_at)
SELECT u.id, c.id, 'Результат почувствовала через 2 недели.', 5, true, NOW()
FROM users u, courses c WHERE u.email = 'user@test.ru' AND c.title = 'Сила интимных мышц'
ON CONFLICT (user_id, course_id) DO UPDATE SET text = EXCLUDED.text;

INSERT INTO reviews (user_id, course_id, text, rating, approved, created_at)
SELECT u.id, c.id, 'Курс помог подготовиться к родам.', 5, true, NOW()
FROM users u, courses c WHERE u.email = 'user@test.ru' AND c.title = 'Легкая беременность'
ON CONFLICT (user_id, course_id) DO UPDATE SET text = EXCLUDED.text;

INSERT INTO reviews (user_id, course_id, text, rating, approved, created_at)
SELECT u.id, c.id, 'Программа восстановления просто спасение.', 5, true, NOW()
FROM users u, courses c WHERE u.email = 'user@test.ru' AND c.title = 'Восстановление после родов'
ON CONFLICT (user_id, course_id) DO UPDATE SET text = EXCLUDED.text;

INSERT INTO reviews (user_id, course_id, text, rating, approved, created_at)
SELECT u.id, c.id, 'Крутой курс.', 4, true, NOW()
FROM users u, courses c WHERE u.email = 'user2@test.ru' AND c.title = 'Женский курс'
ON CONFLICT (user_id, course_id) DO UPDATE SET text = EXCLUDED.text;

INSERT INTO reviews (user_id, course_id, text, rating, approved, created_at)
SELECT u.id, c.id, 'Не понравилось.', 2, true, NOW()
FROM users u, courses c WHERE u.email = 'user3@test.ru' AND c.title = 'Женский курс'
ON CONFLICT (user_id, course_id) DO UPDATE SET text = EXCLUDED.text;

INSERT INTO reviews (user_id, course_id, text, rating, approved, created_at)
SELECT u.id, c.id, 'Можно было и лучше', 3, true, NOW()
FROM users u, courses c WHERE u.email = 'user4@test.ru' AND c.title = 'Женский курс'
ON CONFLICT (user_id, course_id) DO UPDATE SET text = EXCLUDED.text;

INSERT INTO services (title, description, price, duration_minutes) VALUES
(
    'Индивидуальная консультация/занятие',
    'Персональная работа с тренером-реабилитологом для решения конкретных проблем здоровья. Разбор симптомов, диагностика состояния и составление четкого плана действий.|||
     С каким запросом ко мне можно обратиться:|||
     Реабилитация мышц тазового дна. Необходима при таких симптомах, как недержании мочи, опущении тазовых органов, 
     боль во время полового акта, попадание воздуха и др. симптомах дисфункции мышц тазового дна.|||
     Восстановление после родов ЕР и КС. Коррекция диастаза.|||
     Физическая реабилитация после абдоминальных и мочеполовых операций. Работа с рубцом КС и после эпизиотомии.|||
     Синдром хронической тазовой боли (СХТБ) у женщин и мужчин. Физическая терапия при:|||
        - диспареунии (боль во время полового акта),|||
        - вульводинии (боль в области наружных половых органов),|||
        - кокцигодинии (боль в копчике),|||
        - пудендальной невралгии,|||
        - миалгии напряжения мышц тазового дна и абдоминальной области и др.|||
     Физическая реабилитация при РМЖ.
     Необходима при подготовке к лечению, а также после любого вида лечения. 
     При ограничении подвижности в плече со стороны операции, болях, отеках, web-syndrome (тяже), для мобилизации послеоперационного шва.',
    4000,
    60
),
(
    'Бесплатная диагностическая консультация',
    'Короткая онлайн-встреча для знакомства и определения вектора работы. Вы рассказываете о проблеме, я оцениваю ситуацию и предлагаю варианты помощи.
     Идеально, если вы сомневаетесь, какой формат вам подойдет.|||
     Обсуждение ваших симптомов и целей.|||
     Оценка необходимости очной или онлайн-работы.|||
     Первичные рекомендации по облегчению состояния.|||
     Ответы на организационные вопросы.',
    0,
    20
),
(
    'Физическая реабилитация при РМЖ',
    'Комплексная программа физической реабилитации для женщин после лечения рака молочной железы. 
     Поможет восстановить подвижность плеча, справиться с отёками и рубцами, вернуть уверенность в своём теле. 
     Индивидуальный подход с учётом диагноза и этапа лечения.|||
     Составление комплексов упражнений с учетом диагноза и проведенного лечения, а также вашего запроса.|||
     Восстановление функциональной активности с учетом риска лимфостаза.|||
     Работа с самыми частыми жалобами и состояниями на фоне лечения:|||
        - Постмастэктомический синдром (отёки, ограничения движения),|||
        - Web-syndrome (тяж), отеки, грубые рубцы,|||
        - Контрактура плечевого сустава (снижение объёма движения).|||
     Кинезиотейпирование.',
    2700,
    60
),
(
    'Послеродовый патронаж для мамы',
    'Первые дни и недели после родов являются ключевыми для послеродового восстановления матери, 
     поскольку они закладывают прочную основу для долгосрочного здоровья и благополучия.
    Подходит для женщин после родов от 3 дней до 6 месяцев. Как после естественных родов, так и после кесарева сечения.|||
    Что входит в патронаж для матери:|||
        - Мягкие мануальные техники для восстановления положения внутренних органов,|||
        - Дыхательные техники, приводящие в баланс нервную систему мамы,|||
        - физическая реабилитация после родов|||
        - практики, помогающие наладить грудное вскармливание, стимулировать выработку молока, 
          улучшить лимфоток в области груди и грудной клетки, что служит профилактикой и лечением лактостаза.|||
        -  массаж шва после КС|||
        -  обучение самомассажу шва после КС|||
        -  рекомендации по питанию в послеродовый период|||
        -  кинезиотейпирование живота при диастазе, слабости мышц брюшной стенки, тейпирование рубца КС.|||
    Все практики совместимы с грудным вскармливанием, могут выполняться лежа в постели, а также совместно с ребенком.|||
    Если вы хотите грамотно восстановится после родов, но не знаете, как, то я помогу вам в этом.|||
    Услуга доступна с выездом на дом.',
    4000,
    90
);