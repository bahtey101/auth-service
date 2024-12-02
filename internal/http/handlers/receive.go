package handlers

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (h *Handler) Receive(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		logrus.Errorf("failed to parse id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	userIP := ctx.ClientIP()

	logrus.Info("Receice")

	token, err := h.service.Receive(
		ctx.Request.Context(),
		userID,
		userIP,
	)
	if err != nil {
		logrus.Errorf("failed to generate token: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate tokens",
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
