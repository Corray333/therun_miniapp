-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN ref_claimed BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN ref_claimed;
-- +goose StatementEnd
