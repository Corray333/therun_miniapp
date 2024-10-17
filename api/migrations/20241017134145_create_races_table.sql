-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS races (
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    round_id BIGINT NOT NULL REFERENCES rounds(round_id),
    start_time BIGINT NOT NULL DEFAULT 0,
    miles FLOAT NOT NULL DEFAULT 0,
    CONSTRAINT races_pk PRIMARY KEY (user_id, round_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS races;
-- +goose StatementEnd
