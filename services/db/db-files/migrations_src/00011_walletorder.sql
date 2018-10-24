-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE user_wallet ADD COLUMN order_nr integer NOT NULL DEFAULT 0;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE user_wallet DROP COLUMN order_nr;