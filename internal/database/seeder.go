package database

import (
	"library-management-system-go/internal/config"
	"library-management-system-go/internal/domain"
	"library-management-system-go/pkg/utils"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DefaultUser struct {
	Username   string
	Email      string
	Password   string
	FirstName  string
	LastName   string
	Role       domain.UserRole
	EmployeeID *string
	StudentID  *string
}

func SeedDefaultUsers(cfg *config.Config) error {
	// Connect to database
	dsn := cfg.GetDatabaseDSN()

	var gormLogger logger.Interface
	if cfg.Server.Env == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return err
	}

	// Define default users
	adminEmpID := "EMP001"
	admin2EmpID := "EMP000"
	librarianEmpID := "EMP002"
	teacherEmpID := "EMP003"
	studentID := "STU001"

	defaultUsers := []DefaultUser{
		{
			Username:   "admin",
			Email:      "admin@library.com",
			Password:   "admin123",
			FirstName:  "Admin",
			LastName:   "User",
			Role:       domain.RoleAdmin,
			EmployeeID: &adminEmpID,
		},
		{
			Username:   "admin_2",
			Email:      "admin_2@library.com",
			Password:   "admin29",
			FirstName:  "Admin",
			LastName:   "User",
			Role:       domain.RoleAdmin,
			EmployeeID: &admin2EmpID,
		},
		{
			Username:   "librarian",
			Email:      "librarian@library.com",
			Password:   "librarian123",
			FirstName:  "Librarian",
			LastName:   "User",
			Role:       domain.RoleLibrarian,
			EmployeeID: &librarianEmpID,
		},
		{
			Username:   "teacher",
			Email:      "teacher@library.com",
			Password:   "teacher123",
			FirstName:  "Teacher",
			LastName:   "User",
			Role:       domain.RoleTeacher,
			EmployeeID: &teacherEmpID,
		},
		{
			Username:  "student",
			Email:     "student@library.com",
			Password:  "student123",
			FirstName: "Student",
			LastName:  "User",
			Role:      domain.RoleStudent,
			StudentID: &studentID,
		},
	}

	log.Println("Starting to seed default users...")

	for _, defaultUser := range defaultUsers {
		// Check if user already exists
		var existingUser domain.User
		result := db.Where("username = ? OR email = ?", defaultUser.Username, defaultUser.Email).First(&existingUser)

		if result.Error == nil {
			log.Printf("User '%s' already exists, skipping...", defaultUser.Username)
			continue
		}

		// Hash password
		hashedPassword, err := utils.HashPassword(defaultUser.Password)
		if err != nil {
			log.Printf("Failed to hash password for user '%s': %v", defaultUser.Username, err)
			continue
		}

		// Create user
		user := domain.User{
			Username:        defaultUser.Username,
			Email:           defaultUser.Email,
			PasswordHash:    hashedPassword,
			FirstName:       defaultUser.FirstName,
			LastName:        defaultUser.LastName,
			Role:            defaultUser.Role,
			Status:          domain.StatusActive,
			EmployeeID:      defaultUser.EmployeeID,
			StudentID:       defaultUser.StudentID,
			MaxBooksAllowed: getMaxBooksForRole(defaultUser.Role),
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user '%s': %v", defaultUser.Username, err)
			continue
		}

		log.Printf("✓ Created user: %s (Role: %s, Password: %s)",
			defaultUser.Username, defaultUser.Role, defaultUser.Password)
	}

	log.Println("Default users seeding completed!")
	return nil
}

func getMaxBooksForRole(role domain.UserRole) int {
	switch role {
	case domain.RoleAdmin:
		return 20
	case domain.RoleLibrarian:
		return 15
	case domain.RoleTeacher:
		return 10
	case domain.RoleStudent:
		return 5
	default:
		return 5
	}
}
