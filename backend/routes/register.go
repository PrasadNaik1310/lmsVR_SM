package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/lms")
	{
		registerAuthRoutes(api)
		registerCompanyRoutes(api)
	}
}

func registerCompanyRoutes(api *gin.RouterGroup) {
	company := api.Group("/company")
	company.Use(middleware.AuthMiddleWare())
	{
		registerCompanySessionRoutes(company)
		registerCompanyBatchRoutes(company)
	}
}
