-- +goose Up
-- +goose StatementBegin
ALTER TABLE races ADD COLUMN temp_miles FLOAT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE races DROP COLUMN temp_miles; 
-- +goose StatementEnd
