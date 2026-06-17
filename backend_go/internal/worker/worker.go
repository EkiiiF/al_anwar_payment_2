package worker

import (
	"context"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/request"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/service"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"gorm.io/gorm"
)

// StartTokenCleanupWorker berjalan di background untuk menghapus token kadaluarsa dan transaksi pending kadaluarsa.
// Interval: setiap 1 jam.
func StartTokenCleanupWorker(db *gorm.DB) {
	go func() {
		log.Println("[Worker] Token & Payment cleanup worker started")
		cleanupTokens(db)
		cleanupExpiredPayments(db)

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			cleanupTokens(db)
			cleanupExpiredPayments(db)
		}
	}()
}

// StartLogCleanupWorker berjalan di background untuk membersihkan activity logs lama.
// Interval: setiap 24 jam.
func StartLogCleanupWorker(logService service.LogService) {
	go func() {
		log.Println("[Worker] Log cleanup worker started")
		if err := logService.CleanupLogs(); err != nil {
			log.Printf("[Worker] Error cleaning up logs on startup: %v", err)
		} else {
			log.Println("[Worker] Initial log cleanup completed successfully")
		}

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("[Worker] Running scheduled log cleanup...")
			if err := logService.CleanupLogs(); err != nil {
				log.Printf("[Worker] Error cleaning up logs: %v", err)
			} else {
				log.Println("[Worker] Scheduled log cleanup completed successfully")
			}
		}
	}()
}

// cleanupTokens menghapus token blacklist yang sudah kadaluarsa.
func cleanupTokens(db *gorm.DB) {
	log.Println("[Worker] Memulai pembersihan token kadaluarsa...")

	resBlacklist := db.Where("expired_at < ?", time.Now()).Delete(&domain.TokenBlacklist{})
	if resBlacklist.Error != nil {
		log.Printf("[Worker] Gagal menghapus Token Blacklist: %v", resBlacklist.Error)
	} else if resBlacklist.RowsAffected > 0 {
		log.Printf("[Worker] Berhasil membersihkan %d Token Blacklist", resBlacklist.RowsAffected)
	}
}

// StartAutoBillingWorker berjalan di background untuk membuat tagihan otomatis.
// Menggunakan kalender Hijriah untuk menentukan komponen tagihan semester.
// Interval: setiap 24 jam.
func StartAutoBillingWorker(db *gorm.DB, suService service.SuperUserService) {
	go func() {
		log.Println("[Worker] Auto billing worker started (Hijri calendar-based)")

		// Run once on startup regardless of date (it will skip if already exists)
		checkAndGenerateBilling(db, suService, true)

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			// Periodic check only on 1st day
			checkAndGenerateBilling(db, suService, false)
		}
	}()
}

