-- +goose Up
-- +goose StatementBegin
ALTER TABLE apps
    ADD COLUMN app_store_url   TEXT,
    ADD COLUMN play_store_url  TEXT,
    ADD COLUMN ru_store_url    TEXT,
    ADD COLUMN app_gallery_url TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE apps
    DROP COLUMN IF EXISTS app_store_url,
    DROP COLUMN IF EXISTS play_store_url,
    DROP COLUMN IF EXISTS ru_store_url,
    DROP COLUMN IF EXISTS app_gallery_url;
-- +goose StatementEnd
