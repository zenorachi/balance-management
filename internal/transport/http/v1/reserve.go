package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/balance-management/internal/entity"
	"net/http"
)

func (h *Handler) initReserveRoutes(api *gin.RouterGroup) {
	reserves := api.Group("/reserve", h.userIdentity)
	{
		reserves.POST("/create", h.createReserve)
		reserves.POST("/confirm_revenue", h.confirmRevenue)
		reserves.POST("/confirm_refund", h.confirmRefund)
	}
}

type createReserveInput struct {
	OrderID int `json:"order_id" binding:"required"`
}

func (h *Handler) createReserve(c *gin.Context) {
	var (
		input createReserveInput
		id    int
		err   error
	)

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
	}

	id, err = h.services.Reserve.Create(c, entity.Reserve{OrderID: input.OrderID})
	if err != nil {
		if errors.Is(err, entity.ErrOrderDoesNotExist) || errors.Is(err, entity.ErrNotEnoughMoney) ||
			errors.Is(err, entity.ErrOrderCannotBeProcessed) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusOK, "id", id)
}

type confirmOperationInput struct {
	ReserveID int `json:"reserve_id" binding:"required"`
}

func (h *Handler) confirmRevenue(c *gin.Context) {
	var (
		input       confirmOperationInput
		operationId int
		err         error
	)

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	operationId, err = h.services.Reserve.ConfirmRevenueByID(c, input.ReserveID)
	if err != nil {
		if errors.Is(err, entity.ErrReserveDoesNotExist) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusOK, "operation_id", operationId)
}

func (h *Handler) confirmRefund(c *gin.Context) {
	var (
		input       confirmOperationInput
		operationId int
		err         error
	)

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	operationId, err = h.services.Reserve.ConfirmRefundByID(c, input.ReserveID)
	if err != nil {
		if errors.Is(err, entity.ErrReserveDoesNotExist) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusOK, "operation_id", operationId)
}
