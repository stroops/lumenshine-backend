
/* Setup table */
CREATE TABLE public.setup
(
	id SERIAL PRIMARY KEY NOT NULL,
    asset_code character varying(12) NOT NULL,    
    issuer character varying(56) NOT NULL,
	ledger_seq integer NOT NULL CONSTRAINT positive_ledger_seq CHECK (ledger_seq >= 0),
    active boolean NOT NULL,
	wait_time integer NOT NULL CONSTRAINT positive_wait_time CHECK (wait_time >= 0),
    created_at date NOT NULL default CURRENT_DATE,
    updated_at date NOT NULL default CURRENT_DATE,
	updated_by character varying(20) NOT NULL,
    CONSTRAINT unique_asset UNIQUE (asset_code, issuer)
);
ALTER TABLE public.setup OWNER TO postgres;

/* Progress table */
CREATE TABLE public.progress
(
	id SERIAL PRIMARY KEY NOT NULL,
	setup_id integer NOT NULL,
	ledger_seq integer NOT NULL CONSTRAINT positive_ledger_seq CHECK (ledger_seq >= 0),
    created_at date NOT NULL default CURRENT_DATE,
    updated_at date NOT NULL default CURRENT_DATE,
	updated_by character varying(20) NOT NULL,
    CONSTRAINT unique_setup UNIQUE (setup_id),
	CONSTRAINT "fk_setup" FOREIGN KEY (setup_id) REFERENCES public.setup (id)
);
ALTER TABLE public.setup OWNER TO postgres;


/* Active accounts table */
CREATE TABLE public.active_accounts
(
    id SERIAL PRIMARY KEY NOT NULL,
    setup_id integer NOT NULL,
    account_id character varying(56) NOT NULL,
    change_trust_ledger_seq integer NOT NULL CONSTRAINT positive_ledger_seq CHECK (change_trust_ledger_seq >= 0),
    change_trust_tx_id character varying(64) NOT NULL,
	created_at date NOT NULL default CURRENT_DATE,
    updated_at date NOT NULL default CURRENT_DATE,
	updated_by character varying(20) NOT NULL,
    CONSTRAINT "unique_account" UNIQUE (setup_id, account_id),
    CONSTRAINT "fk_setup" FOREIGN KEY (setup_id) REFERENCES public.setup (id) 
);
ALTER TABLE public.active_accounts OWNER TO postgres;	
CREATE INDEX "fki_fk_setup" ON public.active_accounts USING btree (setup_id);

/* Removed accounts table */	
CREATE TABLE public.removed_accounts
(
    id SERIAL PRIMARY KEY NOT NULL,
    setup_id integer NOT NULL,
    account_id character varying(56) NOT NULL,
    change_trust_ledger_seq integer NOT NULL CONSTRAINT positive_ledger_seq CHECK (change_trust_ledger_seq >= 0),
    change_trust_tx_id character varying(64) NOT NULL,
	created_at date NOT NULL default CURRENT_DATE,
    updated_at date NOT NULL default CURRENT_DATE,
	updated_by character varying(20) NOT NULL,
    CONSTRAINT "unique_removed_account" UNIQUE (setup_id, account_id),
    CONSTRAINT "fk_setup" FOREIGN KEY (setup_id) REFERENCES public.setup (id) 
);
ALTER TABLE public.removed_accounts OWNER TO postgres;
CREATE INDEX "fki_fk_removed_setup" ON public.removed_accounts USING btree (setup_id);