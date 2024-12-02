package handlers

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) Refresh(ctx *gin.Context) {
	var request struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		logrus.Errorf("failed to bind json: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid refresh request",
		})
		return
	}

	userIP := ctx.ClientIP()

	decodedRefreshToken, err := base64.StdEncoding.DecodeString(
		request.RefreshToken,
	)
	if err != nil {
		logrus.Errorf("invalid refresh token: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	token, err := h.service.Refresh(
		ctx.Request.Context(),
		userIP,
		request.AccessToken,
		string(decodedRefreshToken),
	)
	if err != nil {
		logrus.Errorf("failed to refresh tokens: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to refresh tokens",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token.AccessToken.Value(),
		"refresh_token": base64.StdEncoding.EncodeToString(
			[]byte(token.RefreshToken.Value()),
		),
	})
}
