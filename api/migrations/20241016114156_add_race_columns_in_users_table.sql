-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN current_miles FLOAT NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN race_start_time FLOAT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN race_start_time;
ALTER TABLE users DROP COLUMN current_miles;
-- +goose StatementEnd
