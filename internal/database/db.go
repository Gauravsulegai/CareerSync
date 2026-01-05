package database

import (
	"fmt"
	"log"

	"github.com/Gauravsulegai/careersync/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a global variable
var DB *gorm.DB

func ConnectDB() {
	// 1. Connection settings
	dsn := "host=localhost user=admin password=password123 dbname=careersync port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	// 2. Open connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// 3. Check for errors
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("🚀 Database Connected Successfully!")

	// 4. MIGRATE THE TABLES (Updated to include all new models)
	fmt.Println("⚙️  Migrating the database...")
	err = DB.AutoMigrate(
		&models.User{}, 
		&models.Company{}, 
		&models.ReferralRequest{}, 
		&models.Notification{},
	)
	
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	fmt.Println("✅ Database Migration Complete!")
}