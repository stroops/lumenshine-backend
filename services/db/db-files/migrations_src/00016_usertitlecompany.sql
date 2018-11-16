-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE user_profile
DROP COLUMN title,
DROP COLUMN company;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE user_profile
ADD COLUMN title character varying(64) NOT NULL DEFAULT '',
ADD COLUMN company character varying(128) NOT NULL DEFAULT '';