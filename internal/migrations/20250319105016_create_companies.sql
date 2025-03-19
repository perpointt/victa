-- +goose Up
-- +goose StatementBegin
CREATE TABLE companies
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT                     NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS companies;
-- +goose StatementEnd
