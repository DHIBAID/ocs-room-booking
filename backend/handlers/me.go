package handlers

import (
	"net/http"

	"ocs-room-booking/middleware"
	"ocs-room-booking/utils"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	user, ok := middleware.GetAuthUser(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}
