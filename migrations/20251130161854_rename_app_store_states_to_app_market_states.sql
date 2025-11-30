-- +goose Up
-- +goose StatementBegin
ALTER TABLE app_store_states
    RENAME TO app_market_states;

ALTER INDEX idx_app_store_states_store_id
    RENAME TO idx_app_market_states_store_id;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE app_market_states
    RENAME TO app_store_states;

ALTER INDEX idx_app_market_states_store_id
    RENAME TO idx_app_store_states_store_id;
-- +goose StatementEnd
