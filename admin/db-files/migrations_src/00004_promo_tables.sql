-- +goose Up
-- SQL in this section is executed when the migration is applied.

/*promo_type*/
CREATE TYPE promo_type AS ENUM ('small','big');

/* admin promo table */
CREATE TABLE admin_promo
(
	id SERIAL PRIMARY KEY NOT NULL,
	name character varying(256) NOT NULL,	
    title character varying(512) NOT NULL,
	promo_text text NOT NULL,	
	promo_type promo_type NOT NULL,	
	buttons text NOT NULL,	
	active boolean NOT NULL default false,
	order_index integer NOT NULL,		
	created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
	updated_by character varying NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

drop table IF EXISTS admin_promo;

drop type IF EXISTS promo_type;