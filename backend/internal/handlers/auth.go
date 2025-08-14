package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tahhh2234/tomomo/backend/internal/auth"
	"github.com/tahhh2234/tomomo/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in RegisterInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email & password required"})
			return
		}
		email := strings.TrimSpace(strings.ToLower(in.Email))
		if email == "" || strings.TrimSpace(in.Password) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email & password required"})
			return
		}

		var exists int64
		db.Model(&models.User{}).Where("email = ?", email).Count(&exists)
		if exists > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		user := models.User{
			Email:        email,
			PasswordHash: string(hash),
			Name:         in.Name,
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
		token, _ := auth.GenerateToken(user.ID)
		c.JSON(http.StatusCreated, gin.H{
			"user":  gin.H{"id": user.ID, "email": user.Email, "name": user.Name},
			"token": token,
		})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in LoginInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email & password required"})
			return
		}
		email := strings.TrimSpace(strings.ToLower(in.Email))

		var user models.User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		token, _ := auth.GenerateToken(user.ID)
		c.JSON(http.StatusOK, gin.H{
			"user":  gin.H{"id": user.ID, "email": user.Email, "name": user.Name},
			"token": token,
		})
	}
}
