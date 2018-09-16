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

CREATE TYPE payment_network AS ENUM('fiat', 'stellar', 'ethereum', 'bitcoin');

CREATE TYPE exchange_currency_type AS ENUM('crypto', 'fiat');

CREATE TABLE exchange_currency (
  id SERIAL PRIMARY KEY NOT null,
  exchange_currency_type exchange_currency_type NOT NULL,

  /* e.g. BTC, ETH, XLM, MOBI, USD, EUR */
  asset_code VARCHAR(12) NOT NULL,

  /* e.g. Wei, Satoshi, Stroop, Cent */
  denom_asset_code VARCHAR(64) NOT NULL,

  payment_network payment_network not null,

  /* only needed for currency/token from the stellar blockchain */
  issuer_pk VARCHAR(56) NOT NULL check(
    (payment_network<>'stellar')
    or
    (payment_network='stellar' and issuer_pk<>'')
  ),
  /* number of max decimals for the currency */
  decimals int not null,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_by VARCHAR NOT NULL
);

/* list of exchange currencies that can be supported by an ICO */
INSERT INTO exchange_currency (id, exchange_currency_type, asset_code, denom_asset_code, payment_network, issuer_pk, decimals, updated_by) VALUES (1, 'crypto', 'BTC', 'Satoshi', 'bitcoin', '', 8, 'chris');
INSERT INTO exchange_currency (id, exchange_currency_type, asset_code, denom_asset_code, payment_network, issuer_pk, decimals, updated_by) VALUES (2, 'crypto', 'ETH', 'Wei', 'ethereum', '', 8, 'chris');
INSERT INTO exchange_currency (id, exchange_currency_type, asset_code, denom_asset_code, payment_network, issuer_pk, decimals, updated_by) VALUES (3, 'crypto', 'XLM', 'Stroop', 'stellar', 'Gxxxxx', 7, 'chris');
INSERT INTO exchange_currency (id, exchange_currency_type, asset_code, denom_asset_code, payment_network, issuer_pk, decimals, updated_by) VALUES (4, 'fiat', 'USD', 'Cent', 'fiat', '', 2, 'chris');
INSERT INTO exchange_currency (id, exchange_currency_type, asset_code, denom_asset_code, payment_network, issuer_pk, decimals, updated_by) VALUES (5, 'fiat', 'EUR', 'Cent', 'fiat', '', 2, 'chris');

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

CREATE TYPE ico_phase_status AS ENUM ('planning', 'ready', 'active', 'finished', 'completed','stopped');

CREATE TABLE ico_phase (
  id SERIAL PRIMARY KEY NOT null,
  ico_id INTEGER NOT NULL REFERENCES ico(id),
  ico_phase_name VARCHAR(255) NOT NULL,
  ico_phase_status ico_phase_status NOT NULL,
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
  token_max_order_amount bigint not null default 0,
  token_min_order_amount bigint not null default 0,
  max_user_orders int not null,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_by VARCHAR NOT NULL
);

create unique index on ico_phase (ico_id, ico_phase_status) where ico_phase_status = 'active'; /* only one active per ico at a time */

CREATE TABLE ico_phase_bank_account (
  id SERIAL PRIMARY KEY NOT NULL,
  account_name VARCHAR NOT NULL,
  recepient_name VARCHAR NOT NULL,
  bank_name VARCHAR NOT NULL,
  iban VARCHAR NOT NULL,
  bic_swift VARCHAR NOT NULL,
  /* this must be in go string-format e.g. 'Payment ID: %s'. %s will be replaced with the correct payment id, needed to idetify the order */
  paymend_usage_string VARCHAR NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_by VARCHAR NOT NULL
);

/* exchange currencies that are currently activated for an ICO Phase */
/* must be supported by the ICO of the ICO Phase */
CREATE TABLE ico_phase_activated_exchange_currency (
  id SERIAL PRIMARY KEY NOT NULL,
  ico_phase_id INTEGER NOT NULL REFERENCES ico_phase(id),
  exchange_currency_id INTEGER NOT NULL REFERENCES exchange_currency(id),  
  exchange_master_key text not null, /* master key for generating the addresses and seeds in the payment network */
  denom_price_per_token BIGINT NOT NULL,
  stellar_starting_balance_denom varchar(64) NOT NULL, /*starting-balannce (in denomination) for creating the stellar-account */
  tokens_released BIGINT NOT NULL,
  tokens_blocked BIGINT NOT NULL,
  /* only needed if the customer wants to transfer fiat to our bank account*/
  ico_phase_bank_account_id INTEGER REFERENCES ico_phase_bank_account(id), 
  created_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL default current_timestamp,
  updated_by VARCHAR NOT NULL,
  CONSTRAINT ico_phase_currency_unique UNIQUE(ico_phase_id, exchange_currency_id)
);

/* create some demo data */
insert into ico(id, ico_name, ico_status, kyc, sales_model, issuer_pk, asset_code, updated_by) values 
  (1, 'Demo-ICO', 'active', true, 'fixed', 'GCCBLT6VFEUODLP36C675TJDNZNHQFD5P6L3BBCYUMU2TIO3UQCVXXX3', 'CaliCoin', 'setup');

