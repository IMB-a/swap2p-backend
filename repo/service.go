package repo

import (
	"context"

	"github.com/IMB-a/swap2p-backend/api"
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

	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		return nil, errors.Wrapf(err, "can't connect to db")
	}

	return &Service{db: db}, nil
}

type Repository interface {
	GetTrades(ctx context.Context, offset, limit int) (api.TradeList, error)
	GetPersonalData(ctx context.Context, chatID string) (*api.PersonalData, error)
	UpsertPersonAddress(ctx context.Context, chatID, address string) error
	UpdatePersonState(ctx context.Context, chatID, state string) error
}
