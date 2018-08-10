
-- +goose Up
-- SQL in this section is executed when the migration is applied.

/* Snapshot table */
CREATE TABLE public.snapshot
(
	id SERIAL PRIMARY KEY NOT NULL,
    asset_code character varying(12) NOT NULL,
    issuer character varying(56) NOT NULL,	
    created_at date NOT NULL default CURRENT_DATE,
    updated_at date NOT NULL default CURRENT_DATE,
	updated_by character varying(20) NOT NULL
);
ALTER TABLE public.snapshot OWNER TO postgres;

/* dividend table */
CREATE TABLE public.dividend
(
	id SERIAL PRIMARY KEY NOT NULL,
	snapshot_id integer NOT NULL,
	account_id character varying(56) NOT NULL,
	balance_limit bigint NOT NULL,
	balance bigint NOT NULL,
    dividend_amount bigint NULL,
	CONSTRAINT dividend_balance_limit_check CHECK (balance_limit > 0),
    CONSTRAINT dividend_balance_check CHECK (balance >= 0),
	CONSTRAINT "fk_snapshot" FOREIGN KEY (snapshot_id) REFERENCES public.snapshot (id)
);
ALTER TABLE public.dividend OWNER TO postgres;
CREATE INDEX "fki_fk_snapshot" ON public.dividend USING btree (snapshot_id);



-- +goose Down
-- SQL in this section1 is executed when the migration is rolled back.