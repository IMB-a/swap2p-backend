-- +goose Up

CREATE TABLE telegram_user
(
    user_id UUID PRIMARY KEY,
    chat_id TEXT UNIQUE NOT NULL,
    state   TEXT        NOT NULL
);

create table address
(
    address TEXT PRIMARY KEY,
    user_id UUID REFERENCES telegram_user
);

create table asset
(
    address  TEXT PRIMARY KEY,
    ticker   TEXT NOT NULL,
    decimals INT  NOT NULL
);

create table balance
(
    asset_address TEXT REFERENCES asset,
    user_id       UUID REFERENCES telegram_user,
    amount        INT NOT NULL DEFAULT 0,
    primary key (user_id, asset_address)
);

-- +goose Down

drop table balance;
drop table asset;
drop table address;
drop table telegram_user;
