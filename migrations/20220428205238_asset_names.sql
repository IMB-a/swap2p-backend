-- +goose Up
alter table asset
    rename column ticker to short_name;
alter table asset
    add column full_name text default '';

-- +goose Down
alter table asset
    drop column full_name;
alter table asset
    rename column short_name to ticker;
