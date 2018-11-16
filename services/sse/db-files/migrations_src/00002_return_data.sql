-- +goose Up

alter table sse_config add COLUMN return_data boolean not null default false;

-- +goose Down
alter table sse_config drop COLUMN return_data;