-- +goose Up
alter table trade
    add trade_type text not null default '';
-- +goose Down
alter table trade
    drop trade_type;

