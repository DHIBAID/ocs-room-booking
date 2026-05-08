package services

import (
	"log"
	"os"
	"strings"

	"ocs-room-booking/db"
	"ocs-room-booking/models"
	"ocs-room-booking/utils"
)

func EnsureAdminUser() {
	username := strings.TrimSpace(os.Getenv("ADMIN_USERNAME"))
	password := strings.TrimSpace(os.Getenv("ADMIN_PASSWORD"))
	if username == "" || password == "" {
		log.Println("ADMIN_USERNAME or ADMIN_PASSWORD not set; skipping admin bootstrap")
		return
	}
	var count int64
	if err := db.DB.Model(&models.User{}).Where("role = ?", utils.RoleAdmin).Count(&count).Error; err != nil {
		log.Println("Failed to check admin users:", err)
		return
	}
	if count > 0 {
		return
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Println("Failed to hash admin password:", err)
		return
	}
	admin := models.User{
		Username:     username,
		PasswordHash: hash,
		Role:         utils.RoleAdmin,
		IsActive:     true,
	}
	if err := db.DB.Create(&admin).Error; err != nil {
		log.Println("Failed to create admin user:", err)
		return
	}
	log.Println("Admin user created")
}
