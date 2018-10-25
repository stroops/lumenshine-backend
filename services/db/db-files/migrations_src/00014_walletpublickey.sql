-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE user_wallet RENAME public_key_0 to public_key;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE user_wallet RENAME public_key to public_key_0;