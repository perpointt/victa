-- +goose Up
-- +goose StatementBegin
/* 1. Добавляем колонку */
ALTER TABLE stores
    ADD COLUMN slug TEXT;

/* 2. Заполняем для уже существующих строк
      (можно заменить выражение на своё правило генерации) */
UPDATE stores
SET slug = lower(regexp_replace(name, '\s+', '_', 'g'))
WHERE slug IS NULL;

/* 3. Делаем NOT NULL + уникальность */
ALTER TABLE stores
    ALTER COLUMN slug SET NOT NULL,
    ADD CONSTRAINT stores_slug_unique UNIQUE (slug);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE stores
    DROP CONSTRAINT IF EXISTS stores_slug_unique,
    DROP COLUMN IF EXISTS slug;
-- +goose StatementEnd
