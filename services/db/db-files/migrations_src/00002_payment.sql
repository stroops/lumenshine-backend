-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE key_value_store (
  key varchar(255) NOT NULL PRIMARY KEY not null,
  str_value varchar(255) NOT NULL,
  int_value bigint NOT NULL
);

INSERT INTO key_value_store (key, str_value, int_value) VALUES ('eth_address_index', '', 0);
INSERT INTO key_value_store (key, str_value, int_value) VALUES ('eth_last_block', '0', 0);

INSERT INTO key_value_store (key, str_value, int_value) VALUES ('btc_address_index', '', 0);
INSERT INTO key_value_store (key, str_value, int_value) VALUES ('btc_last_block', '0', 0);

INSERT INTO key_value_store (key, str_value, int_value) VALUES ('xlm_last_ledger_id', '0', 0);

CREATE TYPE chain AS ENUM ('fiat', 'btc', 'eth', 'xlm');
CREATE TYPE order_status AS ENUM ('waiting_for_payment', 'payment_received', 'waiting_user_tx', 'tx_created', 'payment_error', 'finished', 'error', 'under_pay', 'over_pay', 'no_coins_left', 'phase_expired');

CREATE TABLE ico_phase (
  phase_name varchar(255) PRIMARY key NOT NULL,
  start_time timestamp with time zone NOT NULL,
  end_time timestamp with time zone NOT NULL,
  coin_amount bigint not null,
  is_active boolean not null,
  created_at timestamp with time zone NOT NULL default current_timestamp,
  updated_at timestamp with time zone NOT NULL default current_timestamp  
);
create unique index on ico_phase (is_active) where is_active = true; /* only one active at a time */

insert into ico_phase (phase_name, start_time, end_time, coin_amount, is_active) VALUES 
  ('Phase1', '2018-1-1'::timestamp, '2020-1-1'::timestamp, 10000000, true);
insert into ico_phase (phase_name, start_time, end_time, coin_amount, is_active) VALUES 
  ('Phase2', '2018-1-1'::timestamp, '2020-1-1'::timestamp, 10000000, false);


CREATE TABLE user_order (
  id SERIAL PRIMARY KEY NOT null,
  user_id integer not null REFERENCES user_profile(id),

  /* ========== */
  /* order data */
  order_phase_id varchar(255) not null REFERENCES ico_phase(phase_name),
  order_status order_status not null,
  coin_amount bigint not null,
  chain_amount numeric(18, 8) not null, /* max one billion */
  chain_amount_denom varchar(64) not null, /* denomination for chain_amount */
  
  /* ========== */
  /* chain data */
  chain chain NOT NULL,
  address_index bigint NOT NULL, /* used in eth and btc for generating the address */
  /* bitcoin 34 characters */
  /* ethereum 42 characters */  
  /* stellar 56 characters */  
  chain_address varchar(56) NOT NULL,  
  chain_address_seed varchar(56) NOT NULL, /* used only for stellar accounts */

  /* this is the destination address where the coins must be transfered to */
  /* if we use more than one public_key, default should be key_0, therefore not related via user_id */
  user_stellar_public_key varchar(56) NOT NULL, 

  /* this field is used to save any error message that happened during the client payment */
  payment_error_message text not null,
  payment_tx text not null,

  /* ============== */
  /* default fields */
  created_at timestamp with time zone NOT NULL default current_timestamp,
  updated_at timestamp with time zone NOT NULL default current_timestamp,
  updated_by character varying not null,  

  CONSTRAINT valid_address_index CHECK (address_index >= 0),
  CONSTRAINT valid_amount CHECK (coin_amount > 0)
);
create index idx_user_order_user_profile on user_order(user_id);
create unique index idx_user_order_ix1 on user_order(chain, chain_address);
create index idx_user_order_ix2 on user_order(updated_at); /* need this for fast filtering */
create index idx_user_order_ix3 on user_order(order_status);
create index idx_user_order_ix4 on user_order(user_stellar_public_key);
create index idx_user_order_ix5 on user_order(order_phase_id);


CREATE TYPE transaction_status AS ENUM ('new', 'processed');
CREATE TABLE processed_transaction (
  chain chain NOT NULL,
  status transaction_status not null,

  /* Ethereum: "0x"+hash (so 64+2) */
  transaction_id varchar(66) NOT NULL,
  /* bitcoin 34 characters */
  /* ethereum 42 characters */
   /* stellar 56 characters */  
  receiving_address varchar(56) NOT NULL,
  chain_amount_denom varchar(64) not null, /* max one billion */

  user_order_id integer not null REFERENCES user_order(id),

  created_at timestamp with time zone NOT NULL default current_timestamp,
  updated_at timestamp with time zone NOT NULL default current_timestamp,
  PRIMARY KEY (chain, transaction_id)
);
create unique index idx_processed_transaction_ix1 on processed_transaction(user_order_id);


CREATE TABLE multiple_transaction (
  id SERIAL PRIMARY KEY NOT null,
  chain chain NOT NULL,
  transaction_id varchar(66) NOT NULL,
  receiving_address varchar(56) NOT NULL,
  chain_amount_denom varchar(64) not null,
  user_order_id integer not null REFERENCES user_order(id),
  created_at timestamp with time zone NOT NULL default current_timestamp,
  updated_at timestamp with time zone NOT NULL default current_timestamp  
);
create unique index idx_multiple_transaction_ix1 on multiple_transaction(user_order_id);
create index idx_multiple_transaction_ix2 on multiple_transaction(chain, transaction_id);
create index idx_multiple_transaction_ix3 on multiple_transaction(chain, receiving_address);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table IF EXISTS key_value_store;
drop table if exists processed_transaction;
drop table if exists multiple_transaction;
drop table IF EXISTS user_order;
drop table IF EXISTS  ico_phase;
drop type if exists order_status;
drop type if exists transaction_status;
drop type if exists chain;
