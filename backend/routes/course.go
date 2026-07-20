package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterCourseRoutes(api *gin.RouterGroup) {

	courses := api.Group("/courses")
	{

		courses.POST("", handlers.CreateCourse)
		courses.GET("", handlers.ListCoursesForUser)
		courses.GET("/:id", handlers.GetCourseDetails)
		courses.PUT("/:id", handlers.UpdateCourse)
		courses.PATCH("/:id/publish", handlers.PublishCourse)
		courses.POST("/:id/invite", handlers.GenerateCourseInvite)
	}
}
func RegisterTeacherRoutes(api *gin.RouterGroup) {
	teachers := api.Group("/teachers")
	{
		teachers.GET("", handlers.ListTeachers)
	}
}
