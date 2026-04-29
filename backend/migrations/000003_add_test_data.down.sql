-- 1. Удаляем покупки пользователей
DELETE FROM user_purchases
WHERE user_id IN (SELECT id FROM users WHERE email IN ('admin@test.ru', 'user@test.ru'));

-- 2. Удаляем уроки всех тестовых курсов
DELETE FROM lessons
WHERE module_id IN (
  SELECT id FROM modules WHERE course_id IN (
    SELECT id FROM courses WHERE title IN (
      'Женский курс', 
      'Введение в женское здоровье', 
      'Восстановление после родов',
      'Питание во время беременности',
      'Йога для беременных',
      'Психология материнства'
    )
  )
);

-- 3. Удаляем модули всех тестовых курсов
DELETE FROM modules
WHERE course_id IN (
  SELECT id FROM courses WHERE title IN (
    'Женский курс', 
    'Введение в женское здоровье', 
    'Восстановление после родов',
    'Питание во время беременности',
    'Йога для беременных',
    'Психология материнства'
  )
);

-- 4. Удаляем тестовые курсы (и старые, и новые названия)
DELETE FROM courses
WHERE title IN (
  'Женский курс', 
  'Введение в женское здоровье', 
  'Восстановление после родов',
  'Питание во время беременности',
  'Йога для беременных',
  'Психология материнства'
);

-- 5. Удаляем тестовых пользователей
DELETE FROM users
WHERE email IN ('admin@test.ru', 'user@test.ru');

