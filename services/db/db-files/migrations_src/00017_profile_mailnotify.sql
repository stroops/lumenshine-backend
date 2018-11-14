-- +goose Up

alter table user_profile add column mail_notifications boolean not null default false;

-- +goose Down
alter table user_profile drop column mail_notifications;