package main

import (
	"auth-service/config"
	"auth-service/internal/app"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := new(config.Config)

	if err := env.Parse(cfg); err != nil {
		logrus.Fatal("failed to retrieve env variables: ", err)
	}

	if err := app.Run(cfg); err != nil {
		logrus.Fatal("error running service: ", err)
	}
}
