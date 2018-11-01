-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE source_receiver AS ENUM ('payment', 'notify', 'sse');

CREATE TABLE sse_config
(
    id SERIAL PRIMARY KEY NOT null,
    source_receiver source_receiver not null,
    stellar_account varchar(56) NOT NULL,
    operation_types bigint not null,

    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp
);
CREATE INDEX sse_config_ix1 ON sse_config(source_receiver);
CREATE index sse_config_ix2 on sse_config(source_receiver, stellar_account);

CREATE TABLE sse_data
(
    id SERIAL PRIMARY KEY not null,
    sse_config_id integer not null REFERENCES sse_config (id),
    source_receiver source_receiver not null,
    stellar_account varchar(56) NOT NULL,
    operation_types bigint not null,

    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp
);
CREATE INDEX sse_data_ix1 ON sse_data(source_receiver);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table IF EXISTS sse_data;
drop table IF EXISTS sse_config;

drop type IF EXISTS source_receiver;