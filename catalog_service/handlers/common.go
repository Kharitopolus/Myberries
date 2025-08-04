package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type HandlerImpl struct {
	db *gorm.DB
}

func InitDB(dbUrl string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}

	err = db.AutoMigrate(&Category{}, &Product{})
	if err != nil {
		log.Fatal("cat create tables:", err)
	}

	return db
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	log.Println("migrations complete successfully")

	h := HandlerImpl{db: db}
	r := gin.Default()
	r.GET("/categories", h.ListCategories)
	r.GET("/categories/:id", h.GetCategory)
	r.POST("/categories", h.CreateCategory)
	r.PUT("/categories/:id", h.UpdateCategory)
	r.DELETE("/categories/:id", h.DeleteCategory)

	r.GET("/products", h.ListProducts)
	r.GET("/products/:id", h.GetProduct)
	r.POST("/products", h.CreateProduct)
	r.PUT("/products/:id", h.UpdateProduct)
	r.DELETE("/products/:id", h.DeleteProduct)
	return r
}
