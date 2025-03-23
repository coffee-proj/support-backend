package app

import (
	"os"

	"github.com/coffee/support/internal/controller"
	"github.com/coffee/support/internal/usecase"
	"github.com/gosuit/httper"
	"github.com/gosuit/mongo"
	"github.com/gosuit/sl"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type appConfig struct {
	Mongo      mongo.Config      `yaml:"mongo"`
	Logger     sl.Config         `yaml:"logger"`
	Controller controller.Config `yaml:"controller"`
	Server     httper.ServerCfg  `yaml:"http_server"`
	Usecase    usecase.Config    `yaml:"usecase"`
}

func getAppConfig() (*appConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	path := getConfigPath()

	var cfg appConfig

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func getConfigPath() string {
	path := os.Getenv("CONFIG_PATH")

	if path == "" {
		return "config/local.yaml"
	}

	return path
}
