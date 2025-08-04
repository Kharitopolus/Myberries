package handlers

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID `gorm:"primarykey" json:"id"`
	Name      string    `                  json:"name"`
	CreatedAt time.Time `                  json:"created_at"`
	UpdatedAt time.Time `                  json:"updated_at"`

	Products []Product
}

type Product struct {
	ID          uuid.UUID `gorm:"primarykey"`
	Name        string
	Description string
	PriceUSD    float64 `gorm:"type:numeric(8,2)"`
	CategoryID  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Category Category
}

func CreateTables(dbUrl string) error {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}

	err = db.AutoMigrate(&Category{}, &Product{})
	if err != nil {
		log.Fatal("cat create tables:", err)
	}

	log.Println("migrations complete successfully")

	return nil
}
