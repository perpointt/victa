-- +goose Up
-- +goose StatementBegin
-- Удаляем устаревшую таблицу с открытыми секретами
DROP TABLE IF EXISTS company_integrations;
-- +goose StatementEnd


-- +goose Down
-- намеренно оставлено пустым: миграция необратима
