-- +goose Up
alter table trade
    drop constraint trade_pkey;
alter table trade
    add primary key (trade_id, trade_type);

-- +goose Down
alter table trade
    drop constraint trade_pkey;
alter table trade
    add primary key (trade_id);
