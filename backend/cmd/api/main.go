package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"potapov.me/taskflow/internal/http/router"
	"potapov.me/taskflow/internal/repository/postgres"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка .env
	_ = godotenv.Load()

	// Подключение к БД
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo, err := postgres.New(ctx, os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer repo.Pool.Close()

	// Настройка роутера
	r := router.SetupRouter()

	// Middleware для передачи репозитория
	r.Use(func(c *gin.Context) {
		c.Set("repo", repo)
		c.Next()
	})

	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + os.Getenv("APP_PORT"),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}
	log.Println("Server exited")
}
