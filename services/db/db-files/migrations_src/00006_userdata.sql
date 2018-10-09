-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE user_profile
ADD COLUMN additional_name varchar(255) NOT NULL DEFAULT '',
ADD COLUMN birth_country_code varchar(3) NOT NULL DEFAULT '',
ADD COLUMN bank_account_number varchar(255) NOT NULL DEFAULT '',
ADD COLUMN bank_number varchar(255) NOT NULL DEFAULT '',
ADD COLUMN bank_phone_number  varchar(255) NOT NULL DEFAULT '',
ADD COLUMN tax_id varchar(255) NOT NULL DEFAULT '',
ADD COLUMN tax_id_name varchar(255) NOT NULL DEFAULT '',
ADD COLUMN occupation varchar(5) NOT NULL DEFAULT '',
ADD COLUMN employer_name varchar(500) NOT NULL DEFAULT '',
ADD COLUMN employer_address varchar(500) NOT NULL DEFAULT '',
ADD COLUMN language_code varchar(10) NOT NULL DEFAULT '';


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE user_profile
DROP COLUMN additional_name,
DROP COLUMN birth_country_code,
DROP COLUMN bank_account_number,
DROP COLUMN bank_number,
DROP COLUMN bank_phone_number,
DROP COLUMN tax_id,
DROP COLUMN tax_id_name,
DROP COLUMN occupation,
DROP COLUMN employer_name,
DROP COLUMN employer_address,
DROP COLUMN language_code;