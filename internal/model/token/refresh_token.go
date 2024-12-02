package token

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/json"

	"auth-service/pkg/dbconverter"
)

type RefreshToken struct {
	value string
}

type refreshTokenClaims struct {
	UserIP string `json:"userIP"`
}

type RefreshTokenBuilder func(userIP string) (RefreshToken, error)

func NewRefreshTokenBuilder(
	publicKey *rsa.PublicKey,
) RefreshTokenBuilder {
	return func(userIP string) (RefreshToken, error) {
		payload, err := json.Marshal(refreshTokenClaims{
			UserIP: userIP,
		})
		if err != nil {
			return RefreshToken{}, err
		}

		token, err := rsa.EncryptOAEP(
			sha1.New(),
			rand.Reader,
			publicKey,
			payload,
			nil,
		)
		if err != nil {
			return RefreshToken{}, err
		}

		return RefreshToken{
			value: string(token),
		}, err
	}
}

type RefreshTokenParser func(t string) (string, error)

func NewRefreshTokenParser(
	privateKey *rsa.PrivateKey,
) RefreshTokenParser {
	return func(token string) (string, error) {
		var claims refreshTokenClaims

		payload, err := privateKey.Decrypt(
			nil,
			[]byte(token),
			&rsa.OAEPOptions{Hash: crypto.SHA1},
		)
		if err != nil {
			return claims.UserIP, err
		}

		if err := json.Unmarshal(payload, &claims); err != nil {
			return claims.UserIP, err
		}

		return claims.UserIP, nil
	}
}

func (rt RefreshToken) Value() string {
	return rt.value
}

func (rt *RefreshToken) Scan(srv any) error {
	return dbconverter.ConvertToString(&rt.value, srv)
}
