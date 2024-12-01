package handlers

import (
	"github.com/gin-gonic/gin"

	"auth-service/internal/app/authservice"
)

type Handler struct {
	service *authservice.AuthService
}

func NewHandler(service *authservice.AuthService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		api.GET("/receive/:id", h.Receive)
		api.GET("/refresh", h.Refresh)
	}

	return router
}
