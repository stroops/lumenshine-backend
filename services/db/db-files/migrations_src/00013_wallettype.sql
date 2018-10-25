-- +goose Up
-- SQL in this section is executed when the migration is applied.

/*wallet_type*/
CREATE TYPE wallet_type AS ENUM ('internal','external');

ALTER TABLE user_wallet ADD COLUMN wallet_type wallet_type NOT NULL DEFAULT 'internal';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE user_wallet DROP COLUMN IF EXISTS wallet_type;

drop type IF EXISTS wallet_type;