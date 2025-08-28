package handlers

import (
	"net/http"
	"strings"
	"time"

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
type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
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
		user := models.User{Email: email, PasswordHash: string(hash), Name: in.Name}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
		// issue tokens
		access, accessExp, _ := auth.GenerateAccessToken(user.ID)
		rtPlain, rtExp := auth.NewRefreshToken()
		rtHash := auth.HashRefreshToken(rtPlain)

		db.Create(&models.RefreshToken{
			UserID:    user.ID,
			TokenHash: rtHash,
			ExpiresAt: rtExp,
		})

		c.JSON(http.StatusCreated, gin.H{
			"user":          gin.H{"id": user.ID, "email": user.Email, "name": user.Name},
			"access_token":  access,
			"access_exp":    accessExp,
			"refresh_token": rtPlain, // เก็บฝั่ง client
			"refresh_exp":   rtExp,
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

		access, accessExp, _ := auth.GenerateAccessToken(user.ID)
		rtPlain, rtExp := auth.NewRefreshToken()
		rtHash := auth.HashRefreshToken(rtPlain)

		db.Create(&models.RefreshToken{
			UserID:    user.ID,
			TokenHash: rtHash,
			ExpiresAt: rtExp,
		})

		c.JSON(http.StatusOK, gin.H{
			"user":          gin.H{"id": user.ID, "email": user.Email, "name": user.Name},
			"access_token":  access,
			"access_exp":    accessExp,
			"refresh_token": rtPlain,
			"refresh_exp":   rtExp,
		})
	}
}

// POST /auth/refresh
// ใช้ refresh_token แลก access token ใหม่ และ "rotate" refresh token
func Refresh(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in RefreshInput
		if err := c.ShouldBindJSON(&in); err != nil || strings.TrimSpace(in.RefreshToken) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required"})
			return
		}
		hash := auth.HashRefreshToken(in.RefreshToken)

		var rt models.RefreshToken
		if err := db.Where("token_hash = ?", hash).First(&rt).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		if rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired or revoked"})
			return
		}

		// rotate: revoke current and issue new one
		now := time.Now()
		if err := db.Model(&rt).Update("revoked_at", &now).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rotate token"})
			return
		}

		access, accessExp, _ := auth.GenerateAccessToken(rt.UserID)
		newPlain, newExp := auth.NewRefreshToken()
		newHash := auth.HashRefreshToken(newPlain)

		if err := db.Create(&models.RefreshToken{
			UserID:    rt.UserID,
			TokenHash: newHash,
			ExpiresAt: newExp,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue new refresh token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  access,
			"access_exp":    accessExp,
			"refresh_token": newPlain,
			"refresh_exp":   newExp,
		})
	}
}

// POST /auth/logout
// ยกเลิก refresh token ปัจจุบัน (หรือทั้งหมดของ user ก็ได้)
type LogoutInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in LogoutInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required"})
			return
		}
		hash := auth.HashRefreshToken(in.RefreshToken)

		var rt models.RefreshToken
		if err := db.Where("token_hash = ?", hash).First(&rt).Error; err != nil {
			// idempotent: ตอบ 200 ไปเลย
			c.JSON(http.StatusOK, gin.H{"success": true})
			return
		}
		now := time.Now()
		db.Model(&rt).Update("revoked_at", &now)
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

// GET /auth/me
func Me(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetUint("userID") // เอามาจาก middleware.JWT
		var user models.User
		if err := db.First(&user, uid).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		})
	}
}
