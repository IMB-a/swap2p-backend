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

	c, err := jsonrpc.NewClient(cfg.Worker.JSONRPCClient)
	if err != nil {
		log.WithError(err).Fatal()
		return
	}

	ass := assets.NewService(c, r, r, r, time.Second*1, log)

	srv, err := server.NewServer(&cfg.Server, log,
		server.SetupWithRepo(r),
		server.SetupWithAsset(ass),
	)
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
			switch t.Type {
			case worker.TradeEventTypeCreate:
				err = r.AddTrade(context.Background(), &t.Trade)
				if err != nil {
					if !errors.Is(err, repo.TradeAlreadyExistsErr) {
						log.WithError(err).Error("can't add trade")
					}
				}
			case worker.TradeEventTypeAccept:
				err = r.CloseTrade(context.Background(), t.Id, t.Trade.Type, t.Trade.YAddress)
				if err != nil {
					log.WithError(err).Error("can't close trade")
				}
			}
		}
		fmt.Println("!!!!!!!!!!!CLOSED!!!!!!!!!!!!")
	}()

	go ass.RunBalanceUpdater(context.Background(), time.NewTicker(time.Minute))

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
