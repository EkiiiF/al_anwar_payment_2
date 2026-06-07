package seeder

import (
	"log"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

func SeedSettings(db *gorm.DB) {
	db.Unscoped().Where("`key` = ?", "auto_billing_day").Delete(&domain.Setting{})

	settings := []domain.Setting{
		{
			Key:         "auto_generate_billing",
			Value:       "false",
			Description: "Otomatis buat tagihan Syahriyyah setiap tanggal Hijriah yang ditentukan",
		},
		{
			Key:         "billing_hijri_day",
			Value:       "1",
			Description: "Tanggal Hijriah eksekusi auto billing (1-30)",
		},
		{
			Key:         "billing_hijri_month",
			Value:       "0",
			Description: "Bulan Hijriah target billing (0 = otomatis mengikuti bulan sekarang)",
		},
		{
			Key:         "billing_hijri_year",
			Value:       "0",
			Description: "Tahun Hijriah target billing (0 = otomatis mengikuti tahun sekarang)",
		},
		{
			Key:         "auto_generate_semester_billing",
			Value:       "false",
			Description: "Otomatis buat seluruh tagihan semester di awal semester",
		},
		{
			Key:         "midtrans_environment",
			Value:       "sandbox",
			Description: "Midtrans environment (sandbox/production)",
		},
	}

	for _, s := range settings {
		var existing domain.Setting
		if err := db.Where("`key` = ?", s.Key).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&s)
				log.Printf("Seeded setting: %s", s.Key)
			}
		}
	}
}
