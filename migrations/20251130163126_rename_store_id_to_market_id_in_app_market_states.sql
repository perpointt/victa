-- +goose Up
-- +goose StatementBegin
-- Удаляем старый индекс, который всё ещё висит на старом имени столбца
DROP INDEX IF EXISTS idx_app_store_states_store_id;

-- Переименовываем колонку store_id -> market_id
ALTER TABLE app_market_states
    RENAME COLUMN store_id TO market_id;

-- Создаём новый индекс на market_id
CREATE INDEX idx_app_market_states_market_id
    ON app_market_states (market_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Откатываем изменения: удаляем новый индекс
DROP INDEX IF EXISTS idx_app_market_states_market_id;

-- Переименовываем колонку обратно
ALTER TABLE app_market_states
    RENAME COLUMN market_id TO store_id;

-- Восстанавливаем старый индекс
CREATE INDEX idx_app_store_states_store_id
    ON app_market_states (store_id);
-- +goose StatementEnd
