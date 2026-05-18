package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/PrasadNaik1310/LMSVR_SM/models"
)

func DBInsert() {
	dsn := "postgresql://neondb_owner:npg_ohi6U2LZxXDG@ep-rapid-band-ap5bewf7-pooler.c-7.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"
	//dsn := os.Getenv("db_url")
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
			PasswordHash: "prasad",
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
			PasswordHash: "password2",
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
			PasswordHash: "password3",
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

func main() {
	DBInsert()
}
