-- +goose Up
-- +goose StatementBegin
CREATE TABLE apps
(
    id         BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies (id) ON DELETE CASCADE,
    name       TEXT   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS apps;
-- +goose StatementEnd
