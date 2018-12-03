-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE payment_template drop COLUMN if exists memo_type;
drop type if exists memo_type;
CREATE TYPE memo_type AS ENUM('MEMO_TEXT', 'MEMO_ID', 'MEMO_HASH', 'MEMO_RETURN');
ALTER TABLE payment_template ADD COLUMN memo_type memo_type not null default 'MEMO_TEXT';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE payment_template DROP COLUMN memo_type;