-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE payment_template
ALTER COLUMN recepient_pk TYPE varchar(56),
ALTER COLUMN asset_code TYPE varchar(12),
ALTER COLUMN issuer_pk TYPE varchar(56);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE payment_template
ALTER COLUMN recepient_pk TYPE character(56),
ALTER COLUMN asset_code TYPE character(12),
ALTER COLUMN issuer_pk TYPE character(56);