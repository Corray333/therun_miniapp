-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks (
    task_id BIGINT NOT NULL PRIMARY KEY,
    description TEXT NOT NULL,
    type VARCHAR(255) NOT NULL,
    link TEXT NOT NULL DEFAULT '',
    data JSONB NOT NULL DEFAULT '{}',
    expire_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    points_reward BIGINT NOT NULL DEFAULT 0,
    keys_reward BIGINT NOT NULL DEFAULT 0,
    race_reward BIGINT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd
