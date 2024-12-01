package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) Receive(ctx *gin.Context) {
	userID, err := uuid.FromBytes([]byte(ctx.Param("id")))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	userIP := ctx.ClientIP()

	accessToken, refreshToken, err := h.service.Receive(
		ctx.Request.Context(),
		userID,
		userIP,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate tokens",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken.Value(),
		"refresh_token": refreshToken.Value(),
	})
}
