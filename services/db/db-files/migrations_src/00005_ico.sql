-- +goose Up
-- SQL in this section is executed when the migration is applied.

/* Tabels and types used to manage an ICO and it's phases.
** Currently only fixed exchange sale model is supported.
*/

CREATE TYPE ico_status AS ENUM ('planning', 'ready', 'active', 'finished', 'completed','stopped');

/* other sales model values for later 'dutch' and 'hybrid' */
CREATE TYPE ico_sales_model AS ENUM ('fixed');

CREATE TABLE ico (
  id SERIAL PRIMARY KEY NOT null,
  ico_name VARCHAR(255) NOT NULL,
  ico_status ico_status NOT NULL,
  /* kyc == true for enabled, false for disabled */
  kyc BOOLEAN  NOT NULL, 
  sales_model ico_sales_model NOT NULL,
  /* public key or issuer account - must be internal from portal */
  issuer_pk VARCHAR(56) NOT NULL, 
  /* asset code of the token to distribute - must be internal from portal */
  asset_code VARCHAR(12) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_by VARCHAR NOT NULL
);

/* other_crypto is a crypto currency that is not in the stellar blockchain */
CREATE TYPE exchange_currency_type AS ENUM('stellar', 'other_crypto', 'fiat');

CREATE TABLE exchange_currency (
  id SERIAL PRIMARY KEY NOT null,
  exchange_currency_type exchange_currency_type NOT NULL,
  /* e.g. BTC, ETH, XLM, MOBI, USD, EUR */
  asset_code VARCHAR(12) NOT NULL,
  /* only needed for currency/token from the stellar blockchain */
  issuer_pk VARCHAR(56) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_by VARCHAR NOT NULL
);

/* list of exchange currencies that can be supported by an ICO */
INSERT INTO exchange_currency (exchange_currency_type, asset_code, issuer_pk, updated_by) VALUES ('other_crypto', 'BTC', '', 'chris');
INSERT INTO exchange_currency (exchange_currency_type, asset_code, issuer_pk, updated_by) VALUES ('other_crypto', 'ETH', '', 'chris');
INSERT INTO exchange_currency (exchange_currency_type, asset_code, issuer_pk, updated_by) VALUES ('stellar', 'XLM', '', 'chris');
INSERT INTO exchange_currency (exchange_currency_type, asset_code, issuer_pk, updated_by) VALUES ('fiat', 'USD', '', 'chris');
INSERT INTO exchange_currency (exchange_currency_type, asset_code, issuer_pk, updated_by) VALUES ('fiat', 'EUR', '', 'chris');

/* currencies that are currently supported by an ICO */
CREATE TABLE ico_supported_exchange_currency (
  id SERIAL PRIMARY KEY NOT null,
  ico_id INTEGER NOT NULL REFERENCES ico(id),
  exchange_currency_id INTEGER NOT NULL REFERENCES exchange_currency(id),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_by VARCHAR NOT NULL,
  CONSTRAINT ico_currency_unique UNIQUE(ico_id, exchange_currency_id)
);

CREATE TYPE icophase_status AS ENUM ('planning', 'ready', 'active', 'finished', 'completed','stopped');

CREATE TABLE icophase (
  id SERIAL PRIMARY KEY NOT null,
  ico_id INTEGER NOT NULL REFERENCES ico(id),
  icophase_name VARCHAR(255) NOT NULL,
  icophase_status icophase_status NOT NULL,
  /* public key of the distribution account - must be internal from the portal */
  dist_pk VARCHAR(56) NOT NULL,
  /* used for signing the transaction to be sent to the customer */
  dist_presigner_pk VARCHAR(56) NOT NULL,
  dist_presigner_seed VARCHAR(56) NOT NULL,
  /* used to sign the payment transaction after signed by presigner and customer */
  dist_postsigner_pk VARCHAR(56) NOT NULL,
  dist_postsigner_seed VARCHAR(56) NOT NULL,
  start_time TIMESTAMP with time zone NOT NULL,
  end_time TIMESTAMP with time zone NOT NULL,
  tokens_to_distribute BIGINT NOT NULL,
  tokens_released BIGINT NOT NULL,
  /* tokens blocked because the order was payed by customer, but the token payment transaction not yet executed */
  tokens_blocked BIGINT NOT NULL,
  tokens_left BIGINT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_by VARCHAR NOT NULL
);

CREATE TABLE icophase_bank_account (
  id SERIAL PRIMARY KEY NOT NULL,
  account_name VARCHAR NOT NULL,
  recepient_name VARCHAR NOT NULL,
  bank_name VARCHAR NOT NULL,
  iban VARCHAR NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_by VARCHAR NOT NULL
);

/* exchange currencies that are currently activated for an ICO Phase */
/* must be supported by the ICO of the ICO Phase */
CREATE TABLE icophase_activated_exchange_currency (
  id SERIAL PRIMARY KEY NOT NULL,
  icophase_id INTEGER NOT NULL REFERENCES icophase(id),
  exchange_currency_id INTEGER NOT NULL REFERENCES exchange_currency(id),
  price_per_token BIGINT NOT NULL,
  tokens_released BIGINT NOT NULL,
  tokens_blocked BIGINT NOT NULL,
  /* only needed if the customer wants to transfer fiat to our bank account*/
  icophase_bank_account_id INTEGER REFERENCES icophase_bank_account(id), 
  created_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_by VARCHAR NOT NULL,
  CONSTRAINT icophase_currency_unique UNIQUE(icophase_id, exchange_currency_id)
);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS icophase_activated_exchange_currency;
DROP TABLE IF EXISTS icophase_bank_account;
DROP TABLE IF EXISTS ico_supported_exchange_currency;
DROP TABLE IF EXISTS exchange_currency;
DROP TYPE IF EXISTS exchange_currency_type;
DROP TABLE IF EXISTS icophase;
DROP TYPE IF EXISTS icophase_status;
DROP TABLE IF EXISTS ico;
DROP TYPE IF EXISTS ico_status;
DROP TYPE IF EXISTS ico_sales_model;