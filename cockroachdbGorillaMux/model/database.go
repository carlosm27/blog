package model

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() (*gorm.DB, error) {

	dsn := "DATABASE_URL"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}

	var now time.Time
	db.Raw("SELECT NOW()").Scan(&now)

	log.Println(now)

	if err = db.AutoMigrate(&Expenses{}); err != nil {
		log.Println(err)
	}

	return db, err

}
