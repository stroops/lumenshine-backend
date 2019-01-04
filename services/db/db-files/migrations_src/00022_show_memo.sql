-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE user_profile add COLUMN show_memos boolean not null default true;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE user_profile DROP COLUMN show_memos;