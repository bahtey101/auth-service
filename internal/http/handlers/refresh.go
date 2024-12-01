package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Refresh(ctx *gin.Context) {
	var request struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid refresh request",
		})
		return
	}

	userIP := ctx.ClientIP()

	accessToken, refreshToken, err := h.service.Refresh(
		ctx.Request.Context(),
		userIP,
		request.AccessToken,
		request.RefreshToken,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to refresh tokens",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken.Value(),
		"refresh_token": refreshToken.Value(),
	})
}
