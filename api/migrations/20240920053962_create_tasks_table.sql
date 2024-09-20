-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks (
    task_id BIGINT NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1),
    icon TEXT NOT NULL DEFAULT '',
    type VARCHAR(16) NOT NULL,
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
