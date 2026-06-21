package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("WARNING::PANIC TRIGERRED!!!!")
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC RECOVERED : %v ", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"Error": "Internal server error try again"})

			}
		}()
		c.Next()
	}
}
