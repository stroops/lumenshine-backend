-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Create foreign table trustlines
CREATE FOREIGN TABLE trustlines
(
    accountid character varying(56) NOT NULL,
    assettype integer NOT NULL,
    issuer character varying(56) NOT NULL,
    assetcode character varying(12) NOT NULL,
    tlimit bigint NOT NULL,
    balance bigint NOT NULL,
    flags integer NOT NULL,
    lastmodified integer NOT NULL,
    buyingliabilities bigint,
    sellingliabilities bigint
)
SERVER fdw_stellarcore_server;

-- Create foreign table user_profile
CREATE FOREIGN TABLE user_profile
(
    id integer NOT NULL,
    email character varying NOT NULL,
    forename character varying(64) NOT NULL,
    lastname character varying(64) NOT NULL
)
SERVER fdw_icop_server;

-- Create foreign table user_security
CREATE FOREIGN TABLE user_security
(
    id integer NOT NULL,
    user_id integer NOT NULL,
    public_key_0 character(56) NOT NULL
)
SERVER fdw_icop_server;

--Create VIEW customer_trustlines
CREATE OR REPLACE VIEW customer_trustlines AS
	SELECT CONCAT(user_profile.forename, ' ', user_profile.lastname) as name, 
		   user_security.public_key_0 as public_key, 
		   trustlines.issuer as issuer_public_key, 
		   trustlines.assetcode as asset_code, 
		  CASE
				WHEN (trustlines.flags = 1) THEN 'ok'
				WHEN (trustlines.flags = 0 AND aut.status IS NULL) THEN 'waiting'
				WHEN (aut.status IS NOT NULL) THEN aut.status
		   END as status,
		   aut.reason as reason
	FROM trustlines
		INNER JOIN user_security ON trustlines.accountid = user_security.public_key_0
		INNER JOIN user_profile ON user_profile.id = user_security.user_id
		LEFT JOIN admin_unauthorized_trustline AS aut 
			ON aut.issuer_public_key_id = trustlines.issuer AND aut.asset_code = trustlines.assetcode AND aut.trustor_public_key = trustlines.accountid;

-- Create view admin_trustlines
CREATE OR REPLACE VIEW admin_trustlines AS
	SELECT admin_stellar_account.name as name, 
	       admin_stellar_account.public_key as public_key, 
		   trustlines.issuer as issuer_public_key, 
		   trustlines.assetcode as asset_code, 
		   CASE
				WHEN (trustlines.flags = 1) THEN 'ok'
				WHEN (trustlines.flags = 0 AND aut.status IS NULL) THEN 'waiting'
				WHEN (aut.status IS NOT NULL) THEN aut.status
		   END as status,
		   aut.reason as reason
	FROM trustlines
		INNER JOIN admin_stellar_account ON trustlines.accountid = admin_stellar_account.public_key
		LEFT JOIN admin_unauthorized_trustline AS aut 
			ON aut.issuer_public_key_id = trustlines.issuer AND aut.asset_code = trustlines.assetcode AND aut.trustor_public_key = trustlines.accountid;

-- +goose Down
-- SQL in this section1 is executed when the migration is rolled back.
DROP VIEW IF EXISTS customer_trustlines;
DROP VIEW IF EXISTS admin_trustlines;

drop FOREIGN table IF EXISTS trustlines;
drop FOREIGN table IF EXISTS user_profile;
drop FOREIGN table IF EXISTS user_security;