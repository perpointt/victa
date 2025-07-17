-- +goose Up
-- +goose StatementBegin
CREATE TABLE company_secrets
(
    company_id  BIGINT    NOT NULL REFERENCES companies (id) ON DELETE CASCADE,
    secret_type TEXT      NOT NULL,
    cipher      BYTEA     NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (company_id, secret_type)
);

CREATE INDEX idx_company_secrets_type ON company_secrets (secret_type);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_company_secrets_type;
DROP TABLE IF EXISTS company_secrets;
-- +goose StatementEnd
