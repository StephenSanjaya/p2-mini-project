package config

import (
	"fmt"
	"log"
	"p2-mini-project/src/entity"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load env")
		return nil
	}

	var dbConfig DBEnv
	err = envconfig.Process("DATABASE", &dbConfig)
	if err != nil {
		log.Fatal("Failed to process env")
		return nil
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUsername, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
		return nil
	}

	err = db.AutoMigrate(&entity.PaymentMethod{}, &entity.Coupon{}, &entity.Category{}, &entity.Car{}, &entity.User{}, &entity.Rental{}, &entity.Payment{})
	if err != nil {
		log.Fatal("Failed to auto migrate db: ", err)
		return nil
	}

	fmt.Println("DB Connected")

	return db
}
