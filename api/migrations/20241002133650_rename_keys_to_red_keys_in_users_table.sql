-- +goose Up
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN key_balance TO red_key_balance;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN red_key_balance TO key_balance;
-- +goose StatementEnd
