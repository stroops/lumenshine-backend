-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE user_profile ADD COLUMN tfa_temp_secret character varying COLLATE pg_catalog."default" NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE user_profile DROP COLUMN tfa_temp_secret;