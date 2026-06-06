package seed

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/services"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/PrasadNaik1310/LMSVR_SM/models"
)

func DBInsert() {

	dsn := os.Getenv("db_url")
	if dsn == "" {
		log.Fatal("DB URL not set in environment variable 'db_url'")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	users := []models.User{
		{
			ID:           uuid.New(),
			FirstName:    "Prasad",
			LastName:     "Naik",
			Email:        "prasad@lms.com",
			PasswordHash: mustHashPassword("prasad"),
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Phone:        "91845920496",
		},
		{
			ID:           uuid.New(),
			FirstName:    "User2",
			LastName:     "2user",
			Email:        "user2@lms.com",
			PasswordHash: mustHashPassword("password2"),
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Phone:        "9273966592",
		},
		{
			ID:           uuid.New(),
			FirstName:    "User3",
			LastName:     "3user",
			Email:        "user3@lms.com",
			PasswordHash: mustHashPassword("password3"),
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Phone:        "91759275896",
		},
	}

	for _, user := range users {
		var existing models.User
		err := db.Where("email = ?", user.Email).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to insert user %s: %v", user.Email, err)
			} else {
				log.Printf("Inserted user: %s", user.Email)
			}
		} else if err == nil {
			log.Printf("User already exists: %s", user.Email)
		} else {
			log.Printf("Error checking user %s: %v", user.Email, err)
		}
	}
}

func mustHashPassword(password string) string {
	hash, err := services.HashPassword(password)
	if err != nil {
		log.Fatalf("failed to hash seed password: %v", err)
	}
	return hash
}

func Seed(db *gorm.DB) error {
	if err := seedEnquiries(db); err != nil {
		return fmt.Errorf("Seeding enquires :%w", err)
	}
	log.Println("[seeder]: Seeding Enquires complete !!!")
	log.Println("Trying applications!!!")
	if err := seedApplications(db); err != nil {
		return fmt.Errorf("Seeding applications :%w", err)
	}
	log.Println("[seeder]:Seeding Applications complete !!")
	return nil

}
func seedEnquiries(db *gorm.DB) error {
	for _, enquiry := range EnquirySeeds {
		result := db.Where("id = ?", enquiry.ID).FirstOrCreate(&models.Enquiry{}, enquiry)
		if result.Error != nil {
			return fmt.Errorf("enquiry %s: %w", enquiry.ID, result.Error)
		}
		if result.RowsAffected > 0 {
			log.Printf("[seeder] created enquiry: %s (%s)", enquiry.FullName, enquiry.ID)
		} else {
			log.Printf("[seeder] skipped enquiry (already exists): %s", enquiry.ID)
		}
	}
	return nil
}

func seedApplications(db *gorm.DB) error {
	for _, app := range ApplicationSeeds {
		result := db.Where("id = ?", app.ID).FirstOrCreate(&models.Application{}, app)
		if result.Error != nil {
			return fmt.Errorf("application %s: %w", app.ID, result.Error)
		}
		if result.RowsAffected > 0 {
			log.Printf("[seeder] created application: %s (enquiry: %s)", app.ID, app.EnquiryID)
		} else {
			log.Printf("[seeder] skipped application (already exists): %s", app.ID)
		}
	}
	return nil
}
func MigrateAndSeed(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Enquiry{}, &models.Application{}); err != nil {
		return fmt.Errorf("auto-migrate: %w", err)
	}
	return Seed(db)
}
func main() {
	DBInsert()

}
