-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN daily_check_streak INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN daily_check_last BIGINT DEFAULT EXTRACT(EPOCH FROM NOW());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN daily_check_streak;
ALTER TABLE users DROP COLUMN daily_check_last;
-- +goose StatementEnd
