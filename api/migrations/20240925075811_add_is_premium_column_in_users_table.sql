-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN is_premium BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN is_premium;
-- +goose StatementEnd