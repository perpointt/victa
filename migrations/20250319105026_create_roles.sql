-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles
(
    id   BIGSERIAL PRIMARY KEY,
    slug TEXT UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
