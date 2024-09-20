-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS task_translate (
    task_id INT NOT NULL REFERENCES tasks(task_id) ON DELETE CASCADE,
    lang VARCHAR(2) NOT NULL,
    description TEXT NOT NULL,
    PRIMARY KEY (task_id, lang)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS task_translate;
-- +goose StatementEnd
