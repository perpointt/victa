-- +goose Up
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS company_id;

CREATE TABLE user_companies
(
    user_id    BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    company_id BIGINT    NOT NULL REFERENCES companies (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, company_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_companies;
-- +goose StatementEnd
