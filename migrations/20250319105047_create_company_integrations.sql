-- +goose Up
-- +goose StatementBegin
CREATE TABLE company_integrations
(
    company_id BIGINT NOT NULL REFERENCES companies (id) ON DELETE CASCADE,
    codemagic_api_key    TEXT,
    notification_bot_token TEXT,
    notification_chat_id TEXT,
    PRIMARY KEY (company_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS company_integrations;
-- +goose StatementEnd
