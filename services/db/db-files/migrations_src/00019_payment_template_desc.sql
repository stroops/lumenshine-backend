-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE payment_template ADD COLUMN template_name character varying(128) NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE payment_template DROP COLUMN template_name;