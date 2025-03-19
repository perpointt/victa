-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES companies (id) ON DELETE CASCADE,
    email      TEXT NOT NULL UNIQUE,
    password   TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
