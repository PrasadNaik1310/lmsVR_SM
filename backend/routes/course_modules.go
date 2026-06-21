package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterCourseModuleRoutes(api *gin.RouterGroup) {

	courses := api.Group("/courses")
	{
		courses.POST("/id/modules", handlers.CreateModule)
		courses.GET("/:id/modules", handlers.ListModules)
	}

	modules := api.Group("/modules")
	{
		modules.PUT("/:id", handlers.UpdateModule)
		modules.DELETE("/:id", handlers.DeleteModule)
	}
}
