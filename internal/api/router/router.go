package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Setup(cfg *config.Config, authService *service.AuthService, examService *service.ExamService) *fiber.App {
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Routes
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register(authService))
	auth.Post("/login", handlers.Login(authService))

	// Protected routes
	protected := api.Group("/", middleware.AuthMiddleware([]byte(cfg.JWT.Secret)))

	// Exam routes
	exams := protected.Group("/exams")
	exams.Get("/", handlers.GetExams(examService))
	exams.Post("/", handlers.CreateExam(examService))
	exams.Get("/:id", handlers.GetExam(examService))
	exams.Put("/:id", handlers.UpdateExam(examService))
	exams.Delete("/:id", handlers.DeleteExam(examService))

	return app
}
