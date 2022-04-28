package assets

import (
	"context"
	"time"

	"github.com/Pod-Box/swap2p-backend/repo"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/contract/builtin/erc20"
	"github.com/umbracle/ethgo/jsonrpc"
)

type BalanceUpdater interface {
	UpdateBalance(ctx context.Context, assetAddress, walletAddress string, balance int64) error
}

var _ BalanceUpdater = &repo.Service{}

type Service struct {
	bu   BalanceUpdater
	ug   repo.UserGetter
	aa   repo.AssetRepository
	freq time.Duration
	log  *logger.Logger
	c    *jsonrpc.Client
}

func NewService(c *jsonrpc.Client, bu BalanceUpdater, aa repo.AssetRepository, ug repo.UserGetter, freq time.Duration, log *logger.Logger) *Service {
	return &Service{
		bu:   bu,
		aa:   aa,
		freq: freq,
		log:  log,
		c:    c,
		ug:   ug,
	}
}

func (s *Service) GetAssetData(address string) (string, int, error) {
	e20 := erc20.NewERC20(ethgo.HexToAddress(address), s.c)
	name, err := e20.Name(ethgo.Latest)
	if err != nil {
		return "", 0, errors.Wrap(err, "name")
	}
	decimals, err := e20.Decimals(ethgo.Latest)
	if err != nil {
		return "", 0, errors.Wrap(err, "decimals")
	}
	return name, int(decimals), nil
}

func (s *Service) UpdateBalance(ctx context.Context, e *erc20.ERC20, aa ...ethgo.Address) {
	for _, a := range aa {
		time.Sleep(s.freq)
		b, err := e.BalanceOf(a, ethgo.Latest)
		if err != nil {
			if err.Error() != "empty response" {
				s.log.WithError(err).Error("get balance")
			}
			continue
		}
		if err = s.bu.UpdateBalance(ctx, e.Contract().Addr().String(), a.String(), b.Int64()); err != nil {
			s.log.WithError(err).Error("update balance")
		}
	}
}

func (s *Service) RunBalanceUpdater(ctx context.Context, t *time.Ticker) {
	log := s.log.WithField("actor", "balance-updater")
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			uu, err := s.ug.GetAllUsers(ctx)
			if err != nil {
				log.WithError(err).Error("get all users")
				continue
			}
			for _, u := range uu {
				if err = s.UpdateAllBalances(ctx, ethgo.HexToAddress(u.WalletAddress)); err != nil {
					log.WithError(err).Error("can't update balances")
				}
			}
		}
	}
}

func (s *Service) UpdateAllBalances(ctx context.Context, a ethgo.Address) error {
	aa, err := s.aa.GetAssets(ctx)
	if err != nil {
		return errors.Wrap(err, "get assets")
	}

	for _, as := range aa {
		e20 := erc20.NewERC20(ethgo.HexToAddress(as.Address), s.c)
		time.Sleep(s.freq)
		b, err := e20.BalanceOf(ethgo.HexToAddress(as.Address), ethgo.Latest)
		if err != nil {
			if err.Error() != "empty response" {
				s.log.WithError(err).Error("get balance")
			}
			continue
		}
		if err = s.bu.UpdateBalance(ctx, e20.Contract().Addr().String(), a.String(), b.Int64()); err != nil {
			s.log.WithError(err).Error("update balance")
		}
	}

	return nil
}
