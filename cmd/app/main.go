package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/zenorachi/balance-management/internal/app"
	"github.com/zenorachi/balance-management/internal/config"
	"github.com/zenorachi/balance-management/pkg/logger"
)

const (
	envFile   = ".env"
	directory = "configs"
	ymlFile   = "main"
)

func main() {
	if err := godotenv.Load(envFile); err != nil {
		logger.Fatal("config", ".env initialization failed")
	}

	viper.AddConfigPath(directory)
	viper.SetConfigName(ymlFile)

	app.Run(config.New())
}
