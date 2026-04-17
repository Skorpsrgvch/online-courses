-- +migrate Up

-- Тестовый администратор (id=1)
-- Email: admin@test.ru, Пароль: admin123
INSERT INTO users (email, full_name, role, password_hash, created_at)
VALUES (
    'admin@test.ru',
    'Администратор',
    'admin',
    '$2a$10$JH/Kx8uNH1qus2TfXupBRuVGeEZ90AW7vnYjt0nbDzOalRvuSnHIO',
    NOW()
) ON CONFLICT DO NOTHING;

-- Тестовый пользователь (id=2)
-- Email: user@test.ru, Пароль: user123
INSERT INTO users (email, full_name, role, password_hash, created_at)
VALUES (
    'user@test.ru',
    'Тестовый Пользователь',
    'user',
    '$2a$10$/60ZTbYEJNC5kyDGD0LzNuzAv0zq39I.xYTs5xNuH./n2Wzi6v95a',
    NOW()
) ON CONFLICT DO NOTHING;

-- Бесплатный курс
INSERT INTO courses (title, description, is_public, price, author_id, is_active, cover_image_url)
SELECT 'Введение в женское здоровье', 'Бесплатный вводный курс об основах женского здоровья.', true, 0,
       (SELECT id FROM users WHERE email = 'admin@test.ru'), true, ''
WHERE NOT EXISTS (SELECT 1 FROM courses WHERE title = 'Введение в женское здоровье');

-- Платный курс
INSERT INTO courses (title, description, is_public, price, author_id, is_active, cover_image_url)
SELECT 'Восстановление после родов', 'Комплексная программа восстановления после родов.', false, 4900,
       (SELECT id FROM users WHERE email = 'admin@test.ru'), true, ''
WHERE NOT EXISTS (SELECT 1 FROM courses WHERE title = 'Восстановление после родов');

-- Модуль для бесплатного курса
INSERT INTO modules (course_id, title, "order")
SELECT c.id, 'Основы анатомии', 1
FROM courses c
WHERE c.title = 'Введение в женское здоровье'
  AND NOT EXISTS (SELECT 1 FROM modules m WHERE m.course_id = c.id AND m.title = 'Основы анатомии');

-- Модуль для платного курса
INSERT INTO modules (course_id, title, "order")
SELECT c.id, 'Первые недели после родов', 1
FROM courses c
WHERE c.title = 'Восстановление после родов'
  AND NOT EXISTS (SELECT 1 FROM modules m WHERE m.course_id = c.id AND m.title = 'Первые недели после родов');

-- Уроки для бесплатного курса — модуль «Основы анатомии»
-- Урок 1: видео
INSERT INTO lessons (module_id, title, description, lesson_type, video_embed_id, "order")
SELECT m.id, 'Введение в курс', 'Обзор программы и цели курса', 'video', 'abc123def', 1
FROM modules m
WHERE m.title = 'Основы анатомии'
  AND NOT EXISTS (SELECT 1 FROM lessons l WHERE l.module_id = m.id AND l.title = 'Введение в курс');

-- Урок 2: видео
INSERT INTO lessons (module_id, title, description, lesson_type, video_embed_id, "order")
SELECT m.id, 'Анатомия тазового дна', 'Базовое строение и функции мышц тазового дна', 'video', 'ghi456jkl', 2
FROM modules m
WHERE m.title = 'Основы анатомии'
  AND NOT EXISTS (SELECT 1 FROM lessons l WHERE l.module_id = m.id AND l.title = 'Анатомия тазового дна');

-- Урок 3: статья
INSERT INTO lessons (module_id, title, description, lesson_type, article_content, "order")
SELECT m.id, '10 советов для здоровья', 'Практические рекомендации', 'article',
       '<h2>10 простых советов</h2><p>Следуйте этим правилам каждый день:</p><ul><li>Пейте достаточно воды</li><li>Делайте упражнения Кегеля</li><li>Следите за осанкой</li><li>Высыпайтесь</li></ul>', 3
FROM modules m
WHERE m.title = 'Основы анатомии'
  AND NOT EXISTS (SELECT 1 FROM lessons l WHERE l.module_id = m.id AND l.title = '10 советов для здоровья');

-- Уроки для платного курса — модуль «Первые недели после родов»
INSERT INTO lessons (module_id, title, description, lesson_type, video_embed_id, "order")
SELECT m.id, 'Первые упражнения', 'Безопасные упражнения для первой недели', 'video', 'mno789pqr', 1
FROM modules m
WHERE m.title = 'Первые недели после родов'
  AND NOT EXISTS (SELECT 1 FROM lessons l WHERE l.module_id = m.id AND l.title = 'Первые упражнения');

-- Покупка платного курса пользователем
INSERT INTO user_purchases (user_id, course_id)
SELECT u.id, c.id
FROM users u, courses c
WHERE u.email = 'user@test.ru'
  AND c.title = 'Восстановление после родов'
  AND NOT EXISTS (SELECT 1 FROM user_purchases p WHERE p.user_id = u.id AND p.course_id = c.id);

-- +migrate Down

DELETE FROM user_purchases
WHERE user_id IN (SELECT id FROM users WHERE email IN ('admin@test.ru', 'user@test.ru'));

DELETE FROM lessons
WHERE module_id IN (
  SELECT id FROM modules WHERE course_id IN (
    SELECT id FROM courses WHERE title IN ('Введение в женское здоровье', 'Восстановление после родов')
  )
);

DELETE FROM modules
WHERE course_id IN (
  SELECT id FROM courses WHERE title IN ('Введение в женское здоровье', 'Восстановление после родов')
);

DELETE FROM courses
WHERE title IN ('Введение в женское здоровье', 'Восстановление после родов');

DELETE FROM users
WHERE email IN ('admin@test.ru', 'user@test.ru');
