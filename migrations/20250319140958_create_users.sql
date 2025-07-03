-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    tg_id         TEXT UNIQUE,
    name    TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
