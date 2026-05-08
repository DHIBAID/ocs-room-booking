package handlers

import (
	"net/http"
	"os"
	"time"

	"ocs-room-booking/db"
	"ocs-room-booking/models"
	"ocs-room-booking/utils"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	var user models.User
	if err := db.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.JSONError(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if !user.IsActive {
		utils.JSONError(c, http.StatusUnauthorized, "user is inactive")
		return
	}
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		utils.JSONError(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	ttl := 12 * time.Hour
	token, err := utils.GenerateJWT(user.ID, user.Role, os.Getenv("JWT_SECRET"), ttl)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "token generation failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
