package seeder

import (
	"log"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedMasterData(db *gorm.DB) {
	db.Unscoped().Where("name = ?", "Syahriyyah").Delete(&domain.Category{})
	categories := []domain.Category{
		{
			ID:          uuid.New().String(),
			Name:        "Syahriyyah Pondok",
			BaseAmount:  350000,
			IsFixed:     true,
			IsActive:    true,
			Description: "Tagihan Syahriyyah Pondok Bulanan — operasional pondok",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Syahriyyah Muhadhoroh",
			BaseAmount:  40000,
			IsFixed:     true,
			IsActive:    true,
			Description: "Tagihan Syahriyyah Muhadhoroh Bulanan — sekolah ngaji",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Daftar Ulang",
			BaseAmount:  20000,
			IsFixed:     true,
			IsActive:    true,
			Description: "Biaya daftar ulang — dikenakan di awal tahun ajaran (Syawal)",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Ujian Semester",
			BaseAmount:  35000,
			IsFixed:     true,
			IsActive:    true,
			Description: "Biaya ujian semester — dikenakan di akhir semester (Rabi'ul Awal & Sya'ban)",
		},
	}

	for _, cat := range categories {
		var existing domain.Category
		if err := db.Where("name = ?", cat.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&cat)
				log.Printf("Seeded category: %s (Rp %.0f)", cat.Name, cat.BaseAmount)
			}
		} else {
			existing.BaseAmount = cat.BaseAmount
			existing.Description = cat.Description
			existing.IsActive = cat.IsActive
			db.Save(&existing)
			log.Printf("Updated existing category: %s (Rp %.0f)", cat.Name, cat.BaseAmount)
		}
	}

	statusTypes := []domain.StudentStatusType{
		{
			ID:                 uuid.New().String(),
			Name:               "Reguler",
			DiscountPercentage: 0,
			IsActiveBilling:    true,
			IsDefault:          true,
			Description:        "Santri reguler dengan biaya penuh",
		},
		{
			ID:                 uuid.New().String(),
			Name:               "Abdi Dalem",
			DiscountPercentage: 100,
			IsActiveBilling:    false,
			IsDefault:          false,
			Description:        "Gratis tagihan untuk abdi dalem",
		},
	}

	for _, st := range statusTypes {
		var existing domain.StudentStatusType
		if err := db.Where("name = ?", st.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&st)
				log.Printf("Seeded status type: %s", st.Name)
			}
		} else {
			if st.Name == "Abdi Dalem" {
				existing.IsActiveBilling = false
				existing.DiscountPercentage = 100
				db.Save(&existing)
				log.Printf("Updated existing status type: %s to Gratis Tagihan", st.Name)
			}
		}
	}
}
