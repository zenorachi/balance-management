package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/balance-management/internal/service"
	"github.com/zenorachi/balance-management/pkg/auth"
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

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUserRoutes(v1)
		h.initAccountRoutes(v1)
		h.initProductRoutes(v1)
		h.initOrderRoutes(v1)
		h.initReserveRoutes(v1)
		h.initOperationRoutes(v1)
	}
}
