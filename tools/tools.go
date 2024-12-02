package tools

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func Hash(rt string) string {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(rt),
		bcrypt.DefaultCost,
	)
	if err != nil {
		logrus.Errorf("refresh token is too long: %d", len(rt))
	}

	return string(hash)
}
