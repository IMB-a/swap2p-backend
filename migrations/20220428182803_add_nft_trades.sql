-- +goose Up
alter table trade
    add x_nft_asset text not null default '';
alter table trade
    add y_nft_asset text not null default '';
alter table trade
    add x_nft_index int not null default 0;
alter table trade
    add y_nft_index int not null default 0;

-- +goose Down
alter table trade
    drop x_nft_asset;
alter table trade
    drop y_nft_asset;
alter table trade
    drop x_nft_index;
alter table trade
    drop y_nft_index;


