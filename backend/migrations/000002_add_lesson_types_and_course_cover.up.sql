-- Добавляем тип урока и контент для статей
ALTER TABLE lessons
    ADD COLUMN lesson_type VARCHAR(20) NOT NULL DEFAULT 'video'
        CHECK (lesson_type IN ('video', 'article')),
    ADD COLUMN article_content TEXT;

-- VideoEmbedID теперь может быть NULL (для статей)
ALTER TABLE lessons
    ALTER COLUMN video_embed_id DROP NOT NULL;

-- Добавляем обложку для курса
ALTER TABLE courses
    ADD COLUMN cover_image_url VARCHAR(500);

-- Обновляем существующие данные
UPDATE lessons SET lesson_type = 'video' WHERE lesson_type IS NULL;
