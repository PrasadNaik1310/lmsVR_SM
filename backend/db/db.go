package db

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/seed"
	"github.com/PrasadNaik1310/LMSVR_SM/services"
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
		//	return nil
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
		dsn += "&connect_timeout=10"
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
	if os.Getenv("run_migrations") == "true" {
		log.Printf("Running migrations !!")
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
			&models.CourseSchedule{},
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

		//enrollement sequence if doesnot exists
		if err := DB.Exec(`
		CREATE SEQUENCE IF NOT EXISTS student_enrollment_seq
		START WITH 1
		INCREMENT BY 1;
	`).Error; err != nil {
			log.Printf("Failed to create student enrollment sequence: %v", err)
			return err
		}

		// Ensure any old unique index on batches.course_id is removed so multiple batches per course are allowed.
		if err := DB.Exec(`
DO $$
DECLARE r RECORD;
BEGIN
  FOR r IN SELECT indexname FROM pg_indexes WHERE tablename='batches' AND indexdef LIKE '%(course_id)%' LOOP
    EXECUTE format('DROP INDEX IF EXISTS %I', r.indexname);
  END LOOP;
END$$;
`).Error; err != nil {
			log.Printf("Failed to drop unique index on batches.course_id: %v", err)
			// not a fatal error; continue
		}

		log.Printf("models migration doneeee!!")
	}
	// Create default user if not exists
	defaultEmail := "1@2.com"
	defaultPassword := "123"
	hashedDefaultPassword, err := services.HashPassword(defaultPassword)
	if err != nil {
		log.Printf("Failed to hash default user password: %v", err)
		return err
	}
	var user models.User
	err = DB.Where("email = ?", defaultEmail).First(&user).Error
	if err != nil {
		if err.Error() == "record not found" || err == gorm.ErrRecordNotFound {
			user = models.User{
				ID:        uuid.New(),
				FirstName: "Default",
				LastName:  "User",
				Email:     defaultEmail,
				//PasswordHash: defaultPassword, // In production, hash this!
				PasswordHash: hashedDefaultPassword,
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

	// Seed data after migrations and basic user setup
	if err := SeedData(); err != nil {
		log.Printf("Warning: Seeding failed (non-fatal): %v", err)
		// Continue even if seeding fails - it's not critical for DB init
	}
	if os.Getenv("APP_ENV") == "seeding" {
		log.Println("+++++++++++=================+++++++++====+++++++++")
		log.Println("Starting seeding data for applications and enquiry")
		if err := seed.MigrateAndSeed(DB); err != nil {
			log.Fatalf("Failed at data migration for applicationa and enquiry")
		}
	}

	log.Printf("DB CHALUUUUUUUUU !!!!!")
	return nil
}

// SeedData creates RBAC roles/permissions and sample courses/batches per user.
// This is separate from InitDB to keep concerns separated.
func SeedData() error {
	// Seed manage-company RBAC defaults so protected endpoints do not fail with 403 in a fresh DB.
	adminRoleName := "company_admin"
	var adminRole models.Role
	if err := DB.Where("name = ?", adminRoleName).First(&adminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			adminRole = models.Role{
				ID: uuid.New(),
				//ID:          1,
				Name:        adminRoleName,
				Description: "Admin role for company management module",
			}
			if createErr := DB.Create(&adminRole).Error; createErr != nil {
				log.Printf("Failed to create role %s: %v", adminRoleName, createErr)
				return createErr
			}
		} else {
			log.Printf("Failed to query role %s: %v", adminRoleName, err)
			return err
		}
	}

	permissionNames := []string{
		"company.session.assign",
		"company.session.read",
		"company.batch.create",
		"company.batch.read",
		"admission.enquiry.create",
		"admission.enquiry.read",
		"admission.enquiry.update",
		"admission.application.create",
		"admission.application.read",
		"admission.application.approve",
		"admission.application.reject",
		"course_schedule.create",
		"course_schedule.read",
		"course_schedule.update",
		"course_schedule.delete",
		"course_log.create",
		"course_log.read",
		"course_log.update",
	}

	for _, permName := range permissionNames {
		var perm models.Permission
		if err := DB.Where("name = ?", permName).First(&perm).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				perm = models.Permission{
					ID:          uuid.New(),
					Name:        permName,
					Description: "Auto-seeded permission for manage company module",
				}
				if createErr := DB.Create(&perm).Error; createErr != nil {
					log.Printf("Failed to create permission %s: %v", permName, createErr)
					return createErr
				}
			} else {
				log.Printf("Failed to query permission %s: %v", permName, err)
				return err
			}
		}

		var rolePerm models.RolePermission
		if err := DB.Where("role_id = ? AND permission_id = ?", adminRole.ID, perm.ID).First(&rolePerm).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				rolePerm = models.RolePermission{
					ID:           uuid.New(),
					RoleID:       adminRole.ID,
					PermissionID: perm.ID,
				}
				if createErr := DB.Create(&rolePerm).Error; createErr != nil {
					log.Printf("Failed to create role_permission for %s: %v", permName, createErr)
					return createErr
				}
			} else {
				log.Printf("Failed to query role_permission for %s: %v", permName, err)
				return err
			}
		}
	}

	zeroUUID := uuid.UUID{} // all zeros = same as uuid.Nil, but explicit
	if err := DB.Model(&models.User{}).
		Where("role_id = ? OR role_id IS NULL", zeroUUID).
		Update("role_id", adminRole.ID).Error; err != nil {
		log.Printf("Failed to assign default role to users without role: %v", err)
		return err
	}
	// Seed admission RBAC defaults
	admissionRoleName := "admission_admin"
	var admissionRole models.Role
	if err := DB.Where("name = ?", admissionRoleName).First(&admissionRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			admissionRole = models.Role{
				ID:          uuid.New(),
				Name:        admissionRoleName,
				Description: "Admin role for admission module",
			}
			if createErr := DB.Create(&admissionRole).Error; createErr != nil {
				log.Printf("Failed to create role %s: %v", admissionRoleName, createErr)
				return createErr
			}
		} else {
			log.Printf("Failed to query role %s: %v", admissionRoleName, err)
			return err
		}
	}

	admissionPermissions := []string{
		"admission.enquiry.create",
		"admission.enquiry.read",
		"admission.enquiry.update",
		"admission.application.create",
		"admission.application.read",
		"admission.application.approve",
		"admission.application.reject",
	}

	for _, permName := range admissionPermissions {
		var perm models.Permission
		if err := DB.Where("name = ?", permName).First(&perm).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				perm = models.Permission{
					ID:          uuid.New(),
					Name:        permName,
					Description: "Auto-seeded permission for admission module",
				}
				if createErr := DB.Create(&perm).Error; createErr != nil {
					log.Printf("Failed to create permission %s: %v", permName, createErr)
					return createErr
				}
			} else {
				log.Printf("Failed to query permission %s: %v", permName, err)
				return err
			}
		}

		var rolePerm models.RolePermission
		if err := DB.Where("role_id = ? AND permission_id = ?", admissionRole.ID, perm.ID).First(&rolePerm).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				rolePerm = models.RolePermission{
					ID:           uuid.New(),
					RoleID:       admissionRole.ID,
					PermissionID: perm.ID,
				}
				if createErr := DB.Create(&rolePerm).Error; createErr != nil {
					log.Printf("Failed to create role_permission for %s: %v", permName, createErr)
					return createErr
				}
			} else {
				log.Printf("Failed to query role_permission for %s: %v", permName, err)
				return err
			}
		}
	}

	// Seed sample courses and batches per user so frontend can show account-specific data.
	var users []models.User
	if err := DB.Find(&users).Error; err != nil {
		log.Printf("Failed to load users for seeding courses: %v", err)
		return err
	}
	for i, u := range users {
		// create 1-2 courses per user depending on index
		numCourses := 1
		if i%2 == 0 {
			numCourses = 2
		}
		for j := 0; j < numCourses; j++ {
			title := "Course - " + u.Email + " - " + strconv.Itoa(j+1)
			var existingCourse models.Course
			if err := DB.Where("title = ? AND created_by = ?", title, u.ID).First(&existingCourse).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					start := time.Now().AddDate(0, j, 0)
					end := start.AddDate(0, 3, 0)
					c := models.Course{
						ID:         uuid.New(),
						Title:      title,
						Level:      "Beginner",
						Status:     "PUBLISHED",
						TotalSeats: 30 + j*10,
						StartDate:  &start,
						EndDate:    &end,
						CreatedBy:  u.ID,
					}
					if err := DB.Create(&c).Error; err != nil {
						log.Printf("Failed to create course %s: %v", title, err)
						continue
					}
					existingCourse = c
				} else {
					// other error
					log.Printf("Failed checking course %s: %v", title, err)
					continue
				}
			}
			// create a batch for this course if not exists
			var existingBatch models.Batch
			if err := DB.Where("course_id = ?", existingCourse.ID).First(&existingBatch).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					start := time.Now().AddDate(0, j, 0)
					end := start.AddDate(0, 3, 0)
					b := models.Batch{
						ID:          uuid.New(),
						CourseID:    existingCourse.ID,
						BatchName:   existingCourse.Title + " - Batch 1",
						StartDate:   &start,
						EndDate:     &end,
						MaxStudents: 30,
						Status:      "active",
					}
					if err := DB.Create(&b).Error; err != nil {
						log.Printf("Failed to create batch for course %s: %v", existingCourse.Title, err)
						continue
					}
				} else {
					log.Printf("Failed checking batch for course %s: %v", existingCourse.Title, err)
					continue
				}
			}
		}
	}

	return nil
}
