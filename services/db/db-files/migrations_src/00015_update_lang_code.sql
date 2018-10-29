-- +goose Up
-- SQL in this section is executed when the migration is applied.

UPDATE salutations SET lang_code=lower(lang_code);
UPDATE countries set lang_code=lower(lang_code);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
UPDATE salutations SET lang_code=upper(lang_code);
UPDATE countries set lang_code=upper(lang_code);