package database

import (
	"fmt"
	"log"
	"os"

	"careersync/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// 1. Get the DB_URL from the environment (Render will provide this!)
	dsn := os.Getenv("DB_URL")
	
	// If no DB_URL is found, default to local (for your laptop)
	if dsn == "" {
		dsn = "host=localhost user=postgres password=password123 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("🚀 Database Connected Successfully!")

	// 2. Migrate Tables
	fmt.Println("⚙️  Migrating the database...")
	err = DB.AutoMigrate(
		&models.User{}, 
		&models.Company{}, 
		&models.ReferralRequest{}, 
	)
	
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	fmt.Println("✅ Database Migration Complete!")
}