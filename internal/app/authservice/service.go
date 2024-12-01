package authservice

import (
	"auth-service/config"
	token "auth-service/internal/model/token"
	"auth-service/internal/repository/tokenrepository"
	"auth-service/internal/repository/userrepository"
	"auth-service/pkg/emailsender"

	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	accessTokenBuilder  token.TokenBuilder
	refreshTokenBuilder token.TokenBuilder
	tokenParser         token.TokenParser
	userRepository      *userrepository.UserRepository
	tokenRepository     *tokenrepository.TokenRepository

	emailSender emailsender.EmailSender
}

func NewAuthService(
	cfg *config.Config,
	userRepository *userrepository.UserRepository,
	tokenRepository *tokenrepository.TokenRepository,
) *AuthService {
	accessTokenBuilder := token.NewTokenBuilder(
		cfg.ExpAccessToken,
		[]byte(cfg.SecretKey),
		jwt.SigningMethodHS512,
	)

	refreshTokenBuilder := token.NewTokenBuilder(
		cfg.ExpRefreshToken,
		[]byte(cfg.SecretKey),
		jwt.SigningMethodHS256,
	)

	tokenParser := token.NewtokenParser(
		[]byte(cfg.SecretKey),
	)

	emailSender := emailsender.NewEmailSender(cfg)

	return &AuthService{
		accessTokenBuilder:  accessTokenBuilder,
		refreshTokenBuilder: refreshTokenBuilder,
		tokenParser:         tokenParser,
		userRepository:      userRepository,
		tokenRepository:     tokenRepository,

		emailSender: emailSender,
	}
}
