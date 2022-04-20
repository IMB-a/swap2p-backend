-- +goose Up
alter table balance
    alter column amount type decimal;

-- +goose Down
alter table balance
    alter column amount type int;