func checkAndGenerateBilling(db *gorm.DB, suService service.SuperUserService, force bool) {
	// ─── Hijri Calendar Info ───────────────────────────────
	hijriNow := utils.GetCurrentHijriDate()

	// Ambil setting tanggal auto billing (Hijriah)
	targetDay := 1
	var daySetting domain.Setting
	if err := db.Where("`key` = ?", "billing_hijri_day").First(&daySetting).Error; err == nil {
		if val, err := strconv.Atoi(daySetting.Value); err == nil && val >= 1 && val <= 30 {
			targetDay = val
		}
	}

	// Tentukan bulan dan tahun target billing (Hijriah) - selalu mengikuti bulan/tahun saat ini
	targetMonth := hijriNow.Month
	targetYear := hijriNow.Year

	// Jika tidak force, hanya proses di tanggal Hijriah yang ditentukan
	if !force && hijriNow.Day != targetDay {
		// Cek H-5 Notifikasi Pre-Billing sebelum tanggal otomatis tertagih berikutnya
		var nextMonth, nextYear int
		if hijriNow.Day < targetDay {
			nextMonth = hijriNow.Month
			nextYear = hijriNow.Year
		} else {
			if hijriNow.Month == 12 {
				nextMonth = 1
				nextYear = hijriNow.Year + 1
			} else {
				nextMonth = hijriNow.Month + 1
				nextYear = hijriNow.Year
			}
		}

		todayGreg := utils.HijriToGregorian(hijriNow.Year, hijriNow.Month, hijriNow.Day)
		billingGreg := utils.HijriToGregorian(nextYear, nextMonth, targetDay)
		daysDiff := int(math.Round(billingGreg.Sub(todayGreg).Hours() / 24.0))

		if daysDiff == 5 {
			log.Printf("[Worker] Hari ini H-5 sebelum tagihan otomatis %s %d H. Mengirim notifikasi pre-billing...", utils.GetHijriMonthName(nextMonth), nextYear)
			ctx := context.Background()
			if err := suService.SendPreBillingNotifications(ctx, nextMonth, nextYear, targetDay); err != nil {
				log.Printf("[Worker] Gagal mengirim notifikasi pre-billing: %v", err)
			}
		}
		return
	}

	semesterInfo := utils.GetSemesterInfo(targetMonth, targetYear)

	log.Printf("[Worker] Target Hijri Billing: %s %d H | %s | Exam: %v | Registration: %v",
		utils.GetHijriMonthName(targetMonth), targetYear,
		semesterInfo.Name, semesterInfo.IsExamMonth, semesterInfo.IsRegistration)

	ctx := context.Background()

	// ─── 1. Auto Billing Syahriyyah ───────────────────────────
	var billingMonthlySetting domain.Setting
	if err := db.Where("`key` = ?", "auto_generate_billing").First(&billingMonthlySetting).Error; err == nil {
		if billingMonthlySetting.Value == "true" {
			if !utils.IsBillableMonth(targetMonth) {
				log.Printf("[Worker] Bulan %s — tidak ada tagihan Syahriyyah", utils.GetHijriMonthName(targetMonth))
			} else {
				log.Println("[Worker] Auto billing Syahriyyah ON — generating...")
				// Hitung padanan masehi untuk request payload agar validasi struct lolos
				gregDate := utils.HijriToGregorian(targetYear, targetMonth, 1)
				req := request.GenerateInvoiceRequest{
					Month:      int(gregDate.Month()),
					Year:       gregDate.Year(),
					HijriMonth: targetMonth,
					HijriYear:  targetYear,
				}
				count, err := suService.GenerateInvoices(ctx, req, "SYSTEM", "127.0.0.1", "Worker-AutoBilling")
				if err != nil {
					log.Printf("[Worker] Error generating Syahriyyah invoices: %v", err)
				} else if count > 0 {
					log.Printf("[Worker] Generated %d Syahriyyah invoices for %s %d H",
						count, utils.GetHijriMonthName(targetMonth), targetYear)
				} else {
					log.Println("[Worker] No new Syahriyyah invoices (already exist)")
				}
			}
		} else {
			log.Println("[Worker] Auto billing Syahriyyah OFF")
		}
	}

	// ─── 2. Auto Billing Semester ──────────────────────────
	var billingSemesterSetting domain.Setting
	if err := db.Where("`key` = ?", "auto_generate_semester_billing").First(&billingSemesterSetting).Error; err == nil {
		if billingSemesterSetting.Value == "true" {
			isSemesterStart := (targetMonth == utils.HijriSyawal || targetMonth == utils.HijriRobiulAkhir)
			if isSemesterStart || force {
				log.Printf("[Worker] Auto billing semester ON — generating semester %d...", semesterInfo.Number)
				count, err := suService.GenerateSemesterInvoices(ctx, semesterInfo.Number, targetYear, "SYSTEM", "127.0.0.1", "Worker-AutoSemesterBilling")
				if err != nil {
					log.Printf("[Worker] Error generating semester invoices: %v", err)
				} else if count > 0 {
					log.Printf("[Worker] Generated %d semester invoices for Semester %d %d H",
						count, semesterInfo.Number, targetYear)
				} else {
					log.Println("[Worker] No new semester invoices (already exist)")
				}
			}
		} else {
			log.Println("[Worker] Auto billing semester OFF")
		}
	}
}

// cleanupExpiredPayments mencari pembayaran pending yang sudah kadaluarsa (lebih dari 24 jam)
// dan mengembalikannya ke status unpaid.
func cleanupExpiredPayments(db *gorm.DB) {
	log.Println("[Worker] Memulai pemeriksaan transaksi pending yang kadaluarsa...")

	var expiredPayments []domain.Payment
	expirationThreshold := time.Now().Add(-24 * time.Hour)
	err := db.Preload("Invoice").
		Where("(transaction_status = ? OR transaction_status = ?) AND created_at < ?", "pending", "pending_payment", expirationThreshold).
		Find(&expiredPayments).Error

	if err != nil {
		log.Printf("[Worker] Gagal mencari transaksi kadaluarsa: %v", err)
		return
	}

	if len(expiredPayments) == 0 {
		log.Println("[Worker] Tidak ada transaksi pending yang kadaluarsa.")
		return
	}

	log.Printf("[Worker] Ditemukan %d transaksi pending yang kadaluarsa. Memproses pembatalan...", len(expiredPayments))

	for _, p := range expiredPayments {
		err := db.Transaction(func(tx *gorm.DB) error {
			statusExpired := "expire"
			p.TransactionStatus = &statusExpired
			
			if err := tx.Save(&p).Error; err != nil {
				return err
			}

			if err := tx.Model(&domain.Invoice{}).Where("id = ?", p.InvoiceID).Update("status", "unpaid").Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Printf("[Worker] Gagal membatalkan transaksi %s: %v", p.ID, err)
		} else {
			log.Printf("[Worker] Transaksi %s (Invoice ID %s) telah kadaluarsa. Status dikembalikan ke unpaid.", p.ID, p.InvoiceID)
		}
	}
}
