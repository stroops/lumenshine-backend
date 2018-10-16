-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TYPE payment_state AS ENUM ('open', 'close');
CREATE TYPE kyc_status AS ENUM ('not_supported','waiting_for_data', 'waiting_for_review','in_review', 'pending', 'rejected', 'approved');

CREATE TABLE user_profile
(
    id SERIAL PRIMARY KEY NOT null,
    email character varying NOT NULL,
    forename character varying(64) NOT NULL,
    lastname character varying(64) NOT NULL,
    company character varying(128) NOT NULL,
    salutation character varying(64) NOT NULL,
    title character varying(64) NOT NULL,

    street_address character varying(128) NOT NULL,
    street_number character varying(128) NOT NULL,
    zip_code character varying(32) NOT NULL,
    city character varying(128) NOT NULL,
    state character varying(128) NOT NULL,
    country_code character varying(128) NOT NULL,
    nationality character varying(2) NOT NULL,

    mobile_nr character varying(64) NOT NULL,
    birth_day date NOT NULL,
    birth_place character varying(128) NOT NULL,

	mail_confirmation_key character varying NOT NULL,
	mail_confirmation_expiry_date timestamp with time zone NOT NULL,

    tfa_secret character varying NOT NULL,
    tfa_temp_secret character varying NOT NULL,

    mail_confirmed boolean NOT null default false,
    tfa_confirmed bool NOT NULL default false,
    mnemonic_confirmed boolean NOT null default false,

    message_count integer not null default 0,
    payment_state payment_state not null default 'close',

    kyc_status kyc_status NOT NULL DEFAULT 'not_supported',

    password character varying not null,

    stellar_account_created boolean not null default false,

    reset2fa_by_admin boolean not null default false,

    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
    updated_by character varying not null
);
CREATE INDEX user_ix_profile_email ON user_profile USING gin (email gin_trgm_ops);
CREATE INDEX user_ix_profile_forename ON user_profile USING gin (forename gin_trgm_ops);
CREATE INDEX user_ix_profile_lastname ON user_profile USING gin (lastname gin_trgm_ops);
CREATE UNIQUE index idx_mail_confirmation_key on user_profile(mail_confirmation_key) where mail_confirmation_key<>'';

CREATE TABLE user_security
(
    id SERIAL PRIMARY KEY not null,
    user_id integer not null REFERENCES user_profile (id),
    kdf_salt character(44) not null,
    mnemonic_master_key character(44) not null,
    mnemonic_master_iv character(24) not null,
    wordlist_master_key character(44) not null,
    wordlist_master_iv character(24) not null,
    mnemonic character varying not null,
    mnemonic_iv character(24) not null,
    wordlist character varying not null,
    wordlist_iv character(24) not null,
    public_key_0 character(56) NOT NULL,
    public_key_188 character(56) NOT NULL,

    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
    updated_by character varying not null
);
create UNIQUE index idx_user_security_user_profile on user_security(user_id);
CREATE UNIQUE index idx_user_security_key0 on user_security(public_key_0);
CREATE UNIQUE index idx_user_security_key188 on user_security(public_key_188);

create table salutations
(
    id SERIAL PRIMARY KEY not null,
    lang_code character(2) not null,
    salutation character varying(60) not null
);
CREATE index idx_salutations_lang_code on salutations(lang_code);
CREATE unique index idx_salutations on salutations(lang_code, salutation);
insert into salutations(lang_code, salutation) values
  ('de', 'Herr'),
  ('de', 'Frau'),
  ('de', 'Familie'),
  ('en', 'Mr.'),
  ('en', 'Ms.');


create table countries
(
    id SERIAL PRIMARY KEY not null,
    lang_code character(2) not null,
    country_name character varying(256) not null
);
CREATE index idx_country_lang_code on countries(lang_code);
CREATE unique index idx_country on countries(lang_code, country_name);
insert into countries(lang_code, country_name) values
  ('de', 'Deutschland'),
  ('en', 'Germany'),
  ('de', 'Amerika'),
  ('en', 'America');

