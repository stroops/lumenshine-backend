-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE user_order drop COLUMN if exists address_index;
ALTER TABLE user_order ADD COLUMN btc_src_out_index int NOT NULL DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE user_order ADD COLUMN address_index varchar NOT NULL DEFAULT 0;
ALTER TABLE user_order drop COLUMN if exists btc_src_out_index;