package token

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type accessTokenClaims struct {
	UserIP string

	jwt.StandardClaims
}

type AccessToken struct {
	value string
}

type AccessTokenBuilder func(userIP string) (AccessToken, error)

func NewAccessTokenBuilder(
	expMinute int,
	tokenKey []byte,
) AccessTokenBuilder {
	return func(userIP string) (AccessToken, error) {
		expirationTime := time.Now().Add(time.Duration(expMinute) * time.Minute)

		claims := accessTokenClaims{
			UserIP: userIP,

			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

		tokenString, err := token.SignedString(tokenKey)
		if err != nil {
			return AccessToken{}, err
		}

		return AccessToken{
			value: tokenString,
		}, nil
	}
}

func (at *AccessToken) Value() string {
	return at.value
}
