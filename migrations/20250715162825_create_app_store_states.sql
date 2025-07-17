-- +goose Up
-- +goose StatementBegin
CREATE TABLE app_store_states
(
    app_id         BIGINT NOT NULL REFERENCES apps (id) ON DELETE CASCADE,
    store_id       BIGINT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
    last_version   BIGINT NULL,
    last_review_id BIGINT NULL,
    updated_at     TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (app_id, store_id)
);

CREATE INDEX idx_app_store_states_store_id ON app_store_states (store_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS app_store_states;
-- +goose StatementEnd