insert into ico_phase(id, ico_id, ico_phase_name, ico_phase_status, dist_pk, dist_presigner_pk, dist_presigner_seed, dist_postsigner_pk, dist_postsigner_seed, start_time, end_time, tokens_to_distribute, tokens_released,tokens_blocked,tokens_left,token_max_order_amount,token_min_order_amount,max_user_orders,updated_by) values
  (1,1, 'Phase 1', 'active', 'GA2XZU7BHRI6VWKAR5GYWMGAPQKFFFMYGVZNBS7QD3QSR5K2HKP7Y672', 'GCSC4YWWKTAFSKMYP25THDQNLYLWHZFCLTEXVBLBDHSEUFWDEKWWPUL7', 'SAOBW63XEFF7H5OLX53VZVA56MW2XNY33K4QPJEST4THZMSRV7TRE35U', 'GA2XZU7BHRI6VWKAR5GYWMGAPQKFFFMYGVZNBS7QD3QSR5K2HKP7Y672', 'SDREDOAGCJUHJQYNFVV6IRO5LTMHA3JACCDKRAIXTABWAXRHAUKP6SU6', '2018-1-1'::timestamp, '2020-1-1'::timestamp, 1000000, 0, 0, 1000000, 1000, 1, 10, 'setup'),
  (2,1, 'Phase 2', 'planning', 'GA2XZU7BHRI6VWKAR5GYWMGAPQKFFFMYGVZNBS7QD3QSR5K2HKP7Y672', 'GCSC4YWWKTAFSKMYP25THDQNLYLWHZFCLTEXVBLBDHSEUFWDEKWWPUL7', 'SAOBW63XEFF7H5OLX53VZVA56MW2XNY33K4QPJEST4THZMSRV7TRE35U', 'GA2XZU7BHRI6VWKAR5GYWMGAPQKFFFMYGVZNBS7QD3QSR5K2HKP7Y672', 'SDREDOAGCJUHJQYNFVV6IRO5LTMHA3JACCDKRAIXTABWAXRHAUKP6SU6', '2018-1-1'::timestamp, '2020-1-1'::timestamp, 1000000, 0, 0, 1000000, 2000, 10, 5, 'setup'),
  (3,1, 'Phase 3', 'ready', 'GA2XZU7BHRI6VWKAR5GYWMGAPQKFFFMYGVZNBS7QD3QSR5K2HKP7Y672', 'GCSC4YWWKTAFSKMYP25THDQNLYLWHZFCLTEXVBLBDHSEUFWDEKWWPUL7', 'SAOBW63XEFF7H5OLX53VZVA56MW2XNY33K4QPJEST4THZMSRV7TRE35U', 'GA2XZU7BHRI6VWKAR5GYWMGAPQKFFFMYGVZNBS7QD3QSR5K2HKP7Y672', 'SDREDOAGCJUHJQYNFVV6IRO5LTMHA3JACCDKRAIXTABWAXRHAUKP6SU6', '2018-1-1'::timestamp, '2020-1-1'::timestamp, 1000000, 0, 0, 1000000, 1000, 5, 1, 'setup');

insert into ico_supported_exchange_currency(id, ico_id, exchange_currency_id, updated_by) values
  (1, 1, 1, 'setup'),
  (2, 1, 2, 'setup'),
  (3, 1, 3, 'setup'),
  (4, 1, 4, 'setup'),
  (5, 1, 5, 'setup');

insert into ico_phase_bank_account(id, account_name, recepient_name, bank_name, iban, bic_swift, paymend_usage_string, updated_by) values
  (1, 'Bank-Acc1', 'Udo Polder', 'MyBank', 'DE12344', 'LZ1233', 'Payment for %s', 'system'),
  (2, 'Bank-Acc2', 'Chris Rogobete', 'HisBank', 'DE12366', 'LZ1266', 'HisPayment for %s', 'system');

insert into ico_phase_activated_exchange_currency (id, ico_phase_id, exchange_currency_id, denom_price_per_token, tokens_released, tokens_blocked, ico_phase_bank_account_id, exchange_master_key, stellar_starting_balance_denom, updated_by) values
  (1, 1, 1, 1000000, 0, 0, null, 'xpub6DxSCdWu6jKqr4isjo7bsPeDD6s3J4YVQV1JSHZg12Eagdqnf7XX4fxqyW2sLhUoFWutL7tAELU2LiGZrEXtjVbvYptvTX5Eoa4Mamdjm9u', '30000000', 'setup'), /*  btc 0,10000000*/
  (2, 1, 2, 2000000, 0, 0, null, 'xpub6DxSCdWu6jKqr4isjo7bsPeDD6s3J4YVQV1JSHZg12Eagdqnf7XX4fxqyW2sLhUoFWutL7tAELU2LiGZrEXtjVbvYptvTX5Eoa4Mamdjm9u', '40000000', 'setup'), /*  eth 0,20000000*/
  (3, 1, 3, 3000000, 0, 0, null, '', '15000000', 'setup'), /*  xlm 0,30000000*/
 /* EUR 1,50*/ (4, 1, 5, 150, 0, 0, 1, '', '20000000', 'setup'); 

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS ico_phase_activated_exchange_currency;
DROP TABLE IF EXISTS ico_phase_bank_account;
DROP TABLE IF EXISTS ico_supported_exchange_currency;
DROP TABLE IF EXISTS exchange_currency;
DROP TYPE IF EXISTS exchange_currency_type;
DROP TABLE IF EXISTS ico_phase;
DROP TYPE IF EXISTS ico_phase_status;
DROP TABLE IF EXISTS ico;
DROP TYPE IF EXISTS ico_status;
DROP TYPE IF EXISTS ico_sales_model;
DROP TYPE IF EXISTS payment_network;
