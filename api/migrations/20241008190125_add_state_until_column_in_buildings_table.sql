-- +goose Up
-- +goose StatementBegin
ALTER TABLE buildings ADD COLUMN state_until BIGINT DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE buildings DROP COLUMN state_until;
-- +goose StatementEnd
