-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS buildings (
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    type VARCHAR(255) NOT NULL,
    level INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, type)
);
INSERT INTO buildings (user_id, type, level) SELECT user_id, 'warehouse', 0 FROM users;
INSERT INTO buildings (user_id, type, level) SELECT user_id, 'mine', 0 FROM users;
INSERT INTO buildings (user_id, type, level) SELECT user_id, 'fabric', 0 FROM users;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS buildings CASCADE;
-- +goose StatementEnd