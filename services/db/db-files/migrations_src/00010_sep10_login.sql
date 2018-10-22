-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE user_profile ADD COLUMN public_key_0 varchar(56) not null default '';
update user_profile set public_key_0 = (select public_key_0 from user_security where user_id=user_profile.id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE user_profile drop COLUMN public_key_0;