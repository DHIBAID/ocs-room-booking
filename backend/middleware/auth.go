package middleware

import (
	"net/http"
	"os"
	"strings"

	"ocs-room-booking/db"
	"ocs-room-booking/models"
	"ocs-room-booking/utils"

	"github.com/gin-gonic/gin"
)

type AuthUser struct {
	ID       uint
	Username string
	Role     string
}

const ContextUserKey = "auth_user"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.JSONError(c, http.StatusUnauthorized, "missing or invalid token")
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(tokenStr, os.Getenv("JWT_SECRET"))
		if err != nil {
			utils.JSONError(c, http.StatusUnauthorized, "invalid token")
			return
		}
		var user models.User
		if err := db.DB.First(&user, claims.UserID).Error; err != nil {
			utils.JSONError(c, http.StatusUnauthorized, "user not found")
			return
		}
		if !user.IsActive {
			utils.JSONError(c, http.StatusUnauthorized, "user is inactive")
			return
		}
		c.Set(ContextUserKey, AuthUser{ID: user.ID, Username: user.Username, Role: user.Role})
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := GetAuthUser(c)
		if !ok || user.Role != utils.RoleAdmin {
			utils.JSONError(c, http.StatusForbidden, "admin access required")
			return
		}
		c.Next()
	}
}

func GetAuthUser(c *gin.Context) (AuthUser, bool) {
	value, ok := c.Get(ContextUserKey)
	if !ok {
		return AuthUser{}, false
	}
	user, ok := value.(AuthUser)
	return user, ok
}
