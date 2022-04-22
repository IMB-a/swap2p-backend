package assets

import (
	"context"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/contract/builtin/erc20"
)

type BalanceUpdater interface {
	UpdateBalance(ctx context.Context, assetAddress, walletAddress string, balance int64) error
}

type Service struct {
	bu   BalanceUpdater
	freq time.Duration
	log  *logger.Logger
}

func NewService(bu BalanceUpdater, freq time.Duration, log *logger.Logger) *Service {
	return &Service{
		bu:   bu,
		freq: freq,
		log:  log,
	}
}
func (s *Service) UpdateBalance(ctx context.Context, e *erc20.ERC20, aa ...ethgo.Address) {
	for _, a := range aa {
		time.Sleep(s.freq)
		b, err := e.BalanceOf(a, ethgo.Latest)
		if err != nil {
			s.log.WithError(err).Error("get balance")
			continue
		}
		if err = s.bu.UpdateBalance(ctx, e.Contract().Addr().String(), a.String(), b.Int64()); err != nil {
			s.log.WithError(err).Error("update balance")
		}
	}
}
