package routes

import (
	"github.com/PrasadNaik1310/LMSVR_SM/handlers"
	"github.com/PrasadNaik1310/LMSVR_SM/middleware"
	"github.com/gin-gonic/gin"
)

func registerCompanySessionRoutes(company *gin.RouterGroup) {
	sessions := company.Group("/sessions")
	{
		sessions.POST("", middleware.RequirePermission("company.session.create"), handlers.CreateAcademicSession)
		sessions.GET("", middleware.RequirePermission("company.session.read"), handlers.ListAcademicSessions)
		sessions.GET("/:session_id", middleware.RequirePermission("company.session.read"), handlers.GetAcademicSession)
		sessions.PUT("/:session_id", middleware.RequirePermission("company.session.update"), handlers.UpdateAcademicSession)
		sessions.DELETE("/:session_id", middleware.RequirePermission("company.session.delete"), handlers.DeleteAcademicSession)
		sessions.PUT(":session_id/courses/:course_id/assign", middleware.RequirePermission("company.session.assign"), handlers.AssignCourseToSession)
		sessions.GET(":session_id/courses", middleware.RequirePermission("company.session.read"), handlers.ListCoursesBySession)
	}
}
