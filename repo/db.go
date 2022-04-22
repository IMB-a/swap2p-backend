package repo

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/Pod-Box/swap2p-backend/api"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
)

func (s *Service) GetAllUsers(ctx context.Context) ([]api.PersonalData, error) {
	pd := make([]api.PersonalData, 0)

	q := `
		select tu.state  	as state,
			   a.address	as wallet_address
		from telegram_user tu
				 left join address a on tu.user_id = a.user_id
		where a.address is not null`

	err := s.db.SelectContext(ctx, &pd, q)
	if err != nil {
		return nil, errors.Wrap(err, "select all personal data")
	}

	return pd, nil
}

func (s *Service) UpdateBalance(ctx context.Context, assetAddress, walletAddress string, balance int64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	q := `delete from balance where asset_address = $1 and user_id = (select user_id from address where lower(address) = lower($2))`
	_, err = tx.ExecContext(ctx, q, assetAddress, walletAddress)
	if err != nil {
		tx.Rollback()
		return err
	}
	q = `insert into balance (asset_address, user_id, amount) VALUES ($1, (select user_id from address where lower(address) = lower($2)), $3)`
	_, err = tx.ExecContext(ctx, q, assetAddress, walletAddress, balance)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *Service) GetAssets(ctx context.Context) (api.AssetList, error) {
	q := `
		select address  as asset_address,
			   ticker   as asset_ticker,
			   decimals as asset_decimals
		from asset`

	al := api.AssetList{}
	err := s.db.SelectContext(ctx, &al, q)
	if err != nil {
		return nil, errors.Wrap(err, "get assets")
	}

	return al, nil
}

func (s *Service) CloseTrade(ctx context.Context, tradeID int) error {
	q := `update trade set closed = true where trade_id = $1`

	if exists, err := s.TradeExists(ctx, tradeID); err != nil {
		return err
	} else if !exists {
		q = `insert into trade (trade_id, closed) VALUES ($1,true)`
	}

	_, err := s.db.ExecContext(ctx, q, tradeID)
	if err != nil {
		return errors.Wrap(err, "close trade")
	}

	return nil
}

func (s *Service) TradeExists(ctx context.Context, tradeID int) (bool, error) {
	q := `select trade_id from trade where trade_id = $1`

	var tid int

	err := s.db.GetContext(ctx, &tid, q, tradeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

var TradeAlreadyExistsErr = errors.New("trade already exists")

func (s *Service) AddTrade(ctx context.Context, trade *api.Trade) error {
	if exists, err := s.TradeExists(ctx, trade.Id); err != nil {
		return err
	} else if exists {
		return TradeAlreadyExistsErr
	}

	q := `insert into trade (trade_id, x_address, y_address, x_asset, y_asset, x_amount, y_amount) 
VALUES (:trade_id, :x_address, :y_address, :x_asset, :y_asset, :x_amount, :y_amount)`

	_, err := s.db.NamedExecContext(ctx, q, trade)
	if err != nil {
		return errors.Wrap(err, "add trade")
	}

	return nil
}

func (s *Service) GetTradesByChatID(ctx context.Context, chatID string) (api.TradeList, error) {
	tl := api.TradeList{}

	q := `
		select t.trade_id                as trade_id,
			   x_address,
			   y_address,
			   x_asset,
			   y_asset,
			   x_amount,
			   y_amount,
			   closed,
			   coalesce(xa.decimals, 18) as x_decimals,
			   coalesce(ya.decimals, 18) as y_decimals
		from trade t
				 join address a on a.address = t.x_address or a.address = t.y_address
				 join telegram_user tu on a.user_id = tu.user_id
				 left join asset xa on t.x_asset = xa.address
				 left join asset ya on t.y_asset = xa.address
		where chat_id = $1`

	err := s.db.SelectContext(ctx, &tl, q, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "trades by chat id")
	}

	return tl, nil
}

func (s *Service) GetTrades(ctx context.Context, offset, limit int) (api.TradeList, error) {
	tl := api.TradeList{}

	q := `select distinct t.trade_id                as trade_id,
			   x_address,
			   y_address,
			   x_asset,
			   y_asset,
			   x_amount,
			   y_amount,
			   closed,
			   coalesce(xa.decimals, 18) as x_decimals,
			   coalesce(ya.decimals, 18) as y_decimals
		from trade t
				 left join address a on a.address = t.x_address or a.address = t.y_address
				 left join telegram_user tu on a.user_id = tu.user_id
				 left join asset xa on t.x_asset = xa.address
				 left join asset ya on t.y_asset = xa.address`

	if offset > 0 {
		q += "\n offset " + strconv.Itoa(limit)
	}
	if limit > 0 {
		q += "\n limit " + strconv.Itoa(limit)
	}
	err := s.db.SelectContext(ctx, &tl, q)
	if err != nil {
		return nil, errors.Wrap(err, "all trades")
	}

	return tl, nil
}

func (s *Service) GetPersonalData(ctx context.Context, chatID string) (*api.PersonalData, error) {
	pd := api.PersonalData{}

	q := `
		select tu.state  				as state,
			   coalesce(a.address, '') 	as wallet_address
		from telegram_user tu
				 left join address a on tu.user_id = a.user_id
		where lower(tu.chat_id) = lower($1)`

	err := s.db.Get(&pd, q, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "select personal data")
	}

	bb, err := s.GetBalancesByChatID(ctx, chatID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		bb = api.Balance{}
	}

	pd.Balance = bb

	return &pd, nil
}

func (s *Service) GetBalancesByChatID(ctx context.Context, chatID string) (api.Balance, error) {
	bb := api.Balance{}

	q := `
		select a.ticker   as asset_name,
			   a.address  as asset_address,
			   b.amount   as amount,
			   a.decimals as asset_decimals
		from balance b
				 join telegram_user tu on b.user_id = tu.user_id
				 join asset a on a.address = b.asset_address
		where lower(tu.chat_id) = lower($1)`

	err := s.db.SelectContext(ctx, &bb, q, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "select balance by chat id")
	}

	return bb, err
}

func (s *Service) GetBalancesByAddress(ctx context.Context, address string) (api.Balance, error) {
	bb := api.Balance{}

	q := `
		select a.ticker   as asset_name,
			   a.address  as asset_address,
			   b.amount   as amount,
			   a.decimals as asset_decimals
		from balance b
				 join telegram_user tu on b.user_id = tu.user_id
				 join asset a on a.address = b.asset_address
				 join address ad on ad.user_id = tu.user_id
		where lower(ad.address) = lower($1)`

	err := s.db.SelectContext(ctx, &bb, q, address)
	if err != nil {
		return nil, errors.Wrap(err, "select balance by address")
	}

	return bb, err
}

const (
	InitialUserState = "new"
)

func (s *Service) UpsertPerson(ctx context.Context, chatID string) error {
	// todo on conflict
	q := `
		insert into telegram_user (user_id, chat_id, state)
		VALUES ($1, $2, $3)`

	_, err := s.db.ExecContext(ctx, q, uuid.NewString(), chatID, InitialUserState)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpsertPersonAddress(ctx context.Context, chatID, address string) error {
	// todo on conflict
	q := `
		insert into address (address, user_id)
		VALUES ($1, (select user_id from telegram_user where chat_id=$2))`

	_, err := s.db.ExecContext(ctx, q, address, chatID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdatePersonState(ctx context.Context, chatID, state string) error {
	q := `update telegram_user set state = $2 where lower(chat_id) = lower($1)`

	_, err := s.db.ExecContext(ctx, q, chatID, state)
	if err != nil {
		return err
	}

	return nil
}
