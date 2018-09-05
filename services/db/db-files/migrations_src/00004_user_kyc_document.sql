-- +goose Up
-- SQL in this section is executed when the migration is applied.

/*document_type*/
CREATE TYPE document_type AS ENUM ('passport','drivers_license','id_card','proof_of_residence');
CREATE TYPE document_format AS ENUM('png','pdf','jpg','jpeg');
CREATE TYPE document_side AS ENUM('front','back');

CREATE TABLE user_kyc_document
(
    id SERIAL PRIMARY KEY not null,
    user_id integer not null REFERENCES user_profile(id),
    type document_type not null,
    format document_format not null,
    side document_side not null,
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


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

drop table IF EXISTS user_kyc_document;

drop type IF EXISTS document_type;
drop type IF EXISTS document_format;
drop type IF EXISTS document_side;