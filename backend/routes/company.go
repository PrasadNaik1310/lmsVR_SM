package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"
	"github.com/PrasadNaik1310/LMSVR_SM/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCompanyUserRoutes(api *gin.RouterGroup) {
	company := api.Group("/company")
	company.Use(middleware.AuthMiddleWare())

	company.POST("/users", middleware.RequirePermission("user.create"), handlers.CreateCompanyUser)
	company.GET("/users", middleware.RequirePermission("user.read"), handlers.ListCompanyUsers)
	// Sessions
	/*sessions := company.Group("/sessions")
	{
		sessions.PUT(":session_id/courses/:course_id/assign", middleware.RequirePermission("company.session.assign"), handlers.AssignCourseToSession)
		sessions.GET(":session_id/courses", middleware.RequirePermission("company.session.read"), handlers.ListCoursesBySession)
	}

	// Courses -> Batches
	courses := company.Group("/courses")
	{
		courses.POST(":course_id/batches", middleware.RequirePermission("company.batch.create"), handlers.CreateBatchForCourse)
		courses.GET(":course_id/batches/:batch_id", middleware.RequirePermission("company.batch.read"), handlers.GetBatchDetails)
	}*/
}
