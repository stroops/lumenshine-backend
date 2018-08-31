-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE user_wallet RENAME COLUMN federation_server TO friendly_id;
ALTER TABLE user_wallet ADD COLUMN domain character varying(255) NOT NULL DEFAULT '';

DROP INDEX IF EXISTS idx_user_wallet_fedname;
CREATE unique INDEX idx_user_wallet_fedname ON user_wallet(friendly_id,domain) where friendly_id <> '' and domain <> '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE user_wallet RENAME COLUMN friendly_id TO federation_server;
ALTER TABLE user_wallet DROP COLUMN domain;

DROP INDEX IF EXISTS idx_user_wallet_fedname;
CREATE unique INDEX idx_user_wallet_fedname ON user_wallet(federation_address) where federation_address <> '';