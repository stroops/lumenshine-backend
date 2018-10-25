-- +goose Up
-- SQL in this section is executed when the migration is applied.

UPDATE user_wallet
SET order_nr = indexer.rownum - 1
FROM (SELECT id, ROW_NUMBER() OVER(PARTITION BY user_id ORDER BY id) as rownum FROM user_wallet) indexer
WHERE user_wallet.id = indexer.id;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

UPDATE user_wallet
SET order_nr = 0;