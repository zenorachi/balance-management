package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/balance-management/internal/entity"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) initOperationRoutes(api *gin.RouterGroup) {
	orders := api.Group("/operations", h.userIdentity)
	{
		orders.GET("/", h.getReportForAccounting)
		orders.GET("/:account_id", h.getReportForUser)
	}
}

func (h *Handler) getReportForAccounting(c *gin.Context) {
	var (
		operations []entity.Operation
		err        error
	)

	operations, err = h.services.Operation.GetReportForAccounting(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(operations) == 0 {
		newResponse(c, http.StatusNoContent, "message", "no operations")
		return
	}

	newResponse(c, http.StatusOK, "operations", operations)
}

func (h *Handler) getReportForUser(c *gin.Context) {
	var (
		accountId  int
		operations []entity.Operation
		err        error
	)

	paramId := strings.Trim(c.Param("account_id"), "/")
	accountId, err = strconv.Atoi(paramId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input (id)")
		return
	}

	operations, err = h.services.Operation.GetReportForUser(c, accountId)
	if err != nil {
		if errors.Is(err, entity.ErrAccountDoesNotExist) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if len(operations) == 0 {
		newResponse(c, http.StatusNoContent, "message", "no operations")
		return
	}

	newResponse(c, http.StatusOK, "operations", operations)
}
