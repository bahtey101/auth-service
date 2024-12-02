package main

import (
	"auth-service/config"
	"auth-service/internal/app"
	"crypto/rand"
	"crypto/rsa"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := new(config.Config)

	if err := env.Parse(cfg); err != nil {
		logrus.Fatal("failed to retrieve env variables: ", err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		logrus.Fatal("failed to generate private key")
	}

	cfg.PrivateKey = privateKey

	if err := app.Run(cfg); err != nil {
		logrus.Fatal("error running service: ", err)
	}
}
