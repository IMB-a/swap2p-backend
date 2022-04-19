-- +goose Up

CREATE TABLE telegram_user
(
    user_id UUID PRIMARY KEY,
    chat_id TEXT UNIQUE NOT NULL
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
    address TEXT PRIMARY KEY REFERENCES asset,
    user_id UUID PRIMARY KEY REFERENCES telegram_user,
    amount  INT NOT NULL DEFAULT 0
);

-- +goose Down

drop table balance;
drop table asset;
drop table address;
drop table telegram_user;
