-- +goose Up
-- +goose StatementBegin
CREATE TABLE application_stores
(
    application_id BIGINT NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    store_id       BIGINT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
    url            TEXT   NOT NULL,
    bundle         TEXT   NOT NULL,
    PRIMARY KEY (application_id, store_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS application_stores;
-- +goose StatementEnd
