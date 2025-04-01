-- +goose Up
-- +goose StatementBegin
CREATE TABLE application_integrations
(
    application_id   BIGINT NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    codemagic_app_id TEXT   NOT NULL,
    PRIMARY KEY (application_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS application_integrations;
-- +goose StatementEnd
