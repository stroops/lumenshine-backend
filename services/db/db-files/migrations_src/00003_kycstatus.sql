-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TYPE kyc_status AS ENUM ('not_supported','waiting_for_data', 'waiting_for_review','in_review', 'pending', 'rejected', 'approved');

ALTER TABLE user_profile ADD COLUMN kyc_status kyc_status NOT NULL DEFAULT 'not_supported';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE user_profile DROP COLUMN kyc_status;

drop type IF EXISTS kyc_status;