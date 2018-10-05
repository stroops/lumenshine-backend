-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE order_transaction_log (
  id SERIAL PRIMARY KEY NOT null,
  order_id integer not null REFERENCES user_order(id),
  status bool not null,
  tx text not null,
  tx_hash text not null,
  result_code text null,
  result_xdr text null,
  error_text text null,
  created_at timestamp with time zone NOT NULL default current_timestamp
);
create index order_transaction_log_ix1 on order_transaction_log(order_id);

CREATE TYPE channel_status AS ENUM ('free', 'in_use', 'merge_reserved', 'merged');
CREATE TABLE channels (
  id SERIAL PRIMARY KEY NOT null,
  pk varchar(56) NOT NULL,
  seed varchar(56) NOT NULL,
  status channel_status not null,
  created_at timestamp with time zone NOT NULL default current_timestamp,
  updated_at timestamp with time zone NOT NULL default current_timestamp
);
create unique index idx_channels_ix1 on channels(pk);

alter table user_order add fee_payed bool not null default false;
alter table user_order ALTER COLUMN payment_seed type varchar(500);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table IF EXISTS order_transaction_log;
drop table IF EXISTS channels;
drop type if exists channel_status;
alter table user_order drop COLUMN if exists fee_payed;
alter table user_order ALTER COLUMN payment_seed type varchar(56) USING SUBSTR(payment_seed, 1, 55);
