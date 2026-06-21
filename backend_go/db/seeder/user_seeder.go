package seeder

import (
	"log"
	"os"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	var superuserRole, pengasuhRole domain.Role
	if err := db.Where("name = ?", "super_user").First(&superuserRole).Error; err != nil {
		log.Printf("Super User role not found: %v", err)
		return
	}
	// if err := db.Where("name = ?", "guardian").First(&guardianRole).Error; err != nil {
	// 	log.Printf("Guardian role not found: %v", err)
	// 	return
	// }
	if err := db.Where("name = ?", "pengasuh").First(&pengasuhRole).Error; err != nil {
		log.Printf("Pengasuh role not found: %v", err)
		return
	}

	defaultPass := os.Getenv("SEEDER_DEFAULT_PASSWORD")
	if defaultPass == "" {
		defaultPass = "password123"
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPass), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return
	}

	users := []domain.User{
		{
			ID:          uuid.New().String(),
			Username:    "super_user",
			Password:    string(hashedPassword),
			RoleID:      superuserRole.ID,
			IsActive:    true,
			FirstName:   "Super",
			LastName:    "User",
			Email:       "super_user@alanwar.com",
			PhoneNumber: "081234567890",
			Gender:      "",
		},
		// {
		// 	ID:          uuid.New().String(),
		// 	Username:    "guardian",
		// 	Password:    string(hashedPassword),
		// 	RoleID:      guardianRole.ID,
		// 	IsActive:    true,
		// 	FirstName:   "Wali",
		// 	LastName:    "Santri",
		// 	Email:       "wali@alanwar.com",
		// 	PhoneNumber: "081234567891",
		// 	Gender:      "L",
		// },
		{
			ID:          uuid.New().String(),
			Username:    "pengasuh",
			Password:    string(hashedPassword),
			RoleID:      pengasuhRole.ID,
			IsActive:    true,
			FirstName:   "Pengasuh",
			LastName:    "",
			Email:       "pengasuh@alanwar.com",
			PhoneNumber: "081234567890",
			Gender:      "",
		},
	}

	for _, user := range users {
		var existing domain.User
		if err := db.Where("username = ?", user.Username).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&user).Error; err != nil {
					log.Printf("Error seeding user %s: %v", user.Username, err)
				} else {
					log.Printf("Seeded user: %s", user.Username)
				}
			}
		} else {
			log.Printf("User %s already exists", user.Username)
		}
	}
}
