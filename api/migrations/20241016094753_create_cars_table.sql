-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cars (
    car_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 CACHE 1),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    element VARCHAR(255) NOT NULL,
    acceleration INTEGER NOT NULL,
    handling INTEGER NOT NULL,
    brakes INTEGER NOT NULL,
    strength INTEGER NOT NULL,
    tank INTEGER NOT NULL,
    fuel FLOAT NOT NULL DEFAULT 0 ,
    health FLOAT NOT NULL DEFAULT 0 ,
    is_main BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT cars_pk PRIMARY KEY (car_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cars;
-- +goose StatementEnd
