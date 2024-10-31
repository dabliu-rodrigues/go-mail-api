package database

import (
	"emailn/internal/domain/campaign"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDb() *gorm.DB {
	dsn := "host=localhost user=emailn_dev password=1234 dbname=emailn_dev port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&campaign.Campaign{}, &campaign.Contact{})
	return db
}

//docker run -e POSTGRES_USER=emailn_dev -e POSTGRES_PASSWORD=1234 -p 5432:5432 -d postgres
