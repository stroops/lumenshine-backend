-- +goose Up
-- SQL in this section is executed when the migration is applied.

INSERT INTO admin_group (name, updated_by) VALUES('Administrators', 'System');
INSERT INTO admin_group (name, updated_by) VALUES('Service', 'System');
INSERT INTO admin_group (name, updated_by) VALUES('Developers', 'System');

-- initial pwd Hawaii11
INSERT INTO admin_user (forename, lastname, email, phone, password, active, last_login, updated_by)
  VALUES ('Cristian', 'Gaudi', 'cristian.gaudi@soneso.com', '+40751043864', '$2a$14$tkW/DpLxiOxJ1.z/SU/wiOUwlPPlkN2Y7hWCLul4wXu3qdgeEnsdW', true, '0001-01-01', 'System');

INSERT INTO admin_user (forename, lastname, email, phone, password, active, last_login, updated_by)
  VALUES ('Udo', 'Polder', 'u.polder@nakamilounge.de', '+496668666', '$2a$14$WNc8SzAb8oMq1ELLrbVQ1O790xIgPt6UvPbCpPbvUmcmjCqAFjisO', true, '0001-01-01', 'System');

INSERT INTO admin_user (forename, lastname, email, phone, password, active, last_login, updated_by)
  VALUES ('User1', 'User1L', 'user1@nakamilounge.de', '+496668666', '$2a$14$WNc8SzAb8oMq1ELLrbVQ1O790xIgPt6UvPbCpPbvUmcmjCqAFjisO', true, '0001-01-01', 'System');


INSERT INTO admin_usergroup (user_id, group_id, updated_by) values (
	(SELECT id from admin_user where email = 'cristian.gaudi@soneso.com'),
	(SELECT id from admin_group where name = 'Administrators'),
	'System'
);

INSERT INTO admin_usergroup (user_id, group_id, updated_by) values (
	(SELECT id from admin_user where email = 'u.polder@nakamilounge.de'),
	(SELECT id from admin_group where name = 'Administrators'),
	'System'
);

INSERT INTO admin_usergroup (user_id, group_id, updated_by) values (
	(SELECT id from admin_user where email = 'u.polder@nakamilounge.de'),
	(SELECT id from admin_group where name = 'Service'),
	'System'
);

INSERT INTO admin_usergroup (user_id, group_id, updated_by) values (
	(SELECT id from admin_user where email = 'u.polder@nakamilounge.de'),
	(SELECT id from admin_group where name = 'Developers'),
	'System'
);

INSERT INTO admin_usergroup (user_id, group_id, updated_by) values (
	(SELECT id from admin_user where email = 'user1@nakamilounge.de'),
	(SELECT id from admin_group where name = 'Service'),
	'System'
);

-- add account
INSERT INTO admin_stellar_account(public_key, name, description, type, updated_by)
VALUES('GBFITCMJXJHPPERECLCRUE4I5XTSALLO5RURDU7IGR2TZOF6XIQXGVLK', 'John Musterman', 'John is an issuer', 'issuing', 'cristian.gaudi@soneso.com');

-- add asset codes
INSERT INTO admin_stellar_asset(issuer_public_key_id, asset_code, updated_by)
VALUES('GBFITCMJXJHPPERECLCRUE4I5XTSALLO5RURDU7IGR2TZOF6XIQXGVLK', 'TOKKE', 'cristian.gaudi@soneso.com');

INSERT INTO admin_stellar_asset(issuer_public_key_id, asset_code, updated_by)
VALUES('GBFITCMJXJHPPERECLCRUE4I5XTSALLO5RURDU7IGR2TZOF6XIQXGVLK', 'CALI', 'cristian.gaudi@soneso.com');

-- add signers
INSERT INTO admin_stellar_signer (stellar_account_public_key_id, name, description, signer_public_key, signer_secret_seed, type, updated_by)
VALUES('GBFITCMJXJHPPERECLCRUE4I5XTSALLO5RURDU7IGR2TZOF6XIQXGVLK', 'John Doe', 'first signer ever', 'ABFITCMJXJHPPERECLCRUE4I5XTSALLO5RURDU7IGR2TZOF6XIQXGVLA','', 'allow_trust', 'cristian.gaudi@soneso.com');

INSERT INTO admin_stellar_signer (stellar_account_public_key_id, name, description, signer_public_key, signer_secret_seed, type, updated_by)
VALUES('GBFITCMJXJHPPERECLCRUE4I5XTSALLO5RURDU7IGR2TZOF6XIQXGVLK', 'Max Musterman', 'second signer ever', 'BBFITCMJXJHPPERECLCRUE4I5XTSALLO5RURDU7IGR2TZOF6XIQXGVLB','', 'allow_trust', 'cristian.gaudi@soneso.com');

INSERT INTO admin_stellar_signer (stellar_account_public_key_id, name, description, signer_public_key, signer_secret_seed, type, updated_by)
VALUES('GBFITCMJXJHPPERECLCRUE4I5XTSALLO5RURDU7IGR2TZOF6XIQXGVLK', 'Jack Doe', 'other signer', 'CBFITCMJXJHPPERECLCRUE4I5XTSALLO5RURDU7IGR2TZOF6XIQXGVLC','', 'other', 'cristian.gaudi@soneso.com');

-- +goose Down
-- SQL in this section1 is executed when the migration is rolled back.
delete from admin_usergroup;
delete from admin_user;
delete from admin_group;

delete from admin_stellar_signer;
delete from admin_stellar_asset;
delete from admin_stellar_account;