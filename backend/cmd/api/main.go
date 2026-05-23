package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/interview-question-005/backend/internal/config"
	"example.com/interview-question-005/backend/internal/database"
	"example.com/interview-question-005/backend/internal/handler"
	"example.com/interview-question-005/backend/internal/repository"
	"example.com/interview-question-005/backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	queueRepository := repository.NewGormQueueRepository(db)
	if err := queueRepository.EnsureState(ctx); err != nil {
		log.Fatalf("ensure queue state: %v", err)
	}

	queueService := service.NewQueueService(queueRepository)
	queueHandler := handler.NewQueueHandler(queueService)

	app := fiber.New(fiber.Config{
		AppName: "interview-question-005-api",
	})
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     joinOrigins(cfg.CORSAllowedOrigins),
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization",
		AllowCredentials: false,
	}))

	app.Get("/health", handler.Health)
	app.Get("/api/queue/current", queueHandler.Current)
	app.Post("/api/queue/next", queueHandler.Next)
	app.Post("/api/queue/reset", queueHandler.Reset)

	go func() {
		log.Printf("api listening on :%s", cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Printf("shutdown server: %v", err)
	}
}

func joinOrigins(origins []string) string {
	if len(origins) == 0 {
		return ""
	}

	result := origins[0]
	for _, origin := range origins[1:] {
		result += "," + origin
	}
	return result
}
