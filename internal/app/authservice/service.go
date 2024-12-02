package authservice

import (
	"auth-service/config"
	token "auth-service/internal/model/token"
	"auth-service/internal/repository/tokenrepository"
	"auth-service/internal/repository/userrepository"
	"auth-service/pkg/emailsender"
)

type AuthService struct {
	accessTokenBuilder  token.AccessTokenBuilder
	refreshTokenBuilder token.RefreshTokenBuilder
	refreshTokenParser  token.RefreshTokenParser
	expRefrashToken     uint32
	userRepository      *userrepository.UserRepository
	tokenRepository     *tokenrepository.TokenRepository

	emailSender emailsender.EmailSender
}

func NewAuthService(
	cfg *config.Config,
	userRepository *userrepository.UserRepository,
	tokenRepository *tokenrepository.TokenRepository,
) *AuthService {
	accessTokenBuilder := token.NewAccessTokenBuilder(
		cfg.ExpAccessToken,
		[]byte(cfg.SecretKey),
	)

	refreshTokenBuilder := token.NewRefreshTokenBuilder(
		&cfg.PrivateKey.PublicKey,
	)

	refreshTokenParser := token.NewRefreshTokenParser(
		cfg.PrivateKey,
	)

	emailSender := emailsender.NewEmailSender(cfg)

	return &AuthService{
		accessTokenBuilder:  accessTokenBuilder,
		refreshTokenBuilder: refreshTokenBuilder,
		refreshTokenParser:  refreshTokenParser,
		expRefrashToken:     uint32(cfg.ExpRefreshToken),
		userRepository:      userRepository,
		tokenRepository:     tokenRepository,

		emailSender: emailSender,
	}
}