create table jwt_key
(
    id SERIAL PRIMARY KEY not null,
    key_name character varying not null,
    key_value1 character varying not null,
    key_value2 character varying not null,
    key_description character varying(256) not null,
    valid1_to timestamp with time zone NOT NULL,
    valid2_to timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL default current_timestamp
);
CREATE unique index idx_key_name on jwt_key(key_name);

insert into jwt_key(key_name, key_value1, key_value2, key_description, valid1_to, valid2_to, updated_at) values
  ('simple', 'default simple1', 'default simple2', 'JWT Keys for simple authentication', to_timestamp(0), to_timestamp(0), to_timestamp(0)),
  ('full', 'default full1', 'default full2', 'JWT Key for full authentication', to_timestamp(0), to_timestamp(0), to_timestamp(0)),
  ('lostpwd', 'default lostpwd1', 'default lostpwd2', 'JWT Key for temporary lost pwd authentication', to_timestamp(0), to_timestamp(0), to_timestamp(0));

create table mail
(
    id SERIAL PRIMARY KEY not null,
    mail_from character varying not null,
    mail_to character varying not null,
    mail_subject character varying not null,
    mail_body  character varying not null,
    external_status character varying not null,
    external_status_id character varying not null,
    internal_status bigint not null,

    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_by character varying not null
);

CREATE TABLE user_message
(
    id SERIAL PRIMARY KEY not null,
    user_id integer not null REFERENCES user_profile(id),
    title character varying not null,
    message character varying not null,
    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
    updated_by character varying not null
);
create index idx_user_message_user_profile on user_message(user_id);

-- +goose StatementBegin
CREATE or replace FUNCTION count_user_messages() RETURNS trigger as $$
    begin

	    IF (TG_OP = 'DELETE') then
	        update user_profile set message_count = (select count(id) from user_message where user_id=old.user_id);
            RETURN OLD;
        ELSIF (TG_OP = 'UPDATE') then
            update user_profile set message_count = (select count(id) from user_message where user_id=new.user_id);
            RETURN NEW;
        ELSIF (TG_OP = 'INSERT') then
            update user_profile set message_count = (select count(id) from user_message where user_id=new.user_id);
            RETURN NEW;
        END IF;

        RETURN null;
    END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER count_user_messages after insert or update or delete on user_message for each row execute procedure count_user_messages();


CREATE TABLE user_message_archive
(
    id SERIAL PRIMARY KEY not null,
    user_id integer not null REFERENCES user_profile(id),
    title character varying not null,
    message character varying not null,
    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
    updated_by character varying not null
);
create index idx_user_message_archive_user_profile on user_message_archive(user_id);

CREATE TABLE user_wallet
(
    id SERIAL PRIMARY KEY not null,
    user_id integer not null REFERENCES user_profile(id),
    public_key_0 character(56) NOT NULL,
    wallet_name character varying(500) NOT NULL,
    friendly_id character varying(256) NOT NULL  DEFAULT '',
    domain character varying(256) NOT NULL DEFAULT '',
    show_on_homescreen boolean not null default true,
    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
    updated_by character varying not null
);

create index idx_user_wallet_user_profile on user_wallet(user_id);
CREATE UNIQUE index idx_user_wallet_pub_key on user_wallet(public_key_0);
CREATE INDEX idx_user_wallet_name ON user_wallet USING gin (wallet_name gin_trgm_ops);
CREATE UNIQUE index idx_user_wallet_name_2 on user_wallet(user_id, wallet_name);
CREATE unique INDEX idx_user_wallet_fedname ON user_wallet(friendly_id,domain) where friendly_id <> '' and domain <> '';

/*message_type*/
CREATE TYPE message_type AS ENUM ('ios','android', 'mail');

/*mail_content_type*/
CREATE TYPE mail_content_type AS ENUM ('text','html');

CREATE TYPE notification_status_code AS ENUM ('new', 'success', 'error');

CREATE TABLE notification
(
    id SERIAL PRIMARY KEY not null,
    user_id integer not null REFERENCES user_profile(id),
    push_token character varying NOT NULL,
    type message_type NOT NULL,
    content character varying NOT NULL,
    mail_subject character varying NOT NULL,
    mail_type mail_content_type NOT NULL,
    user_email character varying NOT NULL,
    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
    updated_by character varying not null
);
create index idx_notification_user_profile on notification(user_id);

