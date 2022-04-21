package main

import (
	"flag"
	"os"

	"github.com/Pod-Box/swap2p-backend/repo"
	"github.com/Pod-Box/swap2p-backend/server"
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

	srv.Run()
}

type Config struct {
	DB     repo.Config   `yaml:"db"`
	Server server.Config `yaml:"server"`
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
