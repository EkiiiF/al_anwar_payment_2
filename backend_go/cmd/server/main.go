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
	cfg := config.NewConfig()
	utils.InitJWTSecret(cfg.JWTSecret)

	db := config.NewDatabase(cfg)

	log.Println("Running seeder...")
	seeder.RunSeeders(db)
	log.Println("Seeding completed successfully")

	app := fiber.New(fiber.Config{
		AppName:      "Al-Anwar Payment API",
		ErrorHandler: exception.ErrorHandler,
	})

	suService := router.NewRouter(app, db, cfg)

	worker.StartTokenCleanupWorker(db)

	logRepo := repository.NewLogRepository()
	logService := service.NewLogService(db, logRepo)
	worker.StartLogCleanupWorker(logService)
	worker.StartAutoBillingWorker(db, suService)

	log.Printf("🚀 Server starting on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
