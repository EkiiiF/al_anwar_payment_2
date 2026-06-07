package main

import (
	"log"

	"github.com/EkiiiF/al_anwar_payment_2.git/db/seeder"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/config"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/exception"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/repository"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/router"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/service"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/worker"
	"github.com/gofiber/fiber/v3"
)

func main() {
	// 1. Load konfigurasi dari .env
	cfg := config.NewConfig()

	// 2. Inisialisasi JWT secret dari config
	utils.InitJWTSecret(cfg.JWTSecret)

	// 3. Koneksi database
	db := config.NewDatabase(cfg)

	// 4. Jalankan Seeder
	log.Println("Running seeder...")
	seeder.RunSeeders(db)
	log.Println("Seeding completed successfully")

	// 5. Setup Fiber v3 dengan custom error handler
	app := fiber.New(fiber.Config{
		AppName:      "Al-Anwar Payment API",
		ErrorHandler: exception.ErrorHandler,
	})

	// 6. Setup router dengan dependency injection
	suService := router.NewRouter(app, db, cfg)

	// 7. Jalankan background workers
	worker.StartTokenCleanupWorker(db)

	logRepo := repository.NewLogRepository()
	logService := service.NewLogService(db, logRepo)
	worker.StartLogCleanupWorker(logService)
	worker.StartAutoBillingWorker(db, suService)

	// 8. Start server
	log.Printf("🚀 Server starting on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
