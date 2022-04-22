-- +goose Up
create table trade
(
    trade_id  int primary key,
    x_address text    not null default '',
    y_address text    not null default '',
    x_asset   text    not null default '',
    y_asset   text    not null default '',
    x_amount  decimal not null default 0,
    y_amount  decimal not null default 0,
    closed    bool    not null default false
);

-- +goose Down
drop table trade;
