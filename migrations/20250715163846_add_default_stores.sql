-- +goose Up
-- +goose StatementBegin
INSERT INTO stores (slug, name)
VALUES
    ('play_store', 'Play Store'),
    ('app_store',  'AppStore'),
    ('ru_store',   'RuStore'),
    ('app_gallery','Huawei App Gallery')
ON CONFLICT (slug) DO NOTHING;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM stores
WHERE slug IN ('play_store','app_store','ru_store','app_gallery');
-- +goose StatementEnd
