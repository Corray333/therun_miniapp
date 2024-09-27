-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks ADD COLUMN class VARCHAR(8) NOT NULL DEFAULT 'task';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks DROP COLUMN class;
-- +goose StatementEnd
