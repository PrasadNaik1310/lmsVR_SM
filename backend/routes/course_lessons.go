package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterCourseLessonRoutes(api *gin.RouterGroup) {

	modules := api.Group("/modules")
	{
		modules.POST("/:module_id/lessons", handlers.CreateLesson)
		modules.GET("/:module_id/lessons", handlers.ListLessons)
	}

	lessons := api.Group("/lessons")
	{
		lessons.GET("/:id", handlers.GetLessonDetails)
		lessons.PUT("/:id", handlers.UpdateLesson)
		lessons.DELETE("/:id", handlers.DeleteLesson)
	}
}
