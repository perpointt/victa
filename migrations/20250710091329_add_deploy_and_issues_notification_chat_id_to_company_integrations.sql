-- +goose Up
-- +goose StatementBegin
BEGIN;

-- 1) Добавляем новые поля
ALTER TABLE company_integrations
    ADD COLUMN deploy_notification_chat_id TEXT,
    ADD COLUMN issues_notification_chat_id TEXT;

-- 2) Мигрируем данные из notification_chat_id в deploy_notification_chat_id
UPDATE company_integrations
SET deploy_notification_chat_id = notification_chat_id;

-- 3) Удаляем старое поле
ALTER TABLE company_integrations
    DROP COLUMN notification_chat_id;

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

-- 1) Восстанавливаем старое поле
ALTER TABLE company_integrations
    ADD COLUMN notification_chat_id TEXT;

-- 2) Копируем обратно данные из deploy_notification_chat_id
UPDATE company_integrations
SET notification_chat_id = deploy_notification_chat_id;

-- 3) Удаляем временные поля
ALTER TABLE company_integrations
    DROP COLUMN deploy_notification_chat_id,
    DROP COLUMN issues_notification_chat_id;

COMMIT;
-- +goose StatementEnd
