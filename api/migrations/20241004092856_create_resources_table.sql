-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS resources (
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    type VARCHAR(255) NOT NULL,
    amount INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, type)
);

-- Create triggers
CREATE OR REPLACE FUNCTION add_resources() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO resources (user_id, type, amount) VALUES (NEW.user_id, 'titan', 0);
    INSERT INTO resources (user_id, type, amount) VALUES (NEW.user_id, 'quartz', 0);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER add_resources_trigger
AFTER INSERT ON users
FOR EACH ROW EXECUTE PROCEDURE add_resources();

-- Insert resources for existing users
INSERT INTO resources (user_id, type, amount) SELECT user_id, 'titan', 0 FROM users;
INSERT INTO resources (user_id, type, amount) SELECT user_id, 'quartz', 0 FROM users;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS add_resources_trigger ON users;
DROP FUNCTION IF EXISTS add_resources();
DROP TABLE IF EXISTS resources CASCADE;
-- +goose StatementEnd