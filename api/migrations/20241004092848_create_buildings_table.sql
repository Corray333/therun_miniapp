-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS buildings (
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    type VARCHAR(255) NOT NULL,
    level INTEGER NOT NULL DEFAULT 0,
    state VARCHAR(8) NOT NULL DEFAULT 'idle',
    last_state_change BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, type)
);

-- Create triggers
CREATE OR REPLACE FUNCTION add_buildings() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO buildings (user_id, type, level) VALUES (NEW.user_id, 'warehouse', 0);
    INSERT INTO buildings (user_id, type, level) VALUES (NEW.user_id, 'mine', 0);
    INSERT INTO buildings (user_id, type, level) VALUES (NEW.user_id, 'fabric', 0);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER add_buildings_trigger
AFTER INSERT ON users
FOR EACH ROW EXECUTE PROCEDURE add_buildings();
INSERT INTO buildings (user_id, type, level) SELECT user_id, 'warehouse', 0 FROM users;
INSERT INTO buildings (user_id, type, level) SELECT user_id, 'mine', 0 FROM users;
INSERT INTO buildings (user_id, type, level) SELECT user_id, 'fabric', 0 FROM users;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS add_buildings_trigger ON users;
DROP FUNCTION IF EXISTS add_buildings();
DROP TABLE IF EXISTS buildings CASCADE;
-- +goose StatementEnd