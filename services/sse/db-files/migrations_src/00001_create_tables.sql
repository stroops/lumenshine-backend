-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE source_receiver AS ENUM ('payment', 'notify', 'sse');

CREATE TABLE sse_config
(
    id SERIAL PRIMARY KEY NOT null,
    source_receiver source_receiver not null,
    stellar_account varchar(56) NOT NULL,
    operation_types bigint not null,
    with_resume boolean not null default false,

    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp
);
create unique index sse_config_ix1 ON sse_config(source_receiver, stellar_account);
CREATE INDEX sse_config_ix2 ON sse_config(stellar_account, with_resume);


CREATE TYPE sse_data_status AS ENUM ('new', 'selected');
CREATE TABLE sse_data
(
    id SERIAL PRIMARY KEY not null,
    sse_config_id integer not null REFERENCES sse_config (id),
    source_receiver source_receiver not null,
    status sse_data_status not null,
    stellar_account varchar(56) NOT NULL,
    operation_type int not null,
    operation_data jsonb null,
    transaction_id bigint not null,
    operation_id bigint not null,
    ledger_id bigint not null,

    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp
);
CREATE INDEX sse_data_ix1 ON sse_data(source_receiver);
CREATE INDEX sse_data_ix2 ON sse_data(source_receiver, stellar_account);
CREATE INDEX sse_data_ix3 ON sse_data(status);

CREATE TYPE sse_index_names AS ENUM ('last_ledger_id');
CREATE TABLE sse_index
(
    id SERIAL PRIMARY KEY not null,
    name sse_index_names not null,
    value bigint not null default 0
);
CREATE INDEX sse_index_ix1 ON sse_index(name);

insert into  sse_index(name, value) values ('last_ledger_id', 0);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table IF EXISTS sse_data;
drop table IF EXISTS sse_config;
drop table if exists sse_index;

drop type IF EXISTS source_receiver;
drop type IF EXISTS sse_data_status;
drop type IF EXISTS sse_index_names;