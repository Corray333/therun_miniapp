-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
    user_id BIGINT NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    avatar TEXT NOT NULL DEFAULT '',
    in_app_id BIGINT NOT NULL,
    point_balance BIGINT NOT NULL DEFAULT 0,
    race_balance BIGINT NOT NULL DEFAULT 0,
    key_balance BIGINT NOT NULL DEFAULT 0,
    last_claim BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    max_points BIGINT NOT NULL DEFAULT 0,
    farm_time BIGINT NOT NULL DEFAULT 7200,
    ref_code VARCHAR(8) NOT NULL UNIQUE,
    referer BIGINT REFERENCES users(user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd