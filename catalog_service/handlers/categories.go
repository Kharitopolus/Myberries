package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (h HandlerImpl) ListCategories(c *gin.Context) {
	var categories []Category
	h.db.Order("created_at").Find(&categories)
	c.JSON(http.StatusOK, categories)
}

type CreateCategoryReqBody struct {
	Name string `json:"name" binding:"required"`
}

func (h HandlerImpl) CreateCategory(c *gin.Context) {
	var b CreateCategoryReqBody

	err := c.ShouldBindJSON(&b)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := Category{ID: uuid.New(), Name: b.Name}
	h.db.Create(&category)

	c.JSON(http.StatusCreated, category)
}

func (h HandlerImpl) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := h.db.Delete(&Category{}, id)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "category does not exist"})
		return
	}
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h HandlerImpl) GetCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category Category
	res := h.db.First(&category, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "category does not exist"})
		return
	}

	c.JSON(200, category)
}

type UpdateCategoryReqBody struct {
	Name string `json:"name" binding:"required"`
}

func (h HandlerImpl) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category Category
	res := h.db.First(&category, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "category does not exist"})
		return
	}
	var b UpdateCategoryReqBody

	err = c.ShouldBindJSON(&b)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.Name = b.Name

	res = h.db.Save(&category)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "category does not exist"})
		return
	}

	c.Status(http.StatusOK)
}
