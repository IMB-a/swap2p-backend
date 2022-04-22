package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/Pod-Box/swap2p-backend/repo"
	"github.com/Pod-Box/swap2p-backend/server"
	"github.com/Pod-Box/swap2p-backend/worker"
	logger "github.com/sirupsen/logrus"
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
					log.WithError(err).Error("can't add trade")
				}
			case worker.TradeEventTypeAccept:
				err = r.CloseTrade(context.Background(), t.Id)
				if err != nil {
					log.WithError(err).Error("can't close trade")
				}
			}
			fmt.Println(t)
		}
		fmt.Println("!!!!!!!!!!!CLOSED!!!!!!!!!!!!")
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
