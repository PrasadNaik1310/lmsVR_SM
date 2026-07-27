package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"
	"github.com/PrasadNaik1310/LMSVR_SM/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCourseRoutes(api *gin.RouterGroup) {

	courses := api.Group("/courses")
	{

		courses.POST("", middleware.RequirePermission("course.create"), handlers.CreateCourse)
		courses.GET("", middleware.RequirePermission("course.read"), handlers.ListCoursesForUser)
		courses.GET("/:id", middleware.RequirePermission("course.read"), handlers.GetCourseDetails)
		courses.PUT("/:id", middleware.RequirePermission("course.update"), handlers.UpdateCourse)
		courses.PATCH("/:id/publish", middleware.RequirePermission("course.publish"), handlers.PublishCourse)
		courses.POST("/:id/invite", middleware.RequirePermission("course.invite"), handlers.GenerateCourseInvite)
	}
}
func RegisterTeacherRoutes(api *gin.RouterGroup) {
	teachers := api.Group("/teachers")
	{
		teachers.GET("", handlers.ListTeachers)
	}
}
