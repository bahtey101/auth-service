package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"auth-service/pkg/dbconverter"
)

type jwtClaims struct {
	UserID uuid.UUID
	UserIP string

	jwt.StandardClaims
}

type Token struct {
	value string
}

type TokenBuilder func(userID uuid.UUID, userIP, sessionID string) (Token, error)

func NewTokenBuilder(
	expMinute int,
	tokenKey []byte,
	signingMethod jwt.SigningMethod,
) TokenBuilder {
	return func(userID uuid.UUID, userIP, sessionID string) (Token, error) {
		expirationTime := time.Now().Add(time.Duration(expMinute) * time.Minute)

		claims := jwtClaims{
			UserID: userID,
			UserIP: userIP,

			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
				Id:        sessionID,
			},
		}

		token := jwt.NewWithClaims(signingMethod, claims)

		tokenString, err := token.SignedString(tokenKey)
		if err != nil {
			return Token{}, err
		}

		return Token{
			value: tokenString,
		}, nil
	}
}

type ParserResponse struct {
	UserID    uuid.UUID
	UserIP    string
	SessionID string
	ExpiresAt int64
}

type TokenParser func(token string) (ParserResponse, error)

func NewtokenParser(
	tokenKey []byte,
) TokenParser {
	return func(token string) (ParserResponse, error) {
		// var (
		// 	userIP    string
		// 	sessionID string
		// 	expiresAt int64
		// )

		parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signature method %v", t.Header["alg"])
			}
			return tokenKey, nil
		})
		if err != nil {
			// return userIP, sessionID, expiresAt, err
			return ParserResponse{}, err
		}

		claims, ok := parsedToken.Claims.(jwtClaims)
		if !ok || !parsedToken.Valid {
			//return userIP, sessionID, expiresAt, fmt.Errorf("invalid token")
			return ParserResponse{}, fmt.Errorf("invalid token")
		}

		// return claims.UserIP, claims.Id, claims.ExpiresAt, nil
		return ParserResponse{
			UserID:    claims.UserID,
			UserIP:    claims.UserIP,
			SessionID: claims.Id,
			ExpiresAt: claims.ExpiresAt,
		}, nil
	}
}

func (t *Token) Value() string {
	return t.value
}

func (t *Token) Scan(srv any) error {
	return dbconverter.ConvertToString(&t.value, srv)
}
