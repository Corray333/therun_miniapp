-- +goose Up
-- +goose StatementBegin
ALTER TABLE users 
ALTER COLUMN farming_from 
SET DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN farming_from SET DEFAULT 0;
-- +goose StatementEnd
