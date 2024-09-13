-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN farming_from BIGINT NOT NULL DEFAULT 0;
UPDATE users SET farming_from = last_claim + 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN farming_from;
-- +goose StatementEnd
