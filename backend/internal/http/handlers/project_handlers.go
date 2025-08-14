package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"potapov.me/taskflow/internal/domain/project"
	"potapov.me/taskflow/internal/repository/postgres"
)

func CreateProject(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем userID из контекста (будет добавлено в middleware аутентификации)
	userID, _ := c.Get("userID")

	project := &project.Project{
		Title:       input.Title,
		Description: input.Description,
		OwnerID:     userID.(uuid.UUID),
	}

	repo := c.MustGet("repo").(*postgres.Repository)
	if err := repo.CreateProject(c.Request.Context(), project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// Добавьте остальные обработчики...
