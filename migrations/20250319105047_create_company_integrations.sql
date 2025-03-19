-- +goose Up
-- +goose StatementBegin
CREATE TABLE company_integrations
(
    company_id BIGINT NOT NULL REFERENCES companies (id) ON DELETE CASCADE,
    codemagic_token    TEXT,
    telegram_bot_token TEXT,
    telegram_chat_id TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS company_integrations;
-- +goose StatementEnd
