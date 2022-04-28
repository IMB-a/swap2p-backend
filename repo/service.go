package repo

import (
	"context"

	"github.com/Pod-Box/swap2p-backend/api"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ Repository = &Service{}

type Service struct {
	db *sqlx.DB
}

func NewService(cfg *Config) (*Service, error) {
	connStr, err := connectionString(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "can't create connection string")
	}

	db, err := sqlx.Connect(cfg.Driver, connStr)

	if err != nil {
		return nil, errors.Wrapf(err, "can't connect to db")
	}

	return &Service{db: db}, nil
}

type Repository interface {
	UserRepository
	TradeRepository
	AssetRepository
	BalanceRepository
}

type UserGetter interface {
	GetAllUsers(ctx context.Context) ([]api.PersonalData, error)
}

type UserRepository interface {
	GetPersonalData(ctx context.Context, chatID string) (*api.PersonalData, error)
	UpsertPersonAddress(ctx context.Context, chatID, address string) error
	UpdatePersonState(ctx context.Context, chatID, state string) error
	UpsertPerson(ctx context.Context, chatID string) error
	UserGetter
}

type TradeFilter struct {
	Closed *bool
}

type TradeRepository interface {
	GetTrades(ctx context.Context, offset, limit int, tf *TradeFilter) (api.TradeList, int, error)
	GetTradesByChatID(ctx context.Context, chatID string) (api.TradeList, error)
	AddTrade(ctx context.Context, trade *api.Trade) error
	TradeExists(ctx context.Context, tradeID int) (bool, error)
	CloseTrade(ctx context.Context, tradeID int, yAddress string) error
}

type AssetRepository interface {
	GetAssets(ctx context.Context) (api.AssetList, error)
	UpdateAsset(ctx context.Context, assetAddress, shortName, fullName string, decimal int64) error
	AddAsset(ctx context.Context, assetAddress, name string, decimal int) error
}

type BalanceRepository interface {
	GetBalancesByAddress(ctx context.Context, address string) (api.Balance, error)
	GetBalancesByChatID(ctx context.Context, chatID string) (api.Balance, error)
}
