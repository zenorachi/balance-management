package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/balance-management/internal/entity"
	"net/http"
)

func (h *Handler) initAccountRoutes(api *gin.RouterGroup) {
	accounts := api.Group("/account", h.userIdentity)
	{
		accounts.POST("/create", h.createAccount)
		accounts.POST("/deposit", h.deposit)
		accounts.POST("/transfer", h.transfer)
		accounts.GET("/balance", h.getBalance)
	}
}

func (h *Handler) createAccount(c *gin.Context) {
	userId, isExists := c.Get(userCtx)
	if !isExists {
		newErrorResponse(c, http.StatusInternalServerError, "cannot find user id")
		return
	}

	id, err := h.services.Account.Create(c, userId.(int))
	if err != nil {
		if errors.Is(err, entity.ErrAccountAlreadyExists) {
			newErrorResponse(c, http.StatusConflict, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "cannot create account")
		}
		return
	}

	newResponse(c, http.StatusCreated, "id", id)
}

type depositInput struct {
	AccountID int     `json:"account_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
}

func (h *Handler) deposit(c *gin.Context) {
	var (
		input depositInput
		err   error
	)

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	err = h.services.Account.DepositByID(c, input.AccountID, input.Amount)
	if err != nil {
		if errors.Is(err, entity.ErrAccountDoesNotExist) || errors.Is(err, entity.ErrAmountIsNegative) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusOK, "message", "success")
}

type transferInput struct {
	SrcAccountID int     `json:"src_account_id" binding:"required"`
	DstAccountID int     `json:"dst_account_id" binding:"required"`
	Amount       float64 `json:"amount"         binding:"required"`
}

func (h *Handler) transfer(c *gin.Context) {
	var (
		input transferInput
		err   error
	)

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	err = h.services.Account.Transfer(c, input.SrcAccountID, input.DstAccountID, input.Amount)
	if err != nil {
		if errors.Is(err, entity.ErrAccountDoesNotExist) || errors.Is(err, entity.ErrNotEnoughMoney) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusOK, "message", "success")
}

type getBalanceInput struct {
	AccountID int `json:"account_id" binding:"required"`
}

func (h *Handler) getBalance(c *gin.Context) {
	var (
		input   getBalanceInput
		account entity.Account
		err     error
	)

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	account, err = h.services.Account.GetByID(c, input.AccountID)
	if err != nil {
		if errors.Is(err, entity.ErrAccountDoesNotExist) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusOK, "balance", account.Balance)
}
