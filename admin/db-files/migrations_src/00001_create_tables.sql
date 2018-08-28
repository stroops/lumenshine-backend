-- +goose Up
-- SQL in this section is executed when the migration is applied.

/* user table */
CREATE TABLE admin_user
(
	id SERIAL PRIMARY KEY NOT NULL,
    forename character varying(64) NOT NULL,
    lastname character varying(64) NOT NULL,
	email character varying NOT NULL,
	phone character varying(64) NOT NULL,
	last_login timestamp with time zone NOT NULL,
	password character varying NOT NULL,
	active boolean NOT NULL default true,
    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
	updated_by character varying NOT NULL,
	CONSTRAINT email_unique UNIQUE(email)
);

/* groups table */
CREATE TABLE admin_group
(
	id SERIAL PRIMARY KEY NOT NULL,
	name character varying(56) NOT NULL,
	created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
	updated_by character varying NOT NULL,
	CONSTRAINT name_unique UNIQUE(name)
);

/* usergroup table */
CREATE TABLE admin_usergroup
(
	id SERIAL PRIMARY KEY NOT NULL,
	user_id integer NOT NULL,
	group_id integer NOT NULL,
	created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
	updated_by character varying NOT NULL,
	CONSTRAINT name_usergroup UNIQUE(user_id, group_id),
	CONSTRAINT "fk_user" FOREIGN KEY (user_id) REFERENCES admin_user (id),
	CONSTRAINT "fk_group" FOREIGN KEY (group_id) REFERENCES admin_group (id)
);
CREATE INDEX "fki_fk_user" ON admin_usergroup(user_id);
CREATE INDEX "fki_fk_group" ON admin_usergroup(group_id);


/*stellar_account_type*/
CREATE TYPE stellar_account_type AS ENUM ('funding','issuing', 'worker');

/* stellar account table */
CREATE TABLE admin_stellar_account
(
	id SERIAL PRIMARY KEY NOT NULL,
    public_key character varying(56) NOT NULL,
    name character varying(256) NOT NULL,
	description character varying NOT NULL,
	type stellar_account_type NOT NULL,
	created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
	updated_by character varying NOT NULL,
	CONSTRAINT account_public_key_unique UNIQUE(public_key)
);

/* stellar account asset table */
CREATE TABLE admin_stellar_asset
(
	id SERIAL PRIMARY KEY NOT NULL,
    issuer_public_key_id character varying(56) NOT NULL,
	asset_code character varying(12) NOT NULL,	
	created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
	updated_by character varying NOT NULL,
	CONSTRAINT account_asset_code_unique UNIQUE(issuer_public_key_id, asset_code),
	CONSTRAINT "fk_stellar_asset_account" FOREIGN KEY (issuer_public_key_id) REFERENCES admin_stellar_account (public_key)
);

/*stellar_signer_type*/
CREATE TYPE stellar_signer_type AS ENUM ('allow_trust','other');

/* stellar signer table */
CREATE TABLE admin_stellar_signer
(
	id SERIAL PRIMARY KEY NOT NULL,
    stellar_account_public_key_id character varying(56) NOT NULL,
	name character varying(256) NOT NULL,
	description character varying NOT NULL,	
	signer_public_key character varying(56) NOT NULL,
	signer_secret_seed character varying(56) NOT NULL,
	type stellar_signer_type NOT NULL,
	created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
	updated_by character varying NOT NULL,
	CONSTRAINT "fk_stellar_signer_account" FOREIGN KEY (stellar_account_public_key_id) REFERENCES admin_stellar_account (public_key)
);

/* stellar known currencies table */
CREATE TABLE admin_known_currencies
(
	id SERIAL PRIMARY KEY NOT NULL,
	name character varying(256) NOT NULL,	
    issuer_public_key character varying(56) NOT NULL,
	asset_code character varying(12) NOT NULL,	
	short_description character varying NOT NULL,	
	long_description character varying NOT NULL,	
	order_index integer NOT NULL,	
	created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
	updated_by character varying NOT NULL,
	CONSTRAINT known_currency_unique UNIQUE(issuer_public_key, asset_code)
);

/* stellar known inflation destination table */
CREATE TABLE admin_known_inflation_destinations
(
	id SERIAL PRIMARY KEY NOT NULL,
	name character varying(256) NOT NULL,	
    issuer_public_key character varying(56) NOT NULL,
	short_description character varying NOT NULL,	
	long_description character varying NOT NULL,	
	order_index integer NOT NULL,	
	created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
	updated_by character varying NOT NULL,
	CONSTRAINT known_inflation_destination_unique UNIQUE(issuer_public_key)
);

/*stellar_trustline_status*/
CREATE TYPE stellar_trustline_status AS ENUM ('ok','waiting','denied','revoked');

/* admin_unauthorized_trustline table */
CREATE TABLE admin_unauthorized_trustline
(
	id SERIAL PRIMARY KEY NOT NULL,
    issuer_public_key_id character varying(56) NOT NULL,
	trustor_public_key character varying(56) NOT NULL,	
	asset_code character varying(12) NOT NULL,
	status stellar_trustline_status NOT NULL,
	reason character varying(1000) NOT NULL,
	created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
	updated_by character varying NOT NULL,
	CONSTRAINT trustline_asset_code_unique UNIQUE(trustor_public_key, issuer_public_key_id, asset_code),	
	CONSTRAINT "fk_unauthorized_trustline_issuer_public_key" FOREIGN KEY (issuer_public_key_id) REFERENCES admin_stellar_account (public_key)
);

-- +goose Down
-- SQL in this section1 is executed when the migration is rolled back.
drop table IF EXISTS admin_usergroup;
drop table IF EXISTS admin_user;
drop table IF EXISTS admin_group;

drop table IF EXISTS admin_stellar_signer;
drop table IF EXISTS admin_stellar_asset;
drop table IF EXISTS admin_stellar_account;

drop table IF EXISTS admin_known_currencies;
drop table IF EXISTS admin_known_inflation_destination;

drop table IF EXISTS admin_unauthorized_trustline;

drop type IF EXISTS stellar_account_type;
drop type IF EXISTS stellar_signer_type;

drop type IF EXISTS stellar_trustline_status;