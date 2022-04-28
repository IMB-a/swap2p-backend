-- +goose Up
alter table balance
    drop constraint balance_asset_address_fkey;

-- +goose Down
alter table balance
    add constraint balance_asset_address_fkey foreign key (asset_address) references asset (address);

