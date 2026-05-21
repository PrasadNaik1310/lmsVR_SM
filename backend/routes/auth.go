package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
	}
}
