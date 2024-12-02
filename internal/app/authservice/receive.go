package authservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"auth-service/internal/model/token"
)

func (s *AuthService) Receive(
	ctx context.Context,
	userID uuid.UUID,
	userIP string,
) (token.Token, error) {
	if _, err := s.userRepository.GetByID(ctx, userID); err != nil {
		return token.Token{}, err
	}

	accessToken, err := s.accessTokenBuilder(userIP)
	if err != nil {
		logrus.Error("build access token err")
		return token.Token{}, err
	}

	refreshToken, err := s.refreshTokenBuilder(userIP)
	if err != nil {
		logrus.Error("build refresh token err")
		return token.Token{}, err
	}

	if err := s.tokenRepository.Create(
		ctx,
		userID,
		userIP,
		refreshToken.Value(),
	); err != nil {

		return token.Token{}, err
	}

	return token.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
