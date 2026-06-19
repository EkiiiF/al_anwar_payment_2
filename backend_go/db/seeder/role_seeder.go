package seeder

import (
	"log"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []domain.Role{
		{
			ID:          uuid.New().String(),
			Name:        "super_user",
			Description: "Super User - Bendahara Utama",
		},
		{
			ID:          uuid.New().String(),
			Name:        "guardian",
			Description: "Guardian - Wali Santri",
		},
		{
			ID:          uuid.New().String(),
			Name:        "pengasuh",
			Description: "Pengasuh - Pemimpin Pesantren",
		},
	}

	for _, role := range roles {
		var existing domain.Role
		if err := db.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&role).Error; err != nil {
					log.Printf("Error seeding role %s: %v", role.Name, err)
				} else {
					log.Printf("Seeded role: %s", role.Name)
				}
			}
		} else {
			log.Printf("Role %s already exists", role.Name)
		}
	}
}
