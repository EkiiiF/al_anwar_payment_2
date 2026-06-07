package seeder

import (
	"log"

	"gorm.io/gorm"
)

func RunSeeders(db *gorm.DB) {
	log.Println(">>> [DEBUG] Entering RunSeeders")

	log.Println(">>> [DEBUG] Seeding Roles...")
	SeedRoles(db)

	log.Println(">>> [DEBUG] Migrating SPP to Syahriyyah if exists...")
	MigrateSppToSyahriyyah(db)

	log.Println(">>> [DEBUG] Seeding Master Data...")
	SeedMasterData(db)

	log.Println(">>> [DEBUG] Seeding Users...")
	SeedUsers(db)

	log.Println(">>> [DEBUG] Seeding Settings...")
	SeedSettings(db)

	CleanupDuplicateProfiles(db)

	log.Println(">>> [DEBUG] All Seeders Finished Successfully")
}

func MigrateSppToSyahriyyah(db *gorm.DB) {
	db.Exec("UPDATE categories SET name = 'Syahriyyah', description = 'Tagihan Syahriyyah Bulanan — dikenakan setiap bulan aktif pembelajaran' WHERE name = 'SPP'")
}

func CleanupDuplicateProfiles(db *gorm.DB) {
	log.Println(">>> [DEBUG] Cleaning up duplicate user profiles...")
	var userIDs []string
	db.Raw("SELECT user_id FROM user_profiles WHERE deleted_at IS NULL GROUP BY user_id HAVING COUNT(*) > 1").Scan(&userIDs)

	for _, userID := range userIDs {
		var profiles []struct {
			ID string
		}
		db.Raw("SELECT id FROM user_profiles WHERE user_id = ? AND deleted_at IS NULL ORDER BY updated_at DESC", userID).Scan(&profiles)

		if len(profiles) > 1 {
			for i := 1; i < len(profiles); i++ {
				db.Exec("DELETE FROM user_profiles WHERE id = ?", profiles[i].ID)
			}
			log.Printf(">>> [DEBUG] Deleted %d duplicate profile records for user %s", len(profiles)-1, userID)
		}
	}
}
