package repo

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type Config struct {
	Driver     string `yaml:"driver"`
	Addr       string `yaml:"addr"`
	Port       string `yaml:"port"`
	DB         string `yaml:"db"`
	UserEnvKey string `yaml:"user_env_key"`
	PassEnvKey string `yaml:"pass_env_key"`
}

func connectionString(c *Config) (string, error) {
	username := os.Getenv(c.UserEnvKey)
	password := os.Getenv(c.PassEnvKey)
	if username == "" || password == "" {
		return "", errors.New("can't get db credentials from env")
	}

	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.Addr, c.Port, c.DB, username, password), nil
}
