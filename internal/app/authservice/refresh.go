package authservice

import (
	"context"
	"fmt"
	"time"

	"auth-service/internal/model/token"

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
) (token.Token, error) {
	rtUserIP, err := s.refreshTokenParser(refreshToken)
	if err != nil {
		return token.Token{}, err
	}

	userID, hashedResfreshToken, createdAt, err := s.tokenRepository.GetByIP(
		ctx,
		rtUserIP,
	)
	if err != nil {
		return token.Token{}, err
	}

	if time.Since(createdAt).Minutes() >= float64(s.expRefrashToken) {
		err := s.tokenRepository.Delete(
			ctx,
			rtUserIP,
		)
		if err != nil {
			return token.Token{}, err
		}

		return token.Token{}, fmt.Errorf("refresh token is expired")
	}

	if err = bcrypt.CompareHashAndPassword(
		[]byte(hashedResfreshToken),
		[]byte(refreshToken),
	); err != nil {
		return token.Token{}, err
	}

	newAccessToken, err := s.accessTokenBuilder(userIP)
	if err != nil {
		return token.Token{}, err
	}

	newRefreshToken, err := s.refreshTokenBuilder(userIP)
	if err != nil {
		return token.Token{}, err
	}

	if userIP == rtUserIP {
		if err = s.tokenRepository.Update(
			ctx,
			userID,
			userIP,
			newRefreshToken.Value(),
		); err != nil {
			return token.Token{}, err
		}
	} else {
		// userEmail, err := s.userRepository.GetByID(
		// 	ctx,
		// 	rtUserID,
		// )
		// if err != nil {
		// 	logrus.Error("failed to find user email in users")
		// }

		// err = s.emailSender(
		// 	userEmail,
		// 	fmt.Sprintf("[WARNING] Login from different IP\n Someone logged into your account from IP: %s", userIP),
		// )
		// if err != nil {
		// 	logrus.Errorf("failed to send email to %s", userEmail)
		// }

		if err = s.tokenRepository.Create(
			ctx,
			userID,
			userIP,
			newRefreshToken.Value(),
		); err != nil {
			logrus.Error("failed to create in db")
			return token.Token{}, err
		}
	}

	return token.Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, err
}
