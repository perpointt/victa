-- +goose Up
-- +goose StatementBegin
ALTER TABLE company_integrations
    ADD COLUMN errors_notification_chat_id TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE company_integrations
    DROP COLUMN IF EXISTS errors_notification_chat_id;
-- +goose StatementEnd
