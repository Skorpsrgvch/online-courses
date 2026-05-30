-- Включаем расширение для UUID (уже есть в современных PG, но на всякий случай)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. ENUM типы для строгой типизации (быстрее и надежнее VARCHAR CHECK)
CREATE TYPE user_role AS ENUM ('admin', 'user');
CREATE TYPE payment_status AS ENUM ('pending', 'waiting_for_capture', 'succeeded', 'canceled', 'failed');

-- 2. Справочник пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'user', -- Используем ENUM
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Каталог курсов
CREATE TABLE IF NOT EXISTS courses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    cover_image_url VARCHAR(500), 
    price BIGINT NOT NULL DEFAULT 0, -- BIGINT для хранения копеек (безопаснее для денег)
    is_public BOOLEAN NOT NULL DEFAULT false, -- По умолчанию курс скрыт/платный
    is_active BOOLEAN NOT NULL DEFAULT true,
    
    -- Контент
    contraindications TEXT,                           
    recommendations TEXT,
    target_audience TEXT,
    course_basis TEXT,
    class_basis TEXT,
    bonuses JSONB DEFAULT '[]'::jsonb,               
    
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Модули
CREATE TABLE IF NOT EXISTS modules (
    id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    week_number INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. Услуги
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price BIGINT NOT NULL DEFAULT 0, -- BIGINT
    duration_minutes INTEGER NOT NULL DEFAULT 60, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 6. Уроки
CREATE TABLE IF NOT EXISTS lessons (
    id SERIAL PRIMARY KEY,
    module_id INTEGER NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,              
    video_embed_id VARCHAR(255) NOT NULL, 
    private_key VARCHAR(255),      -- Для защиты видео
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 7. Покупки (Доступы)
CREATE TABLE IF NOT EXISTS user_purchases (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    purchased_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, course_id)
);

-- 8. Прогресс
CREATE TABLE IF NOT EXISTS user_progress (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, lesson_id)
);

-- 9. Отзывы
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    text TEXT NOT NULL CHECK (length(text) > 0), -- Защита от пустых отзывов
    rejection_reason TEXT, 
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5), -- SMALLINT достаточно для 1-5
    approved BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, course_id)
);

-- 10. Платежи
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- Уникальный ID платежа в нашей системе
    yookassa_payment_id VARCHAR(100), -- ID платежа от ЮKassa (может быть null до создания)
    
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    
    amount BIGINT NOT NULL, -- Сумма в копейках
    currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
    
    status payment_status NOT NULL DEFAULT 'pending', -- ENUM
    
    confirmation_url TEXT,
    description TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    paid_at TIMESTAMP WITH TIME ZONE
);

-- 11. Refresh Tokens (Хешируем сам токен!)
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL, -- Сюда пишем SHA256 от токена
    user_agent TEXT, -- Полезно для безопасности (отслеживать устройства)
    ip_address INET, -- IP адрес выдачи
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    revoked_at TIMESTAMP WITH TIME ZONE -- Для мягкого отзыва токена без удаления
);

-- 12. Токены сброса пароля
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE, -- Тоже хешируем!
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 13. Токены смены Email
CREATE TABLE IF NOT EXISTS email_change_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    new_email VARCHAR(255) NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- =======================================================
-- ИНДЕКСЫ (Критически важно для производительности)
-- =======================================================

-- Индексы на Foreign Keys (для быстрых JOIN и каскадного удаления)
CREATE INDEX IF NOT EXISTS idx_courses_author ON courses(author_id);
CREATE INDEX IF NOT EXISTS idx_modules_course ON modules(course_id);
CREATE INDEX IF NOT EXISTS idx_lessons_module ON lessons(module_id);
CREATE INDEX IF NOT EXISTS idx_purchases_user ON user_purchases(user_id);
CREATE INDEX IF NOT EXISTS idx_purchases_course ON user_purchases(course_id); -- Часто ищем доступы по курсу
CREATE INDEX IF NOT EXISTS idx_progress_user ON user_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_progress_lesson ON user_progress(lesson_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_course_approved ON reviews(course_id, approved) WHERE approved = true; -- Частичный индекс для быстрых выборок публичных отзывов
CREATE INDEX IF NOT EXISTS idx_payments_user ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_yookassa_id ON payments(yookassa_payment_id); -- Быстрый поиск по ID от ЮKassa
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires ON refresh_tokens(expires_at); -- Для очистки протухших токенов
CREATE INDEX IF NOT EXISTS idx_reset_tokens_user ON password_reset_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_email_change_user ON email_change_tokens(user_id);

-- Индексы для поиска
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email); -- Логин по email
CREATE INDEX IF NOT EXISTS idx_courses_active_public ON courses(is_active, is_public); -- Фильтрация каталога

-- =======================================================
-- ТРИГГЕРЫ ДЛЯ АВТО-ОБНОВЛЕНИЯ updated_at
-- =======================================================

-- Функция-триггер
CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Применяем триггер ко всем нужным таблицам
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();
CREATE TRIGGER update_courses_updated_at BEFORE UPDATE ON courses FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();
CREATE TRIGGER update_services_updated_at BEFORE UPDATE ON services FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();
CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();