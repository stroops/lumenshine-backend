-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE exchange_currency
ADD COLUMN name varchar NOT NULL DEFAULT '';

UPDATE exchange_currency SET name = 'Bitcoin' WHERE id = 1;
UPDATE exchange_currency SET name = 'Ether' WHERE id = 2;
UPDATE exchange_currency SET name = 'Lumen' WHERE id = 3;
UPDATE exchange_currency SET name = 'US Dollar' WHERE id = 4;
UPDATE exchange_currency SET name = 'Euro' WHERE id = 5;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE exchange_currency
DROP COLUMN name;