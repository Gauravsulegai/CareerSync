package database

import (
	"fmt"
	"log"

	"careersync/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a global variable
var DB *gorm.DB

func ConnectDB() {
	// 1. Connection settings
	dsn := "host=localhost user=postgres password=password123 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	// 2. Open connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// 3. Check for errors
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("🚀 Database Connected Successfully!")

	// 👇 TEMPORARY: This line deletes the old table so you don't have to do it manually!
	// DB.Migrator().DropTable(&models.ReferralRequest{}) 

	// 4. MIGRATE THE TABLES
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