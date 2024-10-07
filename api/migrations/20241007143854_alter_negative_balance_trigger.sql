-- +goose Up
-- +goose StatementBegin

DROP TRIGGER IF EXISTS negative_balance_trigger ON users;
DROP FUNCTION IF EXISTS check_balances();

CREATE OR REPLACE FUNCTION check_balances() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.point_balance < 0 THEN
        RAISE EXCEPTION 'Point balance cannot be negative';
    END IF;
    IF NEW.race_balance < 0 THEN
        RAISE EXCEPTION 'Race balance cannot be negative';
    END IF;
    IF NEW.red_key_balance < 0 THEN
        RAISE EXCEPTION 'Red key balance cannot be negative';
    END IF;
    IF NEW.blue_key_balance < 0 THEN
        RAISE EXCEPTION 'Blue key balance cannot be negative';
    END IF;
    IF NEW.green_key_balance < 0 THEN
        RAISE EXCEPTION 'Green key balance cannot be negative';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER negative_balance_trigger BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION check_balances();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS negative_balance_trigger ON users;
DROP FUNCTION IF EXISTS check_balances();

-- +goose StatementEnd