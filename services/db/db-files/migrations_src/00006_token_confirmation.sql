-- +goose Up
-- SQL in this section is executed when the migration is applied.

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
drop TABLE IF EXISTS token_history;