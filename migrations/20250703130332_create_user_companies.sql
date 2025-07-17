-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_companies
(
    user_id    BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    company_id BIGINT NOT NULL REFERENCES companies (id) ON DELETE CASCADE,
    role_id    BIGINT NOT NULL REFERENCES roles (id),
    PRIMARY KEY (user_id, company_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_companies;
-- +goose StatementEnd
