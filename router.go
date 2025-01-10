package main

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	router := gin.Default()
	publicRoutes := router.Group("/public")
	{
		publicRoutes.GET("/tasks", ListTasks)
	}

	privateRoutes := router.Group("/tasks")
	privateRoutes.Use(AuthMiddleware())
	{
		privateRoutes.POST("", CreateTask)
		privateRoutes.GET("/:id", FetchTask)
		privateRoutes.PUT("/:id", UpdateTask)
		privateRoutes.DELETE("/:id", DeleteTask)
	}

	return router
}
