package authservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"auth-service/internal/model/token"
	"auth-service/tools"
)

func (s *AuthService) Receive(
	ctx context.Context,
	userID uuid.UUID,
	userIP string,
) (token.Token, token.Token, error) {
	// Логика проверки пользователя
	if _, err := s.userRepository.GetByID(ctx, userID); err != nil {
		logrus.Errorf("failed to find user %v", err)
		return token.Token{}, token.Token{}, err
	}

	// Создаём пару новых токенов
	sessionID, err := tools.GenerateSessionID()
	if err != nil {
		return token.Token{}, token.Token{}, err
	}

	accessToken, err := s.accessTokenBuilder(userID, userIP, sessionID)
	if err != nil {
		return token.Token{}, token.Token{}, err
	}

	refreshToken, err := s.refreshTokenBuilder(userID, userIP, sessionID)
	if err != nil {
		return token.Token{}, token.Token{}, err
	}

	// Сохраняем токен
	if err := s.tokenRepository.Create(
		ctx,
		userID,
		userIP,
		refreshToken.Value(),
	); err != nil {
		return token.Token{}, token.Token{}, err
	}
	logrus.Printf("rtolen len: %d", len(refreshToken.Value()))

	return accessToken, refreshToken, nil
}
