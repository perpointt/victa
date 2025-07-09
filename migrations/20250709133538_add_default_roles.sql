-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (slug) VALUES
                             ('admin'),
                             ('developer')
ON CONFLICT (slug) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM roles
WHERE slug IN ('admin', 'developer');
-- +goose StatementEnd