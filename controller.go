package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTask(ctx *gin.Context) {
	var task Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTask := Task{
		Title:       task.Title,
		Description: task.Description,
		Status:      Pending, // move this default logic to model
	}
	db, err := DBConnection()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Db connection failed."})
		return
	}
	if err := db.Create(&newTask).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Faield to create the task."})
	}
	ctx.JSON(http.StatusCreated, newTask)
}

func FectchTask(ctx *gin.Context) {
	id := ctx.Param("id")
	Id, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "Invalid Id format"})
		return
	}
	// TODO: make some common handler
	db, err := DBConnection()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Db connection failed."})
		return
	}
	// var task = Task{ID: int64(Id)} TODO: FInd why this gives error
	var task Task
	result := db.Where("ID = ?", Id).First(&task)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound,
			gin.H{"error": "Task not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

func ListTasks(ctx *gin.Context) {
	var tasks []Task
	db, err := DBConnection()
	if err != nil {
		log.Println(err)
	}
	// TODO: replcae with common paginator like https://medium.com/@michalkowal567/creating-reusable-pagination-in-golang-and-gorm-4b23e179a54b
	db.Scopes(Paginate(ctx.Request)).Find(&tasks)
	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}
