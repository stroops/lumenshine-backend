-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE user_profile
ADD COLUMN occupation_name varchar(256) NOT NULL DEFAULT '',
ADD COLUMN occupation_code88 varchar(8) NOT NULL DEFAULT '';

ALTER TABLE user_profile
RENAME occupation TO occupation_code08;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE user_profile
RENAME occupation_code08 TO occupation;

ALTER TABLE user_profile
DROP COLUMN occupation_name,
DROP COLUMN occupation_code88;