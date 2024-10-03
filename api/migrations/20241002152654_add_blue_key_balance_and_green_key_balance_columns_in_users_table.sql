-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN blue_key_balance INT NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN green_key_balance INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN blue_key_balance;
ALTER TABLE users DROP COLUMN green_key_balance;
-- +goose StatementEnd
