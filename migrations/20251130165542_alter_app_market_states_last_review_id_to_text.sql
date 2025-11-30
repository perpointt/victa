-- +goose Up
-- +goose StatementBegin
ALTER TABLE app_market_states
    ALTER COLUMN last_review_id TYPE TEXT
        USING last_review_id::TEXT;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE app_market_states
    ALTER COLUMN last_review_id TYPE BIGINT
        USING last_review_id::BIGINT;
-- +goose StatementEnd
