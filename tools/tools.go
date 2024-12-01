package tools

import (
	"crypto/rand"
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const lenSessionID = 32

func GenerateSessionID() (string, error) {
	b := make([]byte, lenSessionID)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func Hash(token string) string {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(token),
		bcrypt.DefaultCost,
	)
	if err != nil {
		logrus.Errorf("refresh token is too long: %d", len(token))
	}

	return string(hash)
}
