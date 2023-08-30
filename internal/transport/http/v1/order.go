package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/balance-management/internal/entity"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) initOrderRoutes(api *gin.RouterGroup) {
	orders := api.Group("/orders", h.userIdentity)
	{
		orders.POST("/", h.makeOrder)
		orders.POST("/cancel", h.cancelOrder)
		orders.GET("/:account_id", h.getOrdersByAccountID)
	}
}

type makeOrderInput struct {
	AccountID int     `json:"account_id" binding:"required"`
	Products  []uint8 `json:"products" binding:"required"`
}

func (h *Handler) makeOrder(c *gin.Context) {
	var (
		input makeOrderInput
		id    int
		err   error
	)

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	order := entity.Order{AccountID: input.AccountID, Products: input.Products}

	id, err = h.services.Order.Create(c, order)
	if err != nil {
		if errors.Is(err, entity.ErrAccountDoesNotExist) || errors.Is(err, entity.ErrProductDoesNotExist) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusCreated, "id", id)
}

type cancelOrderInput struct {
	OrderID int `json:"order_id" binding:"required"`
}

func (h *Handler) cancelOrder(c *gin.Context) {
	var input cancelOrderInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	err := h.services.Order.CancelByID(c, input.OrderID)
	if err != nil {
		if errors.Is(err, entity.ErrOrderDoesNotExist) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusOK, "status", "cancelled")
}

func (h *Handler) getOrdersByAccountID(c *gin.Context) {
	var (
		accountId int
		orders    []entity.Order
		err       error
	)

	paramId := strings.Trim(c.Param("account_id"), "/")
	accountId, err = strconv.Atoi(paramId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input (id)")
		return
	}

	orders, err = h.services.Order.GetAllByAccountID(c, accountId)
	if err != nil {
		if errors.Is(err, entity.ErrAccountDoesNotExist) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusOK, "orders", orders)
}
