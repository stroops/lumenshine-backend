
-- +goose Up
-- SQL in this section is executed when the migration is applied.
DELETE FROM salutations;

with salutation_json (doc) as (
   values
    ('
[
{"lang_code":"EN", "salutation":"Mr."},
{"lang_code":"EN", "salutation":"Mrs."},
{"lang_code":"EN", "salutation":"Miss"},
{"lang_code":"EN", "salutation":"Dr."},
{"lang_code":"EN", "salutation":"Ms."},
{"lang_code":"EN", "salutation":"Prof."},
{"lang_code":"EN", "salutation":"Rev."},
{"lang_code":"EN", "salutation":"Lady"},
{"lang_code":"EN", "salutation":"Sir"},
{"lang_code":"EN", "salutation":"Capt."},
{"lang_code":"EN", "salutation":"Major"},
{"lang_code":"EN", "salutation":"Lt.-Col."},
{"lang_code":"EN", "salutation":"Col."},
{"lang_code":"EN", "salutation":"Lt.-Cmdr."},
{"lang_code":"EN", "salutation":"The Hon."},
{"lang_code":"EN", "salutation":"Cmdr."},
{"lang_code":"EN", "salutation":"Flt. Lt."},
{"lang_code":"EN", "salutation":"Brgdr."},
{"lang_code":"EN", "salutation":"Judge"},
{"lang_code":"EN", "salutation":"Lord"},
{"lang_code":"EN", "salutation":"The Hon. Mrs"},
{"lang_code":"EN", "salutation":"Wng. Cmdr."},
{"lang_code":"EN", "salutation":"Group Capt."},
{"lang_code":"EN", "salutation":"Rt. Hon. Lord"},
{"lang_code":"EN", "salutation":"Revd. Father"},
{"lang_code":"EN", "salutation":"Revd Canon"},
{"lang_code":"EN", "salutation":"Maj.-Gen."},
{"lang_code":"EN", "salutation":"Air Cdre."},
{"lang_code":"EN", "salutation":"Viscount"},
{"lang_code":"EN", "salutation":"Dame"},
{"lang_code":"EN", "salutation":"Rear Admrl."}
]
	 '::json)
)

insert into salutations (lang_code, salutation)
select p.lang_code, p.salutation
from salutation_json l
  cross join lateral json_populate_recordset(null::salutations, doc) as p;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DELETE FROM salutations;