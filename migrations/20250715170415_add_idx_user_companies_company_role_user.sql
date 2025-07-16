-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_user_companies_company_role_user
    ON user_companies (company_id, role_id, user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_companies_company_role_user;
-- +goose StatementEnd