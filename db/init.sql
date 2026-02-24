
-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,        -- bcrypt
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Таблица курсов
CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    is_public BOOLEAN NOT NULL DEFAULT true,   -- true = бесплатно, false = платный
    is_active BOOLEAN NOT NULL DEFAULT true,   -- можно скрыть без удаления
    price INT NOT NULL DEFAULT 0,              -- стоимость в рублях (0 = бесплатный)
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Таблица модулей (внутри курса)
CREATE TABLE modules (
    id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    "order" INTEGER NOT NULL
);

-- Таблица уроков
CREATE TABLE lessons (
    id SERIAL PRIMARY KEY,
    module_id INTEGER NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    video_embed_id VARCHAR(255) NOT NULL,  -- Rutube ID
    "order" INTEGER NOT NULL
);

-- Прогресс пользователя по урокам
CREATE TABLE user_progress (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    completed_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, lesson_id)
);

-- Отзывы на курсы (с модерацией и оценкой)
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5), 
    approved BOOLEAN NOT NULL DEFAULT false,  -- модерация
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, course_id)  -- один отзыв на курс
);

-- Покупки курсов (для платных)
CREATE TABLE user_purchases (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    purchased_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, course_id)
);

-- Индексы для производительности
CREATE INDEX idx_courses_author ON courses(author_id);
CREATE INDEX idx_modules_course ON modules(course_id);
CREATE INDEX idx_lessons_module ON lessons(module_id);
CREATE INDEX idx_progress_user ON user_progress(user_id);
CREATE INDEX idx_reviews_course ON reviews(course_id) WHERE approved = true;
CREATE INDEX idx_purchases_user ON user_purchases(user_id);

-- Вставка тестового администратора (мамы)
-- Пароль: "mamazdorovye123" → замени хэш на реальный через bcrypt!
INSERT INTO users (email, password_hash, full_name, role)
VALUES (
    'anna@example.com',
    '$2a$12$K4vZqJ7bNQd3XrGfT9jY6e8U1vWxYzA1B2C3D4E5F6G7H8I9J0K',
    'Анна Ивановна',
    'admin'
);