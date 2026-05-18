package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		log.Printf("auath check for method %s , uri %s , client ip %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		log.Printf("Auth header present for client %v  , %d", authHeader != "", len(authHeader))

		if authHeader == "" {
			log.Printf("Auth Header missing for client %s", c.ClientIP())
			log.Println("Request failed from middleware")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth header missing"})
			c.Abort()
			return
		}
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			log.Printf("Invalid format for request %s from client ip %s", c.Request.Method, c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid format for auth headers"})
			c.Abort()
			return
		}
		tokenString := strings.TrimSpace(parts[1])
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "JWT missing from server"})
			log.Printf("JWT missing from env , coming from middleware")
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "token not valid or not parsed from client"})
			log.Printf("invalid token coming from client %s ", c.ClientIP())
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("Invalid token claims from %s", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		//c.Set("claims",claims)
		//c.Set("UserEmail",claims["user_email"])
		userIDRaw, ok := claims["user_id"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id missing"})
			c.Abort()
			return
		}

		userIDFloat, ok := userIDRaw.(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id type"})
			c.Abort()
			return
		}

		userID := uint(userIDFloat)
		c.Set("user_id", userID)
		c.Set("userId", userID)
		c.Set("id", userID)
		//c.Set("")
		c.Next()
	}
}
