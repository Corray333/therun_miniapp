-- +goose Up
-- +goose StatementBegin

CREATE OR REPLACE FUNCTION check_balances() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.point_balance < 0 THEN
        RAISE EXCEPTION 'Point balance cannot be negative';
    END IF;
    IF NEW.race_balance < 0 THEN
        RAISE EXCEPTION 'Race balance cannot be negative';
    END IF;
    IF NEW.key_balance < 0 THEN
        RAISE EXCEPTION 'Key balance cannot be negative';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER point_balance_trigger BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION check_balances();
CREATE TRIGGER race_balance_trigger BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION check_balances();
CREATE TRIGGER key_balance_trigger BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION check_balances();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS point_balance_trigger ON users;
DROP TRIGGER IF EXISTS race_balance_trigger ON users;
DROP TRIGGER IF EXISTS key_balance_trigger ON users;
-- +goose StatementEnd
