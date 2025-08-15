package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tahhh2234/tomomo/backend/internal/config"
	"github.com/tahhh2234/tomomo/backend/internal/handlers"
	"github.com/tahhh2234/tomomo/backend/internal/middleware"
	"github.com/tahhh2234/tomomo/backend/internal/models"
	"gorm.io/gorm"

	"net/http"
)

var db *gorm.DB

// --- Input Structs ---
type CreateTaskInput struct {
	Title    string `json:"title" binding:"required"`
	Priority int    `json:"priority"`
}

type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	Priority    *int    `json:"priority"`
}

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env")
	}

	// Connect DB
	db = config.ConnectDB()
	if err := db.AutoMigrate(&models.User{}, &models.Task{}, &models.RefreshToken{}); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	r := gin.Default()

	// public
	r.POST("/auth/register", handlers.Register(db))
	r.POST("/auth/login", handlers.Login(db))
	r.POST("/auth/refresh", handlers.Refresh(db))
	r.POST("/auth/logout", handlers.Logout(db))

	// protected
	api := r.Group("/tasks", middleware.AuthRequired())
	{
		api.GET("", GetTasks)
		api.GET("/:id", GetTaskByID)
		api.POST("", CreateTask)
		api.PUT("/:id", UpdateTask)
		api.DELETE("/:id", DeleteTask)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Starting server on port", port)
	r.Run(":" + port)
}

// GET /tasks
func GetTasks(c *gin.Context) {
	var tasks []models.Task
	uid := c.GetUint("userID")
	db.Where("user_id = ?", uid).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

// GET /tasks/:id
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	uid := c.GetUint("userID")

	var task models.Task
	if err := db.Where("id = ? AND user_id = ?", id, uid).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// POST /tasks
func CreateTask(c *gin.Context) {
	var input CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := models.Task{
		Title:    input.Title,
		Priority: input.Priority,
		UserID:   c.GetUint("userID"), // **มาจาก JWT**
	}
	if err := db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, task)
}

// PUT /tasks/:id
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	uid := c.GetUint("userID")

	var task models.Task
	if err := db.Where("id = ? AND user_id = ?", id, uid).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&task).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// DELETE /tasks/:id
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	uid := c.GetUint("userID")

	var task models.Task
	if err := db.Where("id = ? AND user_id = ?", id, uid).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	if err := db.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}
	c.Status(http.StatusNoContent)
}
