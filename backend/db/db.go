package db

import (
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	log.Println("Init DB")
	dsn := os.Getenv("db_url")
	if dsn == "" {
		log.Fatal("Could not get DB URL , INIT DB")
		return nil
	}
	if strings.HasPrefix(dsn, "postgres://") {
		if !strings.Contains(dsn, "connect_timeout") {
			u, err := url.Parse(dsn)
			if err == nil {
				q := u.Query()
				q.Set("connect_timeout", "10")
				u.RawQuery = q.Encode()
				dsn = u.String()

			}
		}
	} else if !strings.Contains(dsn, "connect_timeout") {
		// DSN format (key=value pairs)
		dsn += "&X connect_timeout=10"
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("Nahi ho rha connect ")
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}
	DB = database
	log.Print("Connected to database successfully")

	//configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetMaxIdleConns(102)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	err = DB.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.User{},
		&models.Student{},
		&models.Teacher{},
		&models.Enquiry{},
		&models.Application{},
		&models.AcademicSession{},
		&models.Team{},
		&models.TeamMember{},
		&models.Course{},
		&models.Batch{},
		&models.BatchTeacher{},
		&models.CourseModule{},
		&models.Lesson{},
		&models.LessonFile{},
		&models.Enrollment{},
		&models.Membership{},
		&models.Payment{},
		&models.LiveSession{},
		&models.LiveSessionAttendance{},
		&models.CoursePlanner{},
		&models.CourseLog{},
		&models.Announcement{},
		&models.AnnouncementCourse{},
		&models.LessonSubmission{},
		&models.CourseInvite{},
		&models.Notification{},
	)
	if err != nil {
		log.Printf("Failed to migrate models: %v", err)
		return err
	}
	/*
		defaultEmail := os.Getenv("default_user")
		if defaultEmail == "" {
			fmt.Errorf("Default user not set")
		}
		defaultPassword := os.Getenv("default_password")
		var testUser models.User
		if err := DB.Where("user_email = ?", defaultEmail).First(&testUser).Error; err != nil {
			testUser = models.User{
				Email:        defaultEmail,
				PasswordHash: defaultPassword,
				FirstName:    "Test",
				LastName:     "User",
			}
			if createErr := DB.Create(&testUser).Error; createErr != nil {
				log.Printf("failed to create default test user: %v", createErr)
				return createErr
			}
			log.Printf("default test user created: %s", defaultEmail)
		} else {
			log.Printf("default test user already exists: %s", defaultEmail)
		}
	*/

	log.Printf("models migration doneeee!!")
	// Create default user if not exists
	defaultEmail := "1@2.com"
	defaultPassword := "123"
	var user models.User
	err = DB.Where("email = ?", defaultEmail).First(&user).Error
	if err != nil {
		if err.Error() == "record not found" || err == gorm.ErrRecordNotFound {
			user = models.User{
				ID:           uuid.New(),
				FirstName:    "Default",
				LastName:     "User",
				Email:        defaultEmail,
				PasswordHash: defaultPassword, // In production, hash this!
				IsActive:     true,
			}
			if createErr := DB.Create(&user).Error; createErr != nil {
				log.Printf("Failed to create default user: %v", createErr)
				return createErr
			}
			log.Printf("Default user created: %s", defaultEmail)
		} else {
			log.Printf("Failed to query for default user: %v", err)
			return err
		}
	} else {
		log.Printf("Default user already exists: %s", defaultEmail)
	}
	log.Printf("DB CHALUUUUUUUUU !!!!!")
	return nil
}
