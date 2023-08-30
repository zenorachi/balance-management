package transport

import (
	"net/http"

	"github.com/zenorachi/balance-management/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/zenorachi/balance-management/internal/service"
	v1 "github.com/zenorachi/balance-management/internal/transport/http/v1"
)

const (
	ginMode = gin.ReleaseMode
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(ginMode)

	router := gin.New()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(r *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager)
	api := r.Group("/api")
	{
		handlerV1.Init(api)
	}
}
