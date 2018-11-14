-- +goose Up
alter table sse_data drop COLUMN sse_config_id;

-- +goose Down
alter table sse_data add COLUMN sse_config_id boolean not null default false;