CREATE TABLE notification_archive
(
    id SERIAL PRIMARY KEY not null,
    user_id integer not null REFERENCES user_profile(id),
    push_token character varying NOT NULL,
    type message_type NOT NULL,
    content character varying NOT NULL,
    mail_subject character varying NOT NULL,
    mail_type mail_content_type NOT NULL,
    user_email character varying NOT NULL,
    status notification_status_code NOT NULL,
    internal_error_string character varying NOT NULL,
    external_status_code character varying NOT NULL,
    external_error_string character varying NOT NULL,
    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
    updated_by character varying not null
);
create index idx_notification_archive_user_profile on notification_archive(user_id);

/*device_type*/
CREATE TYPE device_type AS ENUM ('apple','google');

CREATE TABLE user_pushtoken
(
    id SERIAL PRIMARY KEY not null,
    user_id integer not null REFERENCES user_profile(id),
    device_type device_type NOT NULL,
    push_token character varying NOT NULL,
    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
    updated_by character varying not null
);
CREATE UNIQUE index idx_user_pushtoken on user_pushtoken(push_token);
create index idx_user_pushtoken_user_profile on user_pushtoken(user_id);

/*kyc document_type*/
CREATE TYPE kyc_document_type AS ENUM ('passport','drivers_license','id_card','proof_of_residence');
CREATE TYPE kyc_document_format AS ENUM('png','pdf','jpg','jpeg');
CREATE TYPE kyc_document_side AS ENUM('front','back');

CREATE TABLE user_kyc_document
(
    id SERIAL PRIMARY KEY not null,
    user_id integer not null REFERENCES user_profile(id),
    type kyc_document_type not null,
    format kyc_document_format not null,
    side kyc_document_side not null,
    id_country_code character varying not null,
    id_issue_date timestamp NOT NULL,
    id_expiration_date timestamp NOT NULL,
    id_number character varying not null,
    upload_date timestamp with time zone NOT NULL default current_timestamp,
    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
    updated_by character varying not null
);
create index idx_user_kyc_document_user_profile on user_kyc_document(user_id);


CREATE TABLE user_contact
(
    id SERIAL PRIMARY KEY not null,
    user_id integer not null REFERENCES user_profile(id),
    contact_name character varying not null,
    stellar_address character varying(256) not null,
    public_key character varying(56) not null,
    created_at timestamp with time zone NOT NULL default current_timestamp,
    updated_at timestamp with time zone NOT NULL default current_timestamp,
    updated_by character varying not null
);
create index idx_user_contact_user_profile on user_contact(user_id);

CREATE TABLE token_history
(
    id SERIAL PRIMARY KEY NOT null,
    user_id integer not null REFERENCES user_profile (id),
    mail_confirmation_key character varying NOT NULL,
    created_at timestamp with time zone NOT NULL default current_timestamp
);

CREATE UNIQUE index idx_mail_confirmation_key2 on token_history(mail_confirmation_key) where mail_confirmation_key<>'';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table IF EXISTS user_kyc_document;
drop table IF EXISTS user_message_archive;
drop table IF EXISTS user_message;
drop function if exists  count_user_messages;
drop table IF EXISTS user_security;
drop table IF EXISTS salutations;
drop table IF EXISTS countries;
drop table IF EXISTS jwt_key;
drop table IF EXISTS mail;
drop table IF EXISTS user_wallet;
drop TABLE IF EXISTS notification;
drop TABLE IF EXISTS notification_archive;
drop TABLE IF EXISTS user_pushtoken;
drop TABLE IF EXISTS user_contact;
drop TABLE IF EXISTS token_history;
DROP TABLE IF EXISTS user_profile;

drop type IF EXISTS payment_state;
drop type IF EXISTS message_type;
drop type IF EXISTS mail_content_type;
drop type IF EXISTS notification_status_code;
drop type IF EXISTS device_type;
drop type IF EXISTS kyc_status;
drop type IF EXISTS kyc_document_type;
drop type IF EXISTS kyc_document_format;
drop type IF EXISTS kyc_document_side;