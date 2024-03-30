package main

import (
	"fmt"
	"net/http"

	"time"

	"konzek-jun/app"
	"konzek-jun/configs"
	"konzek-jun/loggerx"
	"konzek-jun/middleware"
	"konzek-jun/prometheus"
	"konzek-jun/repository"
	"konzek-jun/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	prometheus.InitPrometheus()
	loggerx.Init()

	go func() {
		if err := http.ListenAndServe(":2222", promhttp.Handler()); err != nil {
			fmt.Println("Prometheus sunucusunu başlatırken hata oluştu:", err)
		}
	}()

	go func() {
		for {
			prometheus.HandleHTTPRequest()
			time.Sleep(30 * time.Second) // 10 saniye bekleyin
		}
	}()
	appRoute := fiber.New()
	db := configs.ConnectDB()

	defer db.Close()

	taskRepository := repository.NewTaskRepository(db)

	td := app.NewTaskHandler(services.NewTaskService(taskRepository), 5)

	authService := services.NewAuthService(repository.NewUserRepo(db))

	jwtService := services.NewJWTService()

	userService := services.NewUserService(repository.NewUserRepo(db))

	authHandler := app.NewAuthHandler(authService, jwtService, userService)

	appRoute.Use(recover.New())

	jwtMiddleware := middleware.NewJWTMiddleware(services.NewJWTService())

	appRoute.Use(limiter.New(limiter.Config{
		Max:        5, // Maximum 5 requests per second
		Expiration: 1, // Expire limiter after 1 second
		KeyGenerator: func(ctx *fiber.Ctx) string {
			return ctx.IP() // Generate unique key based on client IP address
		},
		LimitReached: func(ctx *fiber.Ctx) error {
			return ctx.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
	}))

	appRoute.Use(func(ctx *fiber.Ctx) error {
		// Middleware'i atlamak istediğimiz endpointlerin adları
		skipEndpoints := []string{"/api/register", "/api/login", "/metrics"}

		// Endpoint adını kontrol et
		for _, skipEndpoint := range skipEndpoints {
			if ctx.Path() == skipEndpoint {
				// Middleware'i atla
				fmt.Println("Atladı:", skipEndpoint)
				return ctx.Next()
			}
		}

		// Diğer durumlarda, JWT doğrulamasını yap
		return jwtMiddleware.AuthorizeJWT(ctx)
	})
	appRoute.Use(prometheus.MeasureRequestDuration)

	appRoute.Post("/api/tasks", td.CreateTask)
	appRoute.Get("/api/tasks", td.GetAllTask)
	appRoute.Get("/api/tasks/page", td.GetAllTaskWithPagination)
	appRoute.Delete("/api/tasks/:id", td.DeleteTask)
	appRoute.Get("/api/tasks/:id", td.GetByID)
	appRoute.Put("/api/tasks", td.UpdateTask)
	appRoute.Post("/api/register", authHandler.Register)
	appRoute.Post("/api/login", authHandler.Login)
	appRoute.Listen(":8080")
}
