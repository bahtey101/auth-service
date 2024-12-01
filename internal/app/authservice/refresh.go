package authservice

import (
	"context"
	"fmt"
	"time"

	"auth-service/internal/model/token"
	"auth-service/tools"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type SendActivateKeyCommand struct {
	UserEmail string
	StrangeIP string
}

func (s *AuthService) Refresh(
	ctx context.Context,
	userIP string,
	accessToken string,
	refreshToken string,
) (token.Token, token.Token, error) {
	// Парсим данные из токенов
	accessTokenClaims, err := s.tokenParser(accessToken)
	if err != nil {
		return token.Token{}, token.Token{}, err
	}

	refreshTokenClaims, err := s.tokenParser(refreshToken)
	if err != nil || !time.Unix(refreshTokenClaims.ExpiresAt, 0).After(time.Now()) {
		return token.Token{}, token.Token{}, err
	}

	// Сравниваем SessionID двух токенов
	if accessTokenClaims.SessionID != refreshTokenClaims.SessionID {
		return token.Token{}, token.Token{}, fmt.Errorf("different session id")
	}

	// Получаем RefreshToken из БД
	hashedResfreshToken, err := s.tokenRepository.Get(
		ctx,
		refreshTokenClaims.UserID,
		refreshTokenClaims.UserIP,
	)
	if err != nil {
		return token.Token{}, token.Token{}, err
	}

	// Сравниваем полученный токен с нашим
	if err = bcrypt.CompareHashAndPassword(
		[]byte(hashedResfreshToken),
		[]byte(refreshToken),
	); err != nil {
		return token.Token{}, token.Token{}, err
	}

	// Создаём пару новых токенов
	newSessionID, err := tools.GenerateSessionID()
	if err != nil {
		return token.Token{}, token.Token{}, err
	}

	newAccessToken, err := s.accessTokenBuilder(
		refreshTokenClaims.UserID,
		userIP,
		newSessionID,
	)
	if err != nil {
		return token.Token{}, token.Token{}, err
	}

	newRefreshToken, err := s.refreshTokenBuilder(
		refreshTokenClaims.UserID,
		userIP,
		newSessionID,
	)
	if err != nil {
		return token.Token{}, token.Token{}, err
	}

	// Проверяем IP адрес
	if userIP == refreshTokenClaims.UserIP {
		if err = s.tokenRepository.Update(
			ctx,
			refreshTokenClaims.UserID,
			userIP,
			newRefreshToken.Value(),
		); err != nil {
			return token.Token{}, token.Token{}, err
		}
	} else {
		// Отправляем email warning
		userEmail, err := s.userRepository.GetByID(
			ctx,
			refreshTokenClaims.UserID,
		)
		if err != nil {
			logrus.Error("failed to find user email in users")
		}

		err = s.emailSender(
			userEmail,
			fmt.Sprintf("[WARNING] Login from different IP\n Someone logged into your account from IP: %s", userIP),
		)
		if err != nil {
			logrus.Errorf("failed to send email to %s", userEmail)
		}

		if err = s.tokenRepository.Create(
			ctx,
			refreshTokenClaims.UserID,
			userIP,
			newRefreshToken.Value(),
		); err != nil {
			return token.Token{}, token.Token{}, err
		}
	}

	return newAccessToken, newRefreshToken, nil
}
