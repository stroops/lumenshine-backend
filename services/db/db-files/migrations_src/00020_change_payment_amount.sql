-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE payment_template
    ALTER COLUMN amount TYPE character varying(64),
    ALTER COLUMN amount SET DEFAULT '',
    ALTER COLUMN amount SET NOT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.