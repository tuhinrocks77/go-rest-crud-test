package main

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/tasks", CreateTask)
	router.GET("/public/tasks", ListTasks)

	return router
}
