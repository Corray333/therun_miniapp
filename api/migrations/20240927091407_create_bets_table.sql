-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bets(
    battle_id BIGINT NOT NULL REFERENCES battles(battle_id),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    pick INT NOT NULL DEFAULT 0,
    PRIMARY KEY (battle_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bets;
-- +goose StatementEnd
