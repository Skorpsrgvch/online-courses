-- Откат изменений
ALTER TABLE courses DROP COLUMN IF EXISTS cover_image_url;
ALTER TABLE lessons DROP COLUMN IF EXISTS article_content;
ALTER TABLE lessons DROP COLUMN IF EXISTS lesson_type;

-- Возвращаем NOT NULL для video_embed_id
UPDATE lessons SET video_embed_id = '' WHERE video_embed_id IS NULL;
ALTER TABLE lessons ALTER COLUMN video_embed_id SET NOT NULL;
