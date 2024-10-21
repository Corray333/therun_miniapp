-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS car_modules(
    car_module_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 CACHE 1),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    car_id BIGINT NOT NULL REFERENCES cars(car_id),
    characteristic VARCHAR(16) NOT NULL,
    boost INTEGER NOT NULL,
    name TEXT NOT NULL

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS car_modules;
-- +goose StatementEnd
