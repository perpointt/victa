-- +goose Up
-- +goose StatementBegin
CREATE TABLE app_integrations
(
    app_id   BIGINT NOT NULL REFERENCES apps (id) ON DELETE CASCADE,
    codemagic_app_id TEXT   NOT NULL,
    PRIMARY KEY (app_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS app_integrations;
-- +goose StatementEnd
