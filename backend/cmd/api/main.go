package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"potapov.me/taskflow/internal/domain/project"
	"potapov.me/taskflow/internal/domain/task"
	"potapov.me/taskflow/internal/domain/user"
	"potapov.me/taskflow/internal/http/router"
	"potapov.me/taskflow/internal/repository/postgres"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	repo, err := postgres.New(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer repo.Close()

	// Автомиграция
	if err := repo.AutoMigrate(
		&user.User{},
		&project.Project{},
		&task.Task{},
	); err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	// Создаем тестового пользователя для разработки
	if os.Getenv("APP_ENV") == "dev" {
		createTestUser(repo.DB)
	}

	r := router.SetupRouter()

	r.Use(func(c *gin.Context) {
		c.Set("repo", repo)
		c.Next()
	})

	// ... остальной код (graceful shutdown)
}

func createTestUser(db *gorm.DB) {
	testUser := &user.User{
		Email:        "test@taskflow.dev",
		Name:         "Test User",
		PasswordHash: "hashed_password", // В реальном приложении используйте bcrypt
	}

	if err := db.FirstOrCreate(testUser, "email = ?", testUser.Email).Error; err != nil {
		log.Printf("Failed to create test user: %v", err)
	}
}
