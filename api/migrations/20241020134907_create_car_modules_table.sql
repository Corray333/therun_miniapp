-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS car_modules(
    car_module_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 CACHE 1),
    characteristic VARCHAR(16) NOT NULL,
    boost INTEGER NOT NULL,
    name TEXT NOT NULL,
    is_temp BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT car_modules_pk PRIMARY KEY (car_module_id)

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS car_modules;
-- +goose StatementEnd
