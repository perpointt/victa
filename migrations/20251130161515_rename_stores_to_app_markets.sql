-- +goose Up
-- +goose StatementBegin
ALTER TABLE stores
    RENAME TO app_markets;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE app_markets
    RENAME TO stores;
-- +goose StatementEnd
