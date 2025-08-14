package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"potapov.me/taskflow/internal/domain/project"
	"potapov.me/taskflow/internal/repository/postgres"
)

func CreateProject(c *gin.Context) {
	var p project.Project
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Генерация ID
	p.ID = uuid.New()

	// Получение репозитория из контекста (добавляется в middleware)
	repo := c.MustGet("repo").(*postgres.Repository)

	if err := repo.CreateProject(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusCreated, p)
}

// Добавьте остальные обработчики...
