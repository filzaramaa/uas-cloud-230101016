package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func koneksi() *gorm.DB {
	// PENTING: Load file .env terlebih dahulu
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	// Sekarang os.Getenv baru bisa mendapatkan isinya
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	dsn := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
