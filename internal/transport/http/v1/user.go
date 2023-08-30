package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zenorachi/balance-management/internal/entity"
)

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.signUp)
		users.POST("/sign-in", h.signIn)
		users.GET("/refresh", h.refresh)
	}
}

type signUpInput struct {
	Login    string `json:"login"    binding:"required,min=2,max=64"`
	Email    string `json:"email"    binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input signUpInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	id, err := h.services.User.SignUp(c, input.Login, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, entity.ErrUserAlreadyExists) {
			newErrorResponse(c, http.StatusConflict, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusOK, "id", id)
}

type signInInput struct {
	Login    string `json:"login"    binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	tokens, err := h.services.User.SignIn(c, input.Login, input.Password)
	if err != nil {
		if errors.Is(err, entity.ErrUserDoesNotExist) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", tokens.RefreshToken))
	newResponse(c, http.StatusOK, "token", tokens.AccessToken)
}

func (h *Handler) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh-token")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "refresh-token not found")
		return
	}

	tokens, err := h.services.RefreshTokens(c, refreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", tokens.RefreshToken))
	newResponse(c, http.StatusOK, "token", tokens.AccessToken)
}
