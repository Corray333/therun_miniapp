-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_task_done (
    user_id INT NOT NULL REFERENCES users(user_id),
    task_id INT NOT NULL REFERENCES tasks(task_id),
    PRIMARY KEY (user_id, task_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
