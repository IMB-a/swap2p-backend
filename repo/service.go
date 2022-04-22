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
}

type UserRepository interface {
	GetPersonalData(ctx context.Context, chatID string) (*api.PersonalData, error)
	UpsertPersonAddress(ctx context.Context, chatID, address string) error
	UpdatePersonState(ctx context.Context, chatID, state string) error
	UpsertPerson(ctx context.Context, chatID string) error
}

type TradeRepository interface {
	GetTrades(ctx context.Context, offset, limit int) (api.TradeList, error)
	GetTradesByChatID(ctx context.Context, chatID string) (api.TradeList, error)
	AddTrade(ctx context.Context, trade *api.Trade) error
	TradeExists(ctx context.Context, tradeID int) (bool, error)
	CloseTrade(ctx context.Context, tradeID int) error
}

type AssetRepository interface {
	GetAssets(ctx context.Context) (api.AssetList, error)
}
