-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE currency (
    id SERIAL PRIMARY KEY NOT NULL,
    currency_code character varying(10) NOT NULL,
    currency_name text NOT NULL,
    currency_issuer text NOT NULL,
    created_at date DEFAULT CURRENT_DATE NOT NULL,
    updated_at date DEFAULT CURRENT_DATE NOT NULL,
    updated_by character varying(20) NOT NULL,
    CONSTRAINT "unique_currency" UNIQUE (currency_code, currency_name, currency_issuer)
);

CREATE TABLE current_chart_data_minutely (
    id SERIAL PRIMARY KEY NOT NULL,
    exchange_rate_time timestamp with time zone NOT NULL DEFAULT current_timestamp,
    source_currency_id integer NOT NULL,
    destination_currency_id integer NOT NULL,
    exchange_rate double precision NOT NULL,
    created_at date DEFAULT CURRENT_DATE NOT NULL,
    updated_at date DEFAULT CURRENT_DATE NOT NULL,
    updated_by character varying(20) NOT NULL,
    CONSTRAINT "fk-source_currency" FOREIGN KEY (source_currency_id) REFERENCES currency (id),
    CONSTRAINT "fk-destination_currency" FOREIGN KEY (destination_currency_id) REFERENCES currency (id)
);

CREATE INDEX ON current_chart_data_minutely (source_currency_id);

CREATE INDEX ON current_chart_data_minutely (destination_currency_id);

ALTER TABLE current_chart_data_minutely ADD CONSTRAINT current_chart_data_minutely_un UNIQUE (exchange_rate_time,source_currency_id,destination_currency_id);

CREATE TABLE current_chart_data_hourly (
    id SERIAL PRIMARY KEY NOT NULL,
    exchange_rate_time timestamp with time zone NOT NULL default current_timestamp,
    source_currency_id integer NOT NULL,
    destination_currency_id integer NOT NULL,
    exchange_rate double precision NOT NULL,
    created_at date DEFAULT CURRENT_DATE NOT NULL,
    updated_at date DEFAULT CURRENT_DATE NOT NULL,
    updated_by character varying(20) NOT NULL,
    CONSTRAINT "fk-source_currency" FOREIGN KEY (source_currency_id) REFERENCES currency (id),
    CONSTRAINT "fk-destination_currency" FOREIGN KEY (destination_currency_id) REFERENCES currency (id)
);

CREATE INDEX ON current_chart_data_hourly (source_currency_id);

CREATE INDEX ON current_chart_data_hourly (destination_currency_id);

ALTER TABLE current_chart_data_hourly ADD CONSTRAINT current_chart_data_hourly_un UNIQUE (exchange_rate_time,source_currency_id,destination_currency_id);

CREATE TABLE history_chart_data (
    id SERIAL PRIMARY KEY NOT NULL,
    exchange_rate_date date NOT NULL,
    source_currency_id integer NOT NULL,
    destination_currency_id integer NOT NULL,
    exchange_rate double precision NOT NULL,
    created_at date DEFAULT CURRENT_DATE NOT NULL,
    updated_at date DEFAULT CURRENT_DATE NOT NULL,
    updated_by character varying(20) NOT NULL,
    CONSTRAINT "fk-source_currency" FOREIGN KEY (source_currency_id) REFERENCES currency (id),
    CONSTRAINT "fk-destination_currency" FOREIGN KEY (destination_currency_id) REFERENCES currency (id)
);

CREATE INDEX ON history_chart_data (source_currency_id);

CREATE INDEX ON history_chart_data (destination_currency_id);

ALTER TABLE history_chart_data ADD CONSTRAINT history_chart_data_un UNIQUE (exchange_rate_date,source_currency_id,destination_currency_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS currency CASCADE;
DROP TABLE current_chart_data_minutely CASCADE;
DROP TABLE current_chart_data_hourly CASCADE;
DROP TABLE IF EXISTS history_chart_data CASCADE;