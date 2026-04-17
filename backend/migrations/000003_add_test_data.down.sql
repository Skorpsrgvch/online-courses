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
