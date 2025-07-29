package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yantology/golang_template/internal/api/handlers"
)

func SetupRoutes(r *gin.RouterGroup, h *handlers.Handler) {
	r.GET("/ping", h.HealthCheck)
}