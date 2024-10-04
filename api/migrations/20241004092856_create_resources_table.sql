-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS resources (
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    type VARCHAR(255) NOT NULL,
    amount INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, type)
);
INSERT INTO resources (user_id, type, amount) SELECT user_id, 'titan', 0 FROM users;
INSERT INTO resources (user_id, type, amount) SELECT user_id, 'quartz', 0 FROM users;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS resources CASCADE;
-- +goose StatementEnd