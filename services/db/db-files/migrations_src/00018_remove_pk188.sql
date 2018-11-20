-- +goose Up
alter table user_profile drop column password;
alter table user_security drop column public_key_188;
alter table user_profile add column date_suspended timestamp with time zone NULL;
alter table user_profile add column date_closed timestamp with time zone NULL;

-- +goose Down
alter table user_profile add column password character varying NOT NULL default '';
alter table user_security add column public_key_188 character(56) NOT NULL default '';
alter table user_profile drop column if exists date_suspended;
alter table user_profile drop column if exists date_closed;