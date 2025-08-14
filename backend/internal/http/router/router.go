package router

import (
	"github.com/gin-gonic/gin"
	"potapov.me/taskflow/internal/http/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		projects := api.Group("/projects")
		{
			projects.POST("", handlers.CreateProject)
			//projects.GET("", handlers.ListProjects)
			//projects.GET("/:id", handlers.GetProject)
			//projects.PATCH("/:id", handlers.UpdateProject)
			//projects.DELETE("/:id", handlers.DeleteProject)
			//
			//tasks := projects.Group("/:id/tasks")
			//{
			//	tasks.POST("", handlers.CreateTask)
			//	tasks.GET("", handlers.ListTasks)
			//}
		}

		//tasks := api.Group("/tasks")
		//{
		//	tasks.GET("/:id", handlers.GetTask)
		//	tasks.PATCH("/:id", handlers.UpdateTask)
		//	tasks.DELETE("/:id", handlers.DeleteTask)
		//}
	}

	return r
}
