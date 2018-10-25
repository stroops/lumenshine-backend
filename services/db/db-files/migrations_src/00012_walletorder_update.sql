-- +goose Up
-- SQL in this section is executed when the migration is applied.

WITH Numbered AS
(
    SELECT id, user_id, order_nr, ROW_NUMBER() OVER(PARTITION BY user_id ORDER BY id) as rownum
    FROM user_wallet
)

UPDATE user_wallet
SET order_nr = rownum - 1
FROM Numbered WHERE user_wallet.id = Numbered.id;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

UPDATE user_wallet
SET order_nr = 0;