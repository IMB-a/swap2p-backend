package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Pod-Box/swap2p-backend/repo"
	"github.com/Pod-Box/swap2p-backend/server"
	"github.com/Pod-Box/swap2p-backend/worker"
	"github.com/Pod-Box/swap2p-backend/worker/assets"
	logger "github.com/sirupsen/logrus"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/contract/builtin/erc20"
	"github.com/umbracle/ethgo/jsonrpc"
	"gopkg.in/yaml.v3"
)

var configPath string

func main() {
	log := logger.New()
	flag.StringVar(&configPath, "c", "config-local.yaml", "Set path to config file.")
	flag.Parse()

	cfg, err := ReadConfig(configPath)
	if err != nil {
		log.WithError(err).WithField("config_file_path", configPath).Fatal("can't configure from config file")
	}

	r, err := repo.NewService(&cfg.DB)
	if err != nil {
		log.WithError(err).Fatal("can't setup repo")
	}

	srv, err := server.NewServer(&cfg.Server, log, server.SetupWithRepo(r))
	if err != nil {
		log.WithError(err).Fatal("can't setup server")
	}

	wrk, err := worker.NewService(&cfg.Worker)
	if err != nil {
		log.WithError(err).Fatal("can't setup worker")
	}

	go wrk.Run(context.Background())

	go func() {
		for t := range wrk.TradeChan {
			log.Println(t)
			switch t.Type {
			case worker.TradeEventTypeCreate:
				err = r.AddTrade(context.Background(), &t.Trade)
				if err != nil {
					if !errors.Is(err, repo.TradeAlreadyExistsErr) {
						log.WithError(err).Error("can't add trade")
					}
				}
			case worker.TradeEventTypeAccept:
				err = r.CloseTrade(context.Background(), t.Id, t.Trade.YAddress)
				if err != nil {
					log.WithError(err).Error("can't close trade")
				}
			}
		}
		fmt.Println("!!!!!!!!!!!CLOSED!!!!!!!!!!!!")
	}()

	go func() {
		ctx := context.Background()
		ass := assets.NewService(r, time.Second*1, log)

		c, err := jsonrpc.NewClient(cfg.Worker.JSONRPCClient)
		if err != nil {
			log.WithError(err).Error()
			return
		}

		for {
			time.Sleep(time.Second * 100)
			aa, err := r.GetAssets(ctx)
			if err != nil {
				log.WithError(err).Error("get assets")
				continue
			}
			uu, err := r.GetAllUsers(ctx)
			if err != nil {
				log.WithError(err).Error("get all users")
				continue
			}
			for _, a := range aa {
				e20 := erc20.NewERC20(ethgo.HexToAddress(a.Address), c)
				time.Sleep(time.Second * 1)
				for _, u := range uu {
					ass.UpdateBalance(ctx, e20, ethgo.HexToAddress(u.WalletAddress))
				}
			}
		}
	}()

	srv.Run()
}

type Config struct {
	DB     repo.Config   `yaml:"db"`
	Server server.Config `yaml:"server"`
	Worker worker.Config `yaml:"worker"`
}

func ReadConfig(fileName string) (Config, error) {
	var cnf Config
	data, err := os.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(data, &cnf)
	if err != nil {
		return Config{}, err
	}
	return cnf, nil
}
