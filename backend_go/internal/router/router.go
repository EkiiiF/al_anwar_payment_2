package router

import (
	"strings"
	"time"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/config"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/controller"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/repository"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/service"
	"github.com/EkiiiF/al_anwar_payment_2.git/pkg/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/gorm"
)

func NewRouter(app *fiber.App, db *gorm.DB, cfg *config.Config) service.SuperUserService {
	validate := validator.New()

	// ─── CORS ──────────────────────────────────────────────
	origins := strings.Split(cfg.CORSOrigins, ",")
	app.Use(cors.New(cors.Config{
		AllowOrigins: origins,
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "Cache-Control", "Pragma", "Expires"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	// ─── Repositories ──────────────────────────────────────
	authRepo := repository.NewAuthRepository()
	guardianRepo := repository.NewGuardianRepository()
	superUserRepo := repository.NewSuperUserRepository()
	logRepo := repository.NewLogRepository()
	paymentRepo := repository.NewPaymentRepository()
	reportRepo := repository.NewReportRepository()
	settingRepo := repository.NewSettingRepository()

	// ─── Services ──────────────────────────────────────────
	logService := service.NewLogService(db, logRepo)
	authService := service.NewAuthService(db, validate, authRepo)
	guardianService := service.NewGuardianService(db, guardianRepo)
	notifService := service.NewNotificationService(db)
	superUserService := service.NewSuperUserService(db, validate, superUserRepo, settingRepo, logService, notifService)
	paymentService := service.NewPaymentService(db, paymentRepo, guardianRepo, cfg)
	reportService := service.NewReportService(db, reportRepo)

	// ─── Controllers ───────────────────────────────────────
	authCtrl := controller.NewAuthController(authService)
	guardianCtrl := controller.NewGuardianController(guardianService)
	superUserCtrl := controller.NewSuperUserController(superUserService)
	paymentCtrl := controller.NewPaymentController(paymentService)
	reportCtrl := controller.NewReportController(reportService)
	logCtrl := controller.NewLogController(logService)

	// ─── JWT Middleware ────────────────────────────────────
	jwtAuth := middleware.JWTMiddleware(db, authRepo)

	// ─── Health Check ──────────────────────────────────────
	app.Get("/", func(c fiber.Ctx) error {
		wib := time.FixedZone("WIB", 7*3600)
		return c.JSON(fiber.Map{
			"status":    "OK",
			"message":   "API Al-Anwar Payment is running",
			"timestamp": time.Now().In(wib).Format("02 January 2006 15:04:05 MST"),
		})
	})

	// ─── API v1 ────────────────────────────────────────────
	api := app.Group("/api/v1")

	// Midtrans Webhook (Public - No Auth)
	api.Post("/payments/notification", paymentCtrl.HandleNotification)

	// Auth (public)
	auth := api.Group("/auth")
	auth.Post("/login", authCtrl.Login)

	// Auth (authenticated)
	authProtected := api.Group("/auth", jwtAuth)
	authProtected.Post("/logout", authCtrl.Logout)
	authProtected.Get("/profile", authCtrl.GetProfile)
	authProtected.Put("/profile", authCtrl.UpdateProfile)
	authProtected.Put("/change-password", authCtrl.ChangePassword)

	// Guardian routes (role: guardian)
	guardian := api.Group("/guardian", jwtAuth, middleware.RequireRole("guardian"))
	guardian.Get("/dashboard", guardianCtrl.GetDashboard)
	guardian.Get("/invoices", guardianCtrl.GetInvoices)
	guardian.Get("/payments", guardianCtrl.GetPayments)
	guardian.Get("/notifications", guardianCtrl.GetNotifications)
	guardian.Patch("/notifications/:id/read", guardianCtrl.MarkNotificationRead)

	// Payment routes (role: guardian)
	payment := api.Group("/payments", jwtAuth, middleware.RequireRole("guardian"))
	payment.Post("/create", paymentCtrl.CreateTransaction)
	payment.Get("/status/:order_id", paymentCtrl.CheckStatus)

	// Super User / Bendahara routes
	su := api.Group("/admin", jwtAuth, middleware.RequireRole("super_user"))

	su.Get("/dashboard", superUserCtrl.GetDashboardStats)
	su.Get("/notifications", guardianCtrl.GetNotifications)
	su.Patch("/notifications/:id/read", guardianCtrl.MarkNotificationRead)

	// Student management
	su.Get("/students", superUserCtrl.GetStudents)
	su.Get("/students/paginated", superUserCtrl.GetStudentsPaginated)
	su.Post("/students", superUserCtrl.CreateStudent)
	su.Put("/students/:id", superUserCtrl.UpdateStudent)
	su.Delete("/students/:id", superUserCtrl.DeleteStudent)
	su.Patch("/students/:id/toggle-status", superUserCtrl.ToggleStudentStatus)

	// Invoice management
	su.Get("/invoices", superUserCtrl.GetInvoices)
	su.Get("/invoices/paginated", superUserCtrl.GetInvoicesPaginated)
	su.Get("/invoices/students-paginated", superUserCtrl.GetStudentsWithInvoicesPaginated)
	su.Post("/invoices/generate", superUserCtrl.GenerateInvoices)
	su.Post("/invoices/generate-semester", superUserCtrl.GenerateSemesterInvoices)

	// Status type management
	su.Get("/status-types", superUserCtrl.GetStatusTypes)
	su.Post("/status-types", superUserCtrl.CreateStatusType)
	su.Put("/status-types/:id", superUserCtrl.UpdateStatusType)
	su.Delete("/status-types/:id", superUserCtrl.DeleteStatusType)

	// Category management
	su.Get("/categories", superUserCtrl.GetCategories)
	su.Post("/categories", superUserCtrl.CreateCategory)
	su.Put("/categories/:id", superUserCtrl.UpdateCategory)
	su.Delete("/categories/:id", superUserCtrl.DeleteCategory)

	// Settings management
	su.Get("/settings", superUserCtrl.GetSettings)
	su.Put("/settings", superUserCtrl.UpdateSetting)

	// User management (super_user can manage all users)
	su.Get("/users", superUserCtrl.GetAllUsers)
	su.Get("/users/paginated", superUserCtrl.GetAllUsersPaginated)
	su.Get("/users/roles", superUserCtrl.GetAllRoles)
	su.Patch("/users/:id/role", superUserCtrl.UpdateUserRole)
	su.Patch("/users/:id/toggle-active", superUserCtrl.ToggleUserActive)

	// Reports (super_user + pengasuh)
	reports := api.Group("/reports", jwtAuth, middleware.RequireRole("super_user", "pengasuh"))
	reports.Get("/", reportCtrl.GetReports)
	reports.Get("/export/csv", reportCtrl.ExportReports)

	// Activity logs (super_user only)
	logs := api.Group("/logs", jwtAuth, middleware.RequireRole("super_user"))
	logs.Get("/", logCtrl.GetLogs)
	logs.Get("/paginated", superUserCtrl.GetActivityLogsPaginated)

	// Pengasuh routes (read-only access)
	pengasuh := api.Group("/pengasuh", jwtAuth, middleware.RequireRole("pengasuh"))
	pengasuh.Get("/dashboard", superUserCtrl.GetDashboardStats)
	pengasuh.Get("/students", superUserCtrl.GetStudents)
	pengasuh.Get("/invoices", superUserCtrl.GetInvoices)

	return superUserService
}
