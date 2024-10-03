-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_case(
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    case_type VARCHAR(16) NOT NULL
);
INSERT INTO user_case(user_id, case_type) SELECT user_id, 'red' FROM users;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_case;
-- +goose StatementEnd
