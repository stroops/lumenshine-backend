-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE user_profile
ADD COLUMN additional_name varchar(256) NOT NULL DEFAULT '',
ADD COLUMN birth_country_code varchar(2) NOT NULL DEFAULT '',
ADD COLUMN bank_account_number varchar(256) NOT NULL DEFAULT '',
ADD COLUMN bank_number varchar(256) NOT NULL DEFAULT '',
ADD COLUMN bank_phone_number  varchar(256) NOT NULL DEFAULT '',
ADD COLUMN tax_id varchar(256) NOT NULL DEFAULT '',
ADD COLUMN tax_id_name varchar(256) NOT NULL DEFAULT '',
ADD COLUMN occupation varchar(8) NOT NULL DEFAULT '',
ADD COLUMN employer_name varchar(512) NOT NULL DEFAULT '',
ADD COLUMN employer_address varchar(512) NOT NULL DEFAULT '',
ADD COLUMN language_code varchar(16) NOT NULL DEFAULT '';

ALTER TABLE user_profile DROP COLUMN street_number;
ALTER TABLE user_profile ALTER COLUMN street_address TYPE varchar(512);
ALTER TABLE user_profile RENAME street_address TO address;


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

ALTER TABLE user_profile ADD COLUMN street_number character varying(128) NOT NULL DEFAULT '';
ALTER TABLE user_profile RENAME address TO street_address;
ALTER TABLE user_profile ALTER COLUMN street_address TYPE VARCHAR(128);