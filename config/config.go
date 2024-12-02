package config

import "crypto/rsa"

type Config struct {
	Port            string `env:"PORT" envDefault:"8080"`
	PgDSN           string `env:"AUTH_SERVICE_PG_DSN" envDefault:"postgres://db:db@localhost:23432/db"`
	PgMaxOpenConn   int    `env:"PG_MAX_OPEN_CONN" envDefault:"5"`
	SecretKey       string `env:"AUTH_SERVICE_SECRET_KEY" envDefault:"super-secret-key"`
	PrivateKey      *rsa.PrivateKey
	ExpRefreshToken int `env:"AUTH_SERVICE_EXP_REFRESH_TOKEN" envDefault:"5"`
	ExpAccessToken  int `env:"AUTH_SERVICE_EXP_ACCESS_TOKEN" envDefault:"2"`

	EmailAdress   string `env:"EMAIL_ADRESS" envDefault:"example@gmail.com"`
	EmailPassword string `env:"EMAIL_PASSWORD" envDefault:"password"`
	SMTPHost      string `env:"SMTP_HOST" envDefault:"smtp.gmail.com"`
	SMTPPort      string `env:"SMTP_PORT" envDefault:"587"`
}
