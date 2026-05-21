package middleware

import (
	"log"
	"net/http"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequirePermission checks whether the authenticated user has the given permission name.
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// try multiple context keys that might be present
		var userIDAny interface{}
		var ok bool
		if userIDAny, ok = c.Get("user_id"); !ok {
			if userIDAny, ok = c.Get("userId"); !ok {
				if userIDAny, ok = c.Get("id"); !ok {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
					c.Abort()
					return
				}
			}
		}

		var user models.User

		switch v := userIDAny.(type) {
		case string:
			// try uuid
			if uid, err := uuid.Parse(v); err == nil {
				if err := db.DB.Where("id = ?", uid).First(&user).Error; err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
					c.Abort()
					return
				}
			} else {
				// fallback: try numeric id
				if err := db.DB.First(&user, v).Error; err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
					c.Abort()
					return
				}
			}
		case uuid.UUID:
			if err := db.DB.Where("id = ?", v).First(&user).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				c.Abort()
				return
			}
		default:
			// try numeric
			if err := db.DB.First(&user, v).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				c.Abort()
				return
			}
		}

		// load permission by name
		var perm models.Permission
		if err := db.DB.Where("name = ?", permission).First(&perm).Error; err != nil {
			log.Printf("permission lookup failed: %v", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "permission not found"})
			c.Abort()
			return
		}

		// check role_permission
		var rp models.RolePermission
		if err := db.DB.Where("role_id = ? AND permission_id = ?", user.RoleID, perm.ID).First(&rp).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not authorised."})
			c.Abort()
			return
		}

		c.Next()
	}
}
