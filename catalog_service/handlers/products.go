package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (h HandlerImpl) ListProducts(c *gin.Context) {
	var products []Product
	h.db.Order("created_at").Find(&products)
	c.JSON(http.StatusOK, products)
}

type CreateProductReqBody struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PriceUSD    float64   `json:"price_usd"`
	CategoryID  uuid.UUID `json:"category_id"`
}

func (h HandlerImpl) CreateProduct(c *gin.Context) {
	var b CreateProductReqBody

	err := c.ShouldBindJSON(&b)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := Product{
		ID:          uuid.New(),
		Name:        b.Name,
		Description: b.Description,
		PriceUSD:    b.PriceUSD,
		CategoryID:  b.CategoryID,
	}
	res := h.db.Create(&product)
	if res.Error != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": res.Error.Error()},
		)
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h HandlerImpl) DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := h.db.Delete(&Product{}, id)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "product does not exist"})
		return
	}
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h HandlerImpl) GetProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product Product
	res := h.db.First(&product, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "product does not exist"})
		return
	}

	c.JSON(200, product)
}

type UpdateProductReqBody struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PriceUSD    float64   `json:"price_usd"`
	CategoryID  uuid.UUID `json:"category_id"`
}

func (h HandlerImpl) UpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product Product
	res := h.db.First(&product, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "product does not exist"})
		return
	}
	var b UpdateProductReqBody

	err = c.ShouldBindJSON(&b)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.Name = b.Name

	res = h.db.Save(&product)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "product does not exist"})
		return
	}

	c.Status(http.StatusOK)
}
