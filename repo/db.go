package repo

import (
	"context"
	"strconv"

	"github.com/IMB-a/swap2p-backend/api"
	"github.com/google/uuid"
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

func (s *Service) GetPersonalData(ctx context.Context, chatID string) (*api.PersonalData, error) {
	pd := api.PersonalData{}

	q := `
		select tu.state  as state,
			   a.address as wallet_address
		from telegram_user tu
				 left join address a on tu.user_id = a.user_id
		where tu.chat_id = $1`

	err := s.db.Get(&pd, q, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "select personal data")
	}

	bb, err := s.GetBalances(ctx, chatID)
	if err != nil {
		return nil, err
	}

	pd.Balance = bb

	return &pd, nil
}

func (s *Service) GetBalances(ctx context.Context, chatID string) (api.Balance, error) {
	bb := api.Balance{}

	q := `
		select a.ticker   as asset_name,
			   a.address  as asset_address,
			   b.amount   as amount,
			   a.decimals as asset_decimals
		from balance b
				 join telegram_user tu on b.user_id = tu.user_id
				 join asset a on a.address = b.asset_address
		where tu.chat_id = $1`

	err := s.db.SelectContext(ctx, &bb, q, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "select balance")
	}

	return bb, err
}

const (
	InitialUserState = "start"
)

func (s *Service) UpsertPerson(ctx context.Context, chatID string) error {
	q := `
		insert into telegram_user (user_id, chat_id, state)
		VALUES ($1, $2, $3)
		on conflict on constraint telegram_user_chat_id_key do nothing`

	_, err := s.db.ExecContext(ctx, q, uuid.NewString(), chatID, InitialUserState)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpsertPersonAddress(ctx context.Context, chatID, address string) error {
	q := `
		insert into address (address, user_id)
		VALUES ($1, (select user_id from telegram_user where chat_id=$2))
		on conflict on constraint address_pkey do nothing`

	_, err := s.db.ExecContext(ctx, q, address, chatID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdatePersonState(ctx context.Context, chatID, state string) error {
	q := `update telegram_user set state = $2 where chat_id = $1`

	_, err := s.db.ExecContext(ctx, q, chatID, state)
	if err != nil {
		return err
	}

	return nil
}
