-- +goose Up
-- +goose StatementBegin
CREATE TABLE company_integrations
(
    id                 BIGSERIAL PRIMARY KEY,
    company_id         BIGINT                   NOT NULL REFERENCES companies (id) ON DELETE CASCADE,
    codemagic_token    TEXT,
    telegram_bot_token TEXT,
    telegram_chat_id   TEXT,
    created_at         TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS company_integrations;
-- +goose StatementEnd
