package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/balance-management/internal/entity"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) initProductRoutes(api *gin.RouterGroup) {
	products := api.Group("/products", h.userIdentity)
	{
		products.POST("/create", h.createProduct)
		products.GET("/:product_id", h.getProduct)
		products.GET("/", h.getAllProducts)
	}
}

type createProductInput struct {
	Name  string  `json:"name"  binding:"required,min=2"`
	Price float64 `json:"price" binding:"required"`
}

func (h *Handler) createProduct(c *gin.Context) {
	var (
		input createProductInput
		id    int
		err   error
	)

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	product := entity.Product{Name: input.Name, Price: input.Price}

	id, err = h.services.Product.Create(c, product)
	if err != nil {
		if errors.Is(err, entity.ErrProductAlreadyExists) {
			newErrorResponse(c, http.StatusConflict, err.Error())
		} else if errors.Is(err, entity.ErrPriceIsNegative) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusCreated, "id", id)
}

func (h *Handler) getProduct(c *gin.Context) {
	var (
		id      int
		product entity.Product
		err     error
	)

	paramId := strings.Trim(c.Param("product_id"), "/")
	id, err = strconv.Atoi(paramId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input (id)")
		return
	}

	product, err = h.services.Product.GetByID(c, id)
	if err != nil {
		if errors.Is(err, entity.ErrProductDoesNotExist) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusOK, "product", product)
}

func (h *Handler) getAllProducts(c *gin.Context) {
	products, err := h.services.Product.GetAll(c)
	if err != nil {
		if errors.Is(err, entity.ErrProductDoesNotExist) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusOK, "products", products)
}
