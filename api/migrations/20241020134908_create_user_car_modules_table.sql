-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_car_modules(
    user_module_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 CACHE 1),
    car_module_id BIGINT NOT NULL REFERENCES car_modules(car_module_id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    car_id BIGINT REFERENCES cars(car_id) ON DELETE CASCADE,
    round_id BIGINT REFERENCES rounds(round_id) ON DELETE CASCADE,
    CONSTRAINT user_module_pk PRIMARY KEY (user_module_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_car_modules;
-- +goose StatementEnd
