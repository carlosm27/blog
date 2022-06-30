package model

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() (*gorm.DB, error) {

	dsn := "postgresql://user:password-tier14.aws-us-east-1.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full&options=--cluster%3Dwide-okapi-2719"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	var now time.Time
	db.Raw("SELECT NOW()").Scan(&now)

	log.Println(now)

	if err = db.AutoMigrate(&Expenses{}); err != nil {
		log.Println(err)
	}

	return db, err

}
