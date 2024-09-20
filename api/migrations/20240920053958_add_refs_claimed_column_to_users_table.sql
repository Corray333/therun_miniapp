-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN refs_claimed INTEGER DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN refs_claimed;
-- +goose StatementEnd
