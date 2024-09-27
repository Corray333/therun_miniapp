-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS battles (
    battle_id BIGINT NOT NULL,
    round_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    opponent_id BIGINT NOT NULL,
    user_username TEXT NOT NULL DEFAULT '',
    opponent_username TEXT NOT NULL DEFAULT '',
    user_photo TEXT NOT NULL DEFAULT '',
    opponent_photo TEXT NOT NULL DEFAULT '',
    user_miles FLOAT NOT NULL DEFAULT 0,
    opponent_miles FLOAT NOT NULL DEFAULT 0,
    PRIMARY KEY (battle_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS battles;
-- +goose StatementEnd