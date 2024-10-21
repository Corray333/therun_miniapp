-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rounds (
    round_id INTEGER NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 CACHE 1),
    end_time BIGINT NOT NULL,
    element VARCHAR(16) NOT NULL,
    CONSTRAINT rounds_pk PRIMARY KEY (round_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rounds;
-- +goose StatementEnd
