-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE user_profile add COLUMN last_lockout_time timestamp with time zone NOT NULL default '1970-01-01';
ALTER TABLE user_profile add COLUMN last_lockout_counter int NOT NULL default 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE user_profile DROP COLUMN last_lockout_time;
ALTER TABLE user_profile DROP COLUMN last_lockout_counter;