-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_task_done (
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    task_id BIGINT NOT NULL REFERENCES tasks(task_id),
    PRIMARY KEY (user_id, task_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_task_done;
-- +goose StatementEnd
