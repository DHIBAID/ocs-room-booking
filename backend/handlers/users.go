package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"ocs-room-booking/db"
	"ocs-room-booking/models"
	"ocs-room-booking/utils"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"`
}

type updateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"`
}

func CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Username) == "" || len(req.Password) < 8 {
		utils.JSONError(c, http.StatusBadRequest, "username and password are required")
		return
	}
	role, ok := utils.CanonicalRole(req.Role)
	if !ok {
		utils.JSONError(c, http.StatusBadRequest, "invalid role")
		return
	}
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to hash password")
		return
	}
	user := models.User{
		Username:     req.Username,
		PasswordHash: hash,
		Role:         role,
		IsActive:     true,
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if err := db.DB.Create(&user).Error; err != nil {
		utils.JSONError(c, http.StatusConflict, "username already exists")
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"role":      user.Role,
		"is_active": user.IsActive,
	})
}

func ListUsers(c *gin.Context) {
	var users []models.User
	if err := db.DB.Order("id asc").Find(&users).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch users")
		return
	}
	c.JSON(http.StatusOK, users)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid user id")
		return
	}
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		utils.JSONError(c, http.StatusNotFound, "user not found")
		return
	}
	updates := map[string]interface{}{}
	if strings.TrimSpace(req.Username) != "" {
		updates["username"] = req.Username
	}
	if strings.TrimSpace(req.Password) != "" {
		if len(req.Password) < 8 {
			utils.JSONError(c, http.StatusBadRequest, "password must be at least 8 characters")
			return
		}
		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			utils.JSONError(c, http.StatusInternalServerError, "failed to hash password")
			return
		}
		updates["password_hash"] = hash
	}
	if strings.TrimSpace(req.Role) != "" {
		role, ok := utils.CanonicalRole(req.Role)
		if !ok {
			utils.JSONError(c, http.StatusBadRequest, "invalid role")
			return
		}
		updates["role"] = role
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if len(updates) == 0 {
		utils.JSONError(c, http.StatusBadRequest, "no fields to update")
		return
	}
	if err := db.DB.Model(&user).Updates(updates).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to update user")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"role":      user.Role,
		"is_active": user.IsActive,
	})
}

func DeactivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid user id")
		return
	}
	if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to remove user")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "removed"})
}
