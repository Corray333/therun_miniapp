-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN current_miles FLOAT NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN miles_claimed FLOAT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN miles_claimed;
ALTER TABLE users DROP COLUMN current_miles;
-- +goose StatementEnd
