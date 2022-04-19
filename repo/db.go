package repo

import (
	"context"
	"strconv"

	"github.com/IMB-a/swap2p-backend/api"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
)

func (s *Service) GetTrades(ctx context.Context, offset, limit int) (api.TradeList, error) {
	panic("implement me")
	tl := api.TradeList{}

	q := ``

	if offset > 0 {
		q += "\n offset " + strconv.Itoa(limit)
	}
	if limit > 0 {
		q += "\n limit " + strconv.Itoa(limit)
	}
	s.db.SelectContext(ctx, &tl, q)

	return tl, nil
}

func (s *Service) GetPersonalData(ctx context.Context, chatID string) (api.PersonalData, error) {
	pd := api.PersonalData{}
	se
}

func (s *Service) GetBalances(ctx context.Context, chatID string) (api.Balance, error) {
	bb := api.Balance{}

	q := `select * from `
	err := s.db.SelectContext(ctx, &bb, q, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "select balance")
	}
}

func (s *Service) UpsertPersonAddress(ctx context.Context, chatID, address string) error {
	// TODO implement me
	panic("implement me")
}

func (s *Service) UpdatePersonState(ctx context.Context, chatID, state string) error {
	// TODO implement me
	panic("implement me")
}
