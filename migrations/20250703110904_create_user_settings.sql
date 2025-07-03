-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_settings
(
    user_id    BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    active_company_id  BIGINT NULL
        REFERENCES companies (id) ON DELETE SET NULL,
    PRIMARY KEY (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_settings;
-- +goose StatementEnd